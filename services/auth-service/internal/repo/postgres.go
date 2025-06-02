package repo

import (
	"context"
	"database/sql"
	"fmt"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib" // драйвер для подключения через pgx
	"github.com/khaldeezal/Finplan-structure/services/auth-service/internal/model"
	"go.uber.org/zap"
)

// Реализует сохранение пользователей в PostgreSQL
type PostgresUserRepository struct {
	db     *sql.DB
	logger *zap.Logger
}

// Создаёт новый экземпляр репозитория
func NewPostgresUserRepository(db *sql.DB, logger *zap.Logger) *PostgresUserRepository {
	return &PostgresUserRepository{db: db, logger: logger}
}

// Сохраняет нового пользователя в таблицу users
func (r *PostgresUserRepository) CreateUser(ctx context.Context, user *model.User) (string, error) {
	r.logger.Info("CreateUser called", zap.String("email", user.Email))
	query := `INSERT INTO users (email, password, name) VALUES ($1, $2, $3) RETURNING id`
	err := r.db.QueryRowContext(ctx, query, user.Email, user.Password, user.Name).Scan(&user.ID)
	if err != nil {
		r.logger.Error("failed to insert user", zap.Error(err))
		return "", err
	}
	r.logger.Info("User successfully created", zap.String("userID", user.ID))
	return user.ID, nil
}

// Открывает подключение к PostgreSQL и возвращает *sql.DB
func ConnectDB(logger *zap.Logger) (*sql.DB, error) {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		logger.Error("переменная окружения DATABASE_URL не задана")
		return nil, fmt.Errorf("переменная окружения DATABASE_URL не задана")
	}
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		logger.Error("ошибка подключения к БД", zap.Error(err))
		return nil, fmt.Errorf("ошибка подключения к БД: %w", err)
	}

	// Проверка подключения
	if err := db.Ping(); err != nil {
		logger.Error("БД недоступна", zap.Error(err))
		return nil, fmt.Errorf("БД недоступна: %w", err)
	}

	logger.Info("✅ Успешное подключение к PostgreSQL")
	return db, nil
}
func RunMigrations(db *sql.DB, logger *zap.Logger) error {
	query := `
    CREATE EXTENSION IF NOT EXISTS "pgcrypto";
    CREATE TABLE IF NOT EXISTS users (
        id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
        email TEXT UNIQUE NOT NULL,
        password TEXT NOT NULL,
        name TEXT NOT NULL,
        created_at TIMESTAMP DEFAULT now()
    );`
	_, err := db.Exec(query)
	if err != nil {
		logger.Error("Ошибка миграции БД", zap.Error(err))
		return err
	}
	logger.Info("Миграция БД успешно выполнена")
	return nil
}

type UserRepository interface {
	CreateUser(ctx context.Context, user *model.User) (string, error)
	GetUserByEmail(ctx context.Context, email string) (*model.User, error)
}

var _ UserRepository = (*PostgresUserRepository)(nil)

// Получает пользователя по email из таблицы users
func (r *PostgresUserRepository) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	r.logger.Info("GetUserByEmail called", zap.String("email", email))
	query := `SELECT id, email, password, name FROM users WHERE email = $1`
	row := r.db.QueryRowContext(ctx, query, email)

	var user model.User
	err := row.Scan(&user.ID, &user.Email, &user.Password, &user.Name)
	if err != nil {
		r.logger.Error("failed to get user by email", zap.Error(err))
		return nil, err
	}
	r.logger.Info("User retrieved", zap.String("userID", user.ID))
	return &user, nil
}
