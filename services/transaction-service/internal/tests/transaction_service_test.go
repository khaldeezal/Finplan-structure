package tests

import (
	"context"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"github.com/khaldeezal/Finplan-structure/services/transaction-service/internal/model"
	"github.com/khaldeezal/Finplan-structure/services/transaction-service/internal/repo"
	"github.com/khaldeezal/Finplan-structure/services/transaction-service/internal/service"
)

func TestAddTransaction_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockRepo := NewMockTransactionRepository(ctrl)
	logger := zap.NewNop()
	s := service.NewTransactionService(mockRepo, logger)

	tx := &model.Transaction{
		UserID:      "user-1",
		CategoryID:  "cat-1",
		Type:        model.Income,
		Amount:      100.0,
		Description: "Доход",
		Date:        time.Now(),
	}

	mockRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)
	id, err := s.AddTransaction(context.Background(), tx)
	require.NoError(t, err)
	require.NotEmpty(t, id)
}

func TestAddTransaction_Invalid(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockRepo := NewMockTransactionRepository(ctrl)
	logger := zap.NewNop()
	s := service.NewTransactionService(mockRepo, logger)

	tx := &model.Transaction{
		UserID:      "user-1",
		CategoryID:  "cat-1",
		Type:        model.Income,
		Amount:      0,
		Description: "Ошибка",
		Date:        time.Now(),
	}

	_, err := s.AddTransaction(context.Background(), tx)
	require.Error(t, err)
}

func TestListTransactions_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockRepo := NewMockTransactionRepository(ctrl)
	logger := zap.NewNop()
	s := service.NewTransactionService(mockRepo, logger)

	userID := "user-1"
	want := []*model.Transaction{
		{ID: "tx1", UserID: userID, Amount: 100},
		{ID: "tx2", UserID: userID, Amount: 200},
	}
	mockRepo.EXPECT().ListByUserID(gomock.Any(), userID, gomock.Any(), gomock.Any()).Return(want, nil)
	got, err := s.ListTransactions(context.Background(), userID, 0, 0)
	require.NoError(t, err)
	require.Len(t, got, 2)
}

func TestListTransactions_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockRepo := NewMockTransactionRepository(ctrl)
	logger := zap.NewNop()
	s := service.NewTransactionService(mockRepo, logger)

	userID := "user-2"
	mockRepo.EXPECT().ListByUserID(gomock.Any(), userID, gomock.Any(), gomock.Any()).Return(nil, repo.ErrNotFound)
	got, err := s.ListTransactions(context.Background(), userID, 0, 0)
	require.Error(t, err)
	require.Nil(t, got)
}

func TestDeleteTransaction_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockRepo := NewMockTransactionRepository(ctrl)
	logger := zap.NewNop()
	s := service.NewTransactionService(mockRepo, logger)

	txID := "tx-1"
	userID := "user-1"
	mockRepo.EXPECT().Delete(gomock.Any(), txID, userID).Return(nil)
	err := s.DeleteTransaction(context.Background(), txID, userID)
	require.NoError(t, err)
}

func TestDeleteTransaction_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockRepo := NewMockTransactionRepository(ctrl)
	logger := zap.NewNop()
	s := service.NewTransactionService(mockRepo, logger)

	txID := "tx-x"
	userID := "user-x"
	mockRepo.EXPECT().Delete(gomock.Any(), txID, userID).Return(repo.ErrNotFound)
	err := s.DeleteTransaction(context.Background(), txID, userID)
	require.Error(t, err)
}

func TestGetBalance_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockRepo := NewMockTransactionRepository(ctrl)
	logger := zap.NewNop()
	s := service.NewTransactionService(mockRepo, logger)

	userID := "user-1"
	income, expense, total := 1000.0, 500.0, 500.0
	mockRepo.EXPECT().GetBalance(gomock.Any(), userID).Return(income, expense, total, nil)
	gotIncome, gotExpense, gotTotal, err := s.GetBalance(context.Background(), userID)
	require.NoError(t, err)
	require.Equal(t, income, gotIncome)
	require.Equal(t, expense, gotExpense)
	require.Equal(t, total, gotTotal)
}
