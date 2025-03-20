package models

type SwiftCode struct {
	Address       string        `json:"address"`
	BankName      string        `json:"bankName"`
	CountryISO2   string        `json:"countryISO2"`
	CountryName   string        `json:"countryName"`
	IsHeadquarter bool          `json:"isHeadquarter"`
	SwiftCode     string        `json:"swiftCode"`
	Branches      []SwiftBranch `json:"branches,omitempty"`
}

type SwiftBranch struct {
	Address       string `json:"address"`
	BankName      string `json:"bankName"`
	CountryISO2   string `json:"countryISO2"`
	IsHeadquarter bool   `json:"isHeadquarter"`
	SwiftCode     string `json:"swiftCode"`
}
