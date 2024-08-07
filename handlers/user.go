package handlers

import (
	"backend/firebase"
	"backend/models"
	"bytes"
	"context"
	"fmt"
	"math/rand"
	"net/http"
	"strings"
	"time"

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

// 画像をFirebase Storageにアップロード
func uploadFileToFirebaseStorage(ctx context.Context, fileName string, data []byte) (string, error) {
	bucketName := "care-connect-eba8d.appspot.com"
	bucket := firebase.StorageClient.Bucket(bucketName)
	object := bucket.Object(fileName)

	writer := object.NewWriter(ctx)
	defer writer.Close()

	if _, err := writer.Write(data); err != nil {
		return "", err
	}

	return fmt.Sprintf("gs://%s/%s", bucketName, fileName), nil
}

// Userを作成
func CreateUser(c *gin.Context) {
	var user models.User

	// リクエストのMultipart formをパース
	if err := c.Request.ParseMultipartForm(10 << 20); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse form"})
		return
	}

	// Userデータの取得
	user.UserID = generateRandomID()
	user.UserName = c.PostForm("user_name")
	user.BirthDate = c.PostForm("birth_date")
	user.Age = c.PostForm("age")
	user.Address = c.PostForm("address")
	user.Contact = c.PostForm("contact")
	user.HospitalDestination = c.PostForm("hospital_destination")
	user.PrimaryCareDoctor = c.PostForm("primary_care_doctor")
	user.Specialty = c.PostForm("specialty")
	user.ChronicDisease = c.PostForm("chronic_disease")
	user.DisabilityGrade = c.PostForm("disability_grade")

	// 緊急連絡先の取得
	emergencyContacts := []models.EmergencyContact{
		{
			Name:  c.PostForm("emergency_contacts[0].name"),
			Phone: c.PostForm("emergency_contacts[0].phone"),
		},
		{
			Name:  c.PostForm("emergency_contacts[1].name"),
			Phone: c.PostForm("emergency_contacts[1].phone"),
		},
	}
	user.EmergencyContacts = emergencyContacts

	// 画像のアップロード
	fileHeader, err := c.FormFile("photo")
	if err == nil {
		file, err := fileHeader.Open()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to open photo file"})
			return
		}
		defer file.Close()

		buf := bytes.NewBuffer(nil)
		if _, err := buf.ReadFrom(file); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read photo file"})
			return
		}

		photoPath, err := uploadFileToFirebaseStorage(c.Request.Context(), fileHeader.Filename, buf.Bytes())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload photo to Firebase Storage"})
			return
		}
		user.Photo = photoPath
	}

	ctx := context.Background()
	client, err := firebase.App.Firestore(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to initialize Firestore"})
		return
	}
	defer client.Close()

	docRef := client.Collection("users").Doc(user.UserID)
	_, err = docRef.Set(ctx, user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User created successfully", "user_id": user.UserID})
}
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
	doc.DataTo(&user)
	user.UserID = doc.Ref.ID

	// プロフィール画像のURLを取得
	if user.Photo != "" {
		gsPrefix := "gs://"
		if strings.HasPrefix(user.Photo, gsPrefix) {
			photoPath := strings.TrimPrefix(user.Photo, gsPrefix)
			parts := strings.SplitN(photoPath, "/", 2)
			if len(parts) == 2 {
				bucketName := parts[0]
				objectName := parts[1]

				// 署名付きURLを生成
				signedURL, err := firebase.GenerateSignedURL(bucketName, objectName)
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate signed URL"})
					return
				}
				user.Photo = signedURL
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid photo path format"})
				return
			}
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Photo path must start with gs://"})
			return
		}
	}

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

// 履歴を更新
func UpdateHistory(c *gin.Context) {
	userID := c.Param("id")
	var updatedHistories []models.History
	if err := c.BindJSON(&updatedHistories); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	ctx := context.Background()
	client, err := firebase.App.Firestore(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to initialize Firestore"})
		return
	}
	defer client.Close()

	docRef := client.Collection("users").Doc(userID)
	_, err = docRef.Update(ctx, []firestore.Update{
		{Path: "historys", Value: updatedHistories},
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user history"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true})
}

// User詳細情報更新
func UpdateUserDetails(c *gin.Context) {
	userID := c.Param("id")
	var updates struct {
		MedicationStatus string `json:"medication_status"`
		DoctorComment    string `json:"doctor_comment"`
	}
	if err := c.BindJSON(&updates); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	ctx := context.Background()
	client, err := firebase.App.Firestore(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to initialize Firestore"})
		return
	}
	defer client.Close()

	docRef := client.Collection("users").Doc(userID)
	_, err = docRef.Update(ctx, []firestore.Update{
		{Path: "medication_status", Value: updates.MedicationStatus},
		{Path: "doctor_comment", Value: updates.DoctorComment},
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user details"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true})
}
