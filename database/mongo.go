package database

import (
	"context"
	"fmt"
	"swift-app/internal/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client
var collection *mongo.Collection
var isConnected bool

func InitMongoDB(uri string, dbName string, collectionName string) error {
	// Sprawdzenie, czy połączenie zostało już nawiązane
	if isConnected {
		return nil // Połączenie już istnieje, więc nic nie robimy
	}

	var err error
	client, err = mongo.Connect(context.Background(), options.Client().ApplyURI(uri))
	if err != nil {
		return fmt.Errorf("failed to connect to MongoDB: %v", err)
	}

	err = client.Ping(context.Background(), nil)
	if err != nil {
		return fmt.Errorf("failed to ping MongoDB: %v", err)
	}

	collection = client.Database(dbName).Collection(collectionName)
	isConnected = true // Ustawiamy flagę na true, po udanym połączeniu
	return nil
}

func SaveSwiftCodes(swiftCodes []models.SwiftCode) error {
	var docs []interface{}
	for _, swiftCode := range swiftCodes {
		// Sprawdź, czy dokument z danym swiftCode już istnieje
		var existingDoc models.SwiftCode
		err := collection.FindOne(context.Background(), bson.M{"swiftCode": swiftCode.SwiftCode}).Decode(&existingDoc)

		// Jeśli dokument już istnieje (błąd jest mongo.ErrNoDocuments), to go pomijamy
		if err == nil {
			// Dokument już istnieje, więc go pomijamy
			continue
		} else if err != mongo.ErrNoDocuments {
			// Błąd podczas szukania dokumentu
			return fmt.Errorf("failed to check for existing document: %v", err)
		}

		// Dokument nie istnieje, więc dodajemy go do listy do zapisania
		doc := bson.M{
			"swiftCode":     swiftCode.SwiftCode,
			"bankName":      swiftCode.BankName,
			"address":       swiftCode.Address,
			"countryISO2":   swiftCode.CountryISO2,
			"countryName":   swiftCode.CountryName,
			"isHeadquarter": swiftCode.IsHeadquarter,
		}

		// Dodaj branches tylko, jeśli są
		if swiftCode.Branches != nil && len(*swiftCode.Branches) > 0 {
			doc["branches"] = *swiftCode.Branches
		}

		docs = append(docs, doc)
	}

	if len(docs) > 0 {
		_, err := collection.InsertMany(context.Background(), docs)
		if err != nil {
			return fmt.Errorf("failed to insert documents into MongoDB: %v", err)
		}
	}

	return nil
}
