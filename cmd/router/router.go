package router

import (
	v1 "swift-app/api/v1"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	api := r.Group("/v1/swift-codes")
	{
		api.GET(":swift-code", v1.GetSwiftCode)
		api.GET("/country/:countryISO2code", v1.GetSwiftCodesByCountry)
	}
}
