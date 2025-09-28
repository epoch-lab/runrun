package routes

import (
	"net/http"
	"runrun/internal/handler"

	"github.com/gin-gonic/gin"
)

// InitApiRouter initializes the API routes for the application.
func InitApiRouter(r *gin.Engine) {
	apiRouter := r.Group("/api")
	{
		// A simple health check endpoint
		apiRouter.GET("/ping", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"message": "pong",
			})
		})

		// Route to trigger the run process
		apiRouter.POST("/run", handler.RunHandler)
		
		// Authentication endpoint for user registration and login
		apiRouter.POST("/auth", handler.AuthHandler)
	}
}