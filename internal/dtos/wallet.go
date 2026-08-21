package dtos

import "time"

type CreateWalletRequest struct {
	UserId    int64         `json:"user_id"`
	Balance   float64       `json:"balance"`
	Currency  string        `json:"currency"`
	CreatedAt time.Duration `json:"created_at"`
	UpdatedAt time.Duration `json:"updated_at"`
}

type CreateWalletResponse struct {
	ID        int64              `json:"id"`
	Owner     CreateUserResponse `json:"owner"`
	Balance   float64            `json:"balance"`
	Currency  string             `json:"currency"`
	CreatedAt time.Duration      `json:"created_at"`
	UpdatedAt time.Duration      `json:"updated_at"`
}

type UpdateWalletRequest struct {
	Currency string `json:"currency"`
}

type UpdateWalletResponse struct {
	ID        int64              `json:"id"`
	Owner     CreateUserResponse `json:"owner"`
	Balance   float64            `json:"balance"`
	Currency  string             `json:"currency"`
	CreatedAt time.Duration      `json:"created_at"`
	UpdatedAt time.Duration      `json:"updated_at"`
}
