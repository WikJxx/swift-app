package router

import (
	"net/http"
	v1 "swift-app/api/v1"
	"swift-app/internal/services"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine, swiftService *services.SwiftCodeService) {
	api := r.Group("/v1/swift-codes")
	{
		api.GET("/:swift-code", func(c *gin.Context) {
			v1.GetSwiftCode(c, swiftService)
		})

		api.GET("/country/:countryISO2code", func(c *gin.Context) {
			v1.GetSwiftCodesByCountry(c, swiftService)
		})

		api.POST("/", func(c *gin.Context) {
			v1.AddSwiftCode(c, swiftService)
		})

		api.DELETE("/:swift-code", func(c *gin.Context) {
			v1.DeleteSwiftCode(c, swiftService)
		})
	}

	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{"error": "Endpoint not found. Please check the URL and try again."})
	})
}
