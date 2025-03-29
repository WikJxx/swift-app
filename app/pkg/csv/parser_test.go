// parser_test.go contains unit tests for the LoadSwiftCodes function,
package csv

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
)

// This file contains unit tests for the LoadSwiftCodes function,
func TestLoadSwiftCodes(t *testing.T) {
	data := `SWIFT CODE,COUNTRY ISO2 CODE,NAME,ADDRESS,COUNTRY NAME
AAAABBB1XXX,US,First Bank,123 First St,United States
AAAABBB1XXX,US,First Bank,123 First St,United States
AAAABBB1123,US,Second Bank,456 Second St,United States
AAAABBB1123,US,Second Bank,456 Second St,United States
`

	tmpFile, err := os.CreateTemp("", "test_swift_codes_*.csv")
	if err != nil {
		t.Fatalf("could not create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	_, err = tmpFile.WriteString(data)
	if err != nil {
		t.Fatalf("could not write to temp file: %v", err)
	}
	tmpFile.Close()

	swiftCodes, err := LoadSwiftCodes(tmpFile.Name())
	if err != nil {
		t.Fatalf("LoadSwiftCodes failed: %v", err)
	}

	expectedCount := 2
	if len(swiftCodes) != expectedCount {
		t.Errorf("expected %d unique swift codes, got %d", expectedCount, len(swiftCodes))
	}

	swiftCodeSet := make(map[string]bool)
	for _, swiftCode := range swiftCodes {
		if _, exists := swiftCodeSet[swiftCode.SwiftCode]; exists {
			t.Errorf("duplicate SwiftCode found: %s", swiftCode.SwiftCode)
		}
		swiftCodeSet[swiftCode.SwiftCode] = true

		if len(swiftCode.Branches) > 0 {
			branchSet := make(map[string]bool)
			for _, branch := range swiftCode.Branches {
				if _, exists := branchSet[branch.SwiftCode]; exists {
					t.Errorf("duplicate branch found: %s for SwiftCode %s", branch.SwiftCode, swiftCode.SwiftCode)
				}
				branchSet[branch.SwiftCode] = true
			}
		}
	}

	if _, exists := swiftCodeSet["AAAABBB1XXX"]; !exists {
		t.Errorf("SwiftCode AAAABBB1XXX not found")
	}
}

func TestLoadSwiftCodesWithBadStringCSV(t *testing.T) {
	const testCSV = `COUNTRY ISO2 CODE,SWIFT CODE,TYPE,NAME,ADDRESS,TOWN NAME,COUNTRY NAME,TIMEZONE
,,,,"FOREST ZUBRA 1, FLOOR 1 WARSZAWA, MAZOWIECKIE, 01-066",WARSZAWA,POLAND,Europe/Warsaw
PL,TPEOPLPGXXX,BIC11,PEKAO TOWARZYSTWO FUNDUSZY INWESTYCYJNYCH SPOLKA AKCYJNA,"FOREST ZUBRA 1, FLOOR 1 WARSZAWA, MAZOWIECKIE, 01-066",WARSZAWA,POLAND,Europe/Warsaw
`

	tmpFile, err := os.CreateTemp("", "swiftcodes.csv")
	if err != nil {
		t.Fatalf("failed to create temporary file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	if _, err := tmpFile.WriteString(testCSV); err != nil {
		t.Fatalf("failed to write to temporary file: %v", err)
	}

	tmpFile.Close()

	swiftCodes, err := LoadSwiftCodes(tmpFile.Name())
	if err != nil {
		t.Fatalf("error loading swift codes: %v", err)
	}
	if len(swiftCodes) != 1 {
		t.Fatalf("expected 1 headquarter entry, got %d", len(swiftCodes))
	}
	if swiftCodes[0].BankName != "PEKAO TOWARZYSTWO FUNDUSZY INWESTYCYJNYCH SPOLKA AKCYJNA" {
		t.Errorf("expected bank name 'PEKAO TOWARZYSTWO FUNDUSZY INWESTYCYJNYCH SPOLKA AKCYJNA', got '%s'", swiftCodes[0].BankName)
	}
}

func TestLoadSwiftCodesFileNotFound(t *testing.T) {
	_, err := LoadSwiftCodes("non_existent_file.csv")
	if err == nil {
		t.Fatal("expected an error for non-existent file, but got none")
	}
}

func TestLoadSwiftCodesInvalidFormat(t *testing.T) {
	const invalidCSV = `COUNTRY ISO2 CODE;SWIFT CODE;TYPE;NAME;ADDRESS;TOWN NAME;COUNTRY NAME;TIMEZONE
PL;TPEOPLPWKOP;BIC11;PEKAO;"FOREST ZUBRA 1";WARSZAWA;POLAND;Europe/Warsaw`

	tmpFile, err := os.CreateTemp("", "invalid_format.csv")
	if err != nil {
		t.Fatalf("failed to create temporary file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	if _, err := tmpFile.WriteString(invalidCSV); err != nil {
		t.Fatalf("failed to write to temporary file: %v", err)
	}
	tmpFile.Close()

	_, err = LoadSwiftCodes(tmpFile.Name())
	if err == nil {
		t.Fatal("expected an error due to invalid CSV format, but got none")
	}
}

func TestLoadSwiftCodesEmptyAddress(t *testing.T) {
	const testCSV = `COUNTRY ISO2 CODE,SWIFT CODE,TYPE,NAME,ADDRESS,TOWN NAME,COUNTRY NAME,TIMEZONE
PL,TPEOPLPGXXX,BIC11,PEKAO,,WARSZAWA,POLAND,Europe/Warsaw`

	tmpFile, err := os.CreateTemp("", "empty_address.csv")
	if err != nil {
		t.Fatalf("failed to create temporary file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	if _, err := tmpFile.WriteString(testCSV); err != nil {
		t.Fatalf("failed to write to temporary file: %v", err)
	}
	tmpFile.Close()

	swiftCodes, err := LoadSwiftCodes(tmpFile.Name())
	if err != nil {
		t.Fatalf("error loading swift codes: %v", err)
	}

	if swiftCodes[0].Address != "" {
		t.Errorf("expected empty address, got '%s'", swiftCodes[0].Address)
	}
}
func TestLoadSwiftCodesCorrectData(t *testing.T) {
	const testCSV = `COUNTRY ISO2 CODE,SWIFT CODE,TYPE,NAME,ADDRESS,TOWN NAME,COUNTRY NAME,TIMEZONE
PL,TPEOPLPWKOP,BIC11,PEKAO TOWARZYSTWO FUNDUSZY INWESTYCYJNYCH SPOLKA AKCYJNA,"FOREST ZUBRA 1, FLOOR 1 WARSZAWA, MAZOWIECKIE, 01-066",WARSZAWA,POLAND,Europe/Warsaw
PL,TPEOPLPWXXX,BIC11,PEKAO TOWARZYSTWO FUNDUSZY INWESTYCYJNYCH SPOLKA AKCYJNA,"FOREST ZUBRA 1, FLOOR 1 WARSZAWA, MAZOWIECKIE, 01-066",WARSZAWA,POLAND,Europe/Warsaw`

	tmpFile, err := os.CreateTemp("", "swiftcodes.csv")
	if err != nil {
		t.Fatalf("failed to create temporary file: %v", err)
	}

	if _, err := tmpFile.WriteString(testCSV); err != nil {
		t.Fatalf("failed to write to temporary file: %v", err)
	}
	tmpFile.Close()

	swiftCodes, err := LoadSwiftCodes(tmpFile.Name())
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	if len(swiftCodes) != 2 {
		t.Fatalf("expected 2 swift code, got %d", len(swiftCodes))
	}

	if swiftCodes[0].BankName != "PEKAO TOWARZYSTWO FUNDUSZY INWESTYCYJNYCH SPOLKA AKCYJNA" {
		t.Errorf("expected bank name 'PEKAO TOWARZYSTWO FUNDUSZY INWESTYCYJNYCH SPOLKA AKCYJNA', got '%s'", swiftCodes[0].BankName)
	}

	defer os.Remove(tmpFile.Name())
}

func TestLoadSwiftCodesWithMissingFields(t *testing.T) {
	const testCSV = `COUNTRY ISO2 CODE,SWIFT CODE,TYPE,NAME,ADDRESS,TOWN NAME,COUNTRY NAME,TIMEZONE
PL,TPEOPLPWKOP,,PEKAO TOWARZYSTWO FUNDUSZY INWESTYCYJNYCH SPOLKA AKCYJNA,"FOREST ZUBRA 1, FLOOR 1 WARSZAWA, MAZOWIECKIE, 01-066",WARSZAWA,POLAND,Europe/Warsaw`

	swiftCodes, err := LoadSwiftCodes(testCSV)
	if err == nil {
		t.Fatal("expected error due to missing required fields, but got none")
	}

	if len(swiftCodes) != 0 {
		t.Fatalf("expected no SwiftCodes due to missing fields, got %d", len(swiftCodes))
	}
}

func TestLoadSwiftCodesFromFile(t *testing.T) {
	filePath := "../data/Interns_2025_SWIFT_CODES.csv"
	swiftCodes, err := LoadSwiftCodes(filePath)
	if err != nil {
		t.Fatalf("error reading from file: %v", err)
	}

	if len(swiftCodes) < 1 {
		t.Fatalf("expected at least 1 SwiftCode, got %d", len(swiftCodes))
	}
}

func TestLoadSwiftCodes_HeaderAliases(t *testing.T) {
	_, currentFilePath, _, _ := runtime.Caller(0)
	projectRoot := filepath.Join(filepath.Dir(currentFilePath), "../..")
	resourceDir := filepath.Join(projectRoot, "internal", "resources")
	_ = os.MkdirAll(resourceDir, os.ModePerm)

	countriesCSVPath := filepath.Join(resourceDir, "countries_test.csv")
	err := os.WriteFile(countriesCSVPath, []byte("ISO2,NAME\nUS,UNITED STATES"), 0644)
	assert.NoError(t, err)

	t.Cleanup(func() {
		_ = os.Remove(countriesCSVPath)
	})

	testCases := map[string]string{
		"SWIFT CODE":  "SWIFT CODE,COUNTRY ISO2 CODE,NAME,ADDRESS,COUNTRY NAME\nAAAABBB1XXX,US,First Bank,123 First St,UNITED STATES",
		"SWIFTCODE":   "SWIFTCODE,COUNTRY ISO2 CODE,NAME,ADDRESS,COUNTRY NAME\nAAAABBB1XXX,US,First Bank,123 First St,UNITED STATES",
		"SWIFT_CODE":  "SWIFT_CODE,COUNTRY ISO2 CODE,NAME,ADDRESS,COUNTRY NAME\nAAAABBB1XXX,US,First Bank,123 First St,UNITED STATES",
		"SWIFT CODES": "SWIFT CODES,COUNTRY ISO2 CODE,NAME,ADDRESS,COUNTRY NAME\nAAAABBB1XXX,US,First Bank,123 First St,UNITED STATES",
		"SWIFT C0DE":  "SWIFT C0DE,COUNTRY ISO2 CODE,NAME,ADDRESS,COUNTRY NAME\nAAAABBB1XXX,US,First Bank,123 First St,UNITED STATES",
	}

	for label, content := range testCases {
		t.Run("Alias: "+label, func(t *testing.T) {
			tmpFile, err := os.CreateTemp("", "alias_test_*.csv")
			assert.NoError(t, err)
			defer os.Remove(tmpFile.Name())

			_, err = tmpFile.WriteString(content)
			assert.NoError(t, err)
			tmpFile.Close()

			codes, err := LoadSwiftCodes(tmpFile.Name())
			assert.NoError(t, err)
			assert.Equal(t, 1, len(codes))
			assert.Equal(t, "AAAABBB1XXX", codes[0].SwiftCode)
		})
	}
}
