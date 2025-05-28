// User представляет собой модель пользователя для работы с БД и бизнес-логикой.
// Используется в сервисах, хендлерах и репозиториях.
// Хранит ID, Email, хешированный пароль и имя пользователя.
package models

type User struct {
	ID       string
	Email    string
	Password string
	Name     string
}
