package application

import (
	"fmt"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/config"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/container"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/routes"
	"github.com/sajimenezher_meli/meli-frescos-8/pkg/database"
	"log"
	"net/http"
	"os"
)

type Application struct {
	PortServer     string
	FolderDatabase string
}

func (app *Application) SetEnvironment() {

	if err := os.Setenv("folder_database", app.FolderDatabase); err != nil {
		panic("can't set environment")
	}

	if err := os.Setenv("port_server", app.PortServer); err != nil {
		panic("can't set environment")
	}
}

func (app *Application) InitApplication() {
	// 1. Load configuration
	cfg := config.LoadConfig()

	// 2. Initialize database
	db := database.InitDB(cfg)
	defer db.Close()

	c, err := container.NewContainer(db)

	if err != nil {
		log.Fatal(fmt.Sprintf("error initialized container dependencies %v", err))
	}

	router := routes.SetupRoutes(c)

	log.Println(fmt.Sprintf("Server starting on port http://%s/api/v1", app.PortServer))

	if err := http.ListenAndServe(app.PortServer, router); err != nil {
		panic(fmt.Sprintf("Error starting server: %v", err))
	}
}
