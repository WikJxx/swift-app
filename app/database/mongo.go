package database

import (
	"context"
	"fmt"
	"log"
	"swift-app/internal/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client
var collection *mongo.Collection
var isConnected bool

func InitMongoDB(uri string, dbName string, collectionName string) error {
	if isConnected {
		return nil
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

	indexModel := mongo.IndexModel{
		Keys:    bson.M{"swiftCode": 1},
		Options: options.Index().SetUnique(true),
	}
	_, err = collection.Indexes().CreateOne(context.Background(), indexModel)
	if err != nil {
		return fmt.Errorf("failed to create unique index: %v", err)
	}

	isConnected = true
	return nil
}
func IsCollectionEmpty() (bool, error) {
	count, err := collection.CountDocuments(context.Background(), bson.M{})
	if err != nil {
		return false, fmt.Errorf("failed to count documents: %v", err)
	}
	return count == 0, nil
}
func SaveSwiftCodes(swiftCodes []models.SwiftCode) error {

	var docs []interface{}
	headquarters := make(map[string]*models.SwiftCode)

	for _, swiftCode := range swiftCodes {
		if swiftCode.IsHeadquarter {
			headquarters[swiftCode.SwiftCode] = &swiftCode
		} else {
			if headquarter, exists := headquarters[swiftCode.SwiftCode]; exists {
				headquarter.Branches = append(headquarter.Branches, models.SwiftBranch{
					Address:       swiftCode.Address,
					BankName:      swiftCode.BankName,
					CountryISO2:   swiftCode.CountryISO2,
					IsHeadquarter: swiftCode.IsHeadquarter,
					SwiftCode:     swiftCode.SwiftCode,
				})
			} else {
				docs = append(docs, bson.M{
					"swiftCode":     swiftCode.SwiftCode,
					"bankName":      swiftCode.BankName,
					"address":       swiftCode.Address,
					"countryISO2":   swiftCode.CountryISO2,
					"countryName":   swiftCode.CountryName,
					"isHeadquarter": swiftCode.IsHeadquarter,
				})
			}
		}
	}
	for _, headquarter := range headquarters {
		doc := bson.M{
			"swiftCode":     headquarter.SwiftCode,
			"bankName":      headquarter.BankName,
			"address":       headquarter.Address,
			"countryISO2":   headquarter.CountryISO2,
			"countryName":   headquarter.CountryName,
			"isHeadquarter": headquarter.IsHeadquarter,
			"branches":      headquarter.Branches,
		}

		docs = append(docs, doc)
	}

	if len(docs) > 0 {
		_, err := collection.InsertMany(context.Background(), docs, options.InsertMany().SetOrdered(false))
		if err != nil {
			if bulkWriteErr, ok := err.(mongo.BulkWriteException); ok {
				for _, writeError := range bulkWriteErr.WriteErrors {
					if writeError.Code == 11000 {
						fmt.Printf("Duplicate key error for swiftCode: %v\n", writeError)
					} else {
						return fmt.Errorf("failed to insert documents into MongoDB: %v", err)
					}
				}
			} else {
				return fmt.Errorf("failed to insert documents into MongoDB: %v", err)
			}
		}
	}

	fmt.Println("Data imported successfully.")
	return nil
}

func CloseMongoDB() error {
	if client == nil {
		return fmt.Errorf("MongoDB client is not initialized")
	}

	err := client.Disconnect(context.Background())
	if err != nil {
		return fmt.Errorf("failed to close MongoDB connection: %v", err)
	}

	log.Println("MongoDB connection closed successfully")
	return nil
}

func GetCollection() *mongo.Collection {
	return collection
}
