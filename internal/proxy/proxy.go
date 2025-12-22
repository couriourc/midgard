package proxy

import (
	"bytes"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"context"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/midgard/gateway/internal/collection"
	"github.com/redis/go-redis/v9"
	"github.com/midgard/gateway/internal/database"
	"github.com/midgard/gateway/internal/health"
	"gorm.io/gorm"
)

// ProxyManager manages proxy requests
type ProxyManager struct {
	collectionManager *collection.CollectionManager
	healthChecker     *health.HealthChecker
	redisClient       *redis.Client
	db                *gorm.DB
	ctx               context.Context
}

// NewProxyManager creates a new proxy manager
func NewProxyManager(cm *collection.CollectionManager, hc *health.HealthChecker, redisClient *redis.Client, db *gorm.DB) *ProxyManager {
	return &ProxyManager{
		collectionManager: cm,
		healthChecker:     hc,
		redisClient:       redisClient,
		db:                db,
		ctx:               context.Background(),
	}
}

// HandleProxyRequest handles a proxy request
func (pm *ProxyManager) HandleProxyRequest(c *gin.Context) {
	prefix := c.Param("prefix")
	path := c.Param("path")
	
	// Remove leading slash from path if present
	path = strings.TrimPrefix(path, "/")

	// Get collection by prefix
	coll, err := pm.collectionManager.GetCollectionByPrefix(prefix)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Collection not found"})
		return
	}

	// Check if collection is active
	if !coll.Active {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "Collection is not active"})
		return
	}

	// Check health if configured
	if coll.HealthPath != "" {
		if !pm.healthChecker.IsHealthy(coll.ID) {
			c.JSON(http.StatusServiceUnavailable, gin.H{"error": "Service is unhealthy"})
			return
		}
	}

	// Build target URL
	targetURL := fmt.Sprintf("%s/%s", strings.TrimSuffix(coll.BaseURL, "/"), path)
	if c.Request.URL.RawQuery != "" {
		targetURL += "?" + c.Request.URL.RawQuery
	}

	// Capture request body and params before processing
	requestBody, _ := io.ReadAll(c.Request.Body)
	c.Request.Body = io.NopCloser(bytes.NewBuffer(requestBody))
	requestParamsJSON, _ := json.Marshal(c.Request.URL.Query())
	fromCache := false

	// Generate cache key early (before body is consumed by proxy)
	var cacheKey string
	if coll.CacheEnabled && pm.redisClient != nil {
		cacheKey = pm.generateCacheKey(coll.ID, c.Request, coll.CacheKeyStrategy, requestBody)
		cachedData, err := pm.redisClient.Get(pm.ctx, cacheKey).Result()
		if err != nil && err != redis.Nil {
			log.Printf("Cache GET error for key %s: %v", cacheKey, err)
		}
		if err == nil && len(cachedData) > 0 {
			// Serve from cache
			var cachedResponse map[string]interface{}
			if err := json.Unmarshal([]byte(cachedData), &cachedResponse); err == nil {
				// Set response headers if available
				if headers, ok := cachedResponse["headers"].(map[string]interface{}); ok {
					for k, v := range headers {
						if str, ok := v.(string); ok {
							c.Header(k, str)
						}
					}
				}
				// Set status code and write body
				status := http.StatusOK
				if statusVal, ok := cachedResponse["status"].(float64); ok {
					status = int(statusVal)
				}
				if body, ok := cachedResponse["body"].(string); ok {
					c.Header("X-Cache", "HIT")
					fromCache = true
					// Log cached request
					if coll.LogEnabled {
						pm.logRequest(coll.ID, path, c.Request.Method, targetURL, status, 0, len(requestBody), len(body), c.ClientIP(), c.Request.Header, c.Writer.Header(), coll.LogMaxEntries, coll.LogRolling, string(requestBody), string(requestParamsJSON), fromCache)
					}
					c.Data(status, "application/json", []byte(body))
					return
				}
			}
		}
	}

	// Start timer
	start := time.Now()

	// Create reverse proxy
	target, err := url.Parse(targetURL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Create proxy
	proxy := httputil.NewSingleHostReverseProxy(target)

	// Modify the request to use the correct path
	originalDirector := proxy.Director
	proxy.Director = func(req *http.Request) {
		originalDirector(req)
		// Set the correct path (remove /proxy/{prefix} prefix)
		req.URL.Path = "/" + path
		req.URL.RawPath = ""
		// Preserve query parameters
		if c.Request.URL.RawQuery != "" {
			req.URL.RawQuery = c.Request.URL.RawQuery
		}
	}

	// Capture response
	responseRecorder := &responseRecorder{
		ResponseWriter: c.Writer,
		body:           &bytes.Buffer{},
		status:         http.StatusOK,
	}

	// Serve the request
	proxy.ServeHTTP(responseRecorder, c.Request)

	// Calculate duration
	duration := time.Since(start).Milliseconds()

	// Log the request if enabled
	if coll.LogEnabled {
		pm.logRequest(coll.ID, path, c.Request.Method, targetURL, responseRecorder.status, duration, len(requestBody), responseRecorder.body.Len(), c.ClientIP(), c.Request.Header, responseRecorder.Header(), coll.LogMaxEntries, coll.LogRolling, string(requestBody), string(requestParamsJSON), fromCache)
	}

	// Cache the response if enabled
	if coll.CacheEnabled && pm.redisClient != nil && responseRecorder.status == http.StatusOK {
		// Use the same cache key generated earlier
		if cacheKey == "" {
			cacheKey = pm.generateCacheKey(coll.ID, c.Request, coll.CacheKeyStrategy, requestBody)
		}
		responseData := map[string]interface{}{
			"status":  responseRecorder.status,
			"body":    responseRecorder.body.String(),
			"headers": responseRecorder.Header(),
		}
		if data, err := json.Marshal(responseData); err == nil {
			err := pm.redisClient.Set(pm.ctx, cacheKey, string(data), time.Duration(coll.CacheTTL)*time.Second).Err()
			if err != nil {
				log.Printf("Failed to set cache for key %s: %v", cacheKey, err)
			}
		} else {
			log.Printf("Failed to marshal response data for caching: %v", err)
		}
	}
}

