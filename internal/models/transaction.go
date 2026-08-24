package models

import "time"

type Transactions struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	WalletID     uint      `gorm:"not null" json:"wallet_id"`
	Type         string    `gorm:"size:50;not null" json:"type"`
	Amount       int64     `gorm:"not null" json:"amount"`
	BalanceAfter int64     `gorm:"not null" json:"balance_after"`
	CreatedAt    time.Time `json:"created_at"`

	Wallet *Wallets `json:"wallet,omitempty"`
}
