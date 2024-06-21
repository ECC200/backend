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

	r.POST("/medical_memos", handlers.CreateMedicalMemo)
	r.GET("/medical_memos/:id", handlers.GetMedicalMemo)

	// Hospitals routes
	r.POST("/hospitals", handlers.CreateHospital)
	r.GET("/hospitals/:id", handlers.GetHospital)

	// 認証ルートを追加
	r.POST("/login", handlers.LoginHandler)

	// 障害者番号をチェックするルートを追加
	r.POST("/checkDisabilityId", handlers.CheckDisabilityIdHandler)

	r.Run()
}
