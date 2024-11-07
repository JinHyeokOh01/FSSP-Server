package main

import(
	"github.com/gin-gonic/gin"
	"github.com/JinHyeokOh01/FSSP-Server/lgn"
)

func main() {

	r := gin.Default()
 
	r.GET("/", lgn.googleForm)
	r.GET("/auth/google/login", lgn.googleLoginHandler)
	r.GET("/auth/google/callback", lgn.googleAuthCallback)
 
	r.Run(":5000")
 }