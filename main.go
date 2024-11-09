package main

import(
	"github.com/gin-gonic/gin"
	"server/login"
	"server/db"
	"server/naver"
	"server/geolocation"
	"fmt"
)

func main() {
	client, err := db.ConnectDB()
	if err != nil{
		panic(err)
	}
	defer db.DisconnectDB(client)

	r := gin.Default()
 
	r.GET("/", login.GoogleForm)
	r.GET("/auth/google/login", login.GoogleLoginHandler)
	r.GET("/auth/google/callback", func(c *gin.Context) {
        login.GoogleAuthCallback(c, client)
		db.Mu.Lock() // 뮤텍스 잠금
        if db.CurrentUser != nil {
            fmt.Println("Current User Info:", db.CurrentUser) // 전역 변수 사용
        }
        db.Mu.Unlock() // 뮤텍스 잠금 해제
    })

	r.GET("/search", naver.QuerySearch)
	r.POST("/savelist", func(c *gin.Context){
		db.UpdateListHandler(c, client, db.CurrentUser["email"].(string))
	})
	r.POST("/deletelist", func(c *gin.Context){
		db.DeleteListHandler(c, client, db.CurrentUser["email"].(string))
	})
	r.GET("/getLocation", geolocation.GeoLocationHandler)

	r.Run(":5000")
 }