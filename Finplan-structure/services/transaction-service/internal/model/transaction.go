package model

import (
	"encoding/json"
	"errors"
	"time"
)

type TransactionType string

const (
	Income  TransactionType = "INCOME"
	Expense TransactionType = "EXPENSE"
)

// Transaction представляет доменную модель транзакции
type Transaction struct {
	ID          string          `db:"id" json:"id"`                   // Уникальный идентификатор транзакции
	UserID      string          `db:"user_id" json:"user_id"`         // Идентификатор пользователя
	CategoryID  string          `db:"category_id" json:"category_id"` // Идентификатор категории
	Type        TransactionType `db:"type" json:"type"`               // Тип транзакции: доход или расход
	Amount      float64         `db:"amount" json:"amount"`           // Сумма транзакции
	Description string          `db:"description" json:"description"` // Описание транзакции
	Date        time.Time       `db:"date" json:"date"`               // Дата, к которой относится транзакция
	CreatedAt   time.Time       `db:"created_at" json:"created_at"`   // Дата создания записи
}

type transactionAlias Transaction

func (t *Transaction) MarshalJSON() ([]byte, error) {
	type Alias transactionAlias
	return json.Marshal(&struct {
		Date      string `json:"date"`
		CreatedAt string `json:"created_at"`
		*Alias
	}{
		Date:      t.Date.Format("2006-01-02"),
		CreatedAt: t.CreatedAt.Format(time.RFC3339),
		Alias:     (*Alias)(t),
	})
}

func (t *Transaction) UnmarshalJSON(data []byte) error {
	type Alias transactionAlias
	aux := &struct {
		Date      string `json:"date"`
		CreatedAt string `json:"created_at"`
		*Alias
	}{
		Alias: (*Alias)(t),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	var err error
	t.Date, err = time.Parse("2006-01-02", aux.Date)
	if err != nil {
		return err
	}
	t.CreatedAt, err = time.Parse(time.RFC3339, aux.CreatedAt)
	if err != nil {
		return err
	}
	return nil
}

func (t *Transaction) Validate() error {
	if t.Amount <= 0 {
		return errors.New("amount must be greater than zero")
	}
	if t.Type != Income && t.Type != Expense {
		return errors.New("invalid transaction type")
	}
	if t.Date.After(time.Now()) {
		return errors.New("date cannot be in the future")
	}
	return nil
}
