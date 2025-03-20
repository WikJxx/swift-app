package router

import (
	"net/http"
	"swift-app/internal/services"

	"github.com/gin-gonic/gin"
)

func handleError(c *gin.Context, err error, statusCode int) {
	c.JSON(statusCode, gin.H{"error": err.Error()})
}

func SetupRoutes(r *gin.Engine, swiftService *services.SwiftCodeService) {
	v1 := r.Group("/v1/swift-codes")
	{
		v1.GET("/:swift-code", func(c *gin.Context) {
			swiftCode := c.Param("swift-code")
			swiftCodeDetails, err := swiftService.GetSwiftCodeDetails(swiftCode)
			if err != nil {
				handleError(c, err, http.StatusNotFound)
				return
			}
			c.JSON(http.StatusOK, gin.H{
				"swiftCode": swiftCodeDetails,
			})
		})

		v1.GET("/country/:countryISO2code", func(c *gin.Context) {
			countryISO2 := c.Param("countryISO2code")
			swiftCodesResponse, err := swiftService.GetSwiftCodesByCountry(countryISO2)
			if err != nil {
				handleError(c, err, http.StatusNotFound)
				return
			}

			c.JSON(http.StatusOK, swiftCodesResponse)
		})
	}
}
