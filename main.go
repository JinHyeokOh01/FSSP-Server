package main

import (
    "github.com/JinHyeokOh01/FSSP-Server/db"
    "github.com/JinHyeokOh01/FSSP-Server/routes"
    "github.com/JinHyeokOh01/FSSP-Server/controllers"
    "github.com/joho/godotenv"
    "github.com/gin-gonic/gin"
    "log"
    "fmt"
)

func init() {
    if err := godotenv.Load(); err != nil {
        log.Fatal("Error loading .env file")
    }
}

func main() {
    // MongoDB 연결
    client, err := db.ConnectDB()
    if err != nil {
        panic(err)
    }
    defer db.DisconnectDB(client)

    // Gin 엔진 생성
    r := gin.Default()

    // Auth 초기화 (세션 및 OAuth 설정)
    if err := controllers.InitAuth(r); err != nil {
        log.Fatal("Failed to initialize auth:", err)
    }

    // 라우터 설정
    routes.SetupRoutes(r, client)
    
    port := 5000
    log.Printf("Server is running on port: %d", port)
    
    if err := r.Run(fmt.Sprintf(":%d", port)); err != nil {
        log.Fatal("Server failed to start:", err)
    }
}