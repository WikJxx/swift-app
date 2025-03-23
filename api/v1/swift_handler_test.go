package v1

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"swift-app/internal/models"
	"swift-app/internal/services"
	testutils "swift-app/internal/testutils"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
)

func clearCollection() {
	_, _ = testutils.Collection.DeleteMany(context.Background(), bson.M{})
}

func TestGetSwiftCode(t *testing.T) {
	clearCollection()

	service := services.NewSwiftCodeService(testutils.Collection)

	_, err := testutils.Collection.InsertOne(context.Background(), bson.M{
		"swiftCode":     "AAAABBBXXX",
		"bankName":      "Test Bank",
		"address":       "123 Test St",
		"countryISO2":   "US",
		"countryName":   "United States",
		"isHeadquarter": true,
	})
	assert.NoError(t, err, "Failed to insert test data")
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = []gin.Param{{Key: "swift-code", Value: "AAAABBBXXX"}}

	GetSwiftCode(c, service)

	assert.Equal(t, http.StatusOK, w.Code, "Expected status code 200")
	assert.Contains(t, w.Body.String(), "Test Bank", "Expected bank name in response")
}

func TestGetSwiftCode_NotFound(t *testing.T) {
	clearCollection()

	service := services.NewSwiftCodeService(testutils.Collection)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = []gin.Param{{Key: "swift-code", Value: "NONEXIST"}}

	GetSwiftCode(c, service)

	assert.Equal(t, http.StatusNotFound, w.Code, "Expected status code 404")
	assert.Contains(t, w.Body.String(), "missing headquarter: NONEXISTXXX", "Expected error message in response")
}

func TestGetSwiftCodesByCountry(t *testing.T) {
	clearCollection()

	service := services.NewSwiftCodeService(testutils.Collection)

	_, err := testutils.Collection.InsertMany(context.Background(), []interface{}{
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

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = []gin.Param{{Key: "countryISO2code", Value: "US"}}

	GetSwiftCodesByCountry(c, service)

	assert.Equal(t, http.StatusOK, w.Code, "Expected status code 200")

	var response models.CountrySwiftCodesResponse
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err, "Failed to unmarshal response")
	assert.Equal(t, 2, len(response.SwiftCodes), "Expected 2 SWIFT codes")
}

func TestAddSwiftCode(t *testing.T) {
	clearCollection()

	service := services.NewSwiftCodeService(testutils.Collection)

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
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("POST", "/swift", bytes.NewBuffer(jsonData))
	c.Request.Header.Set("Content-Type", "application/json")

	AddSwiftCode(c, service)

	assert.Equal(t, http.StatusOK, w.Code, "Expected status code 200")

	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err, "Failed to unmarshal response")
	assert.Equal(t, "Headquarter SWIFT code added successfully", response["message"], "Expected success message")
}

func TestDeleteSwiftCode(t *testing.T) {
	clearCollection()

	service := services.NewSwiftCodeService(testutils.Collection)

	_, err := testutils.Collection.InsertOne(context.Background(), bson.M{
		"swiftCode":     "XYZBANKXXX",
		"bankName":      "XYZ Bank",
		"countryISO2":   "UK",
		"countryName":   "United Kingdom",
		"isHeadquarter": true,
	})
	assert.NoError(t, err, "Failed to insert test data")

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = []gin.Param{{Key: "swift-code", Value: "XYZBANKXXX"}}

	DeleteSwiftCode(c, service)

	assert.Equal(t, http.StatusOK, w.Code, "Expected status code 200")

	var response map[string]string
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err, "Failed to unmarshal response")
	assert.Equal(t, "Deleted 1 records", response["message"], "Expected deletion message")
}
