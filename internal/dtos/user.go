package dtos

import "time"

type CreateUserRequest struct {
	Username string `json:"username" binding:"required,min=3,max=255"`
	Email    string `json:"email"    binding:"required"`
	Password string `json:"password" binding:"required,min=6"`
}

type CreateUserResponse struct {
	ID        int64         `json:"id"`
	Username  string        `json:"username"`
	Email     string        `json:"email"`
	CreatedAt time.Duration `json:"created_at"`
}

type LoginRequest struct {
	Email    string `json:"email"    binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	User         CreateUserResponse `json:"user"`
	AccessToken  string             `json:"access_token"`
	RefreshToken string             `json:"refresh_token"`
}
