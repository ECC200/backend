package handlers_test

import (
	"backend/firebase"
	"backend/handlers"
	"backend/models"
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// 新しくuserを作成する処理
func TestCreateUser(t *testing.T) {
	firebase.Initialize()

	router := gin.Default()
	router.POST("/users", handlers.CreateUser)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/users", strings.NewReader(`{"user_name":"kotarou minami","mailaddress":"john@example.com","password":"password123","emergency_contact":"123456789","work_contact":"987654321","blood_type":"O+"}`))
	req.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "User created successfully")
}

// userを参照する処理
func TestGetUser(t *testing.T) {
	firebase.Initialize()

	// 事前にユーザーを作成しておく
	ctx := context.Background()
	client, _ := firebase.App.Firestore(ctx)
	user := models.User{
		UserName:         "John Doe",
		MailAddress:      "john@example.com",
		Password:         "password123",
		BirthDate:        time.Now(),
		EmergencyContact: "123456789",
		WorkContact:      "987654321",
		BloodType:        "O+",
	}
	docRef, _, _ := client.Collection("users").Add(ctx, user)
	user.UserID = docRef.ID

	// FirestoreにUserIDを更新
	_, _ = docRef.Set(ctx, map[string]interface{}{"user_id": user.UserID}, firestore.MergeAll)

	router := gin.Default()
	router.GET("/users/:id", handlers.GetUser)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/users/"+user.UserID, nil)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "John Doe")
}
