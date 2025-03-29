// Package main is the entry point of the Swift App API.
// It initializes the environment, database, imports data from CSV, and starts the server.
package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"swift-app/cmd/server"
	"swift-app/database"
	_ "swift-app/docs"
	"swift-app/initialization"
	"syscall"

	"github.com/joho/godotenv"
)

// main initializes environment variables, connects to the database, imports SWIFT data, and starts the HTTP server.
// @title Swift App API
// @version 1.0
// @description This is a Swift Code management API.
func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, using default values.")
	}

	mongoURI := os.Getenv("MONGO_URI")
	mongoDB := os.Getenv("MONGO_DB")
	mongoCollection := os.Getenv("MONGO_COLLECTION")
	csvPath := os.Getenv("CSV_PATH")

	err = initialization.InitializeDatabase(mongoURI, mongoDB, mongoCollection)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	summary, err := initialization.ImportData(csvPath)
	if err != nil {
		log.Fatalf("Failed to import data: %v", err)
	}

	log.Printf(`
Data import complete.
Headquarters added: %d
Skipped HQs (duplicates): %d

Branches added: %d
Duplicate branches: %d
Branches with missing HQ: %d
All skipped branches: %d
`, summary.HQAdded, summary.HQSkipped, summary.BranchesAdded, summary.BranchesDuplicate, summary.BranchesMissingHQ, summary.BranchesSkipped)

	handleShutdown()
	fmt.Println("Starting application...")
	server.StartServer()
}

// handleShutdown listens for termination signals and gracefully closes the MongoDB connection.
func handleShutdown() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-sigs
		fmt.Println("\nShutdown requested, closing database connection...")
		err := database.CloseMongoDB()
		if err != nil {
			log.Println("Error closing database:", err)
		} else {
			fmt.Println("Database connection closed.")
		}
		os.Exit(0)
	}()
}
