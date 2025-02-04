package main

import (
	"log"

	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	"github.com/riskibarqy/go-template/config"
	"github.com/riskibarqy/go-template/databases"
)

func main() {
	// Get configuration
	cfg, err := config.GetConfiguration()
	if err != nil {
		log.Fatalln("Error getting configuration:", err)
	}

	databases.NewMigrator(cfg.DBConnectionString).Run("up")
}