// generateCacheKey generates a cache key based on strategy
func (pm *ProxyManager) generateCacheKey(collectionID string, r *http.Request, strategy string, requestBody []byte) string {
	key := fmt.Sprintf("%s:%s:%s", collectionID, r.Method, r.URL.Path)

	// Add query params
	if strategy == "params" || strategy == "all" {
		key += fmt.Sprintf(":%s", r.URL.RawQuery)
	}

	// Add body hash (use the provided requestBody instead of reading from r.Body)
	if (strategy == "body" || strategy == "all") && len(requestBody) > 0 {
		hash := md5.Sum(requestBody)
		key += fmt.Sprintf(":%x", hash)
	}

	return key
}

// logRequest logs a request to the database
func (pm *ProxyManager) logRequest(collectionID, path, method, targetURL string, status int, duration int64, requestSize, responseSize int, clientIP string, requestHeaders, responseHeaders http.Header, maxEntries int, rolling bool, requestBody, requestParams string, fromCache bool) {
	reqHeadersJSON, _ := json.Marshal(requestHeaders)
	respHeadersJSON, _ := json.Marshal(responseHeaders)

	log := database.RequestLog{
		CollectionID:    collectionID,
		Path:            path,
		Method:          method,
		TargetURL:       targetURL,
		Status:          status,
		Duration:        duration,
		RequestSize:     requestSize,
		ResponseSize:    responseSize,
		ClientIP:        clientIP,
		RequestHeaders:  string(reqHeadersJSON),
		ResponseHeaders: string(respHeadersJSON),
		RequestBody:     requestBody,
		RequestParams:   requestParams,
		FromCache:       fromCache,
		Timestamp:       time.Now(),
}

	// Insert log
	pm.db.Create(&log)

	// Handle rolling logs
	if rolling {
		var count int64
		pm.db.Model(&database.RequestLog{}).Where("collection_id = ?", collectionID).Count(&count)
		if count > int64(maxEntries) {
			// Delete oldest logs in batch to reduce database locks
			// Use a subquery to find the IDs of logs to delete, then delete them in one operation
			excessCount := count - int64(maxEntries)
			pm.db.Exec(`
				DELETE FROM request_logs 
				WHERE id IN (
					SELECT id FROM request_logs 
					WHERE collection_id = ? 
					ORDER BY timestamp ASC 
					LIMIT ?
				)
			`, collectionID, excessCount)
		}
	}
}

// responseRecorder captures the response for logging and caching
type responseRecorder struct {
	http.ResponseWriter
	body   *bytes.Buffer
	status int
}

func (r *responseRecorder) WriteHeader(status int) {
	r.status = status
	r.ResponseWriter.WriteHeader(status)
}

func (r *responseRecorder) Write(b []byte) (int, error) {
	r.body.Write(b)
	return r.ResponseWriter.Write(b)
}
