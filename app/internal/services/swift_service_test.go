package services_test

import (
	"context"
	"testing"

	"swift-app/internal/models"
	"swift-app/internal/services"
	testutils "swift-app/internal/testutils"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
)

// Package services_test contains integration tests for SwiftCodeService methods.
func TestAddSwiftCode(t *testing.T) {
	_, _ = testutils.Collection.DeleteMany(context.Background(), bson.M{})

	service := services.NewSwiftCodeService(testutils.Collection)

	swiftCode := &models.SwiftCode{
		SwiftCode:     "AAAABBBXXX",
		BankName:      "Test Bank",
		CountryISO2:   "US",
		CountryName:   "United States",
		Address:       "123 Test St",
		IsHeadquarter: true,
	}

	msg, err := service.AddSwiftCode(swiftCode)
	assert.NoError(t, err, "Adding a SWIFT code should not return an error")
	assert.Equal(t, "Headquarter SWIFT code added successfully", msg)

	var result models.SwiftCode
	err = service.DB.FindOne(context.Background(), bson.M{"swiftCode": "AAAABBBXXX"}).Decode(&result)
	assert.NoError(t, err, "SWIFT code should exist in the database")
}

func TestGetSwiftCodeDetails(t *testing.T) {
	_, _ = testutils.Collection.DeleteMany(context.Background(), bson.M{})

	service := services.NewSwiftCodeService(testutils.Collection)

	swiftCode := bson.M{
		"swiftCode":     "AAAABBBXXX",
		"bankName":      "Test Bank",
		"countryName":   "United States",
		"address":       "123 Test St",
		"countryISO2":   "US",
		"isHeadquarter": true,
	}
	_, err := service.DB.InsertOne(context.Background(), swiftCode)
	assert.NoError(t, err, "Inserting SWIFT code into MongoDB should not return an error")

	result, err := service.GetSwiftCodeDetails("AAAABBBXXX")
	assert.NoError(t, err, "Retrieving SWIFT code should not return an error")
	assert.Equal(t, "Test Bank", result.BankName)
}

func TestGetSwiftCodesByCountry(t *testing.T) {
	_, _ = testutils.Collection.DeleteMany(context.Background(), bson.M{})

	service := services.NewSwiftCodeService(testutils.Collection)

	swiftCodes := []interface{}{
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
	}

	_, err := service.DB.InsertMany(context.Background(), swiftCodes)
	assert.NoError(t, err, "Inserting SWIFT codes into MongoDB should not return an error")

	result, err := service.GetSwiftCodesByCountry("US")
	assert.NoError(t, err, "Retrieving SWIFT codes for the country should not return an error")
	assert.Equal(t, 2, len(result.SwiftCodes))
}

func TestDeleteSwiftCode(t *testing.T) {
	service := services.NewSwiftCodeService(testutils.Collection)

	swiftCode := bson.M{
		"swiftCode":     "XYZBANKXXX",
		"bankName":      "XYZ Bank",
		"countryISO2":   "UK",
		"countryName":   "United Kingdom",
		"isHeadquarter": true,
	}

	_, err := service.DB.InsertOne(context.Background(), swiftCode)
	assert.NoError(t, err, "Inserting SWIFT code should not return an error")

	response, err := service.DeleteSwiftCode("XYZBANKXXX")
	assert.NoError(t, err, "Deleting SWIFT code should not return an error")
	assert.Equal(t, "Deleted hadquarter XYZBANKXXX and its branches", response, "Expected deletion message")

	err = service.DB.FindOne(context.Background(), bson.M{"swiftCode": "XYZBANKXXX"}).Decode(&swiftCode)
	assert.Error(t, err, "SWIFT code should be removed from the database")
}
