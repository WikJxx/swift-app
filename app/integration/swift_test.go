// Package integration contains integration tests that verify full flow from HTTP endpoints
// through router and handlers down to MongoDB persistence. These tests validate proper end-to-end
// behavior including request handling, data validation, and interaction with the database.
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
	testutils "swift-app/internal/testutils"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
)

func setupRouter() *gin.Engine {
	r := gin.Default()

	swiftService := services.NewSwiftCodeService(testutils.Collection)

	router.SetupRoutes(r, swiftService)

	return r
}
func TestSaveAndRetrieveSwiftCode(t *testing.T) {
	_, _ = testutils.Collection.DeleteMany(context.Background(), bson.M{})

	service := services.NewSwiftCodeService(testutils.Collection)

	swiftCode := &models.SwiftCode{
		SwiftCode:     "AAAABBB1XXX",
		BankName:      "Test Bank",
		CountryISO2:   "US",
		CountryName:   "United States",
		Address:       "123 Test St",
		IsHeadquarter: true,
	}
	_, err := service.AddSwiftCode(swiftCode)
	assert.NoError(t, err, "Failed to add SWIFT code")

	result, err := service.GetSwiftCodeDetails("AAAABBB1XXX")
	assert.NoError(t, err, "Failed to retrieve SWIFT code")
	assert.Equal(t, "Test Bank", result.BankName, "Expected bank name 'Test Bank'")
}

func TestAddSwiftCodeEndpoint(t *testing.T) {
	_, _ = testutils.Collection.DeleteMany(context.Background(), bson.M{})

	r := setupRouter()

	swiftCode := models.SwiftCode{
		SwiftCode:     "AAAABBB1XXX",
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
	assert.Equal(t, "headquarter SWIFT code added successfully", response["message"], "Expected success message")

	var result models.SwiftCode
	err = testutils.Collection.FindOne(context.Background(), bson.M{"swiftCode": "AAAABBB1XXX"}).Decode(&result)
	assert.NoError(t, err, "Failed to retrieve SWIFT code from database")
	assert.Equal(t, "Test Bank", result.BankName, "Expected bank name 'Test Bank'")
}
func TestGetSwiftCodeEndpoint(t *testing.T) {
	_, _ = testutils.Collection.DeleteMany(context.Background(), bson.M{})

	_, err := testutils.Collection.InsertOne(context.Background(), bson.M{
		"swiftCode":     "AAAABBB1XXX",
		"bankName":      "Test Bank",
		"address":       "123 Test St",
		"countryISO2":   "US",
		"countryName":   "United States",
		"isHeadquarter": true,
	})
	assert.NoError(t, err, "Failed to insert test data")

	r := setupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/v1/swift-codes/AAAABBB1XXX", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code, "Expected status code 200")

	var response models.SwiftCode
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err, "Failed to unmarshal response")
	assert.Equal(t, "Test Bank", response.BankName, "Expected bank name 'Test Bank'")
}
func TestDeleteSwiftCodeEndpoint(t *testing.T) {
	_, _ = testutils.Collection.DeleteMany(context.Background(), bson.M{})

	_, err := testutils.Collection.InsertOne(context.Background(), bson.M{
		"swiftCode":     "XYZBANK1XXX",
		"bankName":      "XYZ Bank",
		"countryISO2":   "UK",
		"countryName":   "United Kingdom",
		"isHeadquarter": true,
	})
	assert.NoError(t, err, "Failed to insert test data")

	r := setupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/v1/swift-codes/XYZBANK1XXX", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code, "Expected status code 200")

	var response map[string]string
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err, "Failed to unmarshal response")
	assert.Equal(t, "deleted hadquarter XYZBANK1XXX and its branches", response["message"], "Expected deletion message")

	err = testutils.Collection.FindOne(context.Background(), bson.M{"swiftCode": "XYZBANK1XXX"}).Decode(&bson.M{})
	assert.Error(t, err, "SWIFT code should be removed from the database")
}
func TestGetSwiftCodesByCountryEndpoint(t *testing.T) {
	_, _ = testutils.Collection.DeleteMany(context.Background(), bson.M{})

	_, err := testutils.Collection.InsertMany(context.Background(), []interface{}{
		bson.M{
			"swiftCode":     "AAAABBB1XXX",
			"bankName":      "Bank A",
			"countryISO2":   "US",
			"countryName":   "United States",
			"isHeadquarter": true,
		},
		bson.M{
			"swiftCode":     "ZZZZPPP1XXX",
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
