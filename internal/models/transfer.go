package models

import "time"

type Transfers struct {
	ID               uint      `gorm:"primaryKey" json:"id"`
	SenderWalletID   uint      `gorm:"not null" json:"sender_wallet_id"`
	ReceiverWalletID uint      `gorm:"not null" json:"receiver_wallet_id"`
	Amount           int64     `gorm:"not null" json:"amount"`
	CreatedAt        time.Time `json:"created_at"`

	SenderWallet   *Wallets `gorm:"foreignKey:SenderWalletID" json:"sender_wallet,omitempty"`
	ReceiverWallet *Wallets `gorm:"foreignKey:ReceiverWalletID" json:"receiver_wallet,omitempty"`
}
