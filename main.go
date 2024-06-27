package main

import (
	"backend/firebase"
	"backend/handlers"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	firebase.Initialize()

	r := gin.Default()

	// CORSミドルウェアを使用
	r.Use(cors.Default())

	r.POST("/users", handlers.CreateUser)
	r.GET("/users/:id", handlers.GetUser)

	// 認証ルートを追加
	r.POST("/login", handlers.LoginHandler)
	r.POST("/checkDisabilityID", handlers.CheckDisabilityIdHandler)

	r.Run()
}
