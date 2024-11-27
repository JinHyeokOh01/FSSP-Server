package main

import (
    "github.com/JinHyeokOh01/FSSP-Server/db"
    "github.com/JinHyeokOh01/FSSP-Server/routes"
    "github.com/JinHyeokOh01/FSSP-Server/controllers"
    "github.com/joho/godotenv"
    "log"
    "fmt"
)

func init() {
    // 환경 변수 로드
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

    // Google OAuth 초기화
    controllers.InitGoogleOAuth()

    // Gin 라우터 설정
    r := routes.Routes(client)
    
    port := 5000
    log.Printf("Server is running on port: %d", port)
    
    if err := r.Run(fmt.Sprintf(":%d", port)); err != nil {
        log.Fatal("Server failed to start:", err)
    }
}