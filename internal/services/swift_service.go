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
