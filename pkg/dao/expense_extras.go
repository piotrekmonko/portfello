package dao

import "time"

func (e *Expense) IsOperation() {}

func (e *Expense) GetID() string {
	return e.ID
}

func (e *Expense) GetWalletID() string {
	return e.WalletID
}

func (e *Expense) GetAmount() float64 {
	return e.Amount
}

func (e *Expense) GetDescription() *string {
	if !e.Description.Valid {
		return nil
	}
	return &e.Description.String
}

func (e *Expense) GetCreatedAt() time.Time {
	return e.CreatedAt.UTC()
}
