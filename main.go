package main

import(
	"github.com/gin-gonic/gin"
	"server/login"
)

func main() {

	r := gin.Default()
 
	r.GET("/", login.googleForm)
	r.GET("/auth/google/login", login.googleLoginHandler)
	r.GET("/auth/google/callback", login.googleAuthCallback)
 
	r.Run(":5000")
 }