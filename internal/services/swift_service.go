package services

import (
	"context"
	"fmt"
	"strings"
	"swift-app/internal/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type SwiftCodeService struct {
	DB *mongo.Collection
}

func NewSwiftCodeService(db *mongo.Collection) *SwiftCodeService {
	return &SwiftCodeService{DB: db}
}

func (s *SwiftCodeService) GetSwiftCodeDetails(swiftCode string) (*models.SwiftCode, error) {
	var swiftCodeDetails models.SwiftCode

	err := s.DB.FindOne(context.Background(), bson.M{"swiftCode": swiftCode}).Decode(&swiftCodeDetails)
	if err == nil {
		return &swiftCodeDetails, nil
	}

	headquarterSwiftCode := swiftCode[:8] + "XXX"
	var headquarter models.SwiftCode
	filter := bson.M{"swiftCode": headquarterSwiftCode, "isHeadquarter": true}
	err = s.DB.FindOne(context.Background(), filter).Decode(&headquarter)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("headquarter not found for SWIFT code %s. Please add a headquarter before adding branches", headquarterSwiftCode)
		}
		return nil, fmt.Errorf("database error while searching for headquarter: %v", err)
	}

	for _, branch := range headquarter.Branches {
		if branch.SwiftCode == swiftCode {

			return &models.SwiftCode{
				Address:       branch.Address,
				BankName:      branch.BankName,
				CountryISO2:   branch.CountryISO2,
				CountryName:   headquarter.CountryName,
				IsHeadquarter: branch.IsHeadquarter,
				SwiftCode:     branch.SwiftCode,
			}, nil
		}
	}

	return nil, fmt.Errorf("no branch found for swiftCode %s", swiftCode)
}

func (s *SwiftCodeService) GetSwiftCodesByCountry(countryISO2 string) (*models.CountrySwiftCodesResponse, error) {
	countryISO2 = strings.ToUpper(countryISO2)

	cursor, err := s.DB.Find(context.Background(), bson.M{"countryISO2": countryISO2})
	if err != nil {
		return nil, fmt.Errorf("error retrieving SWIFT codes for country %s: %v", countryISO2, err)
	}
	defer cursor.Close(context.Background())

	var swiftCodes []models.SwiftCode
	if err = cursor.All(context.Background(), &swiftCodes); err != nil {
		return nil, fmt.Errorf("error decoding SWIFT codes for country %s: %v", countryISO2, err)
	}

	if len(swiftCodes) == 0 {
		return nil, fmt.Errorf("no swift codes found for country %s", countryISO2)
	}

	var countryName string
	if len(swiftCodes) > 0 {
		countryName = swiftCodes[0].CountryName
	}

	var allSwiftCodes []models.SwiftCode
	var allBranchCodes []models.SwiftCode
	swiftCodeSet := make(map[string]bool)

	for _, swiftCode := range swiftCodes {
		if swiftCode.IsHeadquarter {
			if _, exists := swiftCodeSet[swiftCode.SwiftCode]; !exists {
				allSwiftCodes = append(allSwiftCodes, models.SwiftCode{
					Address:       swiftCode.Address,
					BankName:      swiftCode.BankName,
					CountryISO2:   swiftCode.CountryISO2,
					IsHeadquarter: swiftCode.IsHeadquarter,
					SwiftCode:     swiftCode.SwiftCode,
				})
				swiftCodeSet[swiftCode.SwiftCode] = true
			}
		}

		if len(swiftCode.Branches) > 0 {
			for _, branch := range swiftCode.Branches {
				if _, exists := swiftCodeSet[branch.SwiftCode]; !exists {
					allBranchCodes = append(allBranchCodes, models.SwiftCode{
						Address:       branch.Address,
						BankName:      branch.BankName,
						CountryISO2:   branch.CountryISO2,
						IsHeadquarter: branch.IsHeadquarter,
						SwiftCode:     branch.SwiftCode,
					})
					swiftCodeSet[branch.SwiftCode] = true
				}
			}
		}
	}

	allSwiftCodes = append(allSwiftCodes, allBranchCodes...)

	response := &models.CountrySwiftCodesResponse{
		CountryISO2: countryISO2,
		CountryName: countryName,
		SwiftCodes:  allSwiftCodes,
	}

	return response, nil
}

