package domain

import "context"

// User представляет данные профиля пользователя.
type User struct {
	ID       string // Уникальный идентификатор пользователя
	Name     string // Имя пользователя
	Email    string // Email (нельзя менять через UpdateProfile)
	Currency string // Валюта отображения (например, "RUB", "USD")
	Language string // Предпочтительный язык (например, "ru", "en")
}

// UserService описывает бизнес-логику, связанную с профилем пользователя.
type UserService interface {
	// GetUserByID возвращает профиль по userID.
	GetUserByID(ctx context.Context, id string) (*User, error)

	// UpdateUser обновляет имя, язык и валюту пользователя.
	UpdateUser(ctx context.Context, user User) error
}

// UserRepository описывает интерфейс работы с данными пользователя.
type UserRepository interface {
	// FindByID возвращает пользователя по ID.
	FindByID(ctx context.Context, id string) (*User, error)

	// Update обновляет профиль пользователя.
	Update(ctx context.Context, user User) error
}
