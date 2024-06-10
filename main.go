package main

import (
	"backend/firebase"
	"backend/handlers"

	"github.com/gin-gonic/gin"
)

func main() {
	firebase.Initialize()

	r := gin.Default()

	r.POST("/users", handlers.CreateUser)
	r.GET("/users/:id", handlers.GetUser)

	// r.POST("/medical_memos", handlers.CreateMedicalMemo)
	// r.GET("/medical_memos/:id", handlers.GetMedicalMemo)

	r.Run()
}
