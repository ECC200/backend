package handlers_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"backend/handlers"

	"github.com/gin-gonic/gin"
)

func TestCreateHospital(t *testing.T) {
	router := gin.Default()
	router.POST("/hospitals", handlers.CreateHospital)

	// テストペイロード
	payload := `{"hospital_name":"Test Hospital","password":"test123"}`

	req, _ := http.NewRequest("POST", "/hospitals", bytes.NewBuffer([]byte(payload)))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}

	expected := `{"message":"Hospital created successfully","hospital_id":`
	if !bytes.HasPrefix(w.Body.Bytes(), []byte(expected)) {
		t.Errorf("Unexpected response body: got %v", w.Body.String())
	}
}
