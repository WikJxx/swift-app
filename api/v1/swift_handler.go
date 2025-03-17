package v1

import (
	"net/http"
	"strings"
	"swift-app/database"
	"swift-app/internal/models"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetSwiftCode(c *gin.Context) {
	swiftCode := c.Param("swift-code")
	swiftCode = strings.ToUpper(swiftCode)

	collection := database.GetCollection()

	var result models.SwiftCode
	err := collection.FindOne(c, bson.M{"swiftCode": swiftCode}).Decode(&result)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "SWIFT code not found"})
		return
	}
	response := models.SwiftResponse{
		Address:       result.Address,
		BankName:      result.BankName,
		CountryISO2:   result.CountryISO2,
		CountryName:   result.CountryName,
		IsHeadquarter: result.IsHeadquarter,
		SwiftCode:     result.SwiftCode,
	}

	if result.IsHeadquarter && result.Branches != nil && len(*result.Branches) > 0 {
		var branches []models.SwiftBranchResp
		for _, branchSwiftCode := range *result.Branches {
			var branch models.SwiftCode
			err := collection.FindOne(c, bson.M{"swiftCode": branchSwiftCode}).Decode(&branch)
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
	collection := database.GetCollection()

	var countryResult models.SwiftCode
	err := collection.FindOne(c, bson.M{"countryISO2": countryISO2}).Decode(&countryResult)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch country name"})
		return
	}
	countryName := countryResult.CountryName

	cursor, err := collection.Find(c, bson.M{"countryISO2": countryISO2}, options.Find())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch SWIFT codes"})
		return
	}
	defer cursor.Close(c)

	var swiftCodes []models.SwiftBranchResp
	for cursor.Next(c) {
		var swift models.SwiftCode
		if err := cursor.Decode(&swift); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode SWIFT data"})
			return
		}

		swiftCodes = append(swiftCodes, models.SwiftBranchResp{
			Address:       swift.Address,
			BankName:      swift.BankName,
			CountryISO2:   swift.CountryISO2,
			IsHeadquarter: swift.IsHeadquarter,
			SwiftCode:     swift.SwiftCode,
		})
	}

	if len(swiftCodes) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "No SWIFT codes found for this country"})
		return
	}

	response := gin.H{
		"countryISO2": countryISO2,
		"countryName": countryName,
		"swiftCodes":  swiftCodes,
	}

	c.JSON(http.StatusOK, response)
}
