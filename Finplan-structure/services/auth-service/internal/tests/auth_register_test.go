package tests

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

var ErrUserAlreadyExists = errors.New("user already exists")
var ErrInvalidEmail = errors.New("invalid email")

func TestRegister_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockService := NewMockAuthService(ctrl)
	ctx := context.Background()

	mockService.EXPECT().Register(ctx, "test@example.com", "password123", "Test User").
		Return("user-123", nil)

	resp, err := mockService.Register(ctx, "test@example.com", "password123", "Test User")
	require.NoError(t, err)
	require.Equal(t, "user-123", resp)
}

func TestRegister_AlreadyExists(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockService := NewMockAuthService(ctrl)
	ctx := context.Background()

	mockService.EXPECT().Register(ctx, "existing@example.com", "password123", "Existing User").
		Return("", ErrUserAlreadyExists)

	resp, err := mockService.Register(ctx, "existing@example.com", "password123", "Existing User")
	require.Error(t, err)
	require.Equal(t, ErrUserAlreadyExists, err)
	require.Equal(t, "", resp)
}

func TestRegister_InvalidEmail(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockService := NewMockAuthService(ctrl)
	ctx := context.Background()

	mockService.EXPECT().Register(ctx, "invalid-email", "password123", "Invalid Email User").
		Return("", ErrInvalidEmail)

	resp, err := mockService.Register(ctx, "invalid-email", "password123", "Invalid Email User")
	require.Error(t, err)
	require.Equal(t, ErrInvalidEmail, err)
	require.Equal(t, "", resp)
}

func TestRegister_EmptyFields(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockService := NewMockAuthService(ctrl)
	ctx := context.Background()

	mockService.EXPECT().Register(ctx, "", "", "").
		Return("", errors.New("fields cannot be empty"))

	resp, err := mockService.Register(ctx, "", "", "")
	require.Error(t, err)
	require.EqualError(t, err, "fields cannot be empty")
	require.Equal(t, "", resp)
}
