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

	var swiftCodes []models.SwiftCode
	for i, record := range records[1:] {
		if len(record) != 8 {
			return nil, fmt.Errorf("record on line %d: wrong number of fields", i+2)
		}

		swift := models.SwiftCode{
			CountryISO2: strings.ToUpper(record[0]),
			SwiftCode:   record[1],
			CodeType:    record[2],
			BankName:    record[3],
			Address:     record[4],
			TownName:    record[5],
			CountryName: strings.ToUpper(record[6]),
			TimeZone:    record[7],
		}

		swift.IsHeadquarter = strings.HasSuffix(strings.TrimSpace(record[1]), "XXX")

		swiftCodes = append(swiftCodes, swift)
	}

	return swiftCodes, nil
}
