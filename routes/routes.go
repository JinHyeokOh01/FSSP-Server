package routes

import(
    "github.com/gin-gonic/gin"
    "github.com/JinHyeokOh01/FSSP-Server/controllers"
    "github.com/JinHyeokOh01/FSSP-Server/db"
    "go.mongodb.org/mongo-driver/mongo"
    "fmt"
)

func Routes(client *mongo.Client) *gin.Engine{
    r := gin.Default()

    // 기존 인증 관련 라우트
    r.GET("/", controllers.GoogleForm)
    r.GET("/auth/google/login", controllers.GoogleLoginHandler)
    r.GET("/auth/google/callback", func(c *gin.Context) {
        controllers.GoogleAuthCallback(c, client)
        db.Mu.Lock() // 뮤텍스 잠금
        if db.CurrentUser != nil {
            fmt.Println("Current User Info:", db.CurrentUser)
        }
        db.Mu.Unlock() // 뮤텍스 잠금 해제
    })

    // 기존 API 라우트
    apiRoutes := r.Group("/api")
    {
        apiRoutes.GET("/search", controllers.NaverSearchHandler)
        apiRoutes.POST("/chat", controllers.HandleChat)
    }
    
    // 레스토랑 관리 라우트
    restaurantRoutes := r.Group("/api/restaurants")
    {
        // 현재 로그인한 사용자의 레스토랑 목록 조회
        restaurantRoutes.GET("", func(c *gin.Context){
            db.Mu.Lock()
            if db.CurrentUser == nil {
                c.JSON(401, gin.H{"error": "User not logged in"})
                db.Mu.Unlock()
                return
            }
            email := db.CurrentUser["email"].(string)
            db.Mu.Unlock()
            db.GetRestaurantsHandler(c, client, email)
        })

        // 새로운 레스토랑 추가
        restaurantRoutes.POST("", func(c *gin.Context){
            db.Mu.Lock()
            if db.CurrentUser == nil {
                c.JSON(401, gin.H{"error": "User not logged in"})
                db.Mu.Unlock()
                return
            }
            email := db.CurrentUser["email"].(string)
            db.Mu.Unlock()
            db.UpdateListHandler(c, client, email)
        })

        // 레스토랑 삭제
        restaurantRoutes.DELETE("", func(c *gin.Context){
            db.Mu.Lock()
            if db.CurrentUser == nil {
                c.JSON(401, gin.H{"error": "User not logged in"})
                db.Mu.Unlock()
                return
            }
            email := db.CurrentUser["email"].(string)
            db.Mu.Unlock()
            db.DeleteListHandler(c, client, email)
        })
    }

    return r
}