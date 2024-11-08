package db

import (
	"context"
	"fmt"

	"github.com/joho/godotenv"
	"os"
	"log"
  
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/bson"
  )
  
func ConnectDB() (*mongo.Client, error) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	userURI := os.Getenv("MONGO_URI")
	// Use the SetServerAPIOptions() method to set the version of the Stable API on the client
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(userURI).SetServerAPIOptions(serverAPI)
  
	  // Create a new client and connect to the server
	client, err := mongo.Connect(context.TODO(), opts)
	if err != nil {
		panic(err)
	}

	// Send a ping to confirm a successful connection
	if err := client.Database("admin").RunCommand(context.TODO(), bson.D{{"ping", 1}}).Err(); err != nil {
	  panic(err)
	}
	fmt.Println("Pinged your deployment. You successfully connected to MongoDB!")
	return client, nil
}

func DisconnectDB(client *mongo.Client){
	err := client.Disconnect(context.TODO())
    if err != nil {
        panic(err)
    }
    fmt.Println("Successfully disconnected from MongoDB!")
}


type Book struct {
	Title  string
	Author string
}


func UploadDB(client *mongo.Client, a string){
	coll := client.Database("FSSP_DB").Collection("users")
	doc := Book{Title: "Atonement", Author: a}
	result, err := coll.InsertOne(context.TODO(), doc)
	if err != nil{
		panic(err)
	}
	fmt.Printf("Inserted document with _id: %v\n", result.InsertedID)

}