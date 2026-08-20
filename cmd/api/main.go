package main

import (
	"database/sql"
	"expense-tracker/internal/config"
	"expense-tracker/internal/database"
	"expense-tracker/internal/logger"
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
		err := mainDB.Close()
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to close database")
		}
	}(mainDB)

	log.Info().Msg("Server running....")

}
