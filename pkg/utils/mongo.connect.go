package utils

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var collection *mongo.Collection

func MongoConnect(input string)*mongo.Collection {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
	databasePwd := os.Getenv("MONGO_DATABASE_CONECTION")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	clientOptions := options.Client().ApplyURI(databasePwd)

	client, _ := mongo.Connect(ctx, clientOptions)
	fmt.Println("Pinged your deployment. You successfully connected to MongoDB!")
	collection = client.Database("PropertiesData").Collection(input)
	return collection
}