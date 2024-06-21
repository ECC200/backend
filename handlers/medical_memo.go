package handlers

import (
	"backend/firebase"
	"backend/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CreateMedicalMemo(ctx *gin.Context) {
	var memo models.Health_care_staff
	if err := ctx.ShouldBindJSON(&memo); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	client, err := firebase.App.Firestore(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to initialize Firestore"})
		return
	}
	defer client.Close()

	_, _, err = client.Collection("medical_memos").Add(ctx, memo)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create medical memo"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Medical memo created successfully"})
}

func GetMedicalMemo(ctx *gin.Context) {
	userID := ctx.Param("id")
	client, err := firebase.App.Firestore(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to initialize Firestore"})
		return
	}
	defer client.Close()

	doc, err := client.Collection("medical_memos").Doc(userID).Get(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get medical memo"})
		return
	}

	var memo models.Health_care_staff
	doc.DataTo(&memo)
	ctx.JSON(http.StatusOK, memo)
}
