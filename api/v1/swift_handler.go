package v1

import (
	"net/http"
	"strings"
	"swift-app/internal/models"
	"swift-app/internal/services"

	"github.com/gin-gonic/gin"
)

func GetSwiftCode(c *gin.Context) {
	swiftCode := c.Param("swift-code")
	swiftCode = strings.ToUpper(swiftCode)

	swift, err := services.FindSwiftByCode(swiftCode)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	response := models.SwiftResponse{
		Address:       swift.Address,
		BankName:      swift.BankName,
		CountryISO2:   swift.CountryISO2,
		CountryName:   swift.CountryName,
		IsHeadquarter: swift.IsHeadquarter,
		SwiftCode:     swift.SwiftCode,
	}

	if swift.IsHeadquarter && swift.Branches != nil && len(*swift.Branches) > 0 {
		var branches []models.SwiftBranchResp
		for _, branchSwiftCode := range *swift.Branches {
			branch, err := services.FindSwiftByCode(branchSwiftCode)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch branch data"})
				return
			}
			branches = append(branches, models.SwiftBranchResp{
				Address:       branch.Address,
				BankName:      branch.BankName,
				CountryISO2:   branch.CountryISO2,
				IsHeadquarter: branch.IsHeadquarter,
				SwiftCode:     branch.SwiftCode,
			})
		}
		response.Branches = branches
	}

	c.JSON(http.StatusOK, response)
}

func GetSwiftCodesByCountry(c *gin.Context) {
	countryISO2 := c.Param("countryISO2code")
	countryISO2 = strings.ToUpper(countryISO2)

	swiftCodes, countryName, err := services.FindSwiftCodesByCountry(countryISO2)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := gin.H{
		"countryISO2": countryISO2,
		"countryName": countryName,
		"swiftCodes":  swiftCodes,
	}

	c.JSON(http.StatusOK, response)
}
