package models

import "time"

type Wallets struct {
	ID        int64     `gorm:"primaryKey" json:"id"`
	UserID    int64     `gorm:"unique;not null" json:"user_id"`
	Balance   int64     `gorm:"not null;default:0" json:"balance"`
	CreatedAt time.Time `json:"created_at"`

	User         *Users         `json:"user,omitempty"`
	Transactions []Transactions `gorm:"foreignKey:WalletID" json:"transactions,omitempty"`
}
