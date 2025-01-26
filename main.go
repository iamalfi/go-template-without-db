package main

import (
	"gin-project/middleware"
	"gin-project/routes"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}
}

func main() {

	// Set Gin mode to release
	// gin.SetMode(gin.ReleaseMode)

	r := gin.Default()

	r.Use(middleware.ErrorHandler)

	api := r.Group("/api")
	routes.Routes(api)
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Welcome to the Gin server!",
		})
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
