package main

import (
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
	"github.com/riskibarqy/99backend-challenge/config"
	"github.com/riskibarqy/99backend-challenge/internal/data"
	internalhttp "github.com/riskibarqy/99backend-challenge/internal/http"
	"github.com/riskibarqy/99backend-challenge/internal/user"
	userPg "github.com/riskibarqy/99backend-challenge/internal/user/postgres"
	"github.com/riskibarqy/99backend-challenge/models"
)

// InternalServices represents all the internal domain services
type InternalServices struct {
	userService user.ServiceInterface
}

func buildInternalServices(db *sqlx.DB, _ *config.Config) *InternalServices {
	userPostgresStorage := userPg.NewPostgresStorage(
		data.NewPostgresStorage(db, "user", models.User{}),
	)

	userService := user.NewService(userPostgresStorage)
	return &InternalServices{
		userService: userService,
	}
}

func main() {
	config, err := config.GetConfiguration()
	if err != nil {
		log.Fatalln("failed to get configuration: ", err)
	}
	db, err := sqlx.Open("postgres", config.DBConnectionString)
	if err != nil {
		log.Fatalln("failed to open database x: ", err)
	}
	defer db.Close()

	dataManager := data.NewManager(db)
	internalServices := buildInternalServices(db, config)

	fmt.Println("DB CONNECTION = ", config.DBConnectionString)

	s := internalhttp.NewServer(
		config,
		dataManager,
		internalServices.userService,
	)
	s.Serve()
}
