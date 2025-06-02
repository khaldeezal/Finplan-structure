package tests

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"

	"github.com/khaldeezal/Finplan-structure/services/user-service/internal/domain"
	"github.com/khaldeezal/Finplan-structure/services/user-service/internal/service"
)

func TestGetUserByID_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := NewMockUserRepository(ctrl)
	logger := zap.NewNop()
	userService := service.NewUserService(mockRepo, logger)

	expectedUser := domain.User{
		ID:       "user-123",
		Name:     "Андрей",
		Email:    "test@mail.com",
		Currency: "RUB",
		Language: "ru",
	}
	mockRepo.EXPECT().FindByID(gomock.Any(), "user-123").Return(&expectedUser, nil)

	user, err := userService.GetUserByID(context.Background(), "user-123")
	assert.NoError(t, err)
	assert.Equal(t, expectedUser, *user)
}
