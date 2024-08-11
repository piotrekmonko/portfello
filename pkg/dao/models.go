// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0

package dao

import (
	"database/sql"
	"time"
)

type Expense struct {
	ID          string
	WalletID    string
	Amount      float64
	Description sql.NullString
	CreatedAt   time.Time
}

type History struct {
	ID        string
	Namespace string
	Reference string
	Event     string
	Email     string
	CreatedAt time.Time
}

type LocalUser struct {
	ID          string
	Email       string
	DisplayName string
	Roles       string
	Pwdhash     string
	CreatedAt   time.Time
}

type Wallet struct {
	ID        string
	UserID    string
	Balance   float64
	Currency  string
	CreatedAt time.Time
}
