package model

import "time"

type Transaction struct {
	Date          time.Time
	Type          TransactionType
	SortCode      string
	AccountNumber int
	Description   string
	DebitAmount   Amount
	CreditAmount  Amount
	Balance       Amount
}

type Amount float64
type TransactionType string
