package main

import (
	"fmt"
	"log"
	"swift-app/cmd/server"
	"swift-app/database"
	parser "swift-app/pkg/csv"
)

func main() {
	err := database.InitMongoDB("mongodb://mongo:27017", "swiftDB", "swiftCodes")
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	swiftCodes, err := parser.LoadSwiftCodes("./pkg/test_data/Interns_2025_SWIFT_CODES.csv")
	if err != nil {
		log.Fatalf("Error loading swift codes: %v", err)
	}

	err = database.SaveSwiftCodes(swiftCodes)
	if err != nil {
		log.Fatalf("Failed to save swift codes to database: %v", err)
	}

	fmt.Println("Successfully saved SWIFT codes to MongoDB")

	server.StartServer()
}
