package main

import (
	"log"
	"runrun/config"
	"runrun/internal"
	"runrun/internal/scheduler"
	"runrun/routes"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// Load application configuration from resource/application.yaml
	if err := config.Init(); err != nil {
		log.Fatalf("FATAL: Failed to initialize configuration: %v", err)
	}

	// Initialize Database
	internal.InitDB()

	// Initialize and start scheduler
	sched := scheduler.NewScheduler()
	sched.Start()

	// Initialize Gin router
	r := gin.Default()

	// Configure CORS middleware
	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{
			"http://localhost:3000",           // 本地开发
			"https://runrun.1e27.net",        // 你的Vercel域名
			"https://*.vercel.app",           // 所有Vercel预览域名
		},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// Initialize API routes
	routes.InitApiRouter(r)

	// Get port from config, with a fallback to 23450
	port := config.GetString("server.port", "23450")

	log.Printf("Server starting on port %s", port)

	// Start the server
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("FATAL: Failed to start server: %v", err)
	}
}
