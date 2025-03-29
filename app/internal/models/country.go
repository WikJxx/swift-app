package models

// Country represents a country with its ISO2 code and full name.
type Country struct {
	ISO2 string `json:"iso2"`
	Name string `json:"name"`
}
