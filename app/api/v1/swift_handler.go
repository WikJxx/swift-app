package v1

import (
	"net/http"
	"strings"
	"swift-app/internal/errors"
	"swift-app/internal/models"
	"swift-app/internal/services"

	"github.com/gin-gonic/gin"
)

func GetSwiftCode(c *gin.Context, swiftService *services.SwiftCodeService) {
	swiftCode := strings.ToUpper(c.Param("swift-code"))

	swift, err := swiftService.GetSwiftCodeDetails(swiftCode)
	if err != nil {
		status := errors.MapToStatusCode(err)
		c.JSON(status, models.MessageResponse{Message: err.Error()})
		return
	}

	response := models.SwiftCode{
		Address:       swift.Address,
		BankName:      swift.BankName,
		CountryISO2:   swift.CountryISO2,
		CountryName:   swift.CountryName,
		IsHeadquarter: swift.IsHeadquarter,
		SwiftCode:     swift.SwiftCode,
	}

	if swift.IsHeadquarter {
		response.Branches = swift.Branches
	}

	c.JSON(http.StatusOK, response)
}

func GetSwiftCodesByCountry(c *gin.Context, swiftService *services.SwiftCodeService) {
	countryISO2 := strings.ToUpper(c.Param("countryISO2code"))

	swiftCodesResponse, err := swiftService.GetSwiftCodesByCountry(countryISO2)
	if err != nil {
		status := errors.MapToStatusCode(err)
		c.JSON(status, models.MessageResponse{Message: err.Error()})
		return
	}

	response := models.CountrySwiftCodesResponse{
		CountryISO2: countryISO2,
		CountryName: swiftCodesResponse.CountryName,
		SwiftCodes:  swiftCodesResponse.SwiftCodes,
	}

	c.JSON(http.StatusOK, response)
}

func AddSwiftCode(c *gin.Context, swiftService *services.SwiftCodeService) {
	var swiftCodeRequest models.SwiftCode
	if err := c.ShouldBindJSON(&swiftCodeRequest); err != nil {
		c.JSON(http.StatusBadRequest, models.MessageResponse{Message: "Invalid input data"})
		return
	}

	message, err := swiftService.AddSwiftCode(&swiftCodeRequest)
	if err != nil {
		status := errors.MapToStatusCode(err)
		c.JSON(status, models.MessageResponse{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, models.MessageResponse{Message: message})
}

func DeleteSwiftCode(c *gin.Context, swiftService *services.SwiftCodeService) {
	swiftCode := strings.ToUpper(c.Param("swift-code"))

	message, err := swiftService.DeleteSwiftCode(swiftCode)
	if err != nil {
		status := errors.MapToStatusCode(err)
		c.JSON(status, models.MessageResponse{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, models.MessageResponse{Message: message})
}
