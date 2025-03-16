package csv

import (
	"encoding/csv"
	"fmt"
	"os"
	"strings"
	"swift-app/internal/models"
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
		"TOWN NAME":         -1,
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
	pendingBranches := make(map[string][]models.SwiftCode)

	var skippedRecords []int

	for i, record := range records[1:] {
		if record[fieldIndexes["SWIFT CODE"]] == "" {
			skippedRecords = append(skippedRecords, i+2)
			continue
		}

		address := "Unknown Address"
		if fieldIndexes["ADDRESS"] != -1 && record[fieldIndexes["ADDRESS"]] != "" {
			address = record[fieldIndexes["ADDRESS"]]
		}

		townName := "Unknown Town"
		if fieldIndexes["TOWN NAME"] != -1 && record[fieldIndexes["TOWN NAME"]] != "" {
			townName = record[fieldIndexes["TOWN NAME"]]
		}
		if address == "Unknown Address" && townName != "Unknown Town" {
			address = townName
		}

		bankName := "Unknown Bank Name"
		if fieldIndexes["NAME"] != -1 && record[fieldIndexes["NAME"]] != "" {
			bankName = record[fieldIndexes["NAME"]]
		}

		countryName := "Unknown Country Name"
		if fieldIndexes["COUNTRY NAME"] != -1 && record[fieldIndexes["COUNTRY NAME"]] != "" {
			countryName = record[fieldIndexes["COUNTRY NAME"]]
		}

		countryISO2Code := "Unknown Country Code"
		if fieldIndexes["COUNTRY ISO2 CODE"] != -1 && record[fieldIndexes["COUNTRY ISO2 CODE"]] != "" {
			countryISO2Code = record[fieldIndexes["COUNTRY ISO2 CODE"]]
		}

		swift := models.SwiftCode{
			CountryISO2:   strings.ToUpper(countryISO2Code),
			SwiftCode:     strings.ToUpper(record[fieldIndexes["SWIFT CODE"]]),
			BankName:      bankName,
			Address:       address,
			CountryName:   countryName,
			IsHeadquarter: strings.HasSuffix(strings.TrimSpace(record[fieldIndexes["SWIFT CODE"]]), "XXX"),
		}

		if swift.IsHeadquarter {
			headquartersMap[swift.SwiftCode] = &swift

			if branches, exists := pendingBranches[swift.SwiftCode[:8]]; exists {
				swift.Branches = append(swift.Branches, branches...)
				delete(pendingBranches, swift.SwiftCode[:8])
			}

			swiftCodes = append(swiftCodes, swift)
		} else {
			headquarterCode := swift.SwiftCode[:8] + "XXX"
			if headquarter, exists := headquartersMap[headquarterCode]; exists {
				headquarter.Branches = append(headquarter.Branches, swift)
			} else {
				pendingBranches[swift.SwiftCode[:8]] = append(pendingBranches[swift.SwiftCode[:8]], swift)
			}
		}
	}

	if len(skippedRecords) > 0 {
		fmt.Printf("Warning: skipped records with missing SWIFT CODE on lines: %v\n", skippedRecords)
	}

	return swiftCodes, nil
}
