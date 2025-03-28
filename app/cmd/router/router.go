package router

import (
	v1 "swift-app/api/v1"
	"swift-app/internal/errors"
	"swift-app/internal/models"
	"swift-app/internal/services"

	"github.com/gin-gonic/gin"
)

// SetupRoutes registers SWIFT code API endpoints with their corresponding handlers:
// - GET "/:swift-code": Fetch details of a specific SWIFT code.
// - GET "/country/:countryISO2code": List all SWIFT codes for a given country.
// - POST "/": Add a new SWIFT code.
// - DELETE "/:swift-code": Delete an existing SWIFT code.

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
		err := errors.Wrap(errors.ErrNotFound, "endpoint not found: %s. Please try again", c.Request.URL.Path)
		c.JSON(errors.GetStatusCode(err), models.MessageResponse{Message: err.Error()})
	})
}
