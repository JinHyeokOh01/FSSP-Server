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
 
	r.GET("/", controllers.GoogleForm)
	r.GET("/auth/google/login", controllers.GoogleLoginHandler)
	r.GET("/auth/google/callback", func(c *gin.Context) {
		controllers.GoogleAuthCallback(c, client)
		db.Mu.Lock() // 뮤텍스 잠금
		if db.CurrentUser != nil {
			fmt.Println("Current User Info:", db.CurrentUser) // 전역 변수 사용
		}
		db.Mu.Unlock() // 뮤텍스 잠금 해제
	})

	apiRoutes := r.Group("/api")
	{
		apiRoutes.GET("/search", controllers.NaverSearchHandler)
		apiRoutes.POST("/chat", controllers.HandleChat)
	}
	
	dbRoutes := r.Group("/api/db")
	{
		dbRoutes.POST("/savelist", func(c *gin.Context){
			db.UpdateListHandler(c, client, db.CurrentUser["email"].(string))
		})
		dbRoutes.POST("/deletelist", func(c *gin.Context){
			db.DeleteListHandler(c, client, db.CurrentUser["email"].(string))
		})
	}

	return r
}