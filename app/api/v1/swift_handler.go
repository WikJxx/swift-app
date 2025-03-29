package v1

import (
	"net/http"
	"strings"
	"swift-app/internal/errors"
	"swift-app/internal/models"
	"swift-app/internal/services"
	"swift-app/internal/utils"

	"github.com/gin-gonic/gin"
)

// GetSwiftCode godoc
// @Summary Get SWIFT code
// @Description Returns a SWIFT code by its identifier (headquarter)
// @Tags SWIFT Codes
// @Accept json
// @Produce json
// @Param swift-code path string true "SWIFT code"
// @Success 200 {object} models.SwiftCode
// @Failure 404 {object} models.MessageResponse
// @Router /v1/swift-codes/{swift-code} [get]
func GetSwiftCode(c *gin.Context, swiftService *services.SwiftCodeService) {
	swiftCode := strings.ToUpper(c.Param(utils.ParamSwiftCode))

	swift, err := swiftService.GetSwiftCodeDetails(swiftCode)
	if err != nil {
		c.JSON(errors.GetStatusCode(err), models.MessageResponse{Message: err.Error()})
		return
	}

	if swift.IsHeadquarter {
		c.JSON(http.StatusOK, models.SwiftCode{
			Address:       swift.Address,
			BankName:      swift.BankName,
			CountryISO2:   swift.CountryISO2,
			CountryName:   swift.CountryName,
			IsHeadquarter: true,
			SwiftCode:     swift.SwiftCode,
			Branches:      swift.Branches,
		})
		return
	}

	c.JSON(http.StatusOK, models.SwiftBranch{
		Address:       swift.Address,
		BankName:      swift.BankName,
		CountryISO2:   swift.CountryISO2,
		CountryName:   swift.CountryName,
		IsHeadquarter: false,
		SwiftCode:     swift.SwiftCode,
	})
}

// GetSwiftCodesByCountry godoc
// @Summary Get SWIFT codes by country
// @Description Returns a list of SWIFT codes for a given country ISO2 code
// @Tags SWIFT Codes
// @Accept json
// @Produce json
// @Param countryISO2code path string true "Country ISO2 code"
// @Success 200 {array} models.SwiftCode
// @Failure 404 {object} models.MessageResponse
// @Router /v1/swift-codes/country/{countryISO2code} [get]
func GetSwiftCodesByCountry(c *gin.Context, swiftService *services.SwiftCodeService) {
	countryISO2 := strings.ToUpper(c.Param(utils.ParamCountryISO2))

	swiftCodesResponse, err := swiftService.GetSwiftCodesByCountry(countryISO2)
	if err != nil {
		c.JSON(errors.GetStatusCode(err), models.MessageResponse{Message: err.Error()})
		return
	}

	response := models.CountrySwiftCodesResponse{
		CountryISO2: countryISO2,
		CountryName: swiftCodesResponse.CountryName,
		SwiftCodes:  swiftCodesResponse.SwiftCodes,
	}

	c.JSON(http.StatusOK, response)
}

// AddSwiftCode godoc
// @Summary Add a SWIFT code
// @Description Adds a new SWIFT code (headquarter or branch)
// @Tags SWIFT Codes
// @Accept json
// @Produce json
// @Param swiftCode body models.SwiftCode true "SWIFT code object"
// @Success 200 {object} models.MessageResponse
// @Failure 400 {object} models.MessageResponse
// @Router /v1/swift-codes/ [post]
func AddSwiftCode(c *gin.Context, swiftService *services.SwiftCodeService) {
	var swiftCodeRequest models.SwiftCode
	if err := c.ShouldBindJSON(&swiftCodeRequest); err != nil {
		c.JSON(errors.GetStatusCode(errors.ErrBadRequest), models.MessageResponse{
			Message: "Invalid input data or JSON format",
		})
		return
	}

	message, err := swiftService.AddSwiftCode(&swiftCodeRequest)
	if err != nil {
		c.JSON(errors.GetStatusCode(err), models.MessageResponse{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, models.MessageResponse{Message: message})
}

// DeleteSwiftCode godoc
// @Summary Delete SWIFT code
// @Description Deletes a headquarter SWIFT code and its branches or a single branch
// @Tags SWIFT Codes
// @Accept json
// @Produce json
// @Param swift-code path string true "SWIFT code"
// @Success 200 {object} models.MessageResponse
// @Failure 404 {object} models.MessageResponse
// @Router /v1/swift-codes/{swift-code} [delete]
func DeleteSwiftCode(c *gin.Context, swiftService *services.SwiftCodeService) {
	swiftCode := strings.ToUpper(c.Param(utils.ParamSwiftCode))

	message, err := swiftService.DeleteSwiftCode(swiftCode)
	if err != nil {
		c.JSON(errors.GetStatusCode(err), models.MessageResponse{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, models.MessageResponse{Message: message})

}
