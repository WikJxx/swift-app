package utils

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/testcontainers/testcontainers-go/modules/mongodb"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Package testutils sets up a MongoDB test container and initializes a MongoDB client and collection
// for use in integration tests.

var (
	MongoContainer *mongodb.MongoDBContainer
	MongoClient    *mongo.Client
	Collection     *mongo.Collection
)

func init() {
	gin.SetMode(gin.ReleaseMode)
	ctx := context.Background()
	var err error

	MongoContainer, err = mongodb.Run(ctx, "mongo:8")
	if err != nil {
		panic("Failed to start MongoDB container: " + err.Error())
	}

	uri, err := MongoContainer.ConnectionString(ctx)
	if err != nil {
		panic("Failed to retrieve MongoDB URI: " + err.Error())
	}

	MongoClient, err = mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		panic("Failed to connect to MongoDB: " + err.Error())
	}

	Collection = MongoClient.Database("swiftDB_test").Collection("swiftCodes")
}
