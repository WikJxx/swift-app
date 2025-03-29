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

// GetSwiftCode handles GET requests to fetch a SWIFT code (headquarter or branch) by its identifier.
//
// It returns the full details of a headquarter or branch.
// If the code refers to a headquarter, its branches are also included.
//
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

// GetSwiftCodesByCountry handles GET requests to retrieve all SWIFT codes for a given country.
//
// The country is identified using its ISO2 code. Both headquarters and branches are returned.
//
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

// AddSwiftCode handles POST requests to add a new SWIFT code to the system.
//
// It can add both headquarters and branches. Input is validated from JSON.
//
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

// DeleteSwiftCode handles DELETE requests to remove a SWIFT code from the database.
//
// If the provided code is a headquarter, all its branches are also removed.
//
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
