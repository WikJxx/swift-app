package csv

import (
	"os"
	"swift-app/internal/models"
	"testing"
)

const testCSV = `COUNTRY ISO2 CODE,SWIFT CODE,TYPE,NAME,ADDRESS,TOWN NAME,COUNTRY NAME,TIMEZONE
PL,TPEOPLPWKOP,BIC11,PEKAO TOWARZYSTWO FUNDUSZY INWESTYCYJNYCH SPOLKA AKCYJNA,"FOREST ZUBRA 1, FLOOR 1 WARSZAWA, MAZOWIECKIE, 01-066",WARSZAWA,POLAND,Europe/Warsaw
PL,TPEOPLPWMEG,BIC11,PEKAO TOWARZYSTWO FUNDUSZY INWESTYCYJNYCH SPOLKA AKCYJNA,"FOREST ZUBRA 1, FLOOR 1 WARSZAWA, MAZOWIECKIE, 01-066",WARSZAWA,POLAND,Europe/Warsaw
`

func TestParseSwift_ExampleData(t *testing.T) {

	tmpFile, err := os.CreateTemp("", "swift_test.csv")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())
	if _, err := tmpFile.Write([]byte(testCSV)); err != nil {
		t.Fatal(err)
	}
	tmpFile.Close()

	records, err := LoadSwiftCodes(tmpFile.Name())
	if err != nil {
		t.Fatal(err)
	}

	if len(records) != 2 {
		t.Errorf("Expected 2 records, but got %d", len(records))
	}

	expected := models.SwiftCode{
		CountryISO2:   "PL",
		SwiftCode:     "TPEOPLPWKOP",
		BankName:      "PEKAO TOWARZYSTWO FUNDUSZY INWESTYCYJNYCH SPOLKA AKCYJNA",
		Address:       "FOREST ZUBRA 1, FLOOR 1 WARSZAWA, MAZOWIECKIE, 01-066",
		CountryName:   "POLAND",
		IsHeadquarter: false, // Zakładając, że kod kończy się na "XXX"
	}

	if records[0] != expected {
		t.Errorf("Incorrect record: %+v", records[0])
	}
}
func TestParseSwift_ValidData(t *testing.T) {
	filePath := "../test_data/Interns_2025_SWIFT_CODES.csv"

	records, err := LoadSwiftCodes(filePath)
	if err != nil {
		t.Fatal(err)
	}
	if len(records) != 1061 {
		t.Errorf("Expected 1061 records, but got %d", len(records))
	}

	for i, record := range records {

		if len(record.CountryISO2) != 2 {
			t.Errorf("Record %d: Expected 2-letter country code, but got: %s", i+1, record.CountryISO2)
		}
		if len(record.SwiftCode) < 8 {
			t.Errorf("Record %d: Expected SWIFT code with at least 8 characters, but got: %s", i+1, record.SwiftCode)
		}
		if len(record.BankName) == 0 {
			t.Errorf("Record %d: Missing bank name", i+1)
		}
		if len(record.Address) == 0 {
			t.Errorf("Record %d: Missing address", i+1)
		}
		if len(record.CountryName) == 0 {
			t.Errorf("Record %d: Missing country name", i+1)
		}
	}
}
