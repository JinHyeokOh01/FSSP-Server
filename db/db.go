package db

import (
    "context"
    "fmt"
    "net/http"
    "os"
    "time"

    "github.com/gin-gonic/gin"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/bson/primitive"
)

type Restaurant struct {
    ID        primitive.ObjectID `bson:"_id,omitempty"`
    Name      string            `bson:"name"`
    Category  string            `bson:"category"`
    Address   string            `bson:"address"`
    CreatedAt time.Time         `bson:"created_at"`
}

type UserInfo struct {
    ID          primitive.ObjectID `bson:"_id,omitempty"`
    Email       string            `bson:"email"`
    Password    string            `bson:"password"`
    Name        string            `bson:"name"`
    CreatedAt   time.Time         `bson:"created_at"`
    UpdatedAt   time.Time         `bson:"updated_at"`
    Restaurants []Restaurant      `bson:"restaurants"`
}

type ListRequest struct {
    Name     string `json:"name" binding:"required"`
    Category string `json:"category" binding:"required"`
    Address  string `json:"address" binding:"required"`
}

func ConnectDB() (*mongo.Client, error) {
    userURI := os.Getenv("MONGO_URI")
    if userURI == "" {
        return nil, fmt.Errorf("MONGO_URI environment variable is not set")
    }

    serverAPI := options.ServerAPI(options.ServerAPIVersion1)
    opts := options.Client().ApplyURI(userURI).SetServerAPIOptions(serverAPI)

    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    client, err := mongo.Connect(ctx, opts)
    if err != nil {
        return nil, fmt.Errorf("failed to connect to MongoDB: %v", err)
    }

    if err := client.Database("admin").RunCommand(ctx, bson.D{{"ping", 1}}).Err(); err != nil {
        return nil, fmt.Errorf("failed to ping MongoDB: %v", err)
    }

    fmt.Println("Successfully connected to MongoDB")
    return client, nil
}

func DisconnectDB(client *mongo.Client) error {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    if err := client.Disconnect(ctx); err != nil {
        return fmt.Errorf("failed to disconnect from MongoDB: %v", err)
    }

    fmt.Println("Successfully disconnected from MongoDB")
    return nil
}

func UserAddDB(client *mongo.Client, email, hashedPassword, name string) error {
    coll := client.Database("FSSP_DB").Collection("userlist")

    // Check for existing user
    count, err := coll.CountDocuments(context.TODO(), bson.M{"email": email})
    if err != nil {
        return fmt.Errorf("failed to check existing user: %v", err)
    }
    if count > 0 {
        return fmt.Errorf("user with email %s already exists", email)
    }

    now := time.Now()
    doc := UserInfo{
        Email:       email,
        Password:    hashedPassword,
        Name:        name,
        CreatedAt:   now,
        UpdatedAt:   now,
        Restaurants: []Restaurant{},
    }

    _, err = coll.InsertOne(context.TODO(), doc)
    if err != nil {
        return fmt.Errorf("failed to insert user: %v", err)
    }

    return nil
}

func UpdateListHandler(c *gin.Context, client *mongo.Client, email string) {
    var request ListRequest
    if err := c.ShouldBindJSON(&request); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
        return
    }

    if err := UpdateListDB(client, email, request.Name, request.Category, request.Address); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Restaurant added successfully"})
}

func UpdateListDB(client *mongo.Client, email, name, category, address string) error {
    coll := client.Database("FSSP_DB").Collection("users")

    restaurant := Restaurant{
        ID:        primitive.NewObjectID(),
        Name:      name,
        Category:  category,
        Address:   address,
        CreatedAt: time.Now(),
    }

    filter := bson.M{"email": email}
    update := bson.M{
        "$push": bson.M{"restaurants": restaurant},
        "$set":  bson.M{"updated_at": time.Now()},
    }

    result, err := coll.UpdateOne(context.TODO(), filter, update)
    if err != nil {
        return fmt.Errorf("failed to update user restaurants: %v", err)
    }

    if result.MatchedCount == 0 {
        return fmt.Errorf("no user found with email: %s", email)
    }

    return nil
}

func DeleteListHandler(c *gin.Context, client *mongo.Client, email string) {
    var request ListRequest
    if err := c.ShouldBindJSON(&request); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
        return
    }

    if err := DeleteListDB(client, email, request.Name); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Restaurant deleted successfully"})
}

func DeleteListDB(client *mongo.Client, email, name string) error {
    coll := client.Database("FSSP_DB").Collection("users")

    filter := bson.M{"email": email}
    update := bson.M{
        "$pull": bson.M{"restaurants": bson.M{"name": name}},
        "$set":  bson.M{"updated_at": time.Now()},
    }

    result, err := coll.UpdateOne(context.TODO(), filter, update)
    if err != nil {
        return fmt.Errorf("failed to delete restaurant: %v", err)
    }

    if result.MatchedCount == 0 {
        return fmt.Errorf("no user found with email: %s", email)
    }

    return nil
}

func GetRestaurantsHandler(c *gin.Context, client *mongo.Client, email string) {
    restaurants, err := GetRestaurantsDB(client, email)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "message":     "Restaurants retrieved successfully",
        "restaurants": restaurants,
    })
}

func GetRestaurantsDB(client *mongo.Client, email string) ([]Restaurant, error) {
    coll := client.Database("FSSP_DB").Collection("users")
    
    var user UserInfo
    err := coll.FindOne(context.TODO(), bson.M{"email": email}).Decode(&user)
    if err != nil {
        if err == mongo.ErrNoDocuments {
            return nil, fmt.Errorf("no user found with email: %s", email)
        }
        return nil, fmt.Errorf("failed to retrieve user data: %v", err)
    }
    
    return user.Restaurants, nil
}