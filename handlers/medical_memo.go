package handlers

import (
	"backend/firebase"
	"backend/models"
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CreateMedicalMemo(c *gin.Context) {
	var memo models.Health_care_staff
	if err := c.ShouldBindJSON(&memo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	ctx := context.Background()
	client, err := firebase.App.Firestore(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to initialize Firestore"})
		return
	}
	defer client.Close()

	_, _, err = client.Collection("medical_memos").Add(ctx, memo)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create medical memo"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Medical memo created successfully"})
}

func GetMedicalMemo(c *gin.Context) {
	userID := c.Param("id")
	ctx := context.Background()
	client, err := firebase.App.Firestore(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to initialize Firestore"})
		return
	}
	defer client.Close()

	doc, err := client.Collection("medical_memos").Doc(userID).Get(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get medical memo"})
		return
	}

	var memo models.Health_care_staff
	doc.DataTo(&memo)
	c.JSON(http.StatusOK, memo)
}
