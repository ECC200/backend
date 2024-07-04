package main

import (
	"backend/handlers"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// Firestore 初期化

	r := gin.Default()

	// CORSミドルウェアを使用
	r.Use(cors.Default())

	// User関連の処理
	r.POST("/users", handlers.CreateUser)
	r.GET("/users/:id", handlers.GetUser)
	r.PUT("/users/:id/details", handlers.UpdateUserDetails)
	r.PUT("/users/:id/history", handlers.UpdateHistory)

	// Staff関連の処理
	r.GET("/staffs/:id", handlers.Getstaff)

	// Level2のログイン処理
	r.POST("/login", handlers.LoginHandler)

	// 障がい者番号をチェック
	r.POST("/checkDisabilityID", handlers.CheckDisabilityIdHandler)

	r.Run()
}
