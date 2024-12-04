// routes/routes.go
package routes

import (
    "github.com/gin-gonic/gin"
    "github.com/gin-contrib/sessions"
    "go.mongodb.org/mongo-driver/mongo"
    "github.com/JinHyeokOh01/FSSP-Server/controllers"
    "github.com/JinHyeokOh01/FSSP-Server/db"
)

// routes/routes.go
func SetupRoutes(r *gin.Engine, client *mongo.Client) {
    authController := controllers.NewAuthController(client.Database(db.DatabaseName))

    auth := r.Group("/auth")
    {
        auth.POST("/register", authController.Register)
        auth.POST("/login", authController.Login)
        auth.POST("/logout", AuthRequired(), authController.Logout)
        auth.GET("/current-user", AuthRequired(), authController.GetCurrentUser)
    }

    api := r.Group("/api")
    {
        api.GET("/search", controllers.NaverSearchHandler)
        api.POST("/chat", controllers.HandleChat)

        // 레스토랑 관리 라우트 - AuthRequired 미들웨어 추가
        restaurants := api.Group("/restaurants").Use(AuthRequired())
        {
            restaurants.GET("", func(c *gin.Context) {
                session := sessions.Default(c)
                email := session.Get("email")
                if email == nil {
                    c.JSON(401, gin.H{"error": "인증이 필요합니다"})
                    c.Abort()
                    return
                }
                db.GetRestaurantsHandler(c, client, email.(string))
            })

            restaurants.POST("", func(c *gin.Context) {
                session := sessions.Default(c)
                email := session.Get("email")
                if email == nil {
                    c.JSON(401, gin.H{"error": "인증이 필요합니다"})
                    c.Abort()
                    return
                }
                db.UpdateListHandler(c, client, email.(string))
            })

            restaurants.DELETE("", func(c *gin.Context) {
                session := sessions.Default(c)
                email := session.Get("email")
                if email == nil {
                    c.JSON(401, gin.H{"error": "인증이 필요합니다"})
                    c.Abort()
                    return
                }
                db.DeleteListHandler(c, client, email.(string))
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