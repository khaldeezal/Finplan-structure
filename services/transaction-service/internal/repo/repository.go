package repo

import (
	"context"
	"database/sql"
	"errors"
	"github.com/jmoiron/sqlx"
	"github.com/khaldeezal/Finplan-structure/services/transaction-service/internal/model"
	"go.uber.org/zap"
)

// Если не найдено ни одной строки
var ErrNotFound = errors.New("transaction not found")

// Интерфейс для работы с транзакциями
type TransactionRepository interface {
	Create(ctx context.Context, tx *model.Transaction) error
	ListByUserID(ctx context.Context, userID string, limit, offset int) ([]*model.Transaction, error)
	Delete(ctx context.Context, transactionID, userID string) error
	GetBalance(ctx context.Context, userID string) (income, expense, total float64, err error)
}

// Реализация TransactionRepository
type transactionRepo struct {
	db     *sqlx.DB
	logger *zap.Logger
}

// Создает новый экземпляр transactionRepo
func NewTransactionRepo(db *sqlx.DB, logger *zap.Logger) TransactionRepository {
	return &transactionRepo{db: db, logger: logger}
}

func (r *transactionRepo) Create(ctx context.Context, tx *model.Transaction) error {
	r.logger.Info("creating model", zap.String("user_id", tx.UserID), zap.Float64("amount", tx.Amount), zap.String("type", string(tx.Type)))

	query := `
		INSERT INTO transactions (id, user_id, category_id, type, amount, description, date, created_at)
		VALUES (:id, :user_id, :category_id, :type, :amount, :description, :date, :created_at)
	`
	_, err := r.db.NamedExecContext(ctx, query, tx)
	if err != nil {
		r.logger.Error("failed to create model", zap.Error(err))
	} else {
		r.logger.Info("successfully created transaction", zap.String("id", tx.ID))
	}
	return err
}

func (r *transactionRepo) ListByUserID(ctx context.Context, userID string, limit, offset int) ([]*model.Transaction, error) {
	r.logger.Info("listing transactions", zap.String("user_id", userID), zap.Int("limit", limit), zap.Int("offset", offset))

	query := `
		SELECT id, user_id, category_id, type, amount, description, date, created_at
		FROM transactions
		WHERE user_id = $1
		ORDER BY date DESC
		LIMIT $2 OFFSET $3
	`
	var transactions []*model.Transaction
	err := r.db.SelectContext(ctx, &transactions, query, userID, limit, offset)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}
		r.logger.Error("failed to list transactions", zap.Error(err))
	} else {
		r.logger.Info("found transactions", zap.Int("count", len(transactions)))
	}
	return transactions, err
}

func (r *transactionRepo) Delete(ctx context.Context, transactionID, userID string) error {
	r.logger.Info("deleting model", zap.String("transaction_id", transactionID), zap.String("user_id", userID))

	query := `DELETE FROM transactions WHERE id = $1 AND user_id = $2`
	res, err := r.db.ExecContext(ctx, query, transactionID, userID)
	if err != nil {
		r.logger.Error("failed to delete model", zap.Error(err))
	} else {
		rows, _ := res.RowsAffected()
		if rows == 0 {
			r.logger.Warn("no transaction deleted (not found)", zap.String("transaction_id", transactionID), zap.String("user_id", userID))
			return ErrNotFound
		} else {
			r.logger.Info("transaction deleted", zap.String("transaction_id", transactionID))
		}
	}
	return err
}

func (r *transactionRepo) GetBalance(ctx context.Context, userID string) (float64, float64, float64, error) {
	r.logger.Info("getting balance", zap.String("user_id", userID))

	query := `
		SELECT
			COALESCE(SUM(CASE WHEN type = 'INCOME' THEN amount ELSE 0 END), 0) AS income,
			COALESCE(SUM(CASE WHEN type = 'EXPENSE' THEN amount ELSE 0 END), 0) AS expense
		FROM transactions
		WHERE user_id = $1
	`
	var income, expense float64
	err := r.db.QueryRowContext(ctx, query, userID).Scan(&income, &expense)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, 0, 0, ErrNotFound
		}
		r.logger.Error("failed to get balance", zap.Error(err))
	} else {
		r.logger.Info("balance calculated", zap.Float64("income", income), zap.Float64("expense", expense), zap.Float64("total", income-expense))
	}
	return income, expense, income - expense, err
}
