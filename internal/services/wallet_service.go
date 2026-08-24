package services

import (
	"expense-tracker/internal/config"
	"expense-tracker/internal/models"
	"fmt"

	"gorm.io/gorm"
)

type WalletService struct {
	db  *gorm.DB
	cfg *config.Config
}

func NewWalletService(db *gorm.DB, cfg *config.Config) *WalletService {
	return &WalletService{
		db:  db,
		cfg: cfg,
	}
}

func (s *WalletService) GetUserWallet(userID any) (*models.Wallets, error) {
	var wallet models.Wallets
	if err := s.db.Where("user_id = ?", userID).First(&wallet).Error; err != nil {
		return nil, fmt.Errorf("wallet with id %s not found", userID)
	}

	return &wallet, nil
}
