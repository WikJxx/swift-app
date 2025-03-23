package initialization

import (
	"fmt"
	"log"
	"swift-app/database"
	parser "swift-app/pkg/csv"
)

func InitializeDatabase(uri, dbName, collectionName string) error {
	err := database.InitMongoDB(uri, dbName, collectionName)
	if err != nil {
		return fmt.Errorf("failed to initialize MongoDB: %v", err)
	}
	return nil
}

func ImportDataIfNeeded(csvPath string) error {
	isEmpty, err := database.IsCollectionEmpty()
	if err != nil {
		return fmt.Errorf("failed to check if collection is empty: %v", err)
	}

	if isEmpty {
		swiftCodes, err := parser.LoadSwiftCodes(csvPath)
		if err != nil {
			return fmt.Errorf("failed to load swift codes: %v", err)
		}

		err = database.SaveSwiftCodes(swiftCodes)
		if err != nil {
			return fmt.Errorf("failed to save swift codes: %v", err)
		}
		log.Println("Data imported successfully.")
	} else {
		log.Println("Collection is not empty. Skipping data import.")
	}

	return nil
}
