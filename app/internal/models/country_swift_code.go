package models

type CountrySwiftCodesResponse struct {
	CountryISO2 string        `json:"countryISO2"`
	CountryName string        `json:"countryName"`
	SwiftCodes  []SwiftBranch `json:"swiftCodes"`
}
