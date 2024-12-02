package routes

import (
    "github.com/gin-gonic/gin"
    "github.com/gin-contrib/sessions"
    "github.com/JinHyeokOh01/FSSP-Server/controllers"
    "github.com/JinHyeokOh01/FSSP-Server/db"
    "go.mongodb.org/mongo-driver/mongo"
    "net/http"
)

func SetupRoutes(r *gin.Engine, client *mongo.Client) {
    // 기존 인증 관련 라우트
    r.GET("/", controllers.GoogleForm)
    r.GET("/auth/google/login", controllers.GoogleLoginHandler)
    r.GET("/auth/google/callback", func(c *gin.Context) {
        controllers.GoogleAuthCallback(c, client)
    })
    r.GET("/logout", controllers.LogoutHandler)

    // API 라우트 그룹
    apiRoutes := r.Group("/api")
    apiRoutes.Use(controllers.AuthRequired())
    {
        apiRoutes.GET("/search", controllers.NaverSearchHandler)
        apiRoutes.POST("/chat", controllers.HandleChat)
        apiRoutes.GET("/user/email", getUserEmailHandler)
        apiRoutes.GET("/user/name", getUserInfoHandler)
    }

    // 레스토랑 관리 라우트
    restaurantRoutes := r.Group("/api/restaurants")
    restaurantRoutes.Use(controllers.AuthRequired())
    {
        restaurantRoutes.GET("", func(c *gin.Context) {
            handleRestaurantRoute(c, client, db.GetRestaurantsHandler)
        })

        restaurantRoutes.POST("", func(c *gin.Context) {
            handleRestaurantRoute(c, client, db.UpdateListHandler)
        })

        restaurantRoutes.DELETE("", func(c *gin.Context) {
            handleRestaurantRoute(c, client, db.DeleteListHandler)
        })
    }
}

// 레스토랑 라우트 핸들러를 위한 헬퍼 함수
type RestaurantHandlerFunc func(*gin.Context, *mongo.Client, string)

func handleRestaurantRoute(c *gin.Context, client *mongo.Client, handler RestaurantHandlerFunc) {
    email := getCurrentUserEmail(c)
    if email == "" {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "User not logged in"})
        return
    }
    handler(c, client, email)
}

// 현재 사용자의 이메일 가져오기
func getCurrentUserEmail(c *gin.Context) string {
    session := sessions.Default(c)
    user := session.Get("user")
    if user == nil {
        return ""
    }
    userData := user.(map[string]interface{})
    return userData["email"].(string)
}

// 현재 사용자의 이름 가져오기
func getCurrentUserName(c *gin.Context) string {
    session := sessions.Default(c)
    user := session.Get("user")
    if user == nil {
        return ""
    }
    userData := user.(map[string]interface{})
    if name, exists := userData["given_name"]; exists {
        return name.(string)
    }
    return ""
}

func getUserEmailHandler(c *gin.Context) {
    email := getCurrentUserEmail(c)
    if email == "" {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "User not logged in"})
        return
    }
    c.JSON(http.StatusOK, gin.H{"email": email})
}

func getUserInfoHandler(c *gin.Context) {
    name := getCurrentUserName(c)
    if name == "" {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "User not logged in"})
        return
    }
    c.JSON(http.StatusOK, gin.H{"name": name})
}