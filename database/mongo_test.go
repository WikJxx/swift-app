package database

import (
	"context"
	"testing"

	"swift-app/internal/models"
	utils "swift-app/internal/testutils"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
)

func clearCollection() {
	// Wyczyść kolekcję przed każdym testem
	_, _ = utils.Collection.DeleteMany(context.Background(), bson.M{})
}

func TestInitMongoDB(t *testing.T) {
	// Wyczyść kolekcję przed testem
	clearCollection()

	// Pobierz URI do kontenera MongoDB
	uri, err := utils.MongoContainer.ConnectionString(context.Background())
	if err != nil {
		t.Fatalf("Failed to retrieve MongoDB URI: %v", err)
	}

	// Inicjalizacja bazy danych
	err = InitMongoDB(uri, "swiftDB", "swiftCodes")
	assert.NoError(t, err, "InitMongoDB should not return an error")

	// Sprawdzenie, czy kolekcja została poprawnie zainicjalizowana
	assert.NotNil(t, collection, "Collection should not be nil")

	// Sprawdzenie, czy indeks został poprawnie utworzony
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
	// Wyczyść kolekcję przed testem
	clearCollection()

	// Sprawdzenie, czy kolekcja jest pusta
	empty, err := IsCollectionEmpty()
	assert.NoError(t, err, "IsCollectionEmpty should not return an error")
	assert.True(t, empty, "Collection should be empty")
}

func TestSaveSwiftCodes(t *testing.T) {
	// Wyczyść kolekcję przed testem
	clearCollection()

	// Dane testowe
	swiftCodes := []models.SwiftCode{
		{
			SwiftCode:     "AAAABBBXXX",
			BankName:      "Test Bank",
			Address:       "123 Test St",
			CountryISO2:   "US",
			CountryName:   "United States",
			IsHeadquarter: true,
		},
		{
			SwiftCode:     "ZZZZPPP123",
			BankName:      "Bank B",
			Address:       "456 Test St",
			CountryISO2:   "US",
			CountryName:   "United States",
			IsHeadquarter: false,
		},
	}

	// Zapisywanie danych do bazy
	err := SaveSwiftCodes(swiftCodes)
	assert.NoError(t, err, "SaveSwiftCodes should not return an error")

	// Sprawdzenie, czy dane zostały zapisane
	count, err := utils.Collection.CountDocuments(context.Background(), bson.M{})
	assert.NoError(t, err, "Failed to count documents")
	assert.Equal(t, int64(2), count, "Expected 2 documents in the collection")
}
