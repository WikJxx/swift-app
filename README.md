# Swift APP
This project is a Go-based REST API designed to manage SWIFT (Bank Identifier Codes) data. It allows users to parse, store, and query SWIFT codes efficiently, providing endpoints for retrieving, adding, and deleting SWIFT codes. The application is containerized using Docker and uses MongoDB as the database for fast and scalable data storage.

---
## Chosen Technologies
#### Go (Golang) – Go was used to implement the backend API due to its simplicity and high performance. It enabled building a clean and well-structured project, where the logic is separated into packages for services, handlers, parsers, and validation. Thanks to Go’s strong typing and error handling, it was easier to catch issues early and ensure safe processing of financial data such as SWIFT codes.

#### MongoDB – The document-oriented nature of MongoDB was ideal for storing SWIFT codes, especially due to the nested relationship between headquarters and their branches. The flexible schema allowed each headquarter to include its branches directly, making queries more efficient and avoiding the complexity of joins required in relational databases. It also simplified the insertion and validation logic during data import from CSV files.

#### Gin – Gin was selected for its performance and ease of use in building RESTful APIs. It provided built-in routing, parameter binding, and middleware support, which helped quickly set up a well-structured HTTP interface for interacting with SWIFT codes (for example: adding, retrieving, deleting, filtering by country).

#### Testcontainers-Go – This library allowed writing integration tests with real MongoDB containers started dynamically during testing. It ensured that database-related logic (such as inserting HQ and branch codes, or validating collection state) was properly verified without relying on a shared environment or external DB. This improved test reliability and confidence in the application’s behavior across different setups.

#### Docker – Docker was used to containerize the entire application stack, including the Go backend and MongoDB. This ensured the app could be run and tested consistently across machines and environments. It also simplified deployment and eliminated external dependencies during local development.


---

## Table of Contents

