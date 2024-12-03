// main.go
package main

import (
    "context"
    "log"
    "time"
    "github.com/gin-gonic/gin"
    "github.com/gin-contrib/sessions"
    "github.com/gin-contrib/sessions/cookie"
    "github.com/joho/godotenv"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
    "your-project/routes"
)

func main() {
    // 환경 변수 로드
    if err := godotenv.Load(); err != nil {
        log.Fatal("Error loading .env file")
    }

    // MongoDB 연결
    client, err := connectDB()
    if err != nil {
        log.Fatal("Failed to connect to database:", err)
    }
    defer client.Disconnect(context.Background())

    // Gin 엔진 초기화
    r := gin.Default()

    // 세션 미들웨어 설정
    store := cookie.NewStore([]byte("your-secret-key")) // 실제 운영환경에서는 환경변수에서 가져오세요
    store.Options(sessions.Options{
        MaxAge:   60 * 60 * 24 * 7, // 7일
        Path:     "/",
        Secure:   false, // HTTPS 사용 시 true로 설정
        HttpOnly: true,
    })
    r.Use(sessions.Sessions("session", store))

    // CORS 설정
    r.Use(corsMiddleware())

    // 라우트 설정
    routes.SetupRoutes(r, client.Database("your_database_name"))

    // 서버 시작
    if err := r.Run(":8080"); err != nil {
        log.Fatal("Failed to start server:", err)
    }
}

func connectDB() (*mongo.Client, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
    client, err := mongo.Connect(ctx, clientOptions)
    if err != nil {
        return nil, err
    }

    err = client.Ping(ctx, nil)
    if err != nil {
        return nil, err
    }

    return client, nil
}

func corsMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
        c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
        c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
        c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

        if c.Request.Method == "OPTIONS" {
            c.AbortWithStatus(204)
            return
        }

        c.Next()
    }
}