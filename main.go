package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/midgard/gateway/config"
	"github.com/redis/go-redis/v9"
	"github.com/midgard/gateway/internal/api"
	"github.com/midgard/gateway/internal/collection"
	"github.com/midgard/gateway/internal/database"
	"github.com/midgard/gateway/internal/health"
	"github.com/midgard/gateway/internal/proxy"
)

func main() {
	log.Println("Starting Midgard Gateway...")

	// Load configuration
	cfg := config.LoadConfig()

	// Initialize database
	db, err := database.InitDatabase(&cfg.Database)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// Initialize Redis client
	var redisClient *redis.Client
	if cfg.Redis.Host != "" {
		redisClient = redis.NewClient(&redis.Options{
			Addr:     fmt.Sprintf("%s:%d", cfg.Redis.Host, cfg.Redis.Port),
			Password: cfg.Redis.Password,
			DB:       cfg.Redis.DB,
		})
		log.Printf("Redis client initialized for %s:%d", cfg.Redis.Host, cfg.Redis.Port)
		
		// Test connection
		ctx := context.Background()
		testKey := "midgard:test:connection"
		testValue := "test"
		err := redisClient.Set(ctx, testKey, testValue, 10*time.Second).Err()
		if err != nil {
			log.Printf("Warning: Failed to set test key in Redis: %v", err)
		} else {
			val, err := redisClient.Get(ctx, testKey).Result()
			if err == nil && val == testValue {
				log.Printf("Redis connection test successful")
				// Test key will expire automatically after 10 seconds
			} else {
				log.Printf("Warning: Redis connection test failed: %v", err)
			}
		}
	} else {
		log.Println("Redis not configured (host is empty)")
	}

	// Initialize collection manager
	collectionManager := collection.NewCollectionManager(db)

	// Initialize health checker
	healthChecker := health.NewHealthChecker()

	// Start health checks for existing collections
	collections, err := collectionManager.GetAllCollections()
	if err == nil {
		for _, coll := range collections {
			if coll.HealthPath != "" {
				healthChecker.StartHealthCheck(&coll)
			}
		}
	}

	// Initialize proxy manager
	proxyManager := proxy.NewProxyManager(collectionManager, healthChecker, redisClient, db)

	// Check if frontend is enabled (from environment variable or config)
	enableFrontend := cfg.EnableFrontend
	if envFrontend := os.Getenv("ENABLE_FRONTEND"); envFrontend != "" {
		enableFrontend = envFrontend == "true" || envFrontend == "1"
	}

	// Initialize API server
	apiServer := api.NewAPIServer(collectionManager, proxyManager, healthChecker, db, enableFrontend)
	
	if enableFrontend {
		log.Println("Frontend is enabled")
	} else {
		log.Println("Frontend is disabled (API-only mode)")
	}

	// Create HTTP server
	port := cfg.Server.Port
	if envPort := os.Getenv("PORT"); envPort != "" {
		fmt.Sscanf(envPort, "%d", &port)
	}

	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: apiServer.RegisterRoutes(),
	}

	log.Printf("Midgard Gateway started on :%d", port)
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
