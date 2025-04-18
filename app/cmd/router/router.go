// router.go defines and registers API routes for the SWIFT code service.
// It connects the HTTP endpoints to their corresponding handler functions,
// sets up the versioned API group, and configures a fallback for unknown routes.
package router

import (
	v1 "swift-app/api/v1"
	"swift-app/internal/errors"
	"swift-app/internal/models"
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
		err := errors.Wrap(errors.ErrNotFound, "endpoint not found: %s. Please try again", c.Request.URL.Path)
		c.JSON(errors.GetStatusCode(err), models.MessageResponse{Message: err.Error()})
	})
}
