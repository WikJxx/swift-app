package models

type SwiftCode struct {
	CountryISO2   string
	SwiftCode     string
	CodeType      string
	BankName      string
	Address       string
	TownName      string
	CountryName   string
	TimeZone      string
	IsHeadquarter bool
}
