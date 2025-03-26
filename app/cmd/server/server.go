package server

import (
	"fmt"
	"log"
	"swift-app/cmd/router"
	"swift-app/database"
	"swift-app/internal/services"

	"github.com/gin-gonic/gin"
)

func StartServer() {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	swiftService := services.NewSwiftCodeService(database.GetCollection())

	router.SetupRoutes(r, swiftService)

	port := "8080"
	fmt.Printf("Server running on http://localhost:%s\n", port)

	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
