package tests

import (
	"context"
	"database/sql"
	"github.com/jmoiron/sqlx"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/khaldeezal/Finplan-structure/services/transaction-service/internal/model"
	"github.com/khaldeezal/Finplan-structure/services/transaction-service/internal/repo"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func newTestRepo(t *testing.T) (repo.TransactionRepository, sqlmock.Sqlmock) {
	sqlDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	logger := zap.NewNop()
	db := sqlx.NewDb(sqlDB, "sqlmock")
	r := repo.NewTransactionRepo(db, logger)
	return r, mock
}

func TestCreate_Success(t *testing.T) {
	r, mock := newTestRepo(t)
	tx := &model.Transaction{
		ID:          "tx1",
		UserID:      "user-1",
		CategoryID:  "cat-1",
		Type:        "INCOME",
		Amount:      100.0,
		Description: "Доход",
		Date:        time.Now(),
		CreatedAt:   time.Now(),
	}
	mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO transactions`)).
		WithArgs(tx.ID, tx.UserID, tx.CategoryID, tx.Type, tx.Amount, tx.Description, sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err := r.Create(context.Background(), tx)
	require.NoError(t, err)
}

func TestCreate_Error(t *testing.T) {
	r, mock := newTestRepo(t)
	tx := &model.Transaction{ID: "tx2"}
	mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO transactions`)).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnError(sql.ErrConnDone)

	err := r.Create(context.Background(), tx)
	require.Error(t, err)
}

func TestListByUserID_Success(t *testing.T) {
	r, mock := newTestRepo(t)
	userID := "user-1"
	rows := sqlmock.NewRows([]string{"id", "user_id", "category_id", "type", "amount", "description", "date", "created_at"}).
		AddRow("tx1", userID, "cat-1", "INCOME", 100.0, "Доход", time.Now(), time.Now()).
		AddRow("tx2", userID, "cat-2", "EXPENSE", 50.0, "Расход", time.Now(), time.Now())

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT`)).
		WithArgs(userID, 10, 0).
		WillReturnRows(rows)

	txs, err := r.ListByUserID(context.Background(), userID, 10, 0)
	require.NoError(t, err)
	require.Len(t, txs, 2)
}

func TestListByUserID_Empty(t *testing.T) {
	r, mock := newTestRepo(t)
	userID := "user-x"
	rows := sqlmock.NewRows([]string{"id", "user_id", "category_id", "type", "amount", "description", "date", "created_at"})

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT`)).
		WithArgs(userID, 10, 0).
		WillReturnRows(rows)

	txs, err := r.ListByUserID(context.Background(), userID, 10, 0)
	require.NoError(t, err)
	require.Len(t, txs, 0)
}

func TestListByUserID_Error(t *testing.T) {
	r, mock := newTestRepo(t)
	userID := "fail"
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT`)).
		WithArgs(userID, 10, 0).
		WillReturnError(sql.ErrConnDone)

	txs, err := r.ListByUserID(context.Background(), userID, 10, 0)
	require.Error(t, err)
	require.Nil(t, txs)
}

func TestDelete_Success(t *testing.T) {
	r, mock := newTestRepo(t)
	mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM transactions`)).
		WithArgs("tx1", "user-1").
		WillReturnResult(sqlmock.NewResult(1, 1))

	err := r.Delete(context.Background(), "tx1", "user-1")
	require.NoError(t, err)
}

func TestDelete_NotFound(t *testing.T) {
	r, mock := newTestRepo(t)
	mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM transactions`)).
		WithArgs("notx", "user-1").
		WillReturnResult(sqlmock.NewResult(0, 0)) // 0 row affected

	err := r.Delete(context.Background(), "notx", "user-1")
	require.Error(t, err)
}

func TestGetBalance_Error(t *testing.T) {
	r, mock := newTestRepo(t)
	userID := "fail"
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT`)).
		WithArgs(userID).
		WillReturnError(sql.ErrConnDone)

	income, expense, total, err := r.GetBalance(context.Background(), userID)
	require.Error(t, err)
	require.Zero(t, income)
	require.Zero(t, expense)
	require.Zero(t, total)
}
