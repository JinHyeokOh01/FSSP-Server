package main

import(
	"github.com/gin-gonic/gin"
	"github.com/JinHyeokOh01/FSSP-Server"
)

func main() {

	r := gin.Default()
 
	r.GET("/", login.googleForm)
	r.GET("/auth/google/login", login.googleLoginHandler)
	r.GET("/auth/google/callback", login.googleAuthCallback)
 
	r.Run(":5000")
 }