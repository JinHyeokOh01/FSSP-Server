package main

import(
	"github.com/gin-gonic/gin"
	"github.com/JinHyeokOh01/FSSP-Server/login"
)

func main() {

	r := gin.Default()
 
	r.GET("/", googleForm)
	r.GET("/auth/google/login", googleLoginHandler)
	r.GET("/auth/google/callback", googleAuthCallback)
 
	r.Run(":5000")
 }