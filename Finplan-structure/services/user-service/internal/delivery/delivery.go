package delivery

import (
	"context"
	"log"

	"github.com/khaldeezal/Finplan-structure/services/user-service/internal/domain"
	// userpb — gRPC API, сгенерированный из user.proto
	userpb "github.com/khaldeezal/Finplan-structure/proto-definitions/gen/user"
)

// UserHandler реализует grpc интерфейс UserServiceServer и связывает delivery с бизнес-логикой.
type UserHandler struct {
	userpb.UnimplementedUserServiceServer
	userService domain.UserService
}

// NewUserHandler конструктор хендлера.
func NewUserHandler(userService domain.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

// GetUserProfile — grpc-метод, возвращающий профиль пользователя по user_id.
func (h *UserHandler) GetUserProfile(ctx context.Context, req *userpb.GetUserProfileRequest) (*userpb.GetUserProfileResponse, error) {
	user, err := h.userService.GetUserByID(ctx, req.GetUserId())
	if err != nil {
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

// UpdateUserProfile — grpc-метод для обновления профиля пользователя.
func (h *UserHandler) UpdateUserProfile(ctx context.Context, req *userpb.UpdateUserProfileRequest) (*userpb.UpdateUserProfileResponse, error) {
	err := h.userService.UpdateUser(ctx, domain.User{
		ID:       req.GetUserId(),
		Name:     req.GetName(),
		Currency: req.GetCurrency(),
		Language: req.GetLanguage(),
	})
	if err != nil {
		log.Printf("failed to update user profile: %v", err)
		return &userpb.UpdateUserProfileResponse{Success: false}, err
	}

	return &userpb.UpdateUserProfileResponse{Success: true}, nil
}
