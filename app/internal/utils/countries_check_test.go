package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// Package contains unit tests for country data loading and validation from CSV file.
func TestLoadCountries(t *testing.T) {
	countries, err := LoadCountries()
	assert.NoError(t, err, "LoadCountries should not return an error")
	assert.NotEmpty(t, countries, "Countries map should not be empty")

	poland, exists := countries["PL"]
	assert.True(t, exists, "Poland (PL) should exist in countries")
	assert.Equal(t, "POLAND", poland.Name, "Country name should be POLAND")
}

func TestLoadCountries_MissingFile(t *testing.T) {
	originalPath := getCountriesCSVPath
	getCountriesCSVPath = func() string {
		return "nonexistent.csv"
	}
	defer func() { getCountriesCSVPath = originalPath }()

	_, err := LoadCountries()
	assert.Error(t, err, "Should return error when CSV file does not exist")
}
