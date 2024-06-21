package handlers

import (
	"backend/firebase"
	"backend/models"
	"net/http"

	"cloud.google.com/go/firestore"
	"github.com/gin-gonic/gin"
)

// CreateHospital - hospitalを作成
func CreateHospital(ctx *gin.Context) {
	var hospital models.Hospital
	// リクエストのJSONをhospitalモデルにバインド
	if err := ctx.ShouldBindJSON(&hospital); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// 修正必要: Firestoreが自動生成するため空にしておく
	hospital.HospitalID = ""

	// Firestoreクライアントを初期化
	client, err := firebase.App.Firestore(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to initialize Firestore"})
		return
	}
	defer client.Close()

	// Firestoreにユーザーを追加
	docRef, _, err := client.Collection("hospitals").Add(ctx, hospital)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create hospital"})
		return
	}

	hospital.HospitalID = docRef.ID // ドキュメントIDをHospitalIDとしてセット

	// FirestoreにhospitalIDを更新
	_, err = docRef.Set(ctx, map[string]interface{}{"hospital_id": hospital.HospitalID}, firestore.MergeAll)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update hospital_id"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Hospital created successfully", "hospital_id": hospital.HospitalID})
}

// GetHospital - hospital情報取得
func GetHospital(ctx *gin.Context) {
	hospitalID := ctx.Param("id")
	// Firestoreクライアントを初期化
	client, err := firebase.App.Firestore(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to initialize Firestore"})
		return
	}
	defer client.Close()

	// 指定されたhospitalIDのドキュメントを取得
	doc, err := client.Collection("hospitals").Doc(hospitalID).Get(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get hospital data"})
		return
	}

	var hospital models.Hospital
	doc.DataTo(&hospital) // ドキュメントデータをhospitalモデルにマッピング
	hospital.HospitalID = doc.Ref.ID

	ctx.JSON(http.StatusOK, hospital)
}