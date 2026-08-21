package server

import (
	"expense-tracker/internal/config"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

const (
	dbReadWriteTimeout = 15
	idleTimeOut        = 30
)

type Server struct {
	db  *gorm.DB
	cfg *config.Config
	log zerolog.Logger
}

func New(db *gorm.DB, cfg *config.Config, log zerolog.Logger) *Server {
	return &Server{
		db:  db,
		cfg: cfg,
		log: log,
	}
}

func (s *Server) SetupRoutes() *gin.Engine {
	r := gin.New()

	r.Use(gin.Recovery())
	r.Use(gin.Logger())
	r.Use(s.corsMiddleware())

	r.GET("/health", healthCheckHandler)

	// Auth endpoints
	r.POST("/api/v1/auth/register") // register
	r.POST("/api/v1/auth/login")    // login

	// Wallet endpoints
	r.GET("/api/v1/wallet")           // get user wallet
	r.POST("/api/v1/wallet/deposti")  // user make a deposit
	r.POST("/api/v1/wallet/withdraw") // user make a withdraw

	// Transactions endpoints
	r.GET("/api/v1/transactions")                 // get transactions list
	r.GET("/api/v1/transactions/:id")             // get transactions by id
	r.GET("/api/v1/transactions?type=withdrawal") // get transactions by type
	r.GET("/api/v1/transactions?page=1&limit=20") // get paginated transactions

	// Transfers
	r.POST("/api/v1/transfers") // perform transfer between users

	return r
}

func (s *Server) Run() {
	httpServer := &http.Server{
		Addr:         fmt.Sprintf(":%s", s.cfg.Server.Port),
		Handler:      s.SetupRoutes(),
		ReadTimeout:  dbReadWriteTimeout * time.Second,
		WriteTimeout: dbReadWriteTimeout * time.Second,
		IdleTimeout:  idleTimeOut * time.Second,
	}

	if err := httpServer.ListenAndServe(); err != nil {
		s.log.Fatal().Err(err).Msg("failed to start server")
	}
	s.log.Info().Str("port", s.cfg.Server.Port).Msg("Server running on")
}

func healthCheckHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"app": "expense tracker", "status": "Ok"})
}

func (s *Server) corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Methods", "Content-Type, Authorization")

		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}
