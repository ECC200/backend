package handlers

import (
	"backend/firebase"
	"backend/models"
	"context"
	"net/http"

	"cloud.google.com/go/firestore"
	"github.com/gin-gonic/gin"
	"google.golang.org/api/iterator"
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
	user.BirthDate = "" // 簡単にするために現在時刻をセット
	user.UserID = ""    // Firestoreが自動生成するため空にしておく

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

// 障害者番号をチェックするハンドラー関数
func CheckDisabilityIdHandler(c *gin.Context) {
	var req struct {
		UserID string `json:"disabilityId"` // リクエストボディから障害者番号（UserID）を取得
	}

	// リクエストのJSONをバインド
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	ctx := context.Background()
	// Firestoreクライアントを初期化
	client, err := firebase.App.Firestore(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to initialize Firestore"})
		return
	}
	defer client.Close()

	// 障害者番号（UserID）でユーザーを検索
	iter := client.Collection("users").Where("user_id", "==", req.UserID).Documents(ctx)

	var user models.User
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user data"})
			return
		}
		doc.DataTo(&user) // ドキュメントデータをUserモデルにマッピング
		break
	}

	if user.UserID == "" {
		// ユーザーが見つからない場合
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "Disability ID not found"})
		return
	}

	// ユーザーが見つかった場合
	c.JSON(http.StatusOK, gin.H{"success": true})
}
