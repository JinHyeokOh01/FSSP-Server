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

// Restaurant 정보를 담는 구조체
type Restaurant struct {
    Name     string `bson:"name"`
    Category string `bson:"category"`
    Address  string `bson:"address"`
}

// UserInfo 구조체 업데이트
type UserInfo struct {
    Email       string       `bson:"email"`
    Name        string       `bson:"name"`
    Restaurants []Restaurant `bson:"restaurants"` // List를 Restaurant 슬라이스로 변경
}

// ListRequest 구조체 업데이트
type ListRequest struct {
    Name     string `json:"name"`
    Category string `json:"category"`
    Address  string `json:"address"`
}

func ConnectDB() (*mongo.Client, error) {
    userURI := os.Getenv("MONGO_URI")
    serverAPI := options.ServerAPI(options.ServerAPIVersion1)
    opts := options.Client().ApplyURI(userURI).SetServerAPIOptions(serverAPI)

    client, err := mongo.Connect(context.TODO(), opts)
    if err != nil {
        panic(err)
    }

    if err := client.Database("admin").RunCommand(context.TODO(), bson.D{{"ping", 1}}).Err(); err != nil {
        panic(err)
    }
    fmt.Println("Pinged your deployment. You successfully connected to MongoDB!")
    return client, nil
}

func DisconnectDB(client *mongo.Client) {
    err := client.Disconnect(context.TODO())
    if err != nil {
        panic(err)
    }
    fmt.Println("Successfully disconnected from MongoDB!")
}

func UserAddDB(client *mongo.Client, email string, name string) error {
    coll := client.Database("FSSP_DB").Collection("users")
    doc := UserInfo{
        Email:       email,
        Name:        name,
        Restaurants: []Restaurant{}, // 빈 Restaurant 슬라이스로 초기화
    }

    result, err := coll.InsertOne(context.TODO(), doc)
    if err != nil {
        return fmt.Errorf("failed to insert user: %v", err)
    }

    fmt.Printf("Inserted document with _id: %v\n", result.InsertedID)
    return nil
}

func UpdateListHandler(c *gin.Context, client *mongo.Client, email string) {
    var request ListRequest

    if err := c.ShouldBindJSON(&request); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    UpdateListDB(client, email, request.Name, request.Category, request.Address)
    c.JSON(http.StatusOK, gin.H{"status": "success"})
}

func UpdateListDB(client *mongo.Client, email string, name string, category string, address string) {
    coll := client.Database("FSSP_DB").Collection("users")

    filter := bson.M{"email": email}
    update := bson.M{
        "$push": bson.M{
            "restaurants": Restaurant{
                Name:     name,
                Category: category,
                Address:  address,
            },
        },
    }

    _, err := coll.UpdateOne(context.TODO(), filter, update)
    if err != nil {
        panic(err)
    }
}

func DeleteListHandler(c *gin.Context, client *mongo.Client, email string) {
    var request ListRequest

    if err := c.ShouldBindJSON(&request); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    DeleteListDB(client, email, request.Name)
    c.JSON(http.StatusOK, gin.H{"status": "success"})
}

func DeleteListDB(client *mongo.Client, email string, name string) {
    coll := client.Database("FSSP_DB").Collection("users")

    filter := bson.M{"email": email}
    update := bson.M{
        "$pull": bson.M{
            "restaurants": bson.M{"name": name}, // 식당 이름으로 삭제
        },
    }

    _, err := coll.UpdateOne(context.TODO(), filter, update)
    if err != nil {
        panic(err)
    }
}

// GetRestaurantsHandler handles HTTP requests to get all restaurants for a user
func GetRestaurantsHandler(c *gin.Context, client *mongo.Client, email string) {
    restaurants, err := GetRestaurantsDB(client, email)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, gin.H{"restaurants": restaurants})
}

// GetRestaurantsDB retrieves all restaurants for a given user from the database
func GetRestaurantsDB(client *mongo.Client, email string) ([]Restaurant, error) {
    coll := client.Database("FSSP_DB").Collection("users")
    
    // Find the user document
    var user UserInfo
    filter := bson.M{"email": email}
    err := coll.FindOne(context.TODO(), filter).Decode(&user)
    
    if err != nil {
        if err == mongo.ErrNoDocuments {
            return []Restaurant{}, fmt.Errorf("no user found with email: %s", email)
        }
        return nil, fmt.Errorf("database error: %v", err)
    }
    
    return user.Restaurants, nil
}