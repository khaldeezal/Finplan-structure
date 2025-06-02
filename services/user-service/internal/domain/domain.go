package domain

import "context"

// Данные профиля пользователя
type User struct {
	ID       string
	Name     string
	Email    string
	Currency string
	Language string
}

// Бизнес-логика профиля юзера
type UserService interface {
	// Возвращает профиль по userID
	GetUserByID(ctx context.Context, id string) (*User, error)

	// Обновляет имя, язык и валюту пользователя
	UpdateUser(ctx context.Context, user User) error
}

// Интерфейс работы с данными пользователя
type UserRepository interface {
	// FindByID возвращает пользователя по ID.
	FindByID(ctx context.Context, id string) (*User, error)

	// Update обновляет профиль пользователя.
	Update(ctx context.Context, user User) error
}
