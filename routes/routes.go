// routes/routes.go
package routes

import (
    "github.com/gin-gonic/gin"
    "github.com/gin-contrib/sessions"
    "go.mongodb.org/mongo-driver/mongo"
    "github.com/JinHyeokOh01/FSSP-Server/controllers"
)

func SetupRoutes(r *gin.Engine, db *mongo.Database) {
    authController := controllers.NewAuthController(db)

    auth := r.Group("/auth")
    {
        auth.POST("/register", authController.Register)
        auth.POST("/login", authController.Login)
        auth.POST("/logout", AuthRequired(), authController.Logout)
        auth.GET("/current-user", AuthRequired(), authController.GetCurrentUser)
    }

    // ChatGPT API 라우트 추가
    r.POST("/chat", controllers.HandleChat)

    api := r.Group("/api")
    api.Use(AuthRequired())
    {
        // 인증이 필요한 API 엔드포인트
    }
}

func AuthRequired() gin.HandlerFunc {
    return func(c *gin.Context) {
        session := sessions.Default(c)
        userId := session.Get("userId")
        if userId == nil {
            c.JSON(401, gin.H{"error": "인증이 필요합니다"})
            c.Abort()
            return
        }
        c.Next()
    }
}