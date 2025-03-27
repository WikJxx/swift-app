package services

import (
	"context"
	"fmt"
	"strings"
	"swift-app/internal/errors"
	"swift-app/internal/models"
	"swift-app/internal/utils"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type SwiftCodeService struct {
	DB *mongo.Collection
}

func NewSwiftCodeService(db *mongo.Collection) *SwiftCodeService {
	return &SwiftCodeService{DB: db}
}

// Function retrieves details of a specific SWIFT code, including headquarter or branch information.
func (s *SwiftCodeService) GetSwiftCodeDetails(swiftCode string) (*models.SwiftCode, error) {
	swiftCode = strings.ToUpper(swiftCode)
	if err := utils.ValidateSwiftCode(swiftCode); err != nil {
		return nil, err
	}

	var swiftCodeDetails models.SwiftCode
	err := s.DB.FindOne(context.Background(), bson.M{"swiftCode": swiftCode}).Decode(&swiftCodeDetails)
	if err == nil {
		return &swiftCodeDetails, nil
	}

	headquarter, err := utils.GetHeadquarterBySwiftCode(s.DB, swiftCode)
	if err != nil {
		return nil, err
	}

	for _, branch := range headquarter.Branches {
		if branch.SwiftCode == swiftCode {
			return &models.SwiftCode{
				Address:       branch.Address,
				BankName:      branch.BankName,
				CountryISO2:   branch.CountryISO2,
				CountryName:   headquarter.CountryName,
				IsHeadquarter: false,
				SwiftCode:     branch.SwiftCode,
			}, nil
		}
	}

	return nil, fmt.Errorf("%w: no branch found for SWIFT code %s", errors.ErrNotFound, swiftCode)
}

// Function retrieves all SWIFT codes and branches associated with a specified country ISO2 code.
func (s *SwiftCodeService) GetSwiftCodesByCountry(countryISO2 string) (*models.CountrySwiftCodesResponse, error) {
	countryISO2 = strings.ToUpper(countryISO2)
	_, err := utils.LoadAndValidateCountry(countryISO2)
	if err != nil {
		return nil, err
	}

	cursor, err := s.DB.Find(context.Background(), bson.M{"countryISO2": countryISO2})
	if err != nil {
		return nil, fmt.Errorf("%w: error retrieving SWIFT codes for country %s", errors.ErrInternal, countryISO2)
	}
	defer cursor.Close(context.Background())

	var swiftCodes []models.SwiftCode
	if err = cursor.All(context.Background(), &swiftCodes); err != nil {
		return nil, fmt.Errorf("%w: error decoding SWIFT codes for country %s", errors.ErrInternal, countryISO2)
	}
	if len(swiftCodes) == 0 {
		return nil, fmt.Errorf("%w: no SWIFT codes found for country %s", errors.ErrNotFound, countryISO2)
	}

	var countryName = swiftCodes[0].CountryName
	var allSwiftCodes []models.SwiftCode
	var allBranchCodes []models.SwiftCode
	swiftCodeSet := make(map[string]bool)

	for _, swiftCode := range swiftCodes {
		if swiftCode.IsHeadquarter && !swiftCodeSet[swiftCode.SwiftCode] {
			allSwiftCodes = append(allSwiftCodes, models.SwiftCode{
				Address:       swiftCode.Address,
				BankName:      swiftCode.BankName,
				CountryISO2:   swiftCode.CountryISO2,
				IsHeadquarter: true,
				SwiftCode:     swiftCode.SwiftCode,
			})
			swiftCodeSet[swiftCode.SwiftCode] = true
		}

		for _, branch := range swiftCode.Branches {
			if !swiftCodeSet[branch.SwiftCode] {
				allBranchCodes = append(allBranchCodes, models.SwiftCode{
					Address:       branch.Address,
					BankName:      branch.BankName,
					CountryISO2:   branch.CountryISO2,
					IsHeadquarter: false,
					SwiftCode:     branch.SwiftCode,
				})
				swiftCodeSet[branch.SwiftCode] = true
			}
		}
	}

	allSwiftCodes = append(allSwiftCodes, allBranchCodes...)

	return &models.CountrySwiftCodesResponse{
		CountryISO2: countryISO2,
		CountryName: countryName,
		SwiftCodes:  allSwiftCodes,
	}, nil
}

// Function adds a new SWIFT code (headquarter or branch) to the database with proper validation.
func (s *SwiftCodeService) AddSwiftCode(request *models.SwiftCode) (string, error) {
	request.SwiftCode = strings.ToUpper(request.SwiftCode)
	request.CountryISO2 = strings.ToUpper(request.CountryISO2)
	request.CountryName = strings.ToUpper(request.CountryName)

	if err := utils.ValidateSwiftCode(request.SwiftCode); err != nil {
		return "", err
	}
	if err := utils.ValidateSwiftCodeSuffix(request.SwiftCode, request.IsHeadquarter); err != nil {
		return "", err
	}
	_, err := utils.LoadAndValidateCountryWithName(request.CountryISO2, request.CountryName)
	if err != nil {
		return "", err
	}

	doc := bson.M{
		"swiftCode":     request.SwiftCode,
		"bankName":      request.BankName,
		"address":       request.Address,
		"countryISO2":   request.CountryISO2,
		"countryName":   request.CountryName,
		"isHeadquarter": request.IsHeadquarter,
		"branches":      request.Branches,
	}
	if request.IsHeadquarter && request.Branches == nil {
		doc["branches"] = []bson.M{}
	}

	if request.IsHeadquarter {
		if err := s.DB.FindOne(context.Background(), bson.M{"swiftCode": request.SwiftCode}).Err(); err == nil {
			return "", fmt.Errorf("%w: headquarter SWIFT code already exists", errors.ErrBadRequest)
		}
		if _, err := s.DB.InsertOne(context.Background(), doc); err != nil {
			return "", fmt.Errorf("%w: error inserting SWIFT code into the database", errors.ErrInternal)
		}
		return "Headquarter SWIFT code added successfully", nil
	}

	headquarter, err := utils.GetHeadquarterBySwiftCode(s.DB, request.SwiftCode)
	if err != nil {
		return "", err
	}
	if request.CountryISO2 != headquarter.CountryISO2 {
		return "", fmt.Errorf("%w: branch SWIFT code countryISO does not match headquarter countryISO", errors.ErrBadRequest)
	}
	for _, branch := range headquarter.Branches {
		if branch.SwiftCode == request.SwiftCode {
			return "", fmt.Errorf("%w: branch SWIFT code already exists", errors.ErrBadRequest)
		}
	}

	branch := bson.M{
		"swiftCode":     request.SwiftCode,
		"bankName":      request.BankName,
		"address":       request.Address,
		"countryISO2":   request.CountryISO2,
		"isHeadquarter": false,
	}
	if _, err := s.DB.UpdateOne(context.Background(),
		bson.M{"swiftCode": headquarter.SwiftCode},
		bson.M{"$push": bson.M{"branches": branch}}); err != nil {
		return "", fmt.Errorf("%w: error updating headquarter with branch", errors.ErrInternal)
	}

	return "Branch SWIFT code added to headquarter successfully", nil
}

// Function deletes an existing SWIFT code (headquarter and its branches, or single branch) from the database.
func (s *SwiftCodeService) DeleteSwiftCode(swiftCode string) (string, error) {
	swiftCode = strings.ToUpper(swiftCode)
	if err := utils.ValidateSwiftCode(swiftCode); err != nil {
		return "", err
	}
	isHeadquarter := strings.HasSuffix(swiftCode, "XXX")
	if err := utils.ValidateSwiftCodeSuffix(swiftCode, isHeadquarter); err != nil {
		return "", err
	}

	if isHeadquarter {
		var headquarter models.SwiftCode
		err := s.DB.FindOne(context.Background(), bson.M{"swiftCode": swiftCode, "isHeadquarter": true}).Decode(&headquarter)
		if err == mongo.ErrNoDocuments {
			return "", fmt.Errorf("%w: headquarter %s not found, cannot delete", errors.ErrNotFound, swiftCode)
		}
		if err != nil {
			return "", fmt.Errorf("%w: error while checking headquarter %s", errors.ErrInternal, swiftCode)
		}

		result, err := s.DB.DeleteMany(context.Background(), bson.M{
			"$or": []bson.M{
				{"swiftCode": swiftCode},
				{"swiftCode": bson.M{"$regex": fmt.Sprintf("^%s", swiftCode[:8])}},
			},
		})
		if err != nil {
			return "", fmt.Errorf("%w: error deleting headquarter %s and its branches", errors.ErrInternal, swiftCode)
		}
		if result.DeletedCount == 0 {
			return "", fmt.Errorf("%w: headquarter %s was not deleted", errors.ErrInternal, swiftCode)
		}
		return fmt.Sprintf("Deleted hadquarter %s and its branches", swiftCode), nil
	}

	headquarterCode := swiftCode[:8] + "XXX"
	var headquarter models.SwiftCode
	err := s.DB.FindOne(context.Background(), bson.M{"swiftCode": headquarterCode, "isHeadquarter": true}).Decode(&headquarter)
	if err == mongo.ErrNoDocuments {
		return "", fmt.Errorf("%w: branch %s not found and its headquarter %s does not exist", errors.ErrNotFound, swiftCode, headquarterCode)
	}
	if err != nil {
		return "", fmt.Errorf("%w: error checking headquarter for branch %s", errors.ErrInternal, swiftCode)
	}

	update, err := s.DB.UpdateOne(
		context.Background(),
		bson.M{"swiftCode": headquarterCode},
		bson.M{"$pull": bson.M{"branches": bson.M{"swiftCode": swiftCode}}},
	)
	if err != nil {
		return "", fmt.Errorf("%w: error deleting branch %s", errors.ErrInternal, swiftCode)
	}
	if update.ModifiedCount == 0 {
		return "", fmt.Errorf("%w: branch %s not found under headquarter %s", errors.ErrNotFound, swiftCode, headquarterCode)
	}

	return fmt.Sprintf("Branch %s deleted successfully", swiftCode), nil
}
