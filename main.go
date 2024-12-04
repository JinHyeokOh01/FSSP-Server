// main.go
package main

import (
    "log"
    "os"
    "net/http"
    "github.com/gin-gonic/gin"
    "github.com/gin-contrib/sessions"
    "github.com/gin-contrib/sessions/cookie"
    "github.com/joho/godotenv"
    "github.com/JinHyeokOh01/FSSP-Server/routes"
    "github.com/JinHyeokOh01/FSSP-Server/db"
)

func main() {
    if err := godotenv.Load(); err != nil {
        log.Fatal("Error loading .env file:", err)
    }

    sessionSecret := os.Getenv("SESSION_SECRET")
    if sessionSecret == "" {
        log.Fatal("SESSION_SECRET is not set in .env file")
    }

    client, err := db.ConnectDB()
    if err != nil {
        log.Fatal("Failed to connect to database:", err)
    }
    defer db.DisconnectDB(client)

    r := gin.Default()

    store := cookie.NewStore([]byte(sessionSecret))
    store.Options(sessions.Options{
        MaxAge:   60 * 60 * 24 * 7, // 7 days
        Path:     "/",
        Secure:   false,            // Set to true in production with HTTPS
        HttpOnly: true,
        SameSite: http.SameSiteLaxMode,
    })
    r.Use(sessions.Sessions("session", store))

    // CORS configuration
    r.Use(corsMiddleware())

    routes.SetupRoutes(r, client)

    port := ":8080"
    log.Printf("Server is starting on port%s", port)
    if err := r.Run(port); err != nil {
        log.Fatal("Failed to start server:", err)
    }
}

func corsMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Writer.Header().Set("Access-Control-Allow-Origin", "http://localhost:8080")
        c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
        c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With, Cookie")
        c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")
        c.Writer.Header().Set("Access-Control-Expose-Headers", "Set-Cookie")

        if c.Request.Method == "OPTIONS" {
            c.AbortWithStatus(204)
            return
        }

        c.Next()
    }
}