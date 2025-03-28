package models

type ImportSummary struct {
	HQAdded           int
	HQSkipped         int
	BranchesAdded     int
	BranchesDuplicate int
	BranchesMissingHQ int
	BranchesSkipped   int
}
