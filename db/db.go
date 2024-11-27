package db

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
  
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/bson"
  )
  
func ConnectDB() (*mongo.Client, error) {
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

func UserAddDB(client *mongo.Client, email string, name string) error {
    coll := client.Database("FSSP_DB").Collection("users")
    doc := UserInfo{Email: email, Name: name, List: []string{}}
    
    result, err := coll.InsertOne(context.TODO(), doc)
    if err != nil {
        return fmt.Errorf("failed to insert user: %v", err)
    }
    
    fmt.Printf("Inserted document with _id: %v\n", result.InsertedID)
    return nil
}

type ListRequest struct {
    NewItem string `json:"NewItem"`
}

func UpdateListHandler(c *gin.Context, client *mongo.Client, email string) {
    var request ListRequest

    // 요청 바디에서 데이터 바인딩
    if err := c.ShouldBindJSON(&request); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    // 데이터베이스 업데이트
    UpdateListDB(client, email, request.NewItem)

    c.JSON(http.StatusOK, gin.H{"status": "success"})
}

//식당 저장할 때 추가하는 함수
func UpdateListDB(client *mongo.Client, email string, newItem string){
	coll := client.Database("FSSP_DB").Collection("users")

	// 새로운 아이템을 추가할 문서 찾기
    filter := bson.M{"email": email}

    // 업데이트할 내용 정의
    update := bson.M{
        "$push": bson.M{"list": newItem}, // List에 newItem 추가
    }

    // 문서 업데이트
    _, err := coll.UpdateOne(context.TODO(), filter, update)
    if err != nil {
        panic(err)
    }
}

func DeleteListHandler(c *gin.Context, client *mongo.Client, email string) {
    var request ListRequest

    // 요청 바디에서 데이터 바인딩
    if err := c.ShouldBindJSON(&request); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    // 데이터베이스 업데이트
    DeleteListDB(client, email, request.NewItem)

    c.JSON(http.StatusOK, gin.H{"status": "success"})
}

//식당 삭제하는 함수
func DeleteListDB(client *mongo.Client, email string, newItem string){
	coll := client.Database("FSSP_DB").Collection("users")

	// 아이템을 삭제할 문서 찾기
    filter := bson.M{"email": email}

    // 업데이트할 내용 정의
    update := bson.M{
        "$pull": bson.M{"list": newItem}, // List에서 newItem 삭제
    }

    // 문서 업데이트
    _, err := coll.UpdateOne(context.TODO(), filter, update)
    if err != nil {
        panic(err)
    }
}