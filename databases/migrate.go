package databases

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

type Migrator struct {
	DBURL string
}

func NewMigrator(dbURL string) *Migrator {
	return &Migrator{DBURL: dbURL}
}

func (m *Migrator) Run(action string) {
	// Connect to the database
	db, err := sql.Open("postgres", m.DBURL)
	if err != nil {
		log.Fatal("Failed to connect to the database:", err)
	}
	defer db.Close()

	// Create migration driver
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		log.Fatal("Failed to create migration driver:", err)
	}

	migrator, err := migrate.NewWithDatabaseInstance(
		"file://databases/migrations",
		"postgres", driver,
	)
	if err != nil {
		log.Fatal("Failed to initialize migration:", err)
	}

	// Execute the specified migration action
	switch action {
	case "up":
		err = migrator.Up()
		if err != nil && err != migrate.ErrNoChange {
			log.Fatal("Migration up failed:", err)
		}
		fmt.Println("Migrations applied successfully!")

	case "down":
		err = migrator.Steps(-1) // Rollback the last applied migration
		if err != nil && err != migrate.ErrNoChange {
			log.Fatal("Migration down failed:", err)
		}
		fmt.Println("Last migration rolled back!")

	default:
		log.Fatal("Invalid migration action. Use 'up' or 'down'")
	}
}
