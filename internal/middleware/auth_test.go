package middleware

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"

	jwtpkg "mini-issue/pkg/jwt"
)

func TestAuthMiddlewareRejectsMissingAuthorizationHeader(t *testing.T) {
	router := newAuthTestRouter(t)

	recorder := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	router.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusUnauthorized)
	}
}

func TestAuthMiddlewareRejectsMalformedAuthorizationHeader(t *testing.T) {
	router := newAuthTestRouter(t)

	recorder := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Token abc")
	router.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusUnauthorized)
	}
}

func TestAuthMiddlewareRejectsInvalidToken(t *testing.T) {
	router := newAuthTestRouter(t)

	recorder := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Bearer invalid-token")
	router.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusUnauthorized)
	}
}

func TestAuthMiddlewareSetsUserContext(t *testing.T) {
	router := newAuthTestRouter(t)
	tokenString, err := jwtpkg.GenerateToken(7, "bob")
	if err != nil {
		t.Fatalf("GenerateToken() error = %v", err)
	}

	recorder := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Bearer "+tokenString)
	router.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d; body = %s", recorder.Code, http.StatusOK, recorder.Body.String())
	}

	var body struct {
		UserID   float64 `json:"user_id"`
		Username string  `json:"username"`
	}
	if err := json.Unmarshal(recorder.Body.Bytes(), &body); err != nil {
		t.Fatalf("json.Unmarshal() error = %v", err)
	}
	if body.UserID != 7 {
		t.Fatalf("user_id = %v, want 7", body.UserID)
	}
	if body.Username != "bob" {
		t.Fatalf("username = %q, want bob", body.Username)
	}
}

func newAuthTestRouter(t *testing.T) *gin.Engine {
	t.Helper()

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.GET("/protected", AuthMiddleware(), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"user_id":  c.GetInt64(ContextUserIDKey),
			"username": c.GetString(ContextUsernameKey),
		})
	})
	return router
}
