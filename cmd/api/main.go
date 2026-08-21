package main

import (
	"database/sql"
	"expense-tracker/internal/config"
	"expense-tracker/internal/database"
	"expense-tracker/internal/logger"
	"expense-tracker/internal/server"
)

func main() {
	log := logger.New()

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to load configs")
	}

	db, err := database.New(cfg.Database)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to database")
	}

	mainDB, err := db.DB()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to establish database connection")
	}

	defer func(mainDB *sql.DB) {
		if closeErr := mainDB.Close(); closeErr != nil {
			log.Fatal().Err(closeErr).Msg("Failed to close database")
		}
	}(mainDB)

	srv := server.New(db, cfg, log)

	srv.Run()
}
