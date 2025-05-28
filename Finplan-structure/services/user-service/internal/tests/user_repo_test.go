package tests

import (
	"context"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/khaldeezal/Finplan-structure/services/user-service/internal/domain"
	"github.com/khaldeezal/Finplan-structure/services/user-service/internal/repo"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"regexp"
	"testing"
)

func TestFindByID_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	logger := zap.NewNop()
	userRepo := repo.NewPostgresUserRepository(db, logger)

	expectedUser := &domain.User{
		ID:       "user-123",
		Name:     "Вася",
		Email:    "vasya@example.com",
		Currency: "RUB",
		Language: "ru",
	}

	rows := sqlmock.NewRows([]string{"id", "name", "email", "currency", "language"}).
		AddRow(expectedUser.ID, expectedUser.Name, expectedUser.Email, expectedUser.Currency, expectedUser.Language)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT id, name, email, currency, language FROM users WHERE id = $1")).
		WithArgs(expectedUser.ID).
		WillReturnRows(rows)

	user, err := userRepo.FindByID(context.Background(), expectedUser.ID)
	assert.NoError(t, err)
	assert.Equal(t, expectedUser, user)
	assert.NoError(t, mock.ExpectationsWereMet())
}
