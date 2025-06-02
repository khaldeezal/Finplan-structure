package tests

import (
	"context"
	"errors"
	repo2 "github.com/khaldeezal/Finplan-structure/services/auth-service/internal/repo"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/khaldeezal/Finplan-structure/services/auth-service/internal/model"
	"github.com/lib/pq"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestCreateUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	logger := zap.NewNop()
	repo := repo2.NewPostgresUserRepository(db, logger)

	user := &model.User{
		ID:       "123",
		Email:    "test@example.com",
		Password: "hashedpassword",
		Name:     "Test User",
	}

	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO users (email, password, name) VALUES ($1, $2, $3) RETURNING id`)).
		WithArgs(user.Email, user.Password, user.Name).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(user.ID))

	_, err = repo.CreateUser(context.Background(), user)
	require.NoError(t, err)

	require.NoError(t, mock.ExpectationsWereMet())
}

func TestCreateUser_DuplicateEmail(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	logger := zap.NewNop()
	repo := repo2.NewPostgresUserRepository(db, logger)

	user := &model.User{
		ID:       "123",
		Email:    "duplicate@example.com",
		Password: "hashedpassword",
		Name:     "Test User",
	}

	pqErr := &pq.Error{Code: "23505"}

	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO users (email, password, name) VALUES ($1, $2, $3) RETURNING id`)).
		WithArgs(user.Email, user.Password, user.Name).
		WillReturnError(pqErr)

	_, err = repo.CreateUser(context.Background(), user)
	require.Error(t, err)

	require.NoError(t, mock.ExpectationsWereMet())
}

func TestCreateUser_DBError(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	logger := zap.NewNop()
	repo := repo2.NewPostgresUserRepository(db, logger)

	user := &model.User{
		ID:       "123",
		Email:    "error@example.com",
		Password: "hashedpassword",
		Name:     "Test User",
	}

	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO users (email, password, name) VALUES ($1, $2, $3) RETURNING id`)).
		WithArgs(user.Email, user.Password, user.Name).
		WillReturnError(errors.New("some db error"))

	_, err = repo.CreateUser(context.Background(), user)
	require.Error(t, err)

	require.NoError(t, mock.ExpectationsWereMet())
}
