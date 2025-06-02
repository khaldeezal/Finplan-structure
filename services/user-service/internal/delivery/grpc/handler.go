package grpc

import "context"

import (
	userpb "github.com/khaldeezal/Finplan-proto/proto-definitions/gen/user"
	"github.com/khaldeezal/Finplan-structure/services/user-service/internal/domain"
	"github.com/khaldeezal/Finplan-structure/services/user-service/internal/service"
	"go.uber.org/zap"
)

// gRPC интерфейс UserServiceServer, бизнес-логика для работы с юзерами
type UserHandler struct {
	userpb.UnimplementedUserServiceServer
	svc    service.UserService
	logger *zap.Logger
}

// Создает и возвращает новый экземпляр UserHandler со слоем бизнес-логики
func NewHandler(svc service.UserService, logger *zap.Logger) userpb.UserServiceServer {
	return &UserHandler{svc: svc, logger: logger}
}

// Обрабатывает gRPC-запрос на получение профиля пользователя по userID
// Возвращает данные пользователя или ошибку в случае неудачи
func (h *UserHandler) GetUserProfile(ctx context.Context, req *userpb.GetUserProfileRequest) (*userpb.GetUserProfileResponse, error) {
	h.logger.Info("GetUserProfile called", zap.String("userID", req.GetUserId()))
	user, err := h.svc.GetUserByID(ctx, req.GetUserId())
	if err != nil {
		h.logger.Error("failed to get user", zap.Error(err))
		return nil, err
	}
	return &userpb.GetUserProfileResponse{
		UserId:   user.ID,
		Name:     user.Name,
		Email:    user.Email,
		Currency: user.Currency,
		Language: user.Language,
	}, nil
}

// Обрабатывает gRPC-запрос на обновление профиля пользователя
// Обновляет поля Name, Email, Currency и Language. Возвращает успешность операции и ошибку
func (h *UserHandler) UpdateUserProfile(ctx context.Context, req *userpb.UpdateUserProfileRequest) (*userpb.UpdateUserProfileResponse, error) {
	h.logger.Info("UpdateUserProfile called", zap.String("userID", req.GetUserId()))
	err := h.svc.UpdateUser(ctx, domain.User{
		ID:       req.GetUserId(),
		Name:     req.GetName(),
		Currency: req.GetCurrency(),
		Language: req.GetLanguage(),
	})
	if err != nil {
		h.logger.Error("failed to update user", zap.String("userID", req.GetUserId()), zap.Error(err))
		return &userpb.UpdateUserProfileResponse{Success: false}, err
	}
	h.logger.Info("user profile updated", zap.String("userID", req.GetUserId()))
	return &userpb.UpdateUserProfileResponse{Success: true}, nil
}
