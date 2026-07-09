package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	
	"github.com/aaaaarsen/ai-dos/internal/auth"

	"github.com/gin-gonic/gin"
)

func setupRouter(secret string) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(AuthMiddleware(secret))
	r.GET("/test", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})
	return r
}

func TestAuthMiddleware_NoHeader(t *testing.T) {
	r := setupRouter("secret")
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	r.ServeHTTP(w, req)
	if w.Code != 401 {
		t.Fatalf("expected 401, got %d", w.Code)
	}
}

func TestAuthMiddleware_InvalidForm(t *testing.T){
	r := setupRouter("secret")
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "InvalidToken")
	r.ServeHTTP(w, req)
	if w.Code != 401 {
		t.Fatalf("expected 401, got %d", w.Code)
	}
}

func TestAuthMiddleware_ValidToken(t *testing.T) {
    secret := "secret"
    token, _ := auth.GenerateToken(1, secret)
    r := setupRouter(secret)
    w := httptest.NewRecorder()
    req, _ := http.NewRequest("GET", "/test", nil)
    req.Header.Set("Authorization", "Bearer "+token)
    r.ServeHTTP(w, req)
    if w.Code != 200 {
        t.Fatalf("expected 200, got %d", w.Code)
    }
}