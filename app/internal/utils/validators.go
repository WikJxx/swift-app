package utils

import (
	"context"
	"fmt"
	"strings"
	"swift-app/internal/errors"
	"swift-app/internal/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// Function ensures the ISO2 country code has exactly two uppercase letters.
func ValidateCountryISO2(iso2 string) error {
	if len(iso2) != 2 {
		return fmt.Errorf("%w: country ISO2 must be 2 characters", errors.ErrBadRequest)
	}
	for _, r := range iso2 {
		if r < 'A' || r > 'Z' {
			return fmt.Errorf("%w: country ISO2 must contain only letters", errors.ErrBadRequest)
		}
	}
	return nil
}

// Function verifies if the provided ISO2 code exists in the given countries map.
func ValidateCountryExistence(iso2 string, countries map[string]models.Country) error {
	if _, ok := countries[iso2]; !ok {
		return fmt.Errorf("%w: country ISO2 '%s' not found", errors.ErrBadRequest, iso2)
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
		return nil, fmt.Errorf("%w: error loading country data", errors.ErrInternal)
	}

	if err := ValidateCountryExistence(iso2, countries); err != nil {
		return nil, err
	}

	return countries, nil
}

// Function validates the length of the provided SWIFT code.
func ValidateSwiftCode(swiftCode string) error {
	if len(swiftCode) == 0 {
		return fmt.Errorf("%w: missing SWIFT code", errors.ErrBadRequest)
	}
	if len(swiftCode) < 8 || len(swiftCode) > 11 {
		return fmt.Errorf("%w: SWIFT code must be between 8 and 11 characters", errors.ErrBadRequest)
	}
	return nil
}

// Function checks if the provided country name matches the expected name from ISO2 code.
func ValidateCountryNameMatch(iso2 string, inputName string, countries map[string]models.Country) error {
	expected := countries[iso2].Name
	if !strings.EqualFold(inputName, expected) {
		return fmt.Errorf("%w: country name '%s' does not match ISO2 '%s'", errors.ErrBadRequest, inputName, iso2)
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
				return nil, fmt.Errorf("%w: headquarter not found: %s", errors.ErrNotFound, swiftCode)
			}
			return nil, fmt.Errorf("%w: can not do any actions with branch %s because its headquarter %s is missing", errors.ErrNotFound, swiftCode, headquarterCode)
		}
		return nil, fmt.Errorf("%w: database error while searching for headquarter", errors.ErrInternal)
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
		return fmt.Errorf("%w: HQ swift code must end with 'XXX'", errors.ErrBadRequest)
	}
	if !isHeadquarter && strings.HasSuffix(swiftCode, "XXX") {
		return fmt.Errorf("%w: branch swift code cannot end with 'XXX'", errors.ErrBadRequest)
	}
	return nil
}
