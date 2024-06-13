package handlers

import (
	"backend/firebase"
	"backend/models"
	"context"
	"net/http"

	"cloud.google.com/go/firestore"
	"github.com/gin-gonic/gin"
)

// staffを作成
func Createstaff(c *gin.Context) {
	var staff models.Staff
	// リクエストのJSONをstaffモデルにバインド
	if err := c.ShouldBindJSON(&staff); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// 修正必要
	staff.StaffID = "" // Firestoreが自動生成するため空にしておく
	staff.StaffName = ""
	staff.Password = ""

	ctx := context.Background()
	// Firestoreクライアントを初期化
	client, err := firebase.App.Firestore(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to initialize Firestore"})
		return
	}
	defer client.Close()

	// Firestoreにユーザーを追加
	docRef, _, err := client.Collection("staffs").Add(ctx, staff)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create staff"})
		return
	}

	staff.StaffID = docRef.ID // ドキュメントIDをstaffIDとしてセット

	// FirestoreにstaffIDを更新
	_, err = docRef.Set(ctx, map[string]interface{}{"staff_id": staff.StaffID}, firestore.MergeAll)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update staff_id"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "staff created successfully", "staff_id": staff.StaffID})
}

// staff情報取得
func Getstaff(c *gin.Context) {
	staffID := c.Param("id")
	ctx := context.Background()
	// Firestoreクライアントを初期化
	client, err := firebase.App.Firestore(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to initialize Firestore"})
		return
	}
	defer client.Close()

	// 指定されたstaffIDのドキュメントを取得
	doc, err := client.Collection("staffs").Doc(staffID).Get(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get staff data"})
		return
	}

	var staff models.Staff
	doc.DataTo(&staff) // ドキュメントデータをstaffモデルにマッピング
	staff.StaffID = doc.Ref.ID

	c.JSON(http.StatusOK, staff)
}