func (s *SwiftCodeService) AddSwiftCode(swiftCodeRequest *models.SwiftCode) (map[string]string, error) {
	if len(swiftCodeRequest.SwiftCode) > 11 {
		return map[string]string{"message": "SWIFT code cannot be longer than 11 characters"}, fmt.Errorf("SWIFT code cannot be longer than 11 characters")
	}
	swiftCodeRequest.SwiftCode = strings.ToUpper(swiftCodeRequest.SwiftCode)
	swiftCodeRequest.CountryISO2 = strings.ToUpper(swiftCodeRequest.CountryISO2)
	swiftCodeRequest.CountryName = strings.ToUpper(swiftCodeRequest.CountryName)

	swiftCode := bson.M{
		"swiftCode":     swiftCodeRequest.SwiftCode,
		"bankName":      swiftCodeRequest.BankName,
		"address":       swiftCodeRequest.Address,
		"countryISO2":   swiftCodeRequest.CountryISO2,
		"countryName":   swiftCodeRequest.CountryName,
		"isHeadquarter": swiftCodeRequest.IsHeadquarter,
	}

	if swiftCodeRequest.IsHeadquarter {
		if !strings.HasSuffix(swiftCodeRequest.SwiftCode, "XXX") {
			return map[string]string{"message": "Headquarter SWIFT code must end with 'XXX'"}, fmt.Errorf("headquarter SWIFT code must end with 'XXX'")
		}

		existingHeadquarter := s.DB.FindOne(context.Background(), bson.M{"swiftCode": swiftCodeRequest.SwiftCode})
		if existingHeadquarter.Err() == nil {
			return map[string]string{"message": "Headquarter SWIFT code already exists"}, fmt.Errorf("headquarter SWIFT code already exists")
		}

		_, err := s.DB.InsertOne(context.Background(), swiftCode)
		if err != nil {
			return map[string]string{"message": "Error inserting SWIFT code into the database"}, err
		}
		return map[string]string{"message": "Headquarter SWIFT code added successfully"}, nil
	}

	if strings.HasSuffix(swiftCodeRequest.SwiftCode, "XXX") {
		return map[string]string{"message": "Branch SWIFT code cannot end with 'XXX'"}, fmt.Errorf("branch SWIFT code cannot end with 'XXX'")
	}
	headquarterSwiftCode := swiftCodeRequest.SwiftCode[:8] + "XXX"
	var headquarter models.SwiftCode
	err := s.DB.FindOne(context.Background(), bson.M{"swiftCode": headquarterSwiftCode, "isHeadquarter": true}).Decode(&headquarter)
	if err != nil {
		return map[string]string{"message": "Headquarter not found for the branch SWIFT code"}, err
	}

	if swiftCodeRequest.CountryISO2 != headquarter.CountryISO2 {
		return map[string]string{"message": "Branch SWIFT code countryISO does not match headquarter countryISO"}, fmt.Errorf("branch SWIFT code countryISO does not match headquarter countryISO")
	}

	for _, branch := range headquarter.Branches {
		if branch.SwiftCode == swiftCodeRequest.SwiftCode {
			return map[string]string{"message": "Branch SWIFT code already exists"}, fmt.Errorf("branch SWIFT code already exists")
		}
	}

	branch := models.SwiftBranch{
		SwiftCode:     swiftCodeRequest.SwiftCode,
		BankName:      swiftCodeRequest.BankName,
		Address:       swiftCodeRequest.Address,
		CountryISO2:   swiftCodeRequest.CountryISO2,
		IsHeadquarter: false,
	}

	_, err = s.DB.UpdateOne(
		context.Background(),
		bson.M{"swiftCode": headquarter.SwiftCode},
		bson.M{"$push": bson.M{"branches": branch}},
	)
	if err != nil {
		return map[string]string{"message": "Error updating headquarter with branch"}, err
	}

	return map[string]string{"message": "Branch SWIFT code added to headquarter successfully"}, nil
}
