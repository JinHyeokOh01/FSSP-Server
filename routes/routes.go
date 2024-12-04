// routes/routes.go
package routes

import (
    "github.com/gin-gonic/gin"
    "github.com/gin-contrib/sessions"
    "go.mongodb.org/mongo-driver/mongo"
    "github.com/JinHyeokOh01/FSSP-Server/controllers"
    "github.com/JinHyeokOh01/FSSP-Server/db"
)

func SetupRoutes(r *gin.Engine, client *mongo.Client) {
    authController := controllers.NewAuthController(client.Database(db.DatabaseName))

    auth := r.Group("/auth")
    {
        auth.POST("/register", authController.Register)
        auth.POST("/login", authController.Login)
        auth.POST("/logout", AuthRequired(), authController.Logout)
        auth.GET("/current-user", AuthRequired(), authController.GetCurrentUser)
    }

    // API 라우트 그룹 생성
    api := r.Group("/api")
    {
        // 인증 없이 접근 가능한 엔드포인트들
        api.GET("/search", controllers.NaverSearchHandler)
        api.POST("/chat", controllers.HandleChat)

        // 레스토랑 관리 라우트
        restaurants := api.Group("/restaurants")
        {
            // 현재 로그인한 사용자의 레스토랑 목록 조회
            restaurants.GET("", func(c *gin.Context) {
                session := sessions.Default(c)
                email := session.Get("userEmail").(string)
                db.GetRestaurantsHandler(c, client, email)
            })

            // 새로운 레스토랑 추가
            restaurants.POST("", func(c *gin.Context) {
                session := sessions.Default(c)
                email := session.Get("userEmail").(string)
                db.UpdateListHandler(c, client, email)
            })

            // 레스토랑 삭제
            restaurants.DELETE("", func(c *gin.Context) {
                session := sessions.Default(c)
                email := session.Get("userEmail").(string)
                db.DeleteListHandler(c, client, email)
            })
        }
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
