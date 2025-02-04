package main

import (
	"context"
	"fmt"
	"log"

	"github.com/ancalabrese/reload"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/riskibarqy/go-template/config"
	"github.com/riskibarqy/go-template/databases"
	"github.com/riskibarqy/go-template/internal/data"
	internalhttp "github.com/riskibarqy/go-template/internal/http"
	"github.com/riskibarqy/go-template/internal/redis"
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

func initMetadataConfig() {
	ctx := context.Background()
	rc, err := reload.New(ctx)
	if err != nil {
		log.Fatalln(err)
		return
	}

	config.MetadataConfig = &config.Metadata{}

	go func() {
		for {
			select {
			case err := <-rc.GetErrChannel():
				log.Printf("Received err: %v", err)
			case conf := <-rc.GetReloadChan():
				log.Println("Received new config [", conf.FilePath, "]:", conf.Config)
			}
		}
	}()

	err = rc.AddConfiguration("./metadata.json", &config.MetadataConfig)
	if err != nil {
		panic(err)
	}

	<-ctx.Done()
}

func main() {
	go initMetadataConfig()
	config.GetConfiguration()

	databases.Init()
	redis.Init()

	// Print the current mode
	fmt.Printf("Running in %s mode\n", config.AppConfig.AppMode)

	// Example: Conditional logic based on the mode
	if config.AppConfig.AppMode == "development" {
		// Development-specific settings
		fmt.Println("Development settings applied")
	} else {
		// Production-specific settings
		fmt.Println("Production settings applied")
	}

	dataManager := data.NewManager(config.AppConfig.DatabaseClient)
	internalServices := buildInternalServices(config.AppConfig.DatabaseClient, config.AppConfig)

	s := internalhttp.NewServer(
		config.AppConfig,
		dataManager,
		internalServices.userService,
	)

	s.Serve()
}
