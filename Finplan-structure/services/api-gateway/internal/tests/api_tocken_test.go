package tests

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

func TestGenerateAndValidateJWT(t *testing.T) {
	secret := "khaldee0711"
	userID := "test_user"

	// Генерируем токен с помощью jwt
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour).Unix(),
	})
	tokenStr, err := token.SignedString([]byte(secret))
	assert.NoError(t, err, "should generate token without error")
	assert.NotEmpty(t, tokenStr, "token string should not be empty")

	// Валидируем токен
	parsed, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	assert.NoError(t, err, "token should parse without error")
	assert.True(t, parsed.Valid, "token should be valid")

	claims, ok := parsed.Claims.(jwt.MapClaims)
	assert.True(t, ok, "claims should be of type MapClaims")
	assert.Equal(t, userID, claims["user_id"], "user_id should match")
}
