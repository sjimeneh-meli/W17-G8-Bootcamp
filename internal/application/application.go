package application

import (
	"fmt"
	"log"
	"net/http"

	"github.com/sajimenezher_meli/meli-frescos-8/internal/config"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/container"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/routes"
	"github.com/sajimenezher_meli/meli-frescos-8/pkg/database"
)

type Application struct {
	PortServer string
}

func (app *Application) InitApplication() {
	// 1. Load configuration
	cfg := config.LoadConfig()

	// 2. Initialize database
	db := database.InitDB(cfg)
	defer db.Close()

	c, err := container.NewContainer(db)

	if err != nil {
		log.Fatalf("error initialized container dependencies %v", err)
	}

	router := routes.SetupRoutes(c)

	log.Printf("Server starting on port http://%s/api/v1", app.PortServer)

	if err := http.ListenAndServe(app.PortServer, router); err != nil {
		panic(fmt.Sprintf("Error starting server: %v", err))
	}
}
