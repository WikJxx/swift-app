package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"swift-app/cmd/server"
	"swift-app/database"
	initialization "swift-app/initialization"
	"syscall"

	"github.com/joho/godotenv"
)

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

	hqAdded, hqSkipped, branchesAdded, branchesDuplicate, branchesMissingHQ, branchesSkipped, err := initialization.ImportData(csvPath)
	if err != nil {
		log.Fatalf("Failed to import data: %v", err)
	}

	log.Printf("\nData import complete. \nHeadquarters added: %d\nNumber of skipped headquarters due to beeing duplicates with records in database: %d\n\n", hqAdded, hqSkipped)
	log.Printf("\nBranches added: %d\nNumber of skipped branches due to beeing duplicates with records in database: %d\nNumber of branches skipped due to missing headquarters: %d\nNumber of all skipped branches: %d\n", branchesAdded, branchesDuplicate, branchesMissingHQ, branchesSkipped)

	log.Println("Data imported successfully.")
	handleShutdown()
	fmt.Println("Starting application...")
	server.StartServer()
}
