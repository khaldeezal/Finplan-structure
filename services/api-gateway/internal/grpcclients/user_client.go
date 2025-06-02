package grpcclients

import (
	"context"
	"fmt"

	userpb "github.com/khaldeezal/Finplan-proto/proto-definitions/gen/user"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type UserClient struct {
	client userpb.UserServiceClient
	logger *zap.Logger
}

func NewUserClient(conn *grpc.ClientConn, logger *zap.Logger) *UserClient {
	return &UserClient{
		client: userpb.NewUserServiceClient(conn),
		logger: logger,
	}
}

// Получение профиля пользователя
func (c *UserClient) GetUserProfile(ctx context.Context, userID string) (*userpb.GetUserProfileResponse, error) {
	req := &userpb.GetUserProfileRequest{
		UserId: userID,
	}
	resp, err := c.client.GetUserProfile(ctx, req)
	if err != nil {
		c.logger.Error("failed to get user profile", zap.Error(err))
		return nil, fmt.Errorf("failed to get user profile: %w", err)
	}
	return resp, nil
}

// Обновление профиля пользователя
func (c *UserClient) UpdateUserProfile(ctx context.Context, userID, name, email, currency, language string) (bool, error) {
	req := &userpb.UpdateUserProfileRequest{
		UserId:   userID,
		Name:     name,
		Currency: currency,
		Language: language,
	}
	resp, err := c.client.UpdateUserProfile(ctx, req)
	if err != nil {
		c.logger.Error("failed to update user profile", zap.Error(err))
		return false, fmt.Errorf("failed to update user profile: %w", err)
	}
	return resp.Success, nil
}
