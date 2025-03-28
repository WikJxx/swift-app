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
	err := s.DB.FindOne(context.Background(), bson.M{utils.FieldSwiftCode: swiftCode}).Decode(&swiftCodeDetails)
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

	return nil, errors.Wrap(errors.ErrNotFound, "no branch found for SWIFT code %s", swiftCode)
}

// Function retrieves all SWIFT codes and branches associated with a specified country ISO2 code.
func (s *SwiftCodeService) GetSwiftCodesByCountry(countryISO2 string) (*models.CountrySwiftCodesResponse, error) {
	countryISO2 = strings.ToUpper(countryISO2)
	_, err := utils.LoadAndValidateCountry(countryISO2)
	if err != nil {
		return nil, err
	}

	cursor, err := s.DB.Find(context.Background(), bson.M{utils.FieldCountryISO2: countryISO2})
	if err != nil {
		return nil, errors.Wrap(errors.ErrInternal, "error retrieving SWIFT codes for country %s", countryISO2)
	}
	defer cursor.Close(context.Background())

	var swiftCodes []models.SwiftCode
	if err = cursor.All(context.Background(), &swiftCodes); err != nil {
		return nil, errors.Wrap(errors.ErrInternal, "error decoding SWIFT codes for country %s", countryISO2)
	}
	if len(swiftCodes) == 0 {
		return nil, errors.Wrap(errors.ErrNotFound, "no SWIFT codes found for country %s", countryISO2)
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
		utils.FieldSwiftCode:     request.SwiftCode,
		utils.FieldBankName:      request.BankName,
		utils.FieldAddress:       request.Address,
		utils.FieldCountryISO2:   request.CountryISO2,
		utils.FieldCountryName:   request.CountryName,
		utils.FieldIsHeadquarter: true,
		utils.FieldBranches:      request.Branches,
	}
	if request.IsHeadquarter && request.Branches == nil {
		doc[utils.FieldBranches] = []bson.M{}
	}

	if request.IsHeadquarter {
		if err := s.DB.FindOne(context.Background(), bson.M{utils.FieldSwiftCode: request.SwiftCode}).Err(); err == nil {
			return "", errors.Wrap(errors.ErrBadRequest, "headquarter SWIFT code already exists")
		}
		if _, err := s.DB.InsertOne(context.Background(), doc); err != nil {
			return "", errors.Wrap(errors.ErrInternal, "error inserting SWIFT code into the database")
		}
		return "headquarter SWIFT code added successfully", nil
	}

	headquarter, err := utils.GetHeadquarterBySwiftCode(s.DB, request.SwiftCode)
	if err != nil {
		return "", err
	}
	if request.CountryISO2 != headquarter.CountryISO2 {
		return "", errors.Wrap(errors.ErrBadRequest, "branch countryISO does not match headquarter countryISO")
	}
	for _, branch := range headquarter.Branches {
		if branch.SwiftCode == request.SwiftCode {
			return "", errors.Wrap(errors.ErrBadRequest, "branch SWIFT code already exists")
		}
	}

	branch := bson.M{
		utils.FieldSwiftCode:     request.SwiftCode,
		utils.FieldBankName:      request.BankName,
		utils.FieldAddress:       request.Address,
		utils.FieldCountryISO2:   request.CountryISO2,
		utils.FieldIsHeadquarter: false,
	}
	if _, err := s.DB.UpdateOne(context.Background(),
		bson.M{utils.FieldSwiftCode: headquarter.SwiftCode},
		bson.M{"$push": bson.M{utils.FieldBranches: branch}}); err != nil {
		return "", errors.Wrap(errors.ErrInternal, "error updating headquarter with branch")
	}

	return "branch SWIFT code added to headquarter successfully", nil
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
		err := s.DB.FindOne(context.Background(), bson.M{utils.FieldSwiftCode: swiftCode, utils.FieldIsHeadquarter: true}).Decode(&headquarter)
		if err == mongo.ErrNoDocuments {
			return "", errors.Wrap(errors.ErrNotFound, "headquarter %s not found, cannot delete", swiftCode)
		}
		if err != nil {
			return "", errors.Wrap(errors.ErrInternal, "error while checking headquarter %s", swiftCode)
		}

		result, err := s.DB.DeleteMany(context.Background(), bson.M{
			"$or": []bson.M{
				{utils.FieldSwiftCode: swiftCode},
				{utils.FieldSwiftCode: bson.M{"$regex": fmt.Sprintf("^%s", swiftCode[:8])}},
			},
		})
		if err != nil {
			return "", errors.Wrap(errors.ErrInternal, "error deleting headquarter %s and its branches", swiftCode)
		}
		if result.DeletedCount == 0 {
			return "", errors.Wrap(errors.ErrInternal, "headquarter %s was not deleted", swiftCode)
		}
		return fmt.Sprintf("deleted hadquarter %s and its branches", swiftCode), nil
	}

	headquarterCode := swiftCode[:8] + "XXX"
	var headquarter models.SwiftCode
	err := s.DB.FindOne(context.Background(), bson.M{utils.FieldSwiftCode: headquarterCode, utils.FieldIsHeadquarter: true}).Decode(&headquarter)
	if err == mongo.ErrNoDocuments {
		return "", errors.Wrap(errors.ErrNotFound, "branch %s not found and its headquarter %s does not exist", swiftCode, headquarterCode)
	}
	if err != nil {
		return "", errors.Wrap(errors.ErrInternal, "error checking headquarter for branch %s", swiftCode)
	}

	update, err := s.DB.UpdateOne(
		context.Background(),
		bson.M{utils.FieldSwiftCode: headquarterCode},
		bson.M{"$pull": bson.M{utils.FieldBranches: bson.M{utils.FieldSwiftCode: swiftCode}}},
	)
	if err != nil {
		return "", errors.Wrap(errors.ErrInternal, "error deleting branch %s", swiftCode)
	}
	if update.ModifiedCount == 0 {
		return "", errors.Wrap(errors.ErrNotFound, "branch %s not found under headquarter %s", swiftCode, headquarterCode)
	}
	return fmt.Sprintf("branch %s deleted successfully", swiftCode), nil
}
