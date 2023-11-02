// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.23.0

package dao

import (
	"database/sql"
	"time"
)

// Tracks expenses.
type Expense struct {
	// A base57 encoded uuid.
	ID string
	// Reference to the wallet.
	WalletID    string
	Amount      float64
	Description sql.NullString
	CreatedAt   time.Time
}

type History struct {
	ID string
	// Holds the resource name related resource, such as table name or auth provider name.
	Namespace string
	// Holds the resource ID of the related resource, such as table PK or auth provider id.
	Reference string
	// Describes the event.
	Event string
	// Identifies the user who triggered the event.
	Email     string
	CreatedAt time.Time
}

// Holds together a user and different expenses.
type Wallet struct {
	// A base57-encoded uuid.
	ID string
	// User ID reference to auth provider. This is this wallet Owner.
	UserID    string
	Balance   float64
	Currency  string
	CreatedAt time.Time
}
