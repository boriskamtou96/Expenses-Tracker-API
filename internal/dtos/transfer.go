package dtos

import "time"

type CreateTransferRequest struct {
	SenderWalletID   int64         `json:"sender_wallet_id" binding:"required,gt=0"`
	ReceiverWalletID int64         `json:"receiver_wallet_id" binding:"required,gt=0"`
	Amount           float64       `json:"amount" binding:"required,gt=0"`
	CreatedAt        time.Duration `json:"created_at"`
}

type CreateTransferResponse struct {
	ID               int64         `json:"id"`
	SenderWalletID   int64         `json:"sender_wallet_id"`
	ReceiverWalletID int64         `json:"receiver_wallet_id"`
	Amount           float64       `json:"amount"`
	CreatedAt        time.Duration `json:"created_at"`
}
