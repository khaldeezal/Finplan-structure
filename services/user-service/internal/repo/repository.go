package repo

import (
	"context"
	"database/sql"

	"github.com/khaldeezal/Finplan-structure/services/user-service/internal/domain"
	"go.uber.org/zap"
)

// Методы для работы с пользователями
type UserRepository interface {
	FindByID(ctx context.Context, id string) (*domain.User, error)
	Update(ctx context.Context, user domain.User) error
}

// Реализует UserRepository, работает с PostgreSQL
type PostgresUserRepository struct {
	db     *sql.DB
	logger *zap.Logger
}

// Создаёт новый экземпляр репозитория
func NewPostgresUserRepository(db *sql.DB, logger *zap.Logger) *PostgresUserRepository {
	return &PostgresUserRepository{db: db, logger: logger}
}

// Возвращает пользователя по ID из БД
func (r *PostgresUserRepository) FindByID(ctx context.Context, id string) (*domain.User, error) {
	r.logger.Info("FindByID called", zap.String("userID", id))

	query := `SELECT id, name, email, currency, language FROM users WHERE id = $1`
	row := r.db.QueryRowContext(ctx, query, id)

	var user domain.User
	err := row.Scan(&user.ID, &user.Name, &user.Email, &user.Currency, &user.Language)
	if err != nil {
		r.logger.Error("FindByID failed", zap.String("userID", id), zap.Error(err))
		return nil, err
	}

	r.logger.Info("User found", zap.String("userID", user.ID))
	return &user, nil
}

// Обновляет данные профиля пользователя
func (r *PostgresUserRepository) Update(ctx context.Context, user domain.User) error {
	r.logger.Info("Update called", zap.String("userID", user.ID))

	query := `UPDATE users SET name = $1, currency = $2, language = $3 WHERE id = $4`
	_, err := r.db.ExecContext(ctx, query, user.Name, user.Currency, user.Language, user.ID)
	if err != nil {
		r.logger.Error("Update failed", zap.String("userID", user.ID), zap.Error(err))
		return err
	}

	r.logger.Info("User updated", zap.String("userID", user.ID))
	return nil
}
