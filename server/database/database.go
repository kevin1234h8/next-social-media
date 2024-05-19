package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"social/project/initializer"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func Connect() (*mongo.Client, error) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	logger := initializer.InitializeLogger()
	clientOptions := options.Client().ApplyURI("mongodb://127.0.0.1:27017/socialmedianext")
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		initializer.LogMessage(logger, "Connect", fmt.Sprintf("Error connecting to MongoDB: %v", err))
		return nil, err
	}

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		initializer.LogMessage(logger, "Connect", fmt.Sprintf("Error pinging MongoDB: %v", err))
		return nil, err
	}

	initializer.LogMessage(logger, "Connect", "Connected to MongoDB!")
	logger.Print("database url:", os.Getenv("DATABASE_URL"))
	return client, nil
}
