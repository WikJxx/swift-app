package services

import (
	"context"
	"fmt"
	"swift-app/database"
	"swift-app/internal/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func FindSwiftByCode(code string) (*models.SwiftCode, error) {
	collection := database.GetCollection()

	var swift models.SwiftCode
	err := collection.FindOne(context.Background(), bson.M{"swiftCode": code}).Decode(&swift)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("SWIFT Code not found")
		}
		return nil, err
	}

	return &swift, nil
}

func FindSwiftCodesByCountry(countryISO2 string) ([]models.SwiftBranchResp, string, error) {
	collection := database.GetCollection()

	var countryResult models.SwiftCode
	err := collection.FindOne(context.Background(), bson.M{"countryISO2": countryISO2}).Decode(&countryResult)
	if err != nil {
		return nil, "", fmt.Errorf("Failed to fetch country name")
	}

	countryName := countryResult.CountryName

	cursor, err := collection.Find(context.Background(), bson.M{"countryISO2": countryISO2})
	if err != nil {
		return nil, "", fmt.Errorf("Failed to fetch SWIFT codes")
	}
	defer cursor.Close(context.Background())

	var swiftCodes []models.SwiftBranchResp
	for cursor.Next(context.Background()) {
		var swift models.SwiftCode
		if err := cursor.Decode(&swift); err != nil {
			return nil, "", fmt.Errorf("Failed to decode SWIFT data")
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
		return nil, "", fmt.Errorf("No SWIFT codes found for this country")
	}

	return swiftCodes, countryName, nil
}
