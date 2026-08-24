package server

import (
	"context"
	"expense-tracker/internal/dtos"
	"expense-tracker/internal/services"
	"expense-tracker/internal/utils"

	"github.com/gin-gonic/gin"
)

func (s *Server) register(c *gin.Context) {
	var req dtos.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequestResponse(c, "bad request", err)
		return
	}

	authService := services.NewAuthService(s.db, s.cfg)
	response, err := authService.Register(context.Background(), req)
	if err != nil {
		utils.BadRequestResponse(c, "error while creating user", err)
		return
	}

	utils.CreatedResponse(c, "User create successfully", response)
}

func (s *Server) login(c *gin.Context) {
	var req dtos.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequestResponse(c, "bad request", err)
		return
	}

	authService := services.NewAuthService(s.db, s.cfg)
	response, err := authService.Login(context.Background(), req)
	if err != nil {
		utils.BadRequestResponse(c, "error while login user", err)
		return
	}

	utils.SuccessResponse(c, "User logged successfully", response)
}
