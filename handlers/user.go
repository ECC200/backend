package handlers

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"backend/firebase"
	"backend/models"

	"cloud.google.com/go/firestore"
	"github.com/gin-gonic/gin"
	"google.golang.org/api/iterator"
)

// 10桁の英数字を生成する関数
func generateRandomID() string {
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	rand.Seed(time.Now().UnixNano())
	id := make([]byte, 10)
	for i := range id {
		id[i] = charset[rand.Intn(len(charset))]
	}
	return string(id)
}

// Userを作成
func CreateUser(ctx *gin.Context) {
	// Firestoreクライアントを初期化
	Client, _ := firebase.Initialize(ctx)
	defer Client.Close()

	var user models.User
	// 10桁の英数字を生成してユーザーIDとして設定
	user.UserID = generateRandomID()

	// リクエストのJSONをUserモデルにバインド
	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// Firestoreにユーザーを追加する際にドキュメントIDを指定
	docRef := Client.Collection("users").Doc(user.UserID)
	_, err := docRef.Set(ctx, user)
	if err != nil {
		log.Println(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	// ドキュメントIDをUserIDとしてセット
	user.UserID = docRef.ID

	ctx.JSON(http.StatusOK, gin.H{"message": "User created successfully", "user_id": docRef.ID})
}

// User情報取得
func GetUser(ctx *gin.Context) {
	userID := ctx.Param("id")
	// Firestoreクライアントを初期化
	Client, _ := firebase.Initialize(ctx)
	defer Client.Close()

	// 指定されたuserIDのドキュメントを取得
	doc, err := Client.Collection("users").Doc(userID).Get(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user data"})
		return
	}

	var user models.User
	doc.DataTo(&user)
	user.UserID = doc.Ref.ID

	// プロフィール画像のURLを取得
	if user.Photo != "" {
		// gs:// URLをHTTP/HTTPS URLに変換
		gsPrefix := "gs://"
		if strings.HasPrefix(user.Photo, gsPrefix) {
			photoPath := strings.TrimPrefix(user.Photo, gsPrefix)
			parts := strings.SplitN(photoPath, "/", 2)
			if len(parts) == 2 {
				bucketName := parts[0]
				objectName := parts[1]

				// URLを取得し、署名付きURLとして返す（必要に応じてアクセス制御）
				url := fmt.Sprintf("https://storage.googleapis.com/%s/%s", bucketName, objectName)
				user.Photo = url
			} else {
				ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid photo path format"})
				return
			}
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Photo path must start with gs://"})
			return
		}
	}

	ctx.JSON(http.StatusOK, user)
}

// 障害者番号をチェックするハンドラー関数
func CheckDisabilityIdHandler(ctx *gin.Context) {
	// Firestoreクライアントを初期化
	Client, _ := firebase.Initialize(ctx)
	defer Client.Close()
	var user models.User
	var req struct {
		UserID string `json:"disabilityId"` // リクエストボディから障害者番号（UserID）を取得
	}

	// リクエストのJSONをバインド
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// 障害者番号（UserID）でユーザーを検索
	iter := Client.Collection("users").Where("user_id", "==", req.UserID).Documents(ctx)

	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user data"})
			return
		}
		doc.DataTo(&user) // ドキュメントデータをUserモデルにマッピング
		break
	}

	if user.UserID == "" {
		// ユーザーが見つからない場合
		ctx.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "Disability ID not found"})
		return
	}

	// ユーザーが見つかった場合
	ctx.JSON(http.StatusOK, gin.H{"success": true})
}

// 履歴を更新
func UpdateHistory(ctx *gin.Context) {
	userID := ctx.Param("id")
	var updatedHistories []models.History
	if err := ctx.BindJSON(&updatedHistories); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// Firestoreクライアントを初期化
	Client, _ := firebase.Initialize(ctx)
	defer Client.Close()

	docRef := Client.Collection("users").Doc(userID)
	_, err := docRef.Update(ctx, []firestore.Update{
		{Path: "historys", Value: updatedHistories},
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user history"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"success": true})
}

// User詳細情報更新
func UpdateUserDetails(ctx *gin.Context) {
	userID := ctx.Param("id")
	var updates struct {
		MedicationStatus string `json:"medication_status"`
		DoctorComment    string `json:"doctor_comment"`
	}

	if err := ctx.BindJSON(&updates); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// Firestoreクライアントを初期化
	Client, _ := firebase.Initialize(ctx)
	defer Client.Close()

	docRef := Client.Collection("users").Doc(userID)
	_, err := docRef.Update(ctx, []firestore.Update{
		{Path: "medication_status", Value: updates.MedicationStatus},
		{Path: "doctor_comment", Value: updates.DoctorComment},
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user details"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"success": true})
}
