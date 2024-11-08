package main

import(
	"github.com/gin-gonic/gin"
	"server/login"
	"server/db"
	"server/naver"
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
    })
	r.GET("/search", naver.QuerySearch)
 
	r.Run(":8080")
 }