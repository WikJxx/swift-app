package v1

import (
	"net/http"
	"strings"
	"swift-app/internal/models"
	"swift-app/internal/services"

	"github.com/gin-gonic/gin"
)

func GetSwiftCode(c *gin.Context, swiftService *services.SwiftCodeService) {
	swiftCode := c.Param("swift-code")
	swiftCode = strings.ToUpper(swiftCode)

	swift, err := swiftService.GetSwiftCodeDetails(swiftCode)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
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

	if swift.IsHeadquarter && len(swift.Branches) > 0 {
		response.Branches = swift.Branches
	}

	c.JSON(http.StatusOK, response)
}

func GetSwiftCodesByCountry(c *gin.Context, swiftService *services.SwiftCodeService) {
	countryISO2 := strings.ToUpper(c.Param("countryISO2code"))

	swiftCodesResponse, err := swiftService.GetSwiftCodesByCountry(countryISO2)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	var countryName string
	if len(swiftCodesResponse.SwiftCodes) > 0 {
		countryName = swiftCodesResponse.SwiftCodes[0].CountryName
	}

	response := models.CountrySwiftCodesResponse{
		CountryISO2: countryISO2,
		CountryName: countryName,
		SwiftCodes:  swiftCodesResponse.SwiftCodes,
	}

	c.JSON(http.StatusOK, response)
}

func AddSwiftCode(c *gin.Context, swiftService *services.SwiftCodeService) {
	var swiftCodeRequest models.SwiftCode
	if err := c.ShouldBindJSON(&swiftCodeRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input data"})
		return
	}

	message, err := swiftService.AddSwiftCode(&swiftCodeRequest)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": message["message"]}) // 👈 Teraz zwraca Twój komunikat!
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": message["message"]})
}
