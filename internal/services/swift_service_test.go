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

// func setupTestMongo(t *testing.T) (*services.SwiftCodeService, func()) {
// 	t.Helper()

// 	ctx := context.Background()
// 	mongoContainer, err := mongodb.Run(ctx, "mongo:latest")
// 	if err != nil {
// 		t.Fatalf("Failed to start MongoDB container: %v", err)
// 	}

// 	uri, err := mongoContainer.ConnectionString(ctx)
// 	if err != nil {
// 		t.Fatalf("Failed to retrieve MongoDB URI: %v", err)
// 	}

// 	mongoClient, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
// 	if err != nil {
// 		t.Fatalf("Failed to connect to MongoDB: %v", err)
// 	}

// 	collection := mongoClient.Database("swiftDB").Collection("swiftCodes")
// 	service := services.NewSwiftCodeService(collection)

// 	cleanup := func() {
// 		mongoClient.Disconnect(ctx)
// 		mongoContainer.Terminate(ctx)
// 	}

// 	return service, cleanup
// }

func TestAddSwiftCode(t *testing.T) {
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
	assert.Equal(t, "Headquarter SWIFT code added successfully", msg["message"])

	var result models.SwiftCode
	err = service.DB.FindOne(context.Background(), bson.M{"swiftCode": "AAAABBBXXX"}).Decode(&result)
	assert.NoError(t, err, "SWIFT code should exist in the database")
}

func TestGetSwiftCodeDetails(t *testing.T) {
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

	deletedCount, err := service.DeleteSwiftCode("XYZBANKXXX")
	assert.NoError(t, err, "Deleting SWIFT code should not return an error")
	assert.Equal(t, int64(1), deletedCount)

	err = service.DB.FindOne(context.Background(), bson.M{"swiftCode": "XYZBANKXXX"}).Decode(&swiftCode)
	assert.Error(t, err, "SWIFT code should be removed from the database")
}
