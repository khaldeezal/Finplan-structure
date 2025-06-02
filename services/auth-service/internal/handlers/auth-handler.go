package handlers

import (
	"context"
	authpb "github.com/khaldeezal/Finplan-proto/proto-definitions/gen/auth"
	service "github.com/khaldeezal/Finplan-structure/services/auth-service/internal/services"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type AuthHandler struct {
	authpb.UnimplementedAuthServiceServer
	service   service.AuthService
	logger    *zap.Logger
	jwtSecret string
}

// NewAuthHandler создаёт новый AuthHandler с внедрённой бизнес-логикой.
func NewAuthHandler(s service.AuthService, logger *zap.Logger) *AuthHandler {
	jwtSecret := os.Getenv("JWT_SECRET")
	return &AuthHandler{service: s, logger: logger, jwtSecret: jwtSecret}
}

func (h *AuthHandler) generateJWT(userID string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour * 72).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(h.jwtSecret))
}

// Register обрабатывает запрос регистрации нового пользователя.
// Принимает email, пароль и имя, возвращает userID или ошибку.
func (h *AuthHandler) Register(ctx context.Context, req *authpb.RegisterRequest) (*authpb.AuthResponse, error) {
	h.logger.Info("Register called", zap.String("email", req.Email))
	userID, err := h.service.Register(ctx, req.Email, req.Password, req.Name)
	if err != nil {
		h.logger.Error("Registration failed", zap.Error(err))
		if err == service.ErrUserExists {
			return nil, status.Error(codes.AlreadyExists, "user already exists")
		}
		return nil, status.Error(codes.Internal, "internal error")
	}
	jwtToken, err := h.generateJWT(userID)
	if err != nil {
		h.logger.Error("JWT generation failed", zap.Error(err))
		return nil, status.Error(codes.Internal, "failed to generate token")
	}
	return &authpb.AuthResponse{
		AccessToken:  jwtToken,
		RefreshToken: "",
		ExpiresAt:    "",
	}, nil
}

// Login обрабатывает вход пользователя.
// Принимает email и пароль, возвращает JWT-токен или ошибку.
func (h *AuthHandler) Login(ctx context.Context, req *authpb.LoginRequest) (*authpb.AuthResponse, error) {
	h.logger.Info("Login called", zap.String("email", req.Email))
	userID, err := h.service.Login(ctx, req.Email, req.Password)
	if err != nil {
		h.logger.Error("Login failed", zap.Error(err))
		if err == service.ErrInvalidCredentials || err == service.ErrUserNotFound {
			return nil, status.Error(codes.Unauthenticated, "invalid credentials")
		}
		return nil, status.Error(codes.Internal, "internal error")
	}
	jwtToken, err := h.generateJWT(userID)
	if err != nil {
		h.logger.Error("JWT generation failed", zap.Error(err))
		return nil, status.Error(codes.Internal, "failed to generate token")
	}
	return &authpb.AuthResponse{
		AccessToken:  jwtToken,
		RefreshToken: "",
		ExpiresAt:    "",
	}, nil
}

// VerifyToken проверяет валидность переданного JWT-токена.
// Возвращает userID, связанный с токеном, или ошибку.
func (h *AuthHandler) VerifyToken(ctx context.Context, req *authpb.VerifyTokenRequest) (*authpb.VerifyTokenResponse, error) {
	h.logger.Info("VerifyToken called", zap.String("token", req.Token))
	userID, err := h.service.VerifyToken(ctx, req.Token)
	if err != nil {
		h.logger.Error("Token verification failed",
			zap.String("token", req.Token),
			zap.Error(err),
		)
		return nil, err
	}
	h.logger.Info("Token verification successful",
		zap.String("token", req.Token),
		zap.String("userID", userID),
	)
	h.logger.Info("Returning VerifyToken response", zap.String("userID", userID))
	return &authpb.VerifyTokenResponse{
		UserId: userID,
		Valid:  true,
	}, nil
}
