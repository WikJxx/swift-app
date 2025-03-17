package v1

import (
	"net/http"
	"swift-app/database"
	"swift-app/internal/models"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

func GetSwiftCode(c *gin.Context) {
	swiftCode := c.Param("swift-code")

	collection := database.GetCollection()

	// Pobieranie dokumentu głównej siedziby (headquarter)
	var result models.SwiftCode
	err := collection.FindOne(c, bson.M{"swiftCode": swiftCode}).Decode(&result)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "SWIFT code not found"})
		return
	}

	// Tworzenie struktury odpowiedzi
	response := models.SwiftResponse{
		Address:       result.Address,
		BankName:      result.BankName,
		CountryISO2:   result.CountryISO2,
		CountryName:   result.CountryName,
		IsHeadquarter: result.IsHeadquarter,
		SwiftCode:     result.SwiftCode,
	}

	// Dodanie oddziałów (branches), jeśli istnieją
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

	// Zwracamy odpowiedź w wymaganym formacie
	c.JSON(http.StatusOK, response)
}
