package router

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"swift-app/internal/models"
	"swift-app/internal/services"
	utils "swift-app/internal/testutils"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
)

func setupRouter() *gin.Engine {
	r := gin.Default()
	swiftService := services.NewSwiftCodeService(utils.Collection)
	SetupRoutes(r, swiftService)

	return r
}

func TestGetSwiftCode(t *testing.T) {
	_, _ = utils.Collection.DeleteMany(context.Background(), bson.M{})

	_, err := utils.Collection.InsertOne(context.Background(), bson.M{
		"swiftCode":     "AAAABBBXXX",
		"bankName":      "Test Bank",
		"address":       "123 Test St",
		"countryISO2":   "US",
		"countryName":   "United States",
		"isHeadquarter": true,
	})
	assert.NoError(t, err)

	r := setupRouter()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/v1/swift-codes/AAAABBBXXX", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response models.SwiftCode
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "Test Bank", response.BankName)
}

func TestGetSwiftCode_NotFound(t *testing.T) {
	_, _ = utils.Collection.DeleteMany(context.Background(), bson.M{})

	r := setupRouter()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/v1/swift-codes/NONEXISTXXX", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)

	var response models.MessageResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Contains(t, response.Message, "headquarter not found: NONEXISTXXX")
}

func TestAddSwiftCode(t *testing.T) {
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

	assert.Equal(t, http.StatusOK, w.Code)

	var response models.MessageResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "Headquarter SWIFT code added successfully", response.Message)
}

func TestDeleteSwiftCode(t *testing.T) {
	_, _ = utils.Collection.DeleteMany(context.Background(), bson.M{})

	_, err := utils.Collection.InsertOne(context.Background(), bson.M{
		"swiftCode":     "XYZBANKXXX",
		"bankName":      "XYZ Bank",
		"countryISO2":   "UK",
		"countryName":   "United Kingdom",
		"isHeadquarter": true,
	})
	assert.NoError(t, err)

	r := setupRouter()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/v1/swift-codes/XYZBANKXXX", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response models.MessageResponse
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "Deleted hadquarter XYZBANKXXX and its branches", response.Message)
}
