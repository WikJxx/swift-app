package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	router "swift-app/cmd/router"
	"swift-app/internal/models"
	"swift-app/internal/services"
	utils "swift-app/internal/testutils"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
)

// Package integration contains integration tests verifying the complete HTTP request-response cycle for SWIFT code API endpoints.

func setupRouter() *gin.Engine {
	r := gin.Default()

	swiftService := services.NewSwiftCodeService(utils.Collection)

	router.SetupRoutes(r, swiftService)

	return r
}
func TestSaveAndRetrieveSwiftCode(t *testing.T) {
	_, _ = utils.Collection.DeleteMany(context.Background(), bson.M{})

	service := services.NewSwiftCodeService(utils.Collection)

	swiftCode := &models.SwiftCode{
		SwiftCode:     "AAAABBBXXX",
		BankName:      "Test Bank",
		CountryISO2:   "US",
		CountryName:   "United States",
		Address:       "123 Test St",
		IsHeadquarter: true,
	}
	_, err := service.AddSwiftCode(swiftCode)
	assert.NoError(t, err, "Failed to add SWIFT code")

	result, err := service.GetSwiftCodeDetails("AAAABBBXXX")
	assert.NoError(t, err, "Failed to retrieve SWIFT code")
	assert.Equal(t, "Test Bank", result.BankName, "Expected bank name 'Test Bank'")
}

func TestAddSwiftCodeEndpoint(t *testing.T) {
	_, _ = utils.Collection.DeleteMany(context.Background(), bson.M{})

	r := setupRouter()

	swiftCode := models.SwiftCode{
		SwiftCode:     "AAAABBBXXX",
		BankName:      "Test Bank",
		CountryISO2:   "US",
		CountryName:   "United States",
		Address:       "123 Test St",
		IsHeadquarter: true,
	}
	jsonData, _ := json.Marshal(swiftCode)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/v1/swift-codes/", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code, "Expected status code 200")

	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err, "Failed to unmarshal response")
	assert.Equal(t, "Headquarter SWIFT code added successfully", response["message"], "Expected success message")

	var result models.SwiftCode
	err = utils.Collection.FindOne(context.Background(), bson.M{"swiftCode": "AAAABBBXXX"}).Decode(&result)
	assert.NoError(t, err, "Failed to retrieve SWIFT code from database")
	assert.Equal(t, "Test Bank", result.BankName, "Expected bank name 'Test Bank'")
}
func TestGetSwiftCodeEndpoint(t *testing.T) {
	_, _ = utils.Collection.DeleteMany(context.Background(), bson.M{})

	_, err := utils.Collection.InsertOne(context.Background(), bson.M{
		"swiftCode":     "AAAABBBXXX",
		"bankName":      "Test Bank",
		"address":       "123 Test St",
		"countryISO2":   "US",
		"countryName":   "United States",
		"isHeadquarter": true,
	})
	assert.NoError(t, err, "Failed to insert test data")

	r := setupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/v1/swift-codes/AAAABBBXXX", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code, "Expected status code 200")

	var response models.SwiftCode
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err, "Failed to unmarshal response")
	assert.Equal(t, "Test Bank", response.BankName, "Expected bank name 'Test Bank'")
}
func TestDeleteSwiftCodeEndpoint(t *testing.T) {
	_, _ = utils.Collection.DeleteMany(context.Background(), bson.M{})

	_, err := utils.Collection.InsertOne(context.Background(), bson.M{
		"swiftCode":     "XYZBANKXXX",
		"bankName":      "XYZ Bank",
		"countryISO2":   "UK",
		"countryName":   "United Kingdom",
		"isHeadquarter": true,
	})
	assert.NoError(t, err, "Failed to insert test data")

	r := setupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/v1/swift-codes/XYZBANKXXX", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code, "Expected status code 200")

	var response map[string]string
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err, "Failed to unmarshal response")
	assert.Equal(t, "Deleted hadquarter XYZBANKXXX and its branches", response["message"], "Expected deletion message")

	err = utils.Collection.FindOne(context.Background(), bson.M{"swiftCode": "XYZBANKXXX"}).Decode(&bson.M{})
	assert.Error(t, err, "SWIFT code should be removed from the database")
}
func TestGetSwiftCodesByCountryEndpoint(t *testing.T) {
	_, _ = utils.Collection.DeleteMany(context.Background(), bson.M{})

	_, err := utils.Collection.InsertMany(context.Background(), []interface{}{
		bson.M{
			"swiftCode":     "AAAABBBXXX",
			"bankName":      "Bank A",
			"countryISO2":   "US",
			"countryName":   "United States",
			"isHeadquarter": true,
		},
		bson.M{
			"swiftCode":     "ZZZZPPPXXX",
			"bankName":      "Bank B",
			"countryISO2":   "US",
			"countryName":   "United States",
			"isHeadquarter": true,
		},
	})
	assert.NoError(t, err, "Failed to insert test data")

	r := setupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/v1/swift-codes/country/US", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code, "Expected status code 200")

	var response models.CountrySwiftCodesResponse
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err, "Failed to unmarshal response")
	assert.Equal(t, 2, len(response.SwiftCodes), "Expected 2 SWIFT codes")
}
