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

func TestCreateUser(t *testing.T) {
	firebase.Initialize()
	if firebase.Client == nil {
		t.Fatalf("Failed to initialize Firestore")
	}

	router := gin.Default()
	router.POST("/users", handlers.CreateUser)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/users", strings.NewReader(`{"user_name":"kotarou minami","mailaddress":"john@example.com","password":"password123","work_contact":"987654321","blood_type":"O+","emergency_contacts":[{"name":"Mother","phone":"123456789"}]}`))
	req.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "User created successfully")
}

func TestGetUser(t *testing.T) {
	firebase.Initialize()
	if firebase.Client == nil {
		t.Fatalf("Failed to initialize Firestore")
	}

	ctx := context.Background()
	client := firebase.Client
	user := models.User{
		UserName:    "John Doe",
		MailAddress: "john@example.com",
		Password:    "password123",
		BirthDate:   time.Now(),
		EmergencyContacts: []models.EmergencyContact{
			{Name: "Mother", Phone: "090-1234-5678"},
			{Name: "Father", Phone: "090-8765-4321"},
		},
		WorkContact: "987654321",
		BloodType:   "O+",
	}
	docRef, _, err := client.Collection("users").Add(ctx, user)
	if err != nil {
		t.Fatalf("Failed to add user: %v", err)
	}
	user.UserID = docRef.ID

	_, err = docRef.Set(ctx, map[string]interface{}{"user_id": user.UserID}, firestore.MergeAll)
	if err != nil {
		t.Fatalf("Failed to update user ID: %v", err)
	}

	router := gin.Default()
	router.GET("/users/:id", handlers.GetUser)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/users/"+user.UserID, nil)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "John Doe")
}
