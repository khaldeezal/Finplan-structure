package handlers

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	authpb "github.com/khaldeezal/Finplan-structure/proto-definitions/gen/auth"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type AuthHandler struct {
	client authpb.AuthServiceClient
	logger *zap.Logger
}

func NewAuthHandler(conn *grpc.ClientConn, logger *zap.Logger) *AuthHandler {
	return &AuthHandler{
		client: authpb.NewAuthServiceClient(conn),
		logger: logger,
	}
}

// Регистрация
func (h *AuthHandler) Register(c *gin.Context) {
	var req struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
		Name     string `json:"name" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("bad register request", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	grpcReq := &authpb.RegisterRequest{
		Email:    req.Email,
		Password: req.Password,
		Name:     req.Name,
	}
	resp, err := h.client.Register(context.Background(), grpcReq)
	if err != nil {
		h.logger.Error("register failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"token": resp.AccessToken})
}

// Логин
func (h *AuthHandler) Login(c *gin.Context) {
	var req struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("bad login request", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	grpcReq := &authpb.LoginRequest{
		Email:    req.Email,
		Password: req.Password,
	}
	resp, err := h.client.Login(context.Background(), grpcReq)
	if err != nil {
		h.logger.Error("login failed", zap.Error(err))
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"token": resp.AccessToken})
}

// Проверка токена (пример)
func (h *AuthHandler) VerifyToken(c *gin.Context) {
	var req struct {
		Token string `json:"token" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("bad verify-token request", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	grpcReq := &authpb.VerifyTokenRequest{
		Token: req.Token,
	}
	resp, err := h.client.VerifyToken(context.Background(), grpcReq)
	if err != nil {
		h.logger.Error("verify failed", zap.Error(err))
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"valid": resp.Valid, "user_id": resp.UserId})
}
