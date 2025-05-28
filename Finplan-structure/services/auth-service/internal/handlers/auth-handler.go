package handlers

import (
	"context"
	authpb "github.com/khaldeezal/Finplan-structure/proto-definitions/gen/auth"
	service "github.com/khaldeezal/Finplan-structure/services/auth-service/internal/services"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AuthHandler struct {
	authpb.UnimplementedAuthServiceServer
	service service.AuthService
	logger  *zap.Logger
}

// NewAuthHandler создаёт новый AuthHandler с внедрённой бизнес-логикой.
func NewAuthHandler(s service.AuthService, logger *zap.Logger) *AuthHandler {
	return &AuthHandler{service: s, logger: logger}
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
	return &authpb.AuthResponse{
		AccessToken:  userID, // если userID — это твой accessToken, либо измени на нужный токен
		RefreshToken: "",
		ExpiresAt:    "",
	}, nil
}

// Login обрабатывает вход пользователя.
// Принимает email и пароль, возвращает JWT-токен или ошибку.
func (h *AuthHandler) Login(ctx context.Context, req *authpb.LoginRequest) (*authpb.AuthResponse, error) {
	h.logger.Info("Login called", zap.String("email", req.Email))
	token, err := h.service.Login(ctx, req.Email, req.Password)
	if err != nil {
		h.logger.Error("Login failed", zap.Error(err))
		// Можно добавить свои ошибки в сервисе, например ErrInvalidCredentials или ErrUserNotFound
		if err == service.ErrInvalidCredentials || err == service.ErrUserNotFound {
			return nil, status.Error(codes.Unauthenticated, "invalid credentials")
		}
		return nil, status.Error(codes.Internal, "internal error")
	}
	return &authpb.AuthResponse{
		AccessToken:  token,
		RefreshToken: "",
		ExpiresAt:    "",
	}, nil
}

// VerifyToken проверяет валидность переданного JWT-токена.
// Возвращает userID, связанный с токеном, или ошибку.
func (h *AuthHandler) VerifyToken(ctx context.Context, req *authpb.VerifyTokenRequest) (*authpb.VerifyTokenResponse, error) {
	h.logger.Info("VerifyToken called")
	userID, err := h.service.VerifyToken(ctx, req.Token)
	if err != nil {
		h.logger.Error("Token verification failed", zap.Error(err))
		return nil, err
	}
	return &authpb.VerifyTokenResponse{UserId: userID}, nil
}
