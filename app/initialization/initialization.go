package initialization

import (
	"fmt"
	"swift-app/database"
	"swift-app/internal/models"
	parser "swift-app/pkg/csv"
)

// Function sets up the MongoDB connection and initializes the required database collection.
func InitializeDatabase(uri, dbName, collectionName string) error {
	err := database.InitMongoDB(uri, dbName, collectionName)
	if err != nil {
		return fmt.Errorf("failed to initialize MongoDB: %v", err)
	}
	return nil
}

// Function loads SWIFT codes from a CSV file and imports them into the database, managing headquarters and branches separately.
func ImportData(csvPath string) (int, int, int, int, int, int, error) {
	swiftCodes, err := parser.LoadSwiftCodes(csvPath)
	if err != nil {
		return 0, 0, 0, 0, 0, 0, fmt.Errorf("failed to load swift codes: %v", err)
	}

	var hqList, branchList []models.SwiftCode
	for _, code := range swiftCodes {
		if code.IsHeadquarter {
			hqList = append(hqList, code)
		} else {
			branchList = append(branchList, code)
		}
	}

	hqAdded, hqSkipped, err := database.SaveHeadquarters(hqList)
	if err != nil {
		return 0, 0, 0, 0, 0, 0, fmt.Errorf("failed to save HQs: %v", err)
	}

	branchesAdded, branchesDuplicate, branchesMissingHQ, branchesSkipped, err := database.SaveBranches(branchList)
	if err != nil {
		return 0, 0, 0, 0, 0, 0, fmt.Errorf("failed to save branches: %v", err)
	}

	return hqAdded, hqSkipped, branchesAdded, branchesDuplicate, branchesMissingHQ, branchesSkipped, nil
}
