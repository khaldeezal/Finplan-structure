package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/khaldeezal/Finplan-structure/services/transaction-service/internal/model"
	"github.com/khaldeezal/Finplan-structure/services/transaction-service/internal/repo"
	"go.uber.org/zap"
)

// Интерфейс бизнес-логики транзакций
type TransactionService interface {
	AddTransaction(ctx context.Context, tx *model.Transaction) (string, error)
	ListTransactions(ctx context.Context, userID string, limit, offset int) ([]*model.Transaction, error)
	DeleteTransaction(ctx context.Context, transactionID, userID string) error
	GetBalance(ctx context.Context, userID string) (income, expense, balance float64, err error)
}

// Реализует бизнес-логику транзакций
type transactionService struct {
	repo   repo.TransactionRepository
	logger *zap.Logger
	now    func() time.Time
}

// Создаёт новый экземпляр сервиса
func NewTransactionService(r repo.TransactionRepository, logger *zap.Logger) TransactionService {
	return &transactionService{
		repo:   r,
		logger: logger,
		now:    time.Now,
	}
}

func (s *transactionService) AddTransaction(ctx context.Context, tx *model.Transaction) (string, error) {
	if err := tx.Validate(); err != nil {
		s.logger.Error("invalid transaction", zap.Error(err))
		return "", fmt.Errorf("ошибка валидации транзакции: %w", err)
	}
	s.logger.Info("adding model", zap.String("user_id", tx.UserID), zap.Float64("amount", tx.Amount), zap.String("type", string(tx.Type)))
	tx.ID = uuid.New().String()
	tx.CreatedAt = s.now()
	err := s.repo.Create(ctx, tx)
	if err != nil {
		s.logger.Error("failed to add model", zap.Error(err))
	}
	return tx.ID, err
}

func (s *transactionService) ListTransactions(ctx context.Context, userID string, limit, offset int) ([]*model.Transaction, error) {
	s.logger.Info("listing transactions", zap.String("user_id", userID), zap.Int("limit", limit), zap.Int("offset", offset))
	txs, err := s.repo.ListByUserID(ctx, userID, limit, offset)
	if err != nil {
		if err == repo.ErrNotFound {
			return nil, fmt.Errorf("транзакция не найдена: %w", err)
		}
		return nil, err
	}
	return txs, nil
}

func (s *transactionService) DeleteTransaction(ctx context.Context, transactionID, userID string) error {
	s.logger.Info("deleting model", zap.String("transaction_id", transactionID), zap.String("user_id", userID))
	err := s.repo.Delete(ctx, transactionID, userID)
	if err != nil {
		if err == repo.ErrNotFound {
			return fmt.Errorf("транзакция не найдена: %w", err)
		}
		s.logger.Error("failed to delete model", zap.Error(err))
	}
	return err
}

func (s *transactionService) GetBalance(ctx context.Context, userID string) (float64, float64, float64, error) {
	s.logger.Info("getting balance", zap.String("user_id", userID))
	income, expense, balance, err := s.repo.GetBalance(ctx, userID)
	if err != nil {
		if err == repo.ErrNotFound {
			return 0, 0, 0, fmt.Errorf("транзакция не найдена: %w", err)
		}
		s.logger.Error("failed to get balance", zap.Error(err))
	}
	return income, expense, balance, err
}
