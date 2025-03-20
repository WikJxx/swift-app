// router.go
package router

import (
	v1 "swift-app/api/v1" // Importowanie pakietu v1
	"swift-app/internal/services"

	"github.com/gin-gonic/gin"
)

func handleError(c *gin.Context, err error, statusCode int) {
	c.JSON(statusCode, gin.H{"error": err.Error()})
}

func SetupRoutes(r *gin.Engine, swiftService *services.SwiftCodeService) {
	api := r.Group("/v1/swift-codes")
	{
		// Użycie funkcji GetSwiftCode z pakietu v1
		api.GET("/:swift-code", func(c *gin.Context) {
			v1.GetSwiftCode(c, swiftService) // Wywołanie GetSwiftCode z v1
		})

		// Użycie funkcji GetSwiftCodesByCountry z pakietu v1
		api.GET("/country/:countryISO2code", func(c *gin.Context) {
			v1.GetSwiftCodesByCountry(c, swiftService) // Wywołanie GetSwiftCodesByCountry z v1
		})

		// Użycie funkcji AddSwiftCode z pakietu v1
		api.POST("/", func(c *gin.Context) {
			v1.AddSwiftCode(c, swiftService) // Wywołanie AddSwiftCode z v1
		})
	}
}
