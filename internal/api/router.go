package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/midgard/gateway/internal/collection"
	"github.com/midgard/gateway/internal/database"
	"github.com/midgard/gateway/internal/health"
	"github.com/midgard/gateway/internal/proxy"
	"gorm.io/gorm"
)

// APIServer handles API routes
type APIServer struct {
	collectionManager *collection.CollectionManager
	proxyManager      *proxy.ProxyManager
	healthChecker     *health.HealthChecker
	db                *gorm.DB
	enableFrontend    bool
}

// NewAPIServer creates a new API server
func NewAPIServer(cm *collection.CollectionManager, pm *proxy.ProxyManager, hc *health.HealthChecker, db *gorm.DB, enableFrontend bool) *APIServer {
	return &APIServer{
		collectionManager: cm,
		proxyManager:      pm,
		healthChecker:     hc,
		db:                db,
		enableFrontend:    enableFrontend,
	}
}

// RegisterRoutes registers all routes
func (s *APIServer) RegisterRoutes() http.Handler {
	router := gin.Default()

	// Configure CORS
	router.Use(cors.Default())

	// Serve static files (frontend) - only if enabled
	if s.enableFrontend {
		if _, err := os.Stat("./web/dist"); err == nil {
			router.Use(static.Serve("/", static.LocalFile("./web/dist", false)))
			router.NoRoute(func(c *gin.Context) {
				// Only serve frontend for non-API routes
				if !strings.HasPrefix(c.Request.URL.Path, "/api") {
					c.File("./web/dist/index.html")
				} else {
					c.JSON(http.StatusNotFound, gin.H{"error": "Not found"})
				}
			})
		}
	}

	// API routes
	api := router.Group("/api")
	{
		// Collection management
		api.GET("/collections", s.handleGetCollections)
		api.POST("/collections", s.handleCreateCollection)
		api.GET("/collections/check-prefix/:prefix", s.handleCheckPrefix) // Must be before /:id route
		api.GET("/collections/:id", s.handleGetCollection)
		api.PUT("/collections/:id", s.handleUpdateCollection)
		api.DELETE("/collections/:id", s.handleDeleteCollection)
		api.POST("/collections/:id/toggle", s.handleToggleCollection)
		api.POST("/collections/:id/import-openapi", s.handleImportOpenAPI)

		// Logs
		api.GET("/logs", s.handleGetLogs)
		api.GET("/logs/:collectionId", s.handleGetCollectionLogs)
		api.GET("/logs/:collectionId/latest", s.handleGetLatestLog)
		api.DELETE("/logs", s.handleClearLogs)
		api.DELETE("/logs/:collectionId", s.handleClearCollectionLogs)

		// Statistics
		api.GET("/collections/:id/endpoint-stats", s.handleGetEndpointStats)
	}

	// Proxy routes - using prefix instead of collectionID
	router.Any("/proxy/:prefix/*path", func(c *gin.Context) {
		s.proxyManager.HandleProxyRequest(c)
	})

	// Health check
	router.GET("/health", s.handleHealthCheck)

	return router
}

func (s *APIServer) handleGetCollections(c *gin.Context) {
	collections, err := s.collectionManager.GetAllCollections()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, collections)
}

