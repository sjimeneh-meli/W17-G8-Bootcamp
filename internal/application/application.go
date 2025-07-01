package application

import (
	"fmt"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/handlers"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/routes"
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
	router := chi.NewRouter()

	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	employeeHandler := handlers.GetEmployeeHandler()

	router.Route("/api/v1/employees", func(r chi.Router) {
		r.Get("/", employeeHandler.GetAll())
		r.Post("/", employeeHandler.AddEmployee())
	})

	routes.SetupRoutes(router)

	log.Println(fmt.Sprintf("Server starting on port http://%s/api/v1", app.PortServer))

	if err := http.ListenAndServe(app.PortServer, router); err != nil {
		panic(fmt.Sprintf("Error starting server: %v", err))
	}
}
