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
type SwiftResponse struct {
	Address       string            `json:"address"`
	BankName      string            `json:"bankName"`
	CountryISO2   string            `json:"countryISO2"`
	CountryName   string            `json:"countryName"`
	IsHeadquarter bool              `json:"isHeadquarter"`
	SwiftCode     string            `json:"swiftCode"`
	Branches      []SwiftBranchResp `json:"branches,omitempty"`
}

type SwiftBranchResp struct {
	Address       string `json:"address"`
	BankName      string `json:"bankName"`
	CountryISO2   string `json:"countryISO2"`
	IsHeadquarter bool   `json:"isHeadquarter"`
	SwiftCode     string `json:"swiftCode"`
}
