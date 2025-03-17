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

	var skippedRecords []int

	for i, record := range records[1:] {
		if record[fieldIndexes["SWIFT CODE"]] == "" || record[fieldIndexes["COUNTRY ISO2 CODE"]] == "" {
			skippedRecords = append(skippedRecords, i+2)
			continue
		}

		swiftCode := strings.ToUpper(record[fieldIndexes["SWIFT CODE"]])
		countryISO2 := "Unknown"
		if fieldIndexes["COUNTRY ISO2 CODE"] != -1 {
			countryISO2 = strings.ToUpper(record[fieldIndexes["COUNTRY ISO2 CODE"]])
		}

		bankName := "Unknown"
		if fieldIndexes["NAME"] != -1 {
			bankName = record[fieldIndexes["NAME"]]
		}

		address := "Unknown Address"
		if fieldIndexes["ADDRESS"] != -1 {
			address = record[fieldIndexes["ADDRESS"]]
		}

		countryName := "Unknown Country"
		if fieldIndexes["COUNTRY NAME"] != -1 {
			countryName = record[fieldIndexes["COUNTRY NAME"]]
		}

		isHeadquarter := strings.HasSuffix(swiftCode, "XXX")

		swift := models.SwiftCode{
			SwiftCode:     swiftCode,
			CountryISO2:   countryISO2,
			BankName:      bankName,
			Address:       address,
			CountryName:   countryName,
			IsHeadquarter: isHeadquarter,
		}

		if isHeadquarter {
			emptyBranches := []string{}
			swift.Branches = &emptyBranches
			headquartersMap[swiftCode] = &swift
		} else {
			headquarterCode := swiftCode[:8] + "XXX"
			if headquarter, exists := headquartersMap[headquarterCode]; exists {
				*headquarter.Branches = append(*headquarter.Branches, swiftCode)
			}
			swift.Branches = nil
		}

		swiftCodes = append(swiftCodes, swift)

	}

	if len(skippedRecords) > 0 {
		fmt.Printf("Warning: skipped records with missing SWIFT CODE on lines: %v\n", skippedRecords)
	}

	return swiftCodes, nil
}
