package main

import(
	"github.com/gin-gonic/gin"
	"server/login"
)

func main() {

	r := gin.Default()
 
	r.GET("/", login.GoogleForm)
	r.GET("/auth/google/login", login.GoogleLoginHandler)
	r.GET("/auth/google/callback", login.GoogleAuthCallback)
 
	r.Run(":5000")
 }