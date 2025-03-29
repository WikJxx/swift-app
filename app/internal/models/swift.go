// Package models defines data models used throughout the application,
// including SWIFT code structures for headquarters and branches.
package models

// SwiftCode represents a SWIFT headquarter record, including address, bank details,
// and any associated branch information.
type SwiftCode struct {
	Address       string        `json:"address"`
	BankName      string        `json:"bankName"`
	CountryISO2   string        `json:"countryISO2"`
	CountryName   string        `json:"countryName"`
	IsHeadquarter bool          `json:"isHeadquarter"`
	SwiftCode     string        `json:"swiftCode"`
	Branches      []SwiftBranch `json:"branches"`
}

// SwiftBranch represents a branch of a SWIFT headquarter.
type SwiftBranch struct {
	Address       string `json:"address"`
	BankName      string `json:"bankName"`
	CountryISO2   string `json:"countryISO2"`
	CountryName   string `json:"countryName,omitempty"`
	IsHeadquarter bool   `json:"isHeadquarter"`
	SwiftCode     string `json:"swiftCode"`
}
