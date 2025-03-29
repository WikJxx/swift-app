// mongo.go provides functionality for connecting to MongoDB, managing the collection,
// and saving headquarters and branch SWIFT code records. It includes initialization,
// shutdown, document insertion, and helper methods used throughout the application.
package database

import (
	"context"
	"fmt"
	"log"
	"swift-app/internal/models"
	"swift-app/internal/utils"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client
var collection *mongo.Collection
var isConnected bool

// InitMongoDB establishes a connection to the MongoDB instance,
// initializes the target collection, and creates indexes.
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

	indexModels := []mongo.IndexModel{
		{
			Keys:    bson.M{utils.FieldSwiftCode: 1},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys:    bson.M{utils.FieldCountryISO2: 1},
			Options: options.Index(),
		},
	}
	_, err = collection.Indexes().CreateMany(context.Background(), indexModels)
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
func SaveHeadquarters(hqList []models.SwiftCode) (models.ImportSummary, error) {
	summary := models.ImportSummary{}

	for _, hq := range hqList {
		filter := bson.M{utils.FieldSwiftCode: hq.SwiftCode}
		count, err := collection.CountDocuments(context.Background(), filter)
		if err != nil {
			return summary, fmt.Errorf("error checking HQ existence: %v", err)
		}

		if count == 0 {
			_, err := collection.InsertOne(context.Background(), bson.M{
				utils.FieldSwiftCode:     hq.SwiftCode,
				utils.FieldBankName:      hq.BankName,
				utils.FieldAddress:       hq.Address,
				utils.FieldCountryISO2:   hq.CountryISO2,
				utils.FieldCountryName:   hq.CountryName,
				utils.FieldIsHeadquarter: true,
				utils.FieldBranches:      []interface{}{},
			})
			if err != nil {
				return summary, fmt.Errorf("failed to insert HQ: %v", err)
			}
			summary.HQAdded++
		} else {
			summary.HQSkipped++
		}
	}

	return summary, nil
}

func SaveBranches(branches []models.SwiftCode) (models.ImportSummary, error) {
	summary := models.ImportSummary{}

	for _, branch := range branches {
		hqCode := branch.SwiftCode[:8] + "XXX"
		filter := bson.M{utils.FieldSwiftCode: hqCode, utils.FieldIsHeadquarter: true}

		var hq bson.M
		err := collection.FindOne(context.Background(), filter).Decode(&hq)
		if err != nil {
			summary.BranchesMissingHQ++
			summary.BranchesSkipped++
			continue
		}

		branchesField, ok := hq[utils.FieldBranches].(primitive.A)
		if !ok {
			branchesField = primitive.A{}
		}

		duplicate := false
		for _, existing := range branchesField {
			if bmap, ok := existing.(bson.M); ok {
				if bmap[utils.FieldSwiftCode] == branch.SwiftCode {
					duplicate = true
					break
				}
			}
		}
		if duplicate {
			summary.BranchesDuplicate++
			summary.BranchesSkipped++
			continue
		}

		update := bson.M{"$push": bson.M{utils.FieldBranches: bson.M{
			utils.FieldSwiftCode:     branch.SwiftCode,
			utils.FieldBankName:      branch.BankName,
			utils.FieldAddress:       branch.Address,
			utils.FieldCountryISO2:   branch.CountryISO2,
			utils.FieldIsHeadquarter: false,
		}}}

		_, err = collection.UpdateOne(context.Background(), filter, update)
		if err != nil {
			return summary, fmt.Errorf("failed to add branch: %v", err)
		}
		summary.BranchesAdded++
	}

	return summary, nil
}