func (s *APIServer) handleCreateCollection(c *gin.Context) {
	var coll struct {
		Name             string `json:"name" binding:"required"`
		Description      string `json:"description"`
		Prefix           string `json:"prefix" binding:"required"`
		BaseURL          string `json:"base_url" binding:"required"`
		OpenAPIURL       string `json:"openapi_url"`
		HealthPath       string `json:"health_path"`
		HealthInterval   int    `json:"health_interval"`
		LogEnabled       bool   `json:"log_enabled"`
		LogRolling       bool   `json:"log_rolling"`
		LogMaxEntries    int    `json:"log_max_entries"`
		CacheEnabled     bool   `json:"cache_enabled"`
		CacheTTL         int    `json:"cache_ttl"`
		CacheKeyStrategy string `json:"cache_key_strategy"`
	}

	if err := c.ShouldBindJSON(&coll); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	dbColl := &database.Collection{
		Name:             coll.Name,
		Description:      coll.Description,
		Prefix:           coll.Prefix,
		BaseURL:          coll.BaseURL,
		OpenAPIURL:       coll.OpenAPIURL,
		HealthPath:       coll.HealthPath,
		HealthInterval:   coll.HealthInterval,
		LogEnabled:       coll.LogEnabled,
		LogRolling:       coll.LogRolling,
		LogMaxEntries:    coll.LogMaxEntries,
		CacheEnabled:     coll.CacheEnabled,
		CacheTTL:         coll.CacheTTL,
		CacheKeyStrategy: coll.CacheKeyStrategy,
		Active:           true,
	}

	// Check if prefix already exists
	exists, err := s.collectionManager.CheckPrefixExists(coll.Prefix, "")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if exists {
		c.JSON(http.StatusConflict, gin.H{"error": fmt.Sprintf("Prefix '%s' already exists", coll.Prefix)})
		return
	}

	if err := s.collectionManager.CreateCollection(dbColl); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Start health check if configured
	if dbColl.HealthPath != "" {
		s.healthChecker.StartHealthCheck(dbColl)
	}

	c.JSON(http.StatusCreated, dbColl)
}

func (s *APIServer) handleGetCollection(c *gin.Context) {
	id := c.Param("id")
	coll, err := s.collectionManager.GetCollection(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Collection not found"})
		return
	}
	c.JSON(http.StatusOK, coll)
}

func (s *APIServer) handleUpdateCollection(c *gin.Context) {
	id := c.Param("id")
	var coll struct {
		Name             string `json:"name"`
		Description      string `json:"description"`
		Prefix           string `json:"prefix"`
		BaseURL          string `json:"base_url"`
		OpenAPIURL       string `json:"openapi_url"`
		HealthPath       string `json:"health_path"`
		HealthInterval   int    `json:"health_interval"`
		LogEnabled       bool   `json:"log_enabled"`
		LogRolling       bool   `json:"log_rolling"`
		LogMaxEntries    int    `json:"log_max_entries"`
		CacheEnabled     bool   `json:"cache_enabled"`
		CacheTTL         int    `json:"cache_ttl"`
		CacheKeyStrategy string `json:"cache_key_strategy"`
		Active           bool   `json:"active"`
	}

	if err := c.ShouldBindJSON(&coll); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get existing collection
	existing, err := s.collectionManager.GetCollection(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Collection not found"})
		return
	}

	// Update fields
	if coll.Name != "" {
		existing.Name = coll.Name
	}
	if coll.Description != "" {
		existing.Description = coll.Description
	}
	if coll.Prefix != "" {
		// Check if new prefix already exists (excluding current collection)
		exists, err := s.collectionManager.CheckPrefixExists(coll.Prefix, id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if exists {
			c.JSON(http.StatusConflict, gin.H{"error": fmt.Sprintf("Prefix '%s' already exists", coll.Prefix)})
			return
		}
		existing.Prefix = coll.Prefix
	}
	if coll.BaseURL != "" {
		existing.BaseURL = coll.BaseURL
	}
	if coll.OpenAPIURL != "" {
		existing.OpenAPIURL = coll.OpenAPIURL
	}
	existing.HealthPath = coll.HealthPath
	existing.HealthInterval = coll.HealthInterval
	existing.LogEnabled = coll.LogEnabled
	existing.LogRolling = coll.LogRolling
	existing.LogMaxEntries = coll.LogMaxEntries
	existing.CacheEnabled = coll.CacheEnabled
	existing.CacheTTL = coll.CacheTTL
	existing.CacheKeyStrategy = coll.CacheKeyStrategy
	existing.Active = coll.Active

	if err := s.collectionManager.UpdateCollection(id, existing); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Restart health check if configured
	if existing.HealthPath != "" {
		s.healthChecker.StartHealthCheck(existing)
	} else {
		s.healthChecker.StopHealthCheck(existing.ID)
	}

	c.JSON(http.StatusOK, existing)
}

func (s *APIServer) handleDeleteCollection(c *gin.Context) {
	id := c.Param("id")

	// Stop health check
	s.healthChecker.StopHealthCheck(id)

	if err := s.collectionManager.DeleteCollection(id); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Collection not found"})
		return
	}
	c.Status(http.StatusNoContent)
}

