package initialization

import (
	"fmt"
	"swift-app/database"
	"swift-app/internal/models"
	parser "swift-app/pkg/csv"
)

// InitializeDatabase connects to MongoDB and initializes the target collection.
func InitializeDatabase(uri, dbName, collectionName string) error {
	err := database.InitMongoDB(uri, dbName, collectionName)
	if err != nil {
		return fmt.Errorf("failed to initialize MongoDB: %v", err)
	}
	return nil
}

// ImportData loads SWIFT codes from a CSV file and imports them into the database, managing headquarters and branches separately.
func ImportData(csvPath string) (*models.ImportSummary, error) {
	swiftCodes, err := parser.LoadSwiftCodes(csvPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load swift codes: %v", err)
	}

	var hqList, branchList []models.SwiftCode
	for _, code := range swiftCodes {
		if code.IsHeadquarter {
			hqList = append(hqList, code)
		} else {
			branchList = append(branchList, code)
		}
	}

	hqSummary, err := database.SaveHeadquarters(hqList)
	if err != nil {
		return nil, fmt.Errorf("failed to save HQs: %v", err)
	}

	branchSummary, err := database.SaveBranches(branchList)
	if err != nil {
		return nil, fmt.Errorf("failed to save branches: %v", err)
	}

	summary := &models.ImportSummary{
		HQAdded:           hqSummary.HQAdded,
		HQSkipped:         hqSummary.HQSkipped,
		BranchesAdded:     branchSummary.BranchesAdded,
		BranchesDuplicate: branchSummary.BranchesDuplicate,
		BranchesMissingHQ: branchSummary.BranchesMissingHQ,
		BranchesSkipped:   branchSummary.BranchesSkipped,
	}

	return summary, nil
}
