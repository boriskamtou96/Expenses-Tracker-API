package models

import "time"

type Users struct {
	ID        int64     `gorm:"primaryKey" json:"id"`
	Username  string    `gorm:"size:255;unique;not null" json:"username"`
	Email     string    `gorm:"size:255;unique;not null" json:"email"`
	Password  string    `gorm:"size:255;not null" json:"-"`
	CreatedAt time.Time `json:"created_at"`

	Wallet *Wallets `gorm:"foreignKey:UserID" json:"wallet,omitempty"`
}
