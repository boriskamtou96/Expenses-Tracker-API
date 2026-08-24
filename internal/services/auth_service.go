package services

import (
	"context"
	"errors"
	"expense-tracker/internal/config"
	"expense-tracker/internal/dtos"
	"expense-tracker/internal/models"
	"expense-tracker/internal/utils"
	"fmt"

	"gorm.io/gorm"
)

type AuthService struct {
	db     *gorm.DB
	config *config.Config
}

func NewAuthService(db *gorm.DB, config *config.Config) *AuthService {
	return &AuthService{
		db:     db,
		config: config,
	}
}

var ErrUserAlreadyExists = errors.New("user already exists")

func (s *AuthService) Register(ctx context.Context, req dtos.CreateUserRequest) (*dtos.CreateUserResponse, error) {
	var existingUser models.Users
	err := s.db.WithContext(ctx).Where("email = ?", req.Email).First(&existingUser).Error
	switch {
	case err == nil:
		return nil, fmt.Errorf("%w: %s", ErrUserAlreadyExists, req.Email)
	case !errors.Is(err, gorm.ErrRecordNotFound):
		return nil, fmt.Errorf("searching user by email: %w", err)
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, fmt.Errorf("hashing password: %w", err)
	}

	user := &models.Users{
		Email:    req.Email,
		Username: req.Username,
		Password: hashedPassword,
	}

	err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(user).Error; err != nil {
			return fmt.Errorf("creating user: %w", err)
		}
		if err := tx.Create(&models.Wallets{UserID: user.ID, Balance: 0}).Error; err != nil {
			return fmt.Errorf("creating wallet: %w", err)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return &dtos.CreateUserResponse{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
	}, nil
}

func (s *AuthService) Login(context context.Context, req dtos.LoginRequest) (*dtos.LoginResponse, error) {
	var user models.Users
	if err := s.db.WithContext(context).Where("email = ?", req.Email).First(&user).Error; err != nil {
		return nil, fmt.Errorf("user not found with email: %s", req.Email)
	}

	if !utils.CheckPassword(user.Password, req.Password) {
		return nil, fmt.Errorf("invalid credentials")
	}
	return s.generateAuthResponse(&user)
}

func (s *AuthService) generateAuthResponse(user *models.Users) (*dtos.LoginResponse, error) {
	accessToken, refreshToken, err := utils.GenerateTokenPair(
		&s.config.JWT,
		user.ID,
		user.Username,
		user.Email,
	)
	if err != nil {
		return nil, err
	}

	return &dtos.LoginResponse{
		User: dtos.CreateUserResponse{
			ID:       user.ID,
			Email:    user.Email,
			Username: user.Username,
		},
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
