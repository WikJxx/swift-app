package models

type SwiftCode struct {
	CountryISO2   string    `json:"countryISO2"`
	SwiftCode     string    `json:"swiftCode"`
	BankName      string    `json:"bankName"`
	Address       string    `json:"address"`
	CountryName   string    `json:"countryName"`
	IsHeadquarter bool      `json:"isHeadquarter"`
	Branches      *[]string `json:"branches,omitempty"`
}
