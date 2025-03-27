package database

import (
	"context"
	"testing"
	"time"

	"swift-app/internal/models"
	utils "swift-app/internal/testutils"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func clearCollection() {
	_, _ = utils.Collection.DeleteMany(context.Background(), bson.M{})
}

func TestInitMongoDB(t *testing.T) {
	clearCollection()

	uri, err := utils.MongoContainer.ConnectionString(context.Background())
	assert.NoError(t, err, "Failed to retrieve MongoDB URI")

	err = InitMongoDB(uri, "swiftDB", "swiftCodes")
	assert.NoError(t, err, "InitMongoDB should not return an error")

	assert.NotNil(t, collection, "Collection should not be nil")

	indexes, err := utils.Collection.Indexes().List(context.Background())
	assert.NoError(t, err, "Failed to list indexes")

	var indexNames []string
	for indexes.Next(context.Background()) {
		var index bson.M
		err := indexes.Decode(&index)
		assert.NoError(t, err, "Failed to decode index")

		if name, ok := index["name"].(string); ok {
			indexNames = append(indexNames, name)
		}
	}

	assert.Contains(t, indexNames, "swiftCode_1", "Index 'swiftCode_1' should exist")
}

func TestIsCollectionEmpty(t *testing.T) {
	clearCollection()

	empty, err := IsCollectionEmpty()
	assert.NoError(t, err, "IsCollectionEmpty should not return an error")
	assert.True(t, empty, "Collection should be empty")
}

func TestSaveHeadquarters(t *testing.T) {
	clearCollection()

	hqs := []models.SwiftCode{
		{
			SwiftCode:     "HQBANKXXX",
			BankName:      "HQ Bank",
			Address:       "HQ Street",
			CountryISO2:   "US",
			CountryName:   "United States",
			IsHeadquarter: true,
		},
	}

	hqAdded, hqSkipped, err := SaveHeadquarters(hqs)
	assert.NoError(t, err, "SaveHeadquarters should not return an error")

	count, err := utils.Collection.CountDocuments(context.Background(), bson.M{"isHeadquarter": true})
	assert.NoError(t, err)
	assert.Equal(t, int64(1), count, "Expected 1 HQ in the collection")

	assert.Equal(t, 1, hqAdded, "Expected 1 HQ added")
	assert.Equal(t, 0, hqSkipped, "Expected 0 skipped HQ")
}

func TestSaveBranches(t *testing.T) {
	clearCollection()

	hq := models.SwiftCode{
		SwiftCode:     "BRNBANK1XXX",
		BankName:      "BRN Bank",
		Address:       "HQ St",
		CountryISO2:   "PL",
		CountryName:   "Poland",
		IsHeadquarter: true,
	}

	_, _, err := SaveHeadquarters([]models.SwiftCode{hq})
	assert.NoError(t, err, "SaveHeadquarters should not return an error")

	count, err := utils.Collection.CountDocuments(context.Background(), bson.M{"isHeadquarter": true})
	assert.NoError(t, err)
	assert.Equal(t, int64(1), count, "Expected 1 HQ in the collection")

	branch := models.SwiftCode{
		SwiftCode:     "BRNBANK1ABC",
		BankName:      "BRN Branch",
		Address:       "Branch Blvd",
		CountryISO2:   "PL",
		CountryName:   "Poland",
		IsHeadquarter: false,
	}

	branchesAdded, branchesDuplicate, branchesMissingHQ, branchesSkipped, err := SaveBranches([]models.SwiftCode{branch})
	assert.NoError(t, err, "SaveBranches should not return an error")

	var result bson.M
	found := false
	for i := 0; i < 10; i++ {
		err = GetCollection().FindOne(context.Background(), bson.M{"swiftCode": "BRNBANK1XXX"}).Decode(&result)
		if err == nil {
			found = true
			break
		}
		time.Sleep(100 * time.Millisecond)
	}
	assert.True(t, found, "Expected HQ to be found in database")

	rawBranches, ok := result["branches"].(primitive.A)
	assert.True(t, ok, "Expected branches array in HQ document")

	branchFound := false
	for _, b := range rawBranches {
		if doc, ok := b.(bson.M); ok {
			if doc["swiftCode"] == branch.SwiftCode {
				branchFound = true
				break
			}
		}
	}
	assert.True(t, branchFound, "Expected 1 branch under HQ")

	assert.Equal(t, 1, branchesAdded, "Expected 1 branch added")
	assert.Equal(t, 0, branchesDuplicate, "Expected 0 duplicate branches")
	assert.Equal(t, 0, branchesMissingHQ, "Expected 0 missing HQ")
	assert.Equal(t, 0, branchesSkipped, "Expected 0 skipped branches")
}
