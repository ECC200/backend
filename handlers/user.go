package handlers

import (
	"backend/firebase"
	"backend/models"
	"context"
	"net/http"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/gin-gonic/gin"
)

// Userを作成
func CreateUser(c *gin.Context) {
	var user models.User
	// リクエストのJSONをUserモデルにバインド
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// 修正必要
	user.BirthDate = time.Now() // 簡単にするために現在時刻をセット
	user.UserID = ""            // Firestoreが自動生成するため空にしておく

	ctx := context.Background()
	// Firestoreクライアントを初期化
	client, err := firebase.App.Firestore(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to initialize Firestore"})
		return
	}
	defer client.Close()

	// Firestoreにユーザーを追加
	docRef, _, err := client.Collection("users").Add(ctx, user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	user.UserID = docRef.ID // ドキュメントIDをUserIDとしてセット

	// FirestoreにUserIDを更新
	_, err = docRef.Set(ctx, map[string]interface{}{"user_id": user.UserID}, firestore.MergeAll)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user_id"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User created successfully", "user_id": user.UserID})
}

// User情報取得
func GetUser(c *gin.Context) {
	userID := c.Param("id")
	ctx := context.Background()
	// Firestoreクライアントを初期化
	client, err := firebase.App.Firestore(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to initialize Firestore"})
		return
	}
	defer client.Close()

	// 指定されたuserIDのドキュメントを取得
	doc, err := client.Collection("users").Doc(userID).Get(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user data"})
		return
	}

	var user models.User
	doc.DataTo(&user) // ドキュメントデータをUserモデルにマッピング
	user.UserID = doc.Ref.ID

	c.JSON(http.StatusOK, user)
}
