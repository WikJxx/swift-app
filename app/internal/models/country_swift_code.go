package models

// CountrySwiftCodesResponse represents the response format for retrieving SWIFT codes by country.
// It includes the country ISO2 code, full country name, and a list of associated branch SWIFT codes.
type CountrySwiftCodesResponse struct {
	CountryISO2 string        `json:"countryISO2"`
	CountryName string        `json:"countryName"`
	SwiftCodes  []SwiftBranch `json:"swiftCodes"`
}
