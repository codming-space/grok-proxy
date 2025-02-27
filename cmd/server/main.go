package main

import (
	"fmt"
	"grok-proxy/config"
	"grok-proxy/internal/api"
	"grok-proxy/internal/client"
	"grok-proxy/internal/cookie"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

func init() {
	// Seed the random number generator
	rand.Seed(time.Now().UnixNano())
}

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize cookie manager
	cookieManager, err := cookie.NewManager()
	if err != nil {
		log.Fatalf("Failed to create cookie manager: %v", err)
	}

	// Print current cookie info
	log.Printf("Current cookie: %d / %d",
		cookieManager.CurrentCookieIndex(),
		cookieManager.CookieCount())

	// Initialize Grok client
	grokClient := client.NewGrokClient(cookieManager)

	// Initialize handler
	handler := api.NewHandler(grokClient, cookieManager, cfg)

	// Set up Gin router
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	// Register routes
	handler.RegisterRoutes(router)

	// Get port from environment or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	// Start the server
	log.Printf("Starting server on port %s", port)
	if err := router.Run(fmt.Sprintf("0.0.0.0:%s", port)); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
