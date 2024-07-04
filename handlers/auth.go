package handlers

import (
	"backend/firebase"
	"backend/models"
	"net/http"

	"log"

	"github.com/gin-gonic/gin"
	"google.golang.org/api/iterator"
)

type LoginRequest struct {
	StaffID  string `json:"staffId"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
}

func LoginHandler(ctx *gin.Context) {
	var req LoginRequest
	var staff models.Staff

	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("Invalid request: %v\n", err)
		ctx.JSON(http.StatusBadRequest, LoginResponse{Success: false, Message: "Invalid request"})
		return
	}
	log.Printf("Login attempt: staffId=%s\n", req.StaffID)

	// Firestoreクライアントを初期化
	Client, _ := firebase.Initialize(ctx)
	defer Client.Close()

	// staffIdに基づいてユーザーを取得
	iter := Client.Collection("staffs").Where("staff_id", "==", req.StaffID).Documents(ctx)
	defer iter.Stop()

	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Printf("Error fetching user: %v\n", err)
			ctx.JSON(http.StatusInternalServerError, LoginResponse{Success: false, Message: "Error fetching user"})
			return
		}
		if err := doc.DataTo(&staff); err != nil {
			log.Printf("Error decoding user: %v\n", err)
			ctx.JSON(http.StatusInternalServerError, LoginResponse{Success: false, Message: "Error decoding user"})
			return
		}
		log.Printf("User found: %+v\n", staff)
	}

	if staff.StaffID == "" || staff.Password != req.Password {
		log.Printf("Invalid staff ID or password\n")
		ctx.JSON(http.StatusUnauthorized, LoginResponse{Success: false, Message: "Invalid staff ID or password"})
		return
	}

	log.Printf("Login successful: staffId=%s\n", req.StaffID)
	ctx.JSON(http.StatusOK, LoginResponse{Success: true})
}
