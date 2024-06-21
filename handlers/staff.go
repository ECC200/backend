package handlers

import (
	"backend/firebase"
	"backend/models"
	"net/http"

	"cloud.google.com/go/firestore"
	"github.com/gin-gonic/gin"
)

// staffを作成
func Createstaff(ctx *gin.Context) {
	var staff models.Staff
	// リクエストのJSONをstaffモデルにバインド
	if err := ctx.ShouldBindJSON(&staff); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// 修正必要
	staff.StaffID = "" // Firestoreが自動生成するため空にしておく
	staff.StaffName = ""
	staff.Password = ""

	// Firestoreクライアントを初期化
	client, err := firebase.App.Firestore(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to initialize Firestore"})
		return
	}
	defer client.Close()

	// Firestoreにユーザーを追加
	docRef, _, err := client.Collection("staffs").Add(ctx, staff)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create staff"})
		return
	}

	staff.StaffID = docRef.ID // ドキュメントIDをstaffIDとしてセット

	// FirestoreにstaffIDを更新
	_, err = docRef.Set(ctx, map[string]interface{}{"staff_id": staff.StaffID}, firestore.MergeAll)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update staff_id"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "staff created successfully", "staff_id": staff.StaffID})
}

///----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------

// staff情報取得
func Getstaff(ctx *gin.Context) {
	staffID := ctx.Param("id")
	// Firestoreクライアントを初期化
	client, err := firebase.App.Firestore(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to initialize Firestore"})
		return
	}
	defer client.Close()

	// 指定されたstaffIDのドキュメントを取得
	doc, err := client.Collection("staffs").Doc(staffID).Get(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get staff data"})
		return
	}

	var staff models.Staff
	doc.DataTo(&staff) // ドキュメントデータをstaffモデルにマッピング
	staff.StaffID = doc.Ref.ID

	ctx.JSON(http.StatusOK, staff)
}
