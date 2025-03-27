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
	assert.NoError(t, err)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = []gin.Param{{Key: "swift-code", Value: "aaaabbbxxx"}} // lowercase, test ToUpper

	GetSwiftCode(c, service)

	assert.Equal(t, http.StatusOK, w.Code)
	var response models.SwiftCode
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "Test Bank", response.BankName)
	assert.Equal(t, "AAAABBBXXX", response.SwiftCode)
}

func TestGetSwiftCode_NotFound(t *testing.T) {
	clearCollection()
	service := services.NewSwiftCodeService(testutils.Collection)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = []gin.Param{{Key: "swift-code", Value: "NONEXIST"}}

	GetSwiftCode(c, service)

	assert.Equal(t, http.StatusNotFound, w.Code)
	var response models.MessageResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Contains(t, response.Message, "headquarter NONEXISTXXX is missing")
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
	assert.NoError(t, err)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = []gin.Param{{Key: "countryISO2code", Value: "us"}} // test lowercase input

	GetSwiftCodesByCountry(c, service)

	assert.Equal(t, http.StatusOK, w.Code)

	var response models.CountrySwiftCodesResponse
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "US", response.CountryISO2)
	assert.Equal(t, 2, len(response.SwiftCodes))
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

	assert.Equal(t, http.StatusOK, w.Code)

	var response models.MessageResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "Headquarter SWIFT code added successfully", response.Message)
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
	assert.NoError(t, err)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = []gin.Param{{Key: "swift-code", Value: "XYZBANKXXX"}}

	DeleteSwiftCode(c, service)

	assert.Equal(t, http.StatusOK, w.Code)

	var response models.MessageResponse
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "Deleted hadquarter XYZBANKXXX and its branches", response.Message)
}
