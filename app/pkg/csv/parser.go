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
	records, err := readCSV(filePath)
	if err != nil {
		return nil, err
	}

	header := normalizeHeader(records[0])
	fieldIndexes := extractFieldIndexes(header)
	if fieldIndexes["SWIFT CODE"] == -1 {
		return nil, fmt.Errorf("missing required field: SWIFT CODE")
	}

	countries, err := countries_check.LoadCountries()
	if err != nil {
		return nil, fmt.Errorf("error loading country data: %v", err)
	}

	return processRecords(records[1:], fieldIndexes, countries)
}

func readCSV(filePath string) ([][]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.LazyQuotes = true
	return reader.ReadAll()
}

func normalizeHeader(header []string) []string {
	for i, h := range header {
		header[i] = strings.TrimSpace(h)
	}
	return header
}

func extractFieldIndexes(header []string) map[string]int {
	indexes := map[string]int{
		"SWIFT CODE":        -1,
		"COUNTRY ISO2 CODE": -1,
		"NAME":              -1,
		"ADDRESS":           -1,
		"COUNTRY NAME":      -1,
	}
	for i, field := range header {
		upper := strings.ToUpper(field)
		if _, ok := indexes[upper]; ok {
			indexes[upper] = i
		}
	}
	return indexes
}

func processRecords(records [][]string, idx map[string]int, countries map[string]models.Country) ([]models.SwiftCode, error) {
	var swiftCodes []models.SwiftCode
	hqMap := make(map[string]*models.SwiftCode)
	branchQueue := []*models.SwiftBranch{}
	uniqueHQ := make(map[string]bool)
	uniqueBranches := make(map[string]bool)

	for _, record := range records {
		swiftCode, countryISO2 := strings.ToUpper(record[idx["SWIFT CODE"]]), strings.ToUpper(record[idx["COUNTRY ISO2 CODE"]])

		if swiftCode == "" || countryISO2 == "" || len(countryISO2) != 2 || len(swiftCode) < 8 || len(swiftCode) > 11 {
			fmt.Printf("Warning: Invalid entry for code: %s\n", swiftCode)
			continue
		}

		bankName := strings.ToUpper(record[idx["NAME"]])
		address := getOptionalField(record, idx, "ADDRESS")
		countryName := strings.ToUpper(record[idx["COUNTRY NAME"]])
		isHQ := strings.HasSuffix(swiftCode, "XXX")

		country, ok := countries[countryISO2]
		if !ok {
			fmt.Printf("Warning: Invalid country ISO2 for code %s\n", swiftCode)
			continue
		}
		if !strings.EqualFold(countryName, country.Name) {
			fmt.Printf("Warning: Country name '%s' does not match ISO2 '%s'\n", countryName, countryISO2)
			continue
		}

		if isHQ {
			processHeadquarter(swiftCode, countryISO2, bankName, address, countryName, &swiftCodes, hqMap, &branchQueue, uniqueHQ)
		} else {
			processBranch(swiftCode, countryISO2, bankName, address, hqMap, &branchQueue, uniqueBranches)
		}
	}

	reportUnmatchedBranches(branchQueue)
	return swiftCodes, nil
}

func getOptionalField(record []string, idx map[string]int, key string) string {
	if i, ok := idx[key]; ok && i < len(record) {
		return strings.TrimSpace(strings.ToUpper(record[i]))
	}
	return ""
}

func processHeadquarter(code, iso2, name, address, country string,
	swiftCodes *[]models.SwiftCode,
	hqMap map[string]*models.SwiftCode,
	branchQueue *[]*models.SwiftBranch,
	unique map[string]bool) {

	if unique[code] {
		return
	}
	unique[code] = true

	hq := models.SwiftCode{
		SwiftCode:     code,
		CountryISO2:   iso2,
		BankName:      name,
		Address:       address,
		CountryName:   country,
		IsHeadquarter: true,
		Branches:      []models.SwiftBranch{},
	}

	branches := []models.SwiftBranch{}
	hqPrefix := code[:8]
	remainingQueue := []*models.SwiftBranch{}
	for _, b := range *branchQueue {
		if strings.HasPrefix(b.SwiftCode, hqPrefix) {
			branches = append(branches, *b)
		} else {
			remainingQueue = append(remainingQueue, b)
		}
	}
	hq.Branches = branches
	*branchQueue = remainingQueue

	hqMap[code] = &hq
	*swiftCodes = append(*swiftCodes, hq)
}

func processBranch(code, iso2, name, address string,
	hqMap map[string]*models.SwiftCode,
	branchQueue *[]*models.SwiftBranch,
	unique map[string]bool) {

	if unique[code] {
		return
	}
	unique[code] = true

	branch := models.SwiftBranch{
		SwiftCode:     code,
		BankName:      name,
		Address:       address,
		CountryISO2:   iso2,
		IsHeadquarter: false,
	}

	hqCode := code[:8] + "XXX"
	if hq, found := hqMap[hqCode]; found {
		hq.Branches = append(hq.Branches, branch)
	} else {
		*branchQueue = append(*branchQueue, &branch)
	}
}

func reportUnmatchedBranches(queue []*models.SwiftBranch) {
	for _, b := range queue {
		fmt.Printf("Warning: Branch %s does not have a matching headquarter.\n", b.SwiftCode)
	}
}