1. [Features](#features)
2. [Project Structure](#project-structure)
3. [Getting Started](#getting-started)
   - [Environment Setup](#environment-setup)
   - [Run with Docker (Recommended)](#run-with-docker-recommended)
   - [Run Locally](#run-locally)
4. [API Endpoints](#api-endpoints)
5. [Swagger UI & Documentation](#swagger-ui--documentation)
6. [Testing](#testing)
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

- **Documentation**
  - Auto-generated Swagger UI (`/swagger/index.html`).
  - Developer-friendly GoDoc documentation for internal packages.

- **Containerization**:
  - Dockerized for easy deployment and scalability.
  - Includes a `docker-compose.yml` file for running the application and MongoDB together.

- **Testing**:
    - Comprehensive unit and integration tests for all endpoints and core functionality.

---

## Project Structure
```bash
swift-api/
│
├── app/
│   │
│   ├── cmd/                      # Application entry points
│   │   ├── server/               # HTTP server initialization
│   │   │   ├── server.go         # Gin server setup
│   │   ├── router/               # API routing
│   │   │   ├── router.go         # API route definitions
│   │   │   ├── router_test.go    # Integration tests for routing layer
│   │
│   ├── internal/                 # Business logic
│   │   ├── errors/               # Custom application errors with HTTP status mapping
│   │   │   ├── errors.go
│   │   ├── models/               # Data models
│   │   │   ├── country_swift_code.go  # Response model: SWIFT codes grouped by country
│   │   │   ├── country.go             # Model for country ISO2 and name
│   │   │   ├── import_summary.go      # Model summarizing import statistics
│   │   │   ├── response.go            # Generic message response model
│   │   │   ├── swift.go               # SWIFT code and branch model
│   │   ├── services/              # Business logic implementation
│   │   │   ├── swift_service.go        # SWIFT code operations (add, get, delete)
│   │   │   ├── swift_service_test.go  # Unit tests for service layer
│   │   ├── resources/            # Static resources (CSV, data)
│   │   │   ├── countries.csv         # Country name ↔ ISO2 mapping file
│   │   ├── testutils/            # Shared test setup and MongoDB helpers
│   │   │   ├── testmain.go           # Mongo container & collection bootstrap for tests
│   │   ├── utils/                # Utility helpers
│   │   │   ├── countries_check.go     # ISO2/country name matching and lookup
│   │   │   ├── constans.go           # Constants used across the application
│   │   │   ├── countries_check_test.go # Tests for country validation
│   │   │   ├── validators.go         # Validation logic for SWIFT and country fields
│   │   │   ├── validators_test.go    # Unit tests for validation functions
│
│   ├── pkg/                     # General-purpose packages
│   │   ├── csv/                 # CSV parsing logic
│   │   │   ├── parser.go            # SWIFT code CSV parser
│   │   │   ├── parser_test.go       # Tests for CSV parsing
│   │   ├── data/                # Sample CSV files for parser
│   │   │   ├── ...csv
│
│   ├── database/                # MongoDB connection & setup
│   │   ├── mongo.go                 # Database init, index creation, data saving
│   │   ├── mongo_test.go           # MongoDB-related unit tests
│
│   ├── api/                     # HTTP handlers for API
│   │   ├── v1/                  # API versioning (v1)
│   │   │   ├── swift_handler.go       # Endpoint logic for SWIFT codes
│   │   │   ├── swift_handler_test.go # Unit tests for handler logic
│
│   ├── integration/             # High-level integration tests (end-to-end)
│   │   ├── swift_test.go            # Integration tests combining API + DB
│
│   ├── initialization/          # CSV import and DB connection bootstrap
│   │   ├── initialization.go        # Load data & connect DB from .env
│   │   ├── initialization_test.go   # Unit tests for initialization logic
│   │   ├── data/                    # Embedded CSVs for test/import
│   │   │   ├── ...csv
│
│   ├── docs/                    # Swagger documentation files (auto-generated)
│   │   ├── ...                   # swagger.json, swagger.yaml, etc.
│
│   ├── Dockerfile               # Container definition for API app
│   ├── docker-compose.yml       # Docker Compose setup for API + Mongo
│   ├── main.go                  # Application entry point
│   ├── go.mod                   # Go module definitions
│   ├── go.sum                   # Dependency checksums
│   ├── .env                     # Environment variables (used by app)
│   └── README.md                # Project documentation


```

---

## Getting Started

### Prerequisites

#### For Docker Deployment:
- [Docker Desktop](https://www.docker.com/products/docker-desktop/) (required for `docker-compose`)
- `.env` file with correct values (see below)

#### For Local Development:
- [Go](https://go.dev/dl/) 1.21 or higher
- [MongoDB](https://www.mongodb.com/try/download/community) installed locally and running (`mongod`)
- [Git](https://git-scm.com/)

> If you're on Windows and plan to run tests using **Testcontainers**, Docker must support **Linux containers** (not rootless).

---

### Option 1: Docker Deployment (Recommended)

This is the fastest and simplest way to get started — no need to install Go or MongoDB locally.

#### 1. Clone the Repository
```bash
git clone https://github.com/WikJxx/swift-app.git
cd swift-app/app
```
#### 2. .env file

Make sure you have the .env file in the main folder (/app). 

If you don't, you can create it

Example:
```bash
MONGO_URI=mongodb://mongo:27017
MONGO_DB=swiftDB
MONGO_COLLECTION=swiftCodes
CSV_PATH=./pkg/data/Interns_2025_SWIFT_CODES.csv
HOST=localhost
PORT=8080
```

#### 3. Start the App with Docker Compose

```bash
docker compose up
```

The app will be available at:
http://localhost:8080
Swagger UI: http://localhost:8080/swagger/index.html

---

### Option 2: Running locally
#### 1. Clone the repository:
```bash 
https://github.com/WikJxx/swift-app.git
cd swift-api/app
```
#### 2. .env file

Make sure you have the .env file in the main folder (/app). 

If you don't, you can create it

Example:
```bash
MONGO_URI=mongodb://localhost:27017
MONGO_DB=swiftDB
MONGO_COLLECTION=swiftCodes
CSV_PATH=./pkg/data/Interns_2025_SWIFT_CODES.csv
HOST=localhost
PORT=8080
```

#### 3.  Start MongoDB (if not already running)
```bash
mongod
```

You can also run mongo in docker container.
```bash
docker-compose up -d mongo
```
If using docker remember to change  MONGO_URI path in .env, for example: 
```bash
mongodb://mongo:27017
```

#### 4. Install dependencies:
```bash
go mod tidy
```

#### 5. Run the application
```bash
go run main.go
```
The app will be available at:
http://localhost:8080
Swagger UI: http://localhost:8080/swagger/index.html


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
## Swagger UI & Documentation

This project uses [Swaggo](https://github.com/swaggo/swag) to generate interactive API documentation.

### Setup

No manual installation is required — all necessary dependencies are already included in `go.mod`.

> **Note:** Swagger docs are automatically generated using Swaggo and served via `github.com/swaggo/gin-swagger`.

### Accessing Swagger UI

Once the app is running (default on port `8080`), you can access the Swagger UI at: http://localhost:8080/swagger/index.html


The interface provides a complete list of all available API endpoints, expected request/response formats, and example usages.

### GoDoc-style Documentation

To browse Go documentation locally:

1. Run the local Go documentation server:

```bash
go doc -http=:6060
```
2. Open your browser and go to:
```bash
http://localhost:6060/pkg/
```
You can now explore all your Go packages (e.g., internal/services, internal/models, etc.) using a classic GoDoc-style UI.

> **Note:** GoDoc server runs on http://localhost:6060 by default.

To view the documentation of the services package via terminal:
```bash
go doc swift-app/internal/services
```
Or interactively via browser:
```bash
http://localhost:6060/pkg/swift-app/internal/services/
```

---

## Testing

The project includes comprehensive **unit** and **integration tests** for all key components: parsers, services, API endpoints, and database logic.

---

### Requirements

To run the test suite, **Docker must be running** on your machine. This is required because some tests use [Testcontainers](https://github.com/testcontainers/testcontainers-go) to spin up a real MongoDB instance in a container.

---

### Run All Tests

Make sure Docker is running, then execute:

```bash
cd swift-app/app
go mod tidy 
go test ./...
```
This will automatically:

- Start a MongoDB container

- Connect to it via testcontainers-go

- Run all unit and integration tests

- Tear down the container afterwards

### Test Coverage

The project includes a comprehensive suite of both unit and integration tests to ensure correctness, data consistency, and API behavior. 


| Module/Location          | Description                                                              |
|--------------------------|--------------------------------------------------------------------------|
| `pkg/csv`                | Validates CSV parsing and SWIFT data extraction                          |
| `internal/utils`         | Ensures correctness of validators (e.g., ISO2 format, SWIFT format)      |
| `internal/services`      | Verifies business logic and MongoDB operations (insert, find, delete)    |
| `database/`              | Tests low-level MongoDB logic and collection indexing                    |
| `cmd/router`             | Covers API routing and HTTP response handling                            |
| `integration/`           | Full end-to-end HTTP tests of the API, including data storage & retrieval|


---
## Environment Variables
| Variable            | Description                          | Default Value                          |
|---------------------|--------------------------------------|----------------------------------------|
| `MONGO_URI`         | MongoDB connection URI               | `mongodb://mongo:27017`           |
| `MONGO_DB`          | MongoDB database name                | `swiftDB`                             |
| `MONGO_COLLECTION`  | MongoDB collection name              | `swiftCodes`                          |
| `CSV_PATH`          | Path to the CSV file with SWIFT data | `./pkg/data/Interns_2025_SWIFT_CODES.csv` |
| `HOST`              | Default host                         | `localhost`                           |
| `PORT`              | Default port                         | `8080`                               |

> **Note**: All environment variables are loaded from a `.env` file located in the root directory of the project.  
> Make sure this file exists before running the application locally or via Docker.