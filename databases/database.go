package databases

import (
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
	"github.com/riskibarqy/go-template/config"
)

// Init Connect to the database
func Init() {
	fmt.Println("DB Connection String:", config.AppConfig.DBConnectionString)

	// Open new database connection
	db, err := sqlx.Open("postgres", config.AppConfig.DBConnectionString)
	if err != nil {
		log.Fatalf("Failed to reconnect to the database: %v", err)
	}

	// Assign new connection to AppConfig
	config.AppConfig.DatabaseClient = db

	log.Println("Successfully connected to the database")
}
