package health

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/midgard/gateway/internal/database"
)

// HealthChecker manages health checks for collections
type HealthChecker struct {
	checks map[string]*HealthCheck
	mu     sync.RWMutex
}

// HealthCheck represents a health check for a collection
type HealthCheck struct {
	CollectionID string
	URL          string
	Interval     time.Duration
	LastCheck    time.Time
	IsHealthy    bool
	stopChan     chan struct{}
}

// NewHealthChecker creates a new health checker
func NewHealthChecker() *HealthChecker {
	return &HealthChecker{
		checks: make(map[string]*HealthCheck),
	}
}

// StartHealthCheck starts health checking for a collection
func (hc *HealthChecker) StartHealthCheck(coll *database.Collection) {
	hc.mu.Lock()
	defer hc.mu.Unlock()

	// Stop existing check if any
	if existing, exists := hc.checks[coll.ID]; exists {
		select {
		case <-existing.stopChan:
			// Already closed
		default:
			close(existing.stopChan)
		}
		delete(hc.checks, coll.ID)
	}

	if coll.HealthPath == "" {
		return
	}

	healthURL := fmt.Sprintf("%s%s", coll.BaseURL, coll.HealthPath)
	interval := time.Duration(coll.HealthInterval) * time.Second
	if interval == 0 {
		interval = 30 * time.Second
	}

	check := &HealthCheck{
		CollectionID: coll.ID,
		URL:          healthURL,
		Interval:     interval,
		IsHealthy:    true,
		stopChan:     make(chan struct{}),
	}

	hc.checks[coll.ID] = check

	// Start checking
	go hc.runHealthCheck(check)
}

// StopHealthCheck stops health checking for a collection
func (hc *HealthChecker) StopHealthCheck(collectionID string) {
	hc.mu.Lock()
	defer hc.mu.Unlock()

	if check, exists := hc.checks[collectionID]; exists {
		select {
		case <-check.stopChan:
			// Already closed
		default:
			close(check.stopChan)
		}
		delete(hc.checks, collectionID)
	}
}

// IsHealthy checks if a collection is healthy
func (hc *HealthChecker) IsHealthy(collectionID string) bool {
	hc.mu.RLock()
	defer hc.mu.RUnlock()

	if check, exists := hc.checks[collectionID]; exists {
		return check.IsHealthy
	}
	return true // Default to healthy if no check configured
}

// runHealthCheck runs the health check loop
func (hc *HealthChecker) runHealthCheck(check *HealthCheck) {
	ticker := time.NewTicker(check.Interval)
	defer ticker.Stop()

	// Initial check
	check.IsHealthy = hc.performCheck(check.URL)
	check.LastCheck = time.Now()

	for {
		select {
		case <-ticker.C:
			check.IsHealthy = hc.performCheck(check.URL)
			check.LastCheck = time.Now()
		case <-check.stopChan:
			return
		}
	}
}

// performCheck performs a single health check
func (hc *HealthChecker) performCheck(url string) bool {
	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	resp, err := client.Get(url)
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	return resp.StatusCode >= 200 && resp.StatusCode < 300
}

