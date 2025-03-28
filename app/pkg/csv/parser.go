package csv

import (
	"encoding/csv"
	"fmt"
	"os"
	"strings"
	"swift-app/internal/models"
	"swift-app/internal/utils"
)

// Loads and parses a CSV file containing SWIFT code data, validates each record, and returns a list of unique, validated SWIFT codes.
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

	header := sanitizeHeader(records[0])
	fieldIndexes := getFieldIndexes(header)

	if fieldIndexes["SWIFT CODE"] == -1 {
		return nil, fmt.Errorf("missing required field: SWIFT CODE")
	}

	countries, err := utils.LoadCountries()
	if err != nil {
		return nil, fmt.Errorf("error loading country data: %v", err)
	}

	return processRecords(records[1:], fieldIndexes, countries)
}

// Converts all header fields to uppercase and trims whitespace to ensure consistent field matching.
func sanitizeHeader(header []string) []string {
	for i, h := range header {
		header[i] = strings.ToUpper(strings.TrimSpace(h))
	}
	return header
}

//Maps column names (including supported aliases) to their respective indexes in the CSV header row.

func getFieldIndexes(header []string) map[string]int {
	fieldIndexes := map[string]int{
		"SWIFT CODE":        -1,
		"COUNTRY ISO2 CODE": -1,
		"NAME":              -1,
		"ADDRESS":           -1,
		"COUNTRY NAME":      -1,
	}

	aliases := map[string]string{
		"SWIFTCODE":   "SWIFT CODE",
		"SWIFT CODES": "SWIFT CODE",
		"SWIFT_CODE":  "SWIFT CODE",
		"SWIFT C0DE":  "SWIFT CODE",
	}

	for i, field := range header {
		upperField := strings.ToUpper(strings.TrimSpace(field))

		if alias, ok := aliases[upperField]; ok {
			upperField = alias
		}

		if _, exists := fieldIndexes[upperField]; exists {
			fieldIndexes[upperField] = i
		}
	}

	return fieldIndexes
}

// Processes all rows from the CSV file, validates them, and constructs SwiftCode structs while skipping duplicates or invalid entries.
func processRecords(records [][]string, fieldIndexes map[string]int, countries map[string]models.Country) ([]models.SwiftCode, error) {
	swiftCodes := []models.SwiftCode{}
	uniqueCodes := make(map[string]bool)

	for _, record := range records {
		swiftCode, countryISO2, bankName, address, countryName := extractRecordData(record, fieldIndexes)

		if err := validateRecord(swiftCode, countryISO2, countryName, countries); err != nil {
			fmt.Printf("Warning: %s - %v\n", swiftCode, err)
			continue
		}

		if uniqueCodes[swiftCode] {
			continue
		}
		uniqueCodes[swiftCode] = true

		isHeadquarter := strings.HasSuffix(swiftCode, "XXX")
		if err := utils.ValidateSwiftCodeSuffix(swiftCode, isHeadquarter); err != nil {
			fmt.Printf("Warning: %s - %v\n", swiftCode, err)
			continue
		}

		if isHeadquarter {
			swiftCodes = append(swiftCodes, models.SwiftCode{
				SwiftCode:     swiftCode,
				CountryISO2:   countryISO2,
				BankName:      bankName,
				Address:       address,
				CountryName:   countryName,
				IsHeadquarter: true,
				Branches:      []models.SwiftBranch{},
			})
		} else {
			swiftCodes = append(swiftCodes, models.SwiftCode{
				SwiftCode:     swiftCode,
				CountryISO2:   countryISO2,
				BankName:      bankName,
				Address:       address,
				CountryName:   countryName,
				IsHeadquarter: false,
			})
		}
	}

	return swiftCodes, nil
}

// Extracts and normalizes (uppercase/trim) the values for SWIFT code, ISO2, bank name, address, and country name from a CSV row.
func extractRecordData(record []string, fieldIndexes map[string]int) (string, string, string, string, string) {
	swiftCode := strings.TrimSpace(strings.ToUpper(record[fieldIndexes["SWIFT CODE"]]))
	countryISO2 := strings.TrimSpace(strings.ToUpper(record[fieldIndexes["COUNTRY ISO2 CODE"]]))
	bankName := strings.ToUpper(record[fieldIndexes["NAME"]])
	address := strings.ToUpper(record[fieldIndexes["ADDRESS"]])
	countryName := strings.ToUpper(record[fieldIndexes["COUNTRY NAME"]])

	return swiftCode, countryISO2, bankName, address, countryName
}

// Validates the extracted data from a record against SWIFT code rules and the provided country map.
func validateRecord(swiftCode, countryISO2, countryName string, countries map[string]models.Country) error {
	if err := utils.ValidateSwiftCode(swiftCode); err != nil {
		return fmt.Errorf("invalid SWIFT code: %v", err)
	}
	if err := utils.ValidateCountryISO2(countryISO2); err != nil {
		return fmt.Errorf("invalid ISO2 country code: %v", err)
	}
	if err := utils.ValidateCountryExistence(countryISO2, countries); err != nil {
		return fmt.Errorf("invalid country: %v", err)
	}
	if err := utils.ValidateCountryNameMatch(countryISO2, countryName, countries); err != nil {
		return fmt.Errorf("country name mismatch: %v", err)
	}

	return nil
}
