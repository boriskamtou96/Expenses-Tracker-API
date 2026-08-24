package server

import (
	"expense-tracker/internal/services"
	"expense-tracker/internal/utils"
	"fmt"

	"github.com/gin-gonic/gin"
)

func (s *Server) getUserWallet(c *gin.Context) {
	userID, exist := c.Get("user_id")
	if !exist {
		utils.NotFoundResponse(c, fmt.Sprintf("user with id %s not found", userID), nil)
		return
	}

	walletService := services.NewWalletService(s.db, s.cfg)

	wallet, err := walletService.GetUserWallet(userID)
	if err != nil {
		utils.BadRequestResponse(c, fmt.Sprintf("error while fetching wallet data"), err)
		return
	}

	utils.SuccessResponse(c, "wallet retrieving successfully", wallet)
}
