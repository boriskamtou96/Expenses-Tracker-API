package dtos

import "time"

type CreateTransactionRequest struct {
	WalletID     int64         `json:"wallet_id"`
	Type         string        `json:"type"`
	Amount       float64       `json:"amount"`
	BalanceAfter float64       `json:"balance_after"`
	CreatedAt    time.Duration `json:"created_at"`
}

type TransactionResponse struct {
	ID           int64                `json:"id"`
	Wallet       CreateWalletResponse `json:"wallet"`
	Type         string               `json:"type"`
	Amount       float64              `json:"amount"`
	BalanceAfter float64              `json:"balance_after"`
	CreatedAt    time.Duration        `json:"created_at"`
}
