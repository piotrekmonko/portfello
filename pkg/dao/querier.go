// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0

package dao

import (
	"context"
)

type Querier interface {
	ExpenseInsert(ctx context.Context, arg *ExpenseInsertParams) error
	ExpenseListByWallet(ctx context.Context, walletID string) ([]*Expense, error)
	ExpenseListByWalletByUser(ctx context.Context, walletID string, userID string) ([]*Expense, error)
	HistoryInsert(ctx context.Context, arg *HistoryInsertParams) error
	HistoryList(ctx context.Context) ([]*History, error)
	LocalUserGetByEmail(ctx context.Context, email string) (*LocalUser, error)
	LocalUserInsert(ctx context.Context, arg *LocalUserInsertParams) error
	LocalUserList(ctx context.Context) ([]*LocalUser, error)
	LocalUserSetPass(ctx context.Context, pwdhash string, email string) error
	LocalUserUpdate(ctx context.Context, roles string, email string) error
	WalletInsert(ctx context.Context, arg *WalletInsertParams) error
	WalletUpdateBalance(ctx context.Context, balance float64, iD string) error
	WalletsByUser(ctx context.Context, userID string) ([]*Wallet, error)
}

var _ Querier = (*Queries)(nil)
