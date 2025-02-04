package main

import (
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
	"github.com/riskibarqy/go-template/config"
	"github.com/riskibarqy/go-template/internal/data"
	internalhttp "github.com/riskibarqy/go-template/internal/http"
	"github.com/riskibarqy/go-template/internal/user"
	userPg "github.com/riskibarqy/go-template/internal/user/postgres"
	"github.com/riskibarqy/go-template/models"
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
