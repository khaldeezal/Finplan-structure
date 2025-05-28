package handlers

import (
	"github.com/gin-gonic/gin"
	userpb "github.com/khaldeezal/Finplan-structure/proto-definitions/gen/user"
	"go.uber.org/zap"
)

type UserHandler struct {
	userClient userpb.UserServiceClient
	logger     *zap.Logger
}

func NewUserHandler(userClient userpb.UserServiceClient, logger *zap.Logger) *UserHandler {
	return &UserHandler{
		userClient: userClient,
		logger:     logger,
	}
}

// --- Gin Handlers ---
func (h *UserHandler) GetUserProfile(c *gin.Context) {
	id := c.Param("id")
	req := &userpb.GetUserProfileRequest{UserId: id}
	resp, err := h.userClient.GetUserProfile(c.Request.Context(), req)
	if err != nil {
		h.logger.Error("Ошибка получения профиля пользователя", zap.Error(err))
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, resp)
}

func (h *UserHandler) UpdateUserProfile(c *gin.Context) {
	id := c.Param("id")
	var req userpb.UpdateUserProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Ошибка биндера", zap.Error(err))
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	req.UserId = id
	resp, err := h.userClient.UpdateUserProfile(c.Request.Context(), &req)
	if err != nil {
		h.logger.Error("Ошибка обновления профиля пользователя", zap.Error(err))
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, resp)
}
