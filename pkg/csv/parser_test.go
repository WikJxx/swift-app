package csv

import (
	"os"
	"testing"
)

func TestLoadSwiftCodes(t *testing.T) {
	data := `SWIFT CODE,COUNTRY ISO2 CODE,NAME,ADDRESS,COUNTRY NAME
AAAABBBXXX,US,First Bank,123 First St,USA
AAAABBBXXX,US,First Bank,123 First St,USA
AAAABBB123,US,Second Bank,456 Second St,USA
AAAABBB123,US,Second Bank,456 Second St,USA
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

	expectedCount := 1
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

	if _, exists := swiftCodeSet["AAAABBBXXX"]; !exists {
		t.Errorf("SwiftCode AAAABBBXXX not found")
	}
}

func TestLoadSwiftCodesWithStringCSV(t *testing.T) {
	const testCSV = `COUNTRY ISO2 CODE,SWIFT CODE,TYPE,NAME,ADDRESS,TOWN NAME,COUNTRY NAME,TIMEZONE
,,,,"FOREST ZUBRA 1, FLOOR 1 WARSZAWA, MAZOWIECKIE, 01-066",WARSZAWA,POLAND,Europe/Warsaw
PL,TPEOPLPWMEGXXX,BIC11,PEKAO TOWARZYSTWO FUNDUSZY INWESTYCYJNYCH SPOLKA AKCYJNA,"FOREST ZUBRA 1, FLOOR 1 WARSZAWA, MAZOWIECKIE, 01-066",WARSZAWA,POLAND,Europe/Warsaw
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

	if len(swiftCodes) != 1 {
		t.Fatalf("expected 1 swift code, got %d", len(swiftCodes))
	}

	if swiftCodes[0].BankName != "PEKAO TOWARZYSTWO FUNDUSZY INWESTYCYJNYCH SPOLKA AKCYJNA" {
		t.Errorf("expected bank name 'PEKAO TOWARZYSTWO FUNDUSZY INWESTYCYJNYCH SPOLKA AKCYJNA', got '%s'", swiftCodes[0].BankName)
	}

	if len(swiftCodes[0].Branches) != 1 {
		t.Errorf("expected 1 branch for headquarter, got %d", len(swiftCodes[0].Branches))
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
	filePath := "../test_data/Interns_2025_SWIFT_CODES.csv"
	swiftCodes, err := LoadSwiftCodes(filePath)
	if err != nil {
		t.Fatalf("error reading from file: %v", err)
	}

	if len(swiftCodes) < 1 {
		t.Fatalf("expected at least 1 SwiftCode, got %d", len(swiftCodes))
	}
}
