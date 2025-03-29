// validators_test.go contains unit tests for validation functions related to SWIFT codes and country data.
package utils

import (
	"testing"

	"swift-app/internal/models"

	"github.com/stretchr/testify/assert"
)

func TestValidateCountryISO2(t *testing.T) {
	assert.NoError(t, ValidateCountryISO2("US"))
	assert.Error(t, ValidateCountryISO2("USA"))
	assert.Error(t, ValidateCountryISO2("U1"))
	assert.Error(t, ValidateCountryISO2("us"))
}

func TestValidateSwiftCode(t *testing.T) {
	assert.Error(t, ValidateSwiftCode("ABCDEFGHG"))
	assert.NoError(t, ValidateSwiftCode("ABCDEFGHXXX"))
	assert.Error(t, ValidateSwiftCode("ABCDEFG"))
	assert.Error(t, ValidateSwiftCode("ABCDEFGHIJKL"))
}

func TestValidateSwiftCodeSuffix(t *testing.T) {
	assert.NoError(t, ValidateSwiftCodeSuffix("ABCDEFGHXXX", true))
	assert.Error(t, ValidateSwiftCodeSuffix("ABCDEFGHABC", true))

	assert.NoError(t, ValidateSwiftCodeSuffix("ABCDEFGHABC", false))
	assert.Error(t, ValidateSwiftCodeSuffix("ABCDEFGHXXX", false))
}

func TestValidateCountryExistence(t *testing.T) {
	countries := map[string]models.Country{
		"PL": {ISO2: "PL", Name: "POLAND"},
	}

	assert.NoError(t, ValidateCountryExistence("PL", countries))
	err := ValidateCountryExistence("US", countries)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "country ISO2 'US' not found")
}

func TestValidateCountryNameMatch(t *testing.T) {
	countries := map[string]models.Country{
		"PL": {ISO2: "PL", Name: "POLAND"},
	}

	assert.NoError(t, ValidateCountryNameMatch("PL", "POLAND", countries))
	assert.NoError(t, ValidateCountryNameMatch("PL", "poland", countries))

	err := ValidateCountryNameMatch("PL", "Germany", countries)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "country name 'Germany' does not match ISO2 'PL'")
}

func TestLoadAndValidateCountry(t *testing.T) {
	countries, err := LoadAndValidateCountry("PL")
	assert.NoError(t, err)
	assert.NotEmpty(t, countries)

	_, err = LoadAndValidateCountry("XX")
	assert.Error(t, err)
}

func TestLoadAndValidateCountryWithName(t *testing.T) {
	countries, err := LoadAndValidateCountryWithName("PL", "POLAND")
	assert.NoError(t, err)
	assert.NotEmpty(t, countries)

	_, err = LoadAndValidateCountryWithName("PL", "GERMANY")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "country name 'GERMANY' does not match ISO2 'PL'")
}
