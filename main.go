package main

import (
	"github.com/gin-gonic/gin"
	"log"
	"runrun/config"
	"runrun/routes"
)

func main() {
	// Load application configuration from resource/application.yaml
	if err := config.Init(); err != nil {
		log.Fatalf("FATAL: Failed to initialize configuration: %v", err)
	}

	// Initialize Gin router
	r := gin.Default()

	// Initialize API routes
	routes.InitApiRouter(r)

	// Get port from config, with a fallback to 8080
	port := config.GetString("server.port", "8080")

	log.Printf("Server starting on port %s", port)

	// Start the server
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("FATAL: Failed to start server: %v", err)
	}
}
