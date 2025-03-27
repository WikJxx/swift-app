package server

import (
	"fmt"
	"log"
	"os"
	"swift-app/cmd/router"
	"swift-app/database"
	"swift-app/internal/services"

	"github.com/gin-gonic/gin"
)

// Function initializes and runs the HTTP server, setting up routes and services for handling SWIFT code API requests.

func StartServer() {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	swiftService := services.NewSwiftCodeService(database.GetCollection())

	router.SetupRoutes(r, swiftService)

	host := os.Getenv("HOST")
	port := os.Getenv("PORT")

	if port == "" {
		port = "8080"
	}
	if host == "" {
		host = "localhost"
	}

	address := fmt.Sprintf("%s:%s", host, port)
	fmt.Printf("Server running on http://%s\n", address)

	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
