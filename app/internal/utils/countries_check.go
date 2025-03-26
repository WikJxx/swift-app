package utils

import (
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"swift-app/internal/models"
)

func LoadCountries() (map[string]models.Country, error) {
	_, currentFilePath, _, _ := runtime.Caller(0)
	projectRootDir := filepath.Join(filepath.Dir(currentFilePath), "../..")
	filePath := filepath.Join(projectRootDir, "internal", "resources", "countries.csv")

	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	countries := make(map[string]models.Country)

	header := records[0]
	for i, h := range header {
		header[i] = strings.TrimSpace(h)
	}

	fieldIndexes := map[string]int{
		"ISO2": -1,
		"NAME": -1,
	}

	for i, field := range header {
		upperField := strings.ToUpper(field)
		if _, exists := fieldIndexes[upperField]; exists {
			fieldIndexes[upperField] = i
		}
	}
	if fieldIndexes["ISO2"] == -1 || fieldIndexes["NAME"] == -1 {
		return nil, fmt.Errorf("missing required fields: ISO2 or NAME")
	}

	for _, record := range records[1:] {
		iso2 := strings.ToUpper(record[fieldIndexes["ISO2"]])
		name := strings.ToUpper(record[fieldIndexes["NAME"]])

		countries[iso2] = models.Country{
			ISO2: iso2,
			Name: name,
		}
	}

	return countries, nil
}
