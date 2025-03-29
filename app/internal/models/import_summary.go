package models

// ImportSummary holds statistics about the import process.
type ImportSummary struct {
	HQAdded           int
	HQSkipped         int
	BranchesAdded     int
	BranchesDuplicate int
	BranchesMissingHQ int
	BranchesSkipped   int
}
