package csv

import (
	"encoding/csv"
	"fmt"
	"os"
	"strings"
	"swift-app/internal/models"
	countries_check "swift-app/internal/utils"
)

func LoadSwiftCodes(filePath string) ([]models.SwiftCode, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.LazyQuotes = true
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	header := records[0]
	for i, h := range header {
		header[i] = strings.TrimSpace(h)
	}

	fieldIndexes := map[string]int{
		"SWIFT CODE":        -1,
		"COUNTRY ISO2 CODE": -1,
		"NAME":              -1,
		"ADDRESS":           -1,
		"COUNTRY NAME":      -1,
	}

	for i, field := range header {
		upperField := strings.ToUpper(field)
		if _, exists := fieldIndexes[upperField]; exists {
			fieldIndexes[upperField] = i
		}
	}

	if fieldIndexes["SWIFT CODE"] == -1 {
		return nil, fmt.Errorf("missing required field: SWIFT CODE")
	}

	var swiftCodes []models.SwiftCode
	headquartersMap := make(map[string]*models.SwiftCode)
	branchQueue := make([]*models.SwiftBranch, 0)
	uniqueHQ := make(map[string]bool)
	uniqueBranches := make(map[string]bool)
	countries, err := countries_check.LoadCountries()
	if err != nil {
		return nil, fmt.Errorf("error loading country data: %v", err)
	}

	for _, record := range records[1:] {
		swiftCode := strings.ToUpper(record[fieldIndexes["SWIFT CODE"]])
		countryISO2 := strings.ToUpper(record[fieldIndexes["COUNTRY ISO2 CODE"]])
		if len(countryISO2) != 2 {
			fmt.Printf("Warning: Country ISO %s is smaller than 2 characters\n", swiftCode)
			continue
		}
		if record[fieldIndexes["SWIFT CODE"]] == "" || record[fieldIndexes["COUNTRY ISO2 CODE"]] == "" {
			continue
		}

		if len(swiftCode) > 11 {
			fmt.Printf("Warning: SWIFT code %s is longer than 11 characters\n", swiftCode)
			continue
		}
		if len(swiftCode) < 8 {
			fmt.Printf("Warning: SWIFT code %s is smaller than 8 characters\n", swiftCode)
			continue
		}

		bankName := strings.ToUpper(record[fieldIndexes["NAME"]])
		if fieldIndexes["NAME"] == -1 {
			fmt.Printf("Warning: Name %s can not be empty \n", swiftCode)
			continue
		}

		address := ""
		if fieldIndexes["ADDRESS"] != -1 && fieldIndexes["ADDRESS"] < len(record) {
			address = strings.TrimSpace(strings.ToUpper(record[fieldIndexes["ADDRESS"]]))
		}

		countryName := strings.ToUpper(record[fieldIndexes["COUNTRY NAME"]])
		if fieldIndexes["COUNTRY NAME"] == -1 {
			fmt.Printf("Warning: Country name %s can not be empty \n", swiftCode)
		}

		isHeadquarter := strings.HasSuffix(swiftCode, "XXX")

		country, exists := countries[countryISO2]
		if !exists {
			fmt.Printf("Warning: Invalid country ISO2 for SWIFT code %s\n", swiftCode)
			continue
		}

		if !strings.EqualFold(countryName, country.Name) {
			fmt.Printf("Warning: Country name '%s' does not match ISO2 '%s'\n", countryName, countryISO2)
			continue
		}

		if isHeadquarter {
			if _, exists := uniqueHQ[swiftCode]; exists {
				continue
			}

			uniqueHQ[swiftCode] = true

			swift := models.SwiftCode{
				SwiftCode:     swiftCode,
				CountryISO2:   countryISO2,
				BankName:      bankName,
				Address:       address,
				CountryName:   countryName,
				IsHeadquarter: isHeadquarter,
				Branches:      []models.SwiftBranch{},
			}

			headquartersMap[swiftCode] = &swift

			for i := len(branchQueue) - 1; i >= 0; i-- {
				branch := branchQueue[i]
				if branch.SwiftCode[:8]+"XXX" == swiftCode {
					swift.Branches = append(swift.Branches, *branch)
					branchQueue = append(branchQueue[:i], branchQueue[i+1:]...)
				}
			}

			swiftCodes = append(swiftCodes, swift)
		} else {
			if _, exists := uniqueBranches[swiftCode]; exists {
				continue
			}

			uniqueBranches[swiftCode] = true

			branch := models.SwiftBranch{
				SwiftCode:     swiftCode,
				BankName:      bankName,
				Address:       address,
				CountryISO2:   countryISO2,
				IsHeadquarter: false,
			}

			if hq, found := headquartersMap[swiftCode[:8]+"XXX"]; found {
				hq.Branches = append(hq.Branches, branch)
			} else {
				branchQueue = append(branchQueue, &branch)
			}
		}
	}

	if len(branchQueue) > 0 {
		for _, branch := range branchQueue {
			fmt.Printf("Warning: Branch %s does not have a matching headquarter.\n", branch.SwiftCode)
		}
	}

	return swiftCodes, nil
}
