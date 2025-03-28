package utils

import (
	"context"
	"strings"
	"swift-app/internal/errors"
	"swift-app/internal/models"
	"unicode"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// Function ensures the ISO2 country code has exactly two uppercase letters.
func ValidateCountryISO2(iso2 string) error {
	if len(iso2) != 2 {
		return errors.Wrap(errors.ErrBadRequest, "country ISO2 must be 2 characters")
	}
	for _, r := range iso2 {
		if r < 'A' || r > 'Z' {
			return errors.Wrap(errors.ErrBadRequest, "country ISO2 must contain only letters")
		}
	}
	return nil
}

// Function verifies if the provided ISO2 code exists in the given countries map.
func ValidateCountryExistence(iso2 string, countries map[string]models.Country) error {
	if _, ok := countries[iso2]; !ok {
		return errors.Wrap(errors.ErrBadRequest, "country ISO2 '%s' not found", iso2)
	}
	return nil
}

// Function loads countries from file and checks if provided ISO2 is valid.
func LoadAndValidateCountry(iso2 string) (map[string]models.Country, error) {
	if err := ValidateCountryISO2(iso2); err != nil {
		return nil, err
	}

	countries, err := LoadCountries()
	if err != nil {
		return nil, errors.Wrap(errors.ErrInternal, "error loading country data")
	}

	if err := ValidateCountryExistence(iso2, countries); err != nil {
		return nil, err
	}

	return countries, nil
}

// Function validates the length of the provided SWIFT code.
func ValidateSwiftCode(swiftCode string) error {
	if len(swiftCode) == 0 {
		return errors.Wrap(errors.ErrBadRequest, "missing SWIFT code")
	}
	if len(swiftCode) != 8 && len(swiftCode) != 11 {
		return errors.Wrap(errors.ErrBadRequest, "SWIFT code must be 8 or 11 characters")
	}
	for _, r := range swiftCode {
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) {
			return errors.Wrap(errors.ErrBadRequest, "SWIFT code can only contain letters and digits (no spaces or special characters)")
		}
	}
	return nil
}

// Function checks if the provided country name matches the expected name from ISO2 code.
func ValidateCountryNameMatch(iso2 string, inputName string, countries map[string]models.Country) error {
	expected := countries[iso2].Name
	if !strings.EqualFold(inputName, expected) {
		return errors.Wrap(errors.ErrBadRequest, "country name '%s' does not match ISO2 '%s'", inputName, iso2)
	}
	return nil
}

// Function retrieves the headquarter SWIFT entry for a given SWIFT code.
func GetHeadquarterBySwiftCode(db *mongo.Collection, swiftCode string) (*models.SwiftCode, error) {
	headquarterCode := swiftCode[:8] + "XXX"
	var headquarter models.SwiftCode
	err := db.FindOne(context.Background(), bson.M{
		FieldSwiftCode:     headquarterCode,
		FieldIsHeadquarter: true,
	}).Decode(&headquarter)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			if strings.HasSuffix(swiftCode, "XXX") {
				return nil, errors.Wrap(errors.ErrNotFound, "headquarter not found: %s", swiftCode)
			}
			return nil, errors.Wrap(errors.ErrNotFound, "cannot perform action with branch '%s' because its headquarter '%s' is missing", swiftCode, headquarterCode)
		}
		return nil, errors.Wrap(errors.ErrInternal, "database error while searching for headquarter")
	}

	return &headquarter, nil
}

// Function loads and validates a country by ISO2 and verifies the provided country name.
func LoadAndValidateCountryWithName(iso2, inputName string) (map[string]models.Country, error) {
	countries, err := LoadAndValidateCountry(iso2)
	if err != nil {
		return nil, err
	}
	if err := ValidateCountryNameMatch(iso2, inputName, countries); err != nil {
		return nil, err
	}
	return countries, nil
}

// Function checks whether SWIFT code suffix matches expected format for HQ or branch.
func ValidateSwiftCodeSuffix(swiftCode string, isHeadquarter bool) error {
	if isHeadquarter && !strings.HasSuffix(swiftCode, "XXX") {
		return errors.Wrap(errors.ErrBadRequest, "HQ SWIFT code must end with 'XXX'")
	}
	if !isHeadquarter && strings.HasSuffix(swiftCode, "XXX") {
		return errors.Wrap(errors.ErrBadRequest, "branch SWIFT code cannot end with 'XXX'")
	}
	return nil
}
