package main

import (
	"log"

	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	"github.com/riskibarqy/99backend-challenge/config"
	"github.com/riskibarqy/99backend-challenge/databases"
)

func main() {
	// Get configuration
	cfg, err := config.GetConfiguration()
	if err != nil {
		log.Fatalln("Error getting configuration:", err)
	}

	databases.NewMigrator(cfg.DBConnectionString).Run("up")
}
