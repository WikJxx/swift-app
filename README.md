# Swift APP
This project is a Go-based REST API designed to manage SWIFT (Bank Identifier Codes) data. It allows users to parse, store, and query SWIFT codes efficiently, providing endpoints for retrieving, adding, and deleting SWIFT codes. The application is containerized using Docker and uses MongoDB as the database for fast and scalable data storage.

---

## Table of Contents

1. [Features](#features)
2. [Project Structure](#project-structure)
3. [Getting Started](#getting-started)
   - [Prerequisites](#prerequisites)
   - [Installation](#installation)
   - [Running the Application](#running-the-application)
4. [API Endpoints](#api-endpoints)
5. [Testing](#testing)
6. [Docker Deployment](#docker-deployment)
7. [Environment Variables](#environment-variables)

---

## Features
- **SWIFT Code Parsing**:
  - Parses SWIFT codes from a CSV file.
  - Identifies headquarters (codes ending with "XXX") and branches.
  - Associates branches with their respective headquarters.

- **Database Storage**:
  - Uses MongoDB for efficient storage and retrieval of SWIFT codes.
  - Supports fast querying by SWIFT code or country ISO2 code.

- **REST API**:
  - Provides endpoints for retrieving, adding, and deleting SWIFT codes.
  - Supports querying SWIFT codes by country.

- **Containerization**:
  - Dockerized for easy deployment and scalability.
  - Includes a `docker-compose.yml` file for running the application and MongoDB together.

- **Testing**:
    - Comprehensive unit and integration tests for all endpoints and core functionality.

---

## Project Structure
```bash
swift-api/
│── cmd/               # Application entry points
│   ├── server/        # HTTP server initialization
│   │   ├── server.go  # Gin server setup
│   ├── router/        # API routing
│   │   ├── router.go  # API route definitions
│   │   ├── router_test.go  # API route definitions tests
│── internal/          # Business logic
│   ├── models/        # Data models
│   │   ├── swift.go   # SWIFT code structure
│   ├── services/      # Business logic implementation
│   │   ├── swift_service.go  # SWIFT code operations
│   ├── data/          # Storage
│   │   ├── countries.csv   # Data connectiong coutry names with iso2 codes
│   ├── testutils/      
│   │   ├── testmain.go # Helper for tests
│   ├── utils/        
│   │   ├── countries_check.go   # Functionality to check coutry names and iso2 codes
│── pkg/               # Utilities and helpers
│   ├── csv/           # CSV parsing
│   │   ├── parser.go  # CSV parsing logic
│   │   ├── parser_test.go  # CSV parsing tests
│   ├── test_data/     # CSV files
│   │   ├── ...csv          # CSV file for parser
│── database/          # Database connection
│   ├── mongo.go       # MongoDB initialization
│   ├── mongo_test.go  # MongoDB initialization test
│── api/               # API handlers
│   ├── v1/            # API versioning
│   │   ├── swift_handler.go  # SWIFT code handlers
│   │   ├── swift_handler_test.go  # SWIFT code handlers tests
│── integration/       # Integration tests
│   ├── integration.go # Integration test setup
│── initialization/    # Initialization storage
│   ├── initialization.go # Initialization of database connection
│── Dockerfile         # Dockerfile for containerization
│── docker-compose.yml # Docker Compose setup
│── main.go            # Application entry point
│── go.mod             # Go module dependencies
│── go.sum             # Go dependency checksums
│── .env               # Environment variables
│── README.md          # Project documentation

```

---

## Getting Started

### Prerequisites
- Go 1.24.1 or higher (or Docker for containerized deployment)

- MongoDB (or Docker for containerized deployment)

- Docker (optional, for containerized deployment)

### Installation
#### 1. Clone the repository:
```bash 
git clone https://github.com/your-username/swift-api.git
cd swift-api
```
#### 2. Install dependencies:
```bash
go mod tidy
```
### Running the Application
If you run it locally make sure you have correct variables in .env file:

```bash
MONGO_URI=mongodb://mongo:27017
MONGO_DB=swiftDB
MONGO_COLLECTION=swiftCodes
CSV_PATH=./pkg/test_data/Interns_2025_SWIFT_CODES.csv

```

#### 1. Start MongoDB:
- If using Docker:
    ```bash
    docker-compose up -d mongo
    ```
- If running locally (Make sure you have mongo added to path):
    ```bash
    mongod
    ```
#### 2. Run the application:
```bash
go run main.go
```
The APP will be available at http://localhost:8080.

---

## API Endpoints
### 1. Retrieve Details of a Single SWIFT Code
#### - GET /v1/swift-codes/{swift-code}:

- Retrieves details of a specific SWIFT code (headquarters or branch).

- #### Response Structure:
    ```bash
    {
      "address": "string",
      "bankName": "string",
      "countryISO2": "string",
      "countryName": "string",
      "isHeadquarter": bool,
      "swiftCode": "string",
      "branches": [
        {
          "address": "string",
          "bankName": "string",
          "countryISO2": "string",
          "isHeadquarter": bool,
          "swiftCode": "string"
        }
      ]
    }
    ```
### 2. Retrieve All SWIFT Codes for a Specific Country
#### - GET /v1/swift-codes/country/{countryISO2code}:

- Retrieves all SWIFT codes (headquarters and branches) for a specific country.

- #### Response Structure:
    ```bash
    {
      "countryISO2": "string",
      "countryName": "string",
      "swiftCodes": [
        {
          "address": "string",
          "bankName": "string",
          "countryISO2": "string",
          "isHeadquarter": bool,
          "swiftCode": "string"
        }
      ]
    }
    ```
### 3. Add a New SWIFT Code
#### - POST /v1/swift-codes/:

- Adds a new SWIFT code to the database.

- #### Request Structure:
    ```bash
    {
    "address": "string",
    "bankName": "string",
    "countryISO2": "string",
    "countryName": "string",
    "isHeadquarter": bool,
    "swiftCode": "string"
    }
    ```
- #### Response Structure:
    ```bash
    {
    "message": "string"
    }
    ```

    ---

### 4. Delete a SWIFT Code
#### - DELETE /v1/swift-codes/{swift-code}:

- Deletes a SWIFT code from the database.

- #### Response Structure:
    ```bash
    {
    "message": "string"
    }
    ```
    
---

## Testing
To run the tests, use the following command:
```bash
go test ./...
```
The project includes:

- Unit Tests: For CSV parsing, database operations, and service logic.

- Integration Tests: For API endpoints and database interactions.

---

## Docker Deployment
### 1. Build and start the containers:
```bash
docker-compose up 
```
### 2. The API will be available at http://localhost:8080. 

---
## Environment Variables
| Variable            | Description                          | Default Value                          |
|---------------------|--------------------------------------|----------------------------------------|
| `MONGO_URI`         | MongoDB connection URI               | `mongodb://mongo:27017`           |
| `MONGO_DB`          | MongoDB database name                | `swiftDB`                             |
| `MONGO_COLLECTION`  | MongoDB collection name              | `swiftCodes`                          |
| `CSV_PATH`          | Path to the CSV file with SWIFT data | `./pkg/test_data/Interns_2025_SWIFT_CODES.csv` |

