package grpc

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/khaldeezal/Finplan-structure/services/user-service/internal/domain"
)

// Поддерживаемые валюты и языки (для примера)
var supportedCurrencies = map[string]bool{
	"RUB": true,
	"USD": true,
	"EUR": true,
}

var supportedLanguages = map[string]bool{
	"ru": true,
	"en": true,
	"de": true,
}

// Определяем ошибку для отсутствующего пользователя
var ErrUserNotFound = errors.New("user not found")

// userService — реализация бизнес-логики для работы с профилем пользователя.
type userService struct {
	repo domain.UserRepository
	// callback для событий обновления профиля (можно использовать для событий, уведомлений и т.п.)
	updateCallback func(user domain.User)
}

// NewUserService конструктор, принимает реализацию репозитория.
func NewUserService(repo domain.UserRepository) domain.UserService {
	return &userService{repo: repo}
}

// SetUpdateCallback позволяет установить callback для обработки событий обновления профиля.
func (s *userService) SetUpdateCallback(cb func(user domain.User)) {
	s.updateCallback = cb
}

// GetUserByID возвращает профиль пользователя по ID.
func (s *userService) GetUserByID(ctx context.Context, id string) (*domain.User, error) {
	user, err := s.repo.FindByID(ctx, id)
	if err != nil {
		// Обработка ошибки отсутствия пользователя
		if errors.Is(err, ErrUserNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return user, nil
}

// UpdateUser обновляет имя, язык и валюту (но не email) пользователя.
func (s *userService) UpdateUser(ctx context.Context, user domain.User) error {
	if user.ID == "" {
		return errors.New("missing user ID")
	}

	// Получаем текущие данные пользователя
	existingUser, err := s.repo.FindByID(ctx, user.ID)
	if err != nil {
		return err
	}

	// Проверяем, что email не меняется
	if user.Email != existingUser.Email {
		return errors.New("email change is not allowed")
	}

	// Проверяем валюту
	if !supportedCurrencies[user.Currency] {
		return fmt.Errorf("unsupported currency: %s", user.Currency)
	}

	// Проверяем язык
	if !supportedLanguages[user.Language] {
		return fmt.Errorf("unsupported language: %s", user.Language)
	}

	// Логируем изменения
	log.Printf("Updating user %s: name='%s', currency='%s', language='%s' at %s",
		user.ID, user.Name, user.Currency, user.Language, time.Now().Format(time.RFC3339))

	// Обновляем в репозитории
	err = s.repo.Update(ctx, user)
	if err != nil {
		return err
	}

	// Вызываем callback после успешного обновления (если установлен)
	if s.updateCallback != nil {
		s.updateCallback(user)
	}

	return nil
}
