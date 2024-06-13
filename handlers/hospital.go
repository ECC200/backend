package handlers

import (
	"backend/firebase"
	"backend/models"
	"context"
	"net/http"

	"cloud.google.com/go/firestore"
	"github.com/gin-gonic/gin"
)

// hospitalを作成
func CreateHospital(c *gin.Context) {
	var hospital models.Hosoital
	// リクエストのJSONをhospitalモデルにバインド
	if err := c.ShouldBindJSON(&hospital); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// 修正必要
	hospital.HospitalID = "" // Firestoreが自動生成するため空にしておく
	hospital.HospitalName = ""
	hospital.Password = ""

	ctx := context.Background()
	// Firestoreクライアントを初期化
	client, err := firebase.App.Firestore(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to initialize Firestore"})
		return
	}
	defer client.Close()

	// Firestoreにユーザーを追加
	docRef, _, err := client.Collection("hospitals").Add(ctx, hospital)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create staff"})
		return
	}

	hospital.HospitalID = docRef.ID // ドキュメントIDをHospitalIDとしてセット

	// FirestoreにhospitalIDを更新
	_, err = docRef.Set(ctx, map[string]interface{}{"hospital_id": hospital.HospitalID}, firestore.MergeAll)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update hospital_id"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "hospital created successfully", "hospital_id": hospital.HospitalID})
}

// hospital情報取得
func Gethospital(c *gin.Context) {
	hospitalID := c.Param("id")
	ctx := context.Background()
	// Firestoreクライアントを初期化
	client, err := firebase.App.Firestore(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to initialize Firestore"})
		return
	}
	defer client.Close()

	// 指定されたhospitalIDのドキュメントを取得
	doc, err := client.Collection("hospitals").Doc(hospitalID).Get(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get staff data"})
		return
	}

	var hospital models.Hosoital
	doc.DataTo(&hospital) // ドキュメントデータをhospitalモデルにマッピング
	hospital.HospitalID = doc.Ref.ID

	c.JSON(http.StatusOK, hospital)
}
