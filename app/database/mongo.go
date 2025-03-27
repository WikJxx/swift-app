package database

import (
	"context"
	"fmt"
	"log"
	"swift-app/internal/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
	hqMap := make(map[string]*models.SwiftCode)
	branches := make([]models.SwiftCode, 0)
	branchesSkippedMissingHQ := []string{}
	branchesDuplicates := 0
	hqAdded := 0
	hqAlreadyExists := 0
	branchesAdded := 0

	for _, code := range swiftCodes {
		if code.IsHeadquarter {
			hqMap[code.SwiftCode] = &code
		} else {
			branches = append(branches, code)
		}
	}

	for _, hq := range hqMap {
		filter := bson.M{"swiftCode": hq.SwiftCode}
		count, err := collection.CountDocuments(context.Background(), filter)
		if err != nil {
			return fmt.Errorf("error checking HQ existence: %v", err)
		}

		if count == 0 {
			_, err := collection.InsertOne(context.Background(), bson.M{
				"swiftCode":     hq.SwiftCode,
				"bankName":      hq.BankName,
				"address":       hq.Address,
				"countryISO2":   hq.CountryISO2,
				"countryName":   hq.CountryName,
				"isHeadquarter": true,
				"branches":      hq.Branches,
			})
			if err != nil {
				return fmt.Errorf("failed to insert HQ: %v", err)
			}
			hqAdded++
		} else {
			hqAlreadyExists++
		}
	}

	for _, branch := range branches {
		hqCode := branch.SwiftCode[:8] + "XXX"
		filter := bson.M{"swiftCode": hqCode, "isHeadquarter": true}

		var hq models.SwiftCode
		err := collection.FindOne(context.Background(), filter).Decode(&hq)
		if err != nil {
			branchesSkippedMissingHQ = append(branchesSkippedMissingHQ, branch.SwiftCode)
			continue
		}

		branchExists := false
		for _, existing := range hq.Branches {
			if existing.SwiftCode == branch.SwiftCode {
				branchExists = true
				break
			}
		}
		if branchExists {
			branchesDuplicates++
			continue
		}

		update := bson.M{"$push": bson.M{"branches": bson.M{
			"swiftCode":     branch.SwiftCode,
			"bankName":      branch.BankName,
			"address":       branch.Address,
			"countryISO2":   branch.CountryISO2,
			"isHeadquarter": false,
		}}}

		result, err := collection.UpdateOne(context.Background(), filter, update)
		if err != nil {
			return fmt.Errorf("failed to add branch: %v", err)
		}
		if result.ModifiedCount > 0 {
			branchesAdded++
		} else {
			branchesDuplicates++
		}
	}

	fmt.Printf("Data import complete. HQs added: %d, already existing HQs: %d\n", hqAdded, hqAlreadyExists)
	fmt.Printf("Branches added: %d\n", branchesAdded)
	fmt.Printf("Branches skipped due to duplicates: %d\n", branchesDuplicates)

	if len(branchesSkippedMissingHQ) > 0 {
		fmt.Printf("Branches not added due to missing HQs (%d): %s\n", len(branchesSkippedMissingHQ), joinBranchCodes(branchesSkippedMissingHQ))
	}

	return nil
}

func joinBranchCodes(codes []string) string {
	return fmt.Sprintf("%s", codes)
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
func SaveHeadquarters(hqList []models.SwiftCode) (int, int, int, error) {
	hqAdded := 0
	hqExisting := 0
	hqSkipped := 0

	for _, hq := range hqList {
		filter := bson.M{"swiftCode": hq.SwiftCode}
		count, err := collection.CountDocuments(context.Background(), filter)
		if err != nil {
			return hqAdded, hqExisting, hqSkipped, fmt.Errorf("error checking HQ existence: %v", err)
		}

		if count == 0 {
			_, err := collection.InsertOne(context.Background(), bson.M{
				"swiftCode":     hq.SwiftCode,
				"bankName":      hq.BankName,
				"address":       hq.Address,
				"countryISO2":   hq.CountryISO2,
				"countryName":   hq.CountryName,
				"isHeadquarter": true,
				"branches":      []interface{}{},
			})
			if err != nil {
				return hqAdded, hqExisting, hqSkipped, fmt.Errorf("failed to insert HQ: %v", err)
			}
			hqAdded++
		} else {
			hqExisting++
			hqSkipped++ // Przypisujemy, że te HQ zostały pominięte z powodu ich wcześniejszego istnienia
		}
	}

	return hqAdded, hqExisting, hqSkipped, nil
}

func SaveBranches(branches []models.SwiftCode) (int, int, int, int, error) {
	branchesAdded := 0
	branchesDuplicate := 0
	branchesMissingHQ := 0
	branchesSkipped := 0

	for _, branch := range branches {
		hqCode := branch.SwiftCode[:8] + "XXX"
		filter := bson.M{"swiftCode": hqCode, "isHeadquarter": true}

		var hq bson.M
		err := collection.FindOne(context.Background(), filter).Decode(&hq)
		if err != nil {
			branchesMissingHQ++
			branchesSkipped++
			continue
		}

		branchesField, ok := hq["branches"].(primitive.A)
		if !ok {
			branchesField = primitive.A{}
		}

		duplicate := false
		for _, existing := range branchesField {
			if bmap, ok := existing.(bson.M); ok {
				if bmap["swiftCode"] == branch.SwiftCode {
					duplicate = true
					break
				}
			}
		}
		if duplicate {
			branchesDuplicate++
			branchesSkipped++
			continue
		}

		update := bson.M{"$push": bson.M{"branches": bson.M{
			"swiftCode":     branch.SwiftCode,
			"bankName":      branch.BankName,
			"address":       branch.Address,
			"countryISO2":   branch.CountryISO2,
			"isHeadquarter": false,
		}}}

		_, err = collection.UpdateOne(context.Background(), filter, update)
		if err != nil {
			return branchesAdded, branchesDuplicate, branchesMissingHQ, branchesSkipped, fmt.Errorf("failed to add branch: %v", err)
		}
		branchesAdded++
	}

	return branchesAdded, branchesDuplicate, branchesMissingHQ, branchesSkipped, nil
}
