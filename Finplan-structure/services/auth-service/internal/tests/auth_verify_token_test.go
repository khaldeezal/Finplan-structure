package tests

import (
	"context"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/khaldeezal/Finplan-structure/services/auth-service/internal/services"
	"github.com/stretchr/testify/require"
	"testing"
)

//go:generate mockgen -destination=mock_auth_service.go -package=tests github.com/finplan/Finplan-structure/services/auth-service/internal/services AuthService

func TestVerifyToken_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockAuth := NewMockAuthService(ctrl)
	ctx := context.Background()
	const (
		token  = "valid-token"
		userID = "user123"
	)
	mockAuth.EXPECT().VerifyToken(ctx, token).Return(userID, nil)
	gotUserID, err := mockAuth.VerifyToken(ctx, token)
	require.NoError(t, err)
	require.Equal(t, userID, gotUserID)
}

func TestVerifyToken_InvalidToken(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockAuth := NewMockAuthService(ctrl)
	ctx := context.Background()
	const token = "invalid-token"
	mockAuth.EXPECT().VerifyToken(ctx, token).Return("", services.ErrInvalidToken)
	gotUserID, err := mockAuth.VerifyToken(ctx, token)
	require.ErrorIs(t, err, services.ErrInvalidToken)
	require.Empty(t, gotUserID)
}

func TestVerifyToken_UserNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockAuth := NewMockAuthService(ctrl)
	ctx := context.Background()
	const token = "no-userid-token"
	mockAuth.EXPECT().VerifyToken(ctx, token).Return("", services.ErrUserNotFound)
	gotUserID, err := mockAuth.VerifyToken(ctx, token)
	require.ErrorIs(t, err, services.ErrUserNotFound)
	require.Empty(t, gotUserID)
}

func TestVerifyToken_InternalError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockAuth := NewMockAuthService(ctrl)
	ctx := context.Background()
	const token = "any-token"
	internalErr := errors.New("unexpected fail")
	mockAuth.EXPECT().VerifyToken(ctx, token).Return("", internalErr)
	gotUserID, err := mockAuth.VerifyToken(ctx, token)
	require.ErrorIs(t, err, internalErr)
	require.Empty(t, gotUserID)
}
