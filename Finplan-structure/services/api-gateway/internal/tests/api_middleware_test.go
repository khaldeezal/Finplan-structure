package tests

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/khaldeezal/Finplan-structure/services/api-gateway/internal/middleware"
	"github.com/stretchr/testify/assert"
)

func generateTestJWT(secret string, userID string, expOffset time.Duration) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(expOffset).Unix(),
	})
	tokenStr, _ := token.SignedString([]byte(secret))
	return tokenStr
}

func setupTestRouter(secret string) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(middleware.JWTMiddleware(secret))
	r.GET("/protected", func(c *gin.Context) {
		userID, ok := c.Get("user_id")
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "no claims"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"user_id": userID})
	})
	return r
}

func TestJWTMiddleware_NoToken(t *testing.T) {
	secret := "khaldee0711"
	router := setupTestRouter(secret)

	req, _ := http.NewRequest("GET", "/protected", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusUnauthorized, resp.Code)
	assert.Contains(t, resp.Body.String(), "missing auth header")
}

func TestJWTMiddleware_InvalidToken(t *testing.T) {
	secret := "khaldee0711"
	router := setupTestRouter(secret)

	req, _ := http.NewRequest("GET", "/protected", nil)
	req.Header.Set("Authorization", "Bearer invalid.token.string")
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusUnauthorized, resp.Code)
	assert.Contains(t, resp.Body.String(), "invalid or expired token")
}

func TestJWTMiddleware_ExpiredToken(t *testing.T) {
	secret := "khaldee0711"
	router := setupTestRouter(secret)

	expiredToken := generateTestJWT(secret, "user42", -time.Hour)
	req, _ := http.NewRequest("GET", "/protected", nil)
	req.Header.Set("Authorization", "Bearer "+expiredToken)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusUnauthorized, resp.Code)
	assert.Contains(t, resp.Body.String(), "invalid or expired token")
}

func TestJWTMiddleware_ValidToken(t *testing.T) {
	secret := "khaldee0711"
	router := setupTestRouter(secret)

	validToken := generateTestJWT(secret, "user42", time.Hour)
	req, _ := http.NewRequest("GET", "/protected", nil)
	req.Header.Set("Authorization", "Bearer "+validToken)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
	assert.Contains(t, resp.Body.String(), "user42")
}
