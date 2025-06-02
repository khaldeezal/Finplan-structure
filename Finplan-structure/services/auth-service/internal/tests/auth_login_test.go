package tests

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/khaldeezal/Finplan-structure/services/auth-service/internal/services"
	"github.com/stretchr/testify/require"
)

func TestLogin_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAuth := NewMockAuthService(ctrl)
	email := "user@example.com"
	password := "correctpassword"
	ctx := context.Background()
	expectedToken := "jwt-token-123"

	mockAuth.EXPECT().Login(ctx, email, password).Return(expectedToken, nil)

	resp, err := mockAuth.Login(ctx, email, password)
	require.NoError(t, err)
	require.Equal(t, expectedToken, resp)
}

func TestLogin_InvalidCredentials(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAuth := NewMockAuthService(ctrl)
	email := "user@example.com"
	password := "wrongpassword"
	ctx := context.Background()

	mockAuth.EXPECT().Login(ctx, email, password).Return("", services.ErrInvalidCredentials)

	resp, err := mockAuth.Login(ctx, email, password)
	require.ErrorIs(t, err, services.ErrInvalidCredentials)
	require.Equal(t, "", resp)
}

func TestLogin_UserNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAuth := NewMockAuthService(ctrl)
	email := "noone@example.com"
	password := "any"
	ctx := context.Background()

	mockAuth.EXPECT().Login(ctx, email, password).Return("", services.ErrUserNotFound)

	resp, err := mockAuth.Login(ctx, email, password)
	require.ErrorIs(t, err, services.ErrUserNotFound)
	require.Equal(t, "", resp)
}

func TestLogin_EmptyFields(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAuth := NewMockAuthService(ctrl)
	ctx := context.Background()

	mockAuth.EXPECT().Login(ctx, "", "").Return("", services.ErrInvalidCredentials)

	resp, err := mockAuth.Login(ctx, "", "")
	require.ErrorIs(t, err, services.ErrInvalidCredentials)
	require.Equal(t, "", resp)
}

func TestLogin_InternalError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAuth := NewMockAuthService(ctrl)
	email := "user@example.com"
	password := "any"
	ctx := context.Background()
	internalErr := errors.New("db connection failed")

	mockAuth.EXPECT().Login(ctx, email, password).Return("", internalErr)

	resp, err := mockAuth.Login(ctx, email, password)
	require.ErrorIs(t, err, internalErr)
	require.Equal(t, "", resp)
}
