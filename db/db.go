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

type UserInfo struct {
	Email string
	Name  string
	List []string
}

func UserAddDB(client *mongo.Client, email string, name string){
	coll := client.Database("FSSP_DB").Collection("users")
	doc := UserInfo{Email: email, Name: name, List: []string{}}
	result, err := coll.InsertOne(context.TODO(), doc)
	if err != nil{
		panic(err)
	}
	fmt.Printf("Inserted document with _id: %v\n", result.InsertedID)
}

//식당 저장할 때 추가하는 함수
func UpdateListDB(client *mongo.Client, email string, newItem string){
	coll := client.Database("FSSP_DB").Collection("users")

	// 새로운 아이템을 추가할 문서 찾기
    filter := bson.M{"Email": email}

    // 업데이트할 내용 정의
    update := bson.M{
        "$push": bson.M{"List": newItem}, // List에 newItem 추가
    }

    // 문서 업데이트
    _, err := coll.UpdateOne(context.TODO(), filter, update)
    if err != nil {
        panic(err)
    }
}