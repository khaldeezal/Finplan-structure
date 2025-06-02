package service

import (
	"context"
	"errors"

	"github.com/khaldeezal/Finplan-structure/services/user-service/internal/domain"
	"go.uber.org/zap"
)

// Бизнес-логика работы с профилем юзера
type UserService interface {
	GetUserByID(ctx context.Context, id string) (*domain.User, error)
	UpdateUser(ctx context.Context, user domain.User) error
}

// Реализация
type userService struct {
	repo   domain.UserRepository
	logger *zap.Logger
}

// Конструктор, принимает реализацию репозитория
func NewUserService(repo domain.UserRepository, logger *zap.Logger) domain.UserService {
	return &userService{repo: repo, logger: logger}
}

// Возвращает профиль пользователя по ID
func (s *userService) GetUserByID(ctx context.Context, id string) (*domain.User, error) {
	s.logger.Info("GetUserByID called", zap.String("userID", id))

	user, err := s.repo.FindByID(ctx, id)
	if err != nil {
		s.logger.Error("failed to find user", zap.String("userID", id), zap.Error(err))
		return nil, err
	}
	s.logger.Info("user found", zap.String("userID", user.ID))
	return user, nil
}

// Обновляет имя, язык и валюту (но не email) пользователя
func (s *userService) UpdateUser(ctx context.Context, user domain.User) error {
	s.logger.Info("UpdateUser called", zap.String("userID", user.ID))

	if user.ID == "" {
		err := errors.New("missing user ID")
		s.logger.Error("update failed: empty ID", zap.Error(err))
		return err
	}
	if err := s.repo.Update(ctx, user); err != nil {
		s.logger.Error("failed to update user", zap.String("userID", user.ID), zap.Error(err))
		return err
	}
	s.logger.Info("user updated", zap.String("userID", user.ID))
	return nil
}
