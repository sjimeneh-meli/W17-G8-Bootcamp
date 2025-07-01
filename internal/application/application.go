package application

import (
	"fmt"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/handlers"
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
	sectionHandler := handlers.GetSectionHandler()
	router := chi.NewRouter()

	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	router.Route("/sections", func(rt chi.Router) {
		rt.Get("/", sectionHandler.GetAll())
		rt.Get("/{id}", sectionHandler.GetByID())
		rt.Post("/", sectionHandler.Create())
		rt.Delete("/{id}", sectionHandler.DeleteByID())
	})

	routes.SetupRoutes(router)

	fmt.Printf("Server starting on port %s\n", app.PortServer)
	if err := http.ListenAndServe(app.PortServer, router); err != nil {
		panic(fmt.Sprintf("Error starting server: %v", err))
	}
}
