package services

import (
	"context"
	"errors"
	"strings"
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/khaldeezal/Finplan-structure/services/auth-service/internal/model"
	"github.com/khaldeezal/Finplan-structure/services/auth-service/internal/repo"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUserExists         = errors.New("user already exists")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrInvalidEmail       = errors.New("invalid email format")
	ErrUserNotFound       = errors.New("user not found")
	ErrInvalidToken       = errors.New("invalid token")
)

// Бизнес-логика аутентификации
type AuthService interface {
	Register(ctx context.Context, email, password, name string) (string, error)
	Login(ctx context.Context, email, password string) (string, error)
	VerifyToken(ctx context.Context, token string) (string, error)
}

type authService struct {
	repo      repo.UserRepository
	jwtSecret string
	logger    *zap.Logger
}

// Новый экземпляр сервиса аутентификации
func NewAuthService(r repo.UserRepository, secret string, logger *zap.Logger) AuthService {
	if secret == "" {
		logger.Fatal("JWT_SECRET is required")
		return nil
	}
	return &authService{
		repo:      r,
		jwtSecret: secret,
		logger:    logger,
	}
}

// Хеширует пароль, создаёт нового пользователя и сохраняет его через репо
// Возвращает ID нового пользователя или ошибку.
func (s *authService) Register(ctx context.Context, email, password, name string) (string, error) {
	s.logger.Info("Register called", zap.String("email", email))
	if !strings.Contains(email, "@") {
		return "", ErrInvalidEmail
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		s.logger.Error("failed to hash password", zap.Error(err))
		return "", err
	}
	user := &model.User{
		Email:    email,
		Password: string(hashedPassword),
		Name:     name,
	}
	id, err := s.repo.CreateUser(ctx, user)
	if err != nil {
		return "", err
	}
	s.logger.Info("user created", zap.String("email", email))
	return id, nil
}

// Создаёт JWT-токен с userID, живёт 72 часа
func (s *authService) generateJWT(userID string) (string, error) {
	s.logger.Info("Generating JWT", zap.String("userID", userID))
	s.logger.Info("JWT Secret at generateJWT", zap.String("secret", s.jwtSecret))
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour * 72).Unix(), // 72 часа действия токена
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(s.jwtSecret))
	if err != nil {
		s.logger.Error("failed to sign JWT", zap.Error(err))
		return "", err
	}
	return signedToken, nil
}

// Проверяет email и пароль, возвращает userID или ошибку
func (s *authService) Login(ctx context.Context, email, password string) (string, error) {
	s.logger.Info("Login called", zap.String("email", email))
	user, err := s.repo.GetUserByEmail(ctx, email)
	if err != nil {
		return "", err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		s.logger.Error("invalid credentials", zap.Error(err))
		return "", ErrInvalidCredentials
	}
	s.logger.Info("user logged in", zap.String("userID", user.ID))
	return user.ID, nil
}

// Проверяет валидность JWT-токена и возвращает userID из токена или ошибку
func (s *authService) VerifyToken(ctx context.Context, tokenStr string) (string, error) {
	s.logger.Info("Verifying token")
	s.logger.Info("VerifyToken called", zap.String("secret", s.jwtSecret), zap.String("token", tokenStr))
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(s.jwtSecret), nil
	})
	if err != nil || !token.Valid {
		s.logger.Error("invalid token", zap.Error(err))
		return "", ErrInvalidToken
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", ErrInvalidToken
	}
	userID, ok := claims["user_id"].(string)
	if !ok {
		return "", ErrUserNotFound
	}
	return userID, nil
}