func (s *APIServer) handleToggleCollection(c *gin.Context) {
	id := c.Param("id")
	if err := s.collectionManager.ToggleCollection(id); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Collection not found"})
		return
	}
	coll, _ := s.collectionManager.GetCollection(id)
	c.JSON(http.StatusOK, coll)
}

func (s *APIServer) handleImportOpenAPI(c *gin.Context) {
	id := c.Param("id")

	var req struct {
		OpenAPIURL  string          `json:"openapi_url"`
		OpenAPIJSON json.RawMessage `json:"openapi_json"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var openAPIJSON []byte
	if req.OpenAPIJSON != nil {
		openAPIJSON = []byte(req.OpenAPIJSON)
	}

	if err := s.collectionManager.ImportOpenAPI(id, req.OpenAPIURL, openAPIJSON); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	coll, _ := s.collectionManager.GetCollection(id)
	c.JSON(http.StatusOK, coll)
}

func (s *APIServer) handleHealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func (s *APIServer) handleCheckPrefix(c *gin.Context) {
	prefix := c.Param("prefix")
	excludeID := c.Query("exclude_id")
	
	exists, err := s.collectionManager.CheckPrefixExists(prefix, excludeID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"exists": exists})
}

func (s *APIServer) handleGetLogs(c *gin.Context) {
	var logs []database.RequestLog
	var total int64

	// Pagination parameters
	page := c.DefaultQuery("page", "1")
	pageSize := c.DefaultQuery("pageSize", "5")
	pageInt := 1
	pageSizeInt := 5
	fmt.Sscanf(page, "%d", &pageInt)
	fmt.Sscanf(pageSize, "%d", &pageSizeInt)

	// Filter parameters
	pathFilter := c.Query("path")
	collectionIDFilter := c.Query("collection_id")

	// Build query
	query := s.db.Model(&database.RequestLog{})

	// Apply filters
	if pathFilter != "" {
		query = query.Where("path LIKE ?", "%"+pathFilter+"%")
	}
	if collectionIDFilter != "" {
		query = query.Where("collection_id = ?", collectionIDFilter)
	}

	// Get total count
	if err := query.Count(&total).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Apply pagination and ordering
	offset := (pageInt - 1) * pageSizeInt
	if err := query.Order("timestamp DESC").Offset(offset).Limit(pageSizeInt).Find(&logs).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return paginated response
	c.JSON(http.StatusOK, gin.H{
		"data":  logs,
		"total": total,
		"page":  pageInt,
		"pageSize": pageSizeInt,
	})
}

func (s *APIServer) handleGetCollectionLogs(c *gin.Context) {
	collectionID := c.Param("collectionId")
	var logs []database.RequestLog
	limit := c.DefaultQuery("limit", "100")

	query := s.db.Where("collection_id = ?", collectionID).Order("timestamp DESC")
	if limit != "" {
		query = query.Limit(100)
	}

	if err := query.Find(&logs).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, logs)
}

func (s *APIServer) handleGetLatestLog(c *gin.Context) {
	collectionID := c.Param("collectionId")
	var log database.RequestLog

	if err := s.db.Where("collection_id = ?", collectionID).Order("timestamp DESC").First(&log).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusOK, nil)
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, log)
}

func (s *APIServer) handleClearLogs(c *gin.Context) {
	if err := s.db.Exec("DELETE FROM request_logs").Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

func (s *APIServer) handleClearCollectionLogs(c *gin.Context) {
	collectionID := c.Param("collectionId")
	if err := s.db.Where("collection_id = ?", collectionID).Delete(&database.RequestLog{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

func (s *APIServer) handleGetEndpointStats(c *gin.Context) {
	collectionID := c.Param("id")

	type EndpointStat struct {
		Path         string  `json:"path"`
		Method       string  `json:"method"`
		RequestCount int64   `json:"request_count"`
		AvgDuration  float64 `json:"avg_duration"`
	}

	var stats []EndpointStat

	// Group by path and method, calculate count and average duration
	if err := s.db.Model(&database.RequestLog{}).
		Select("path, method, COUNT(*) as request_count, AVG(duration) as avg_duration").
		Where("collection_id = ?", collectionID).
		Group("path, method").
		Order("request_count DESC").
		Scan(&stats).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, stats)
}
