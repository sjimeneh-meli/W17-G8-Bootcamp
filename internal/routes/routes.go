package routes

import (
	"fmt"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/handlers"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/models"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/repositories"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/services"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/validations"
	"github.com/sajimenezher_meli/meli-frescos-8/pkg/loader"
)

func SetupRoutes(router *chi.Mux) {
	sectionLoader := loader.NewJSONStorage[models.Section](fmt.Sprintf("%s/%s", os.Getenv("folder_database"), "sections.json"))
	sectionRepository, sectionLoaderErr := repositories.GetSectionRepository(sectionLoader)
	if sectionLoaderErr != nil {
		panic(sectionLoaderErr.Error())
	}
	sectionService := services.GetSectionService(sectionRepository)
	sectionValidation := validations.GetSectionValidation()
	sectionHandler := handlers.GetSectionHandler(sectionService, sectionValidation)

	router.Get("/hello-world", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"message": "Hello World from Chi Router!"}`))
	})

	router.Route("/api/v1", func(r chi.Router) {
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"message": "API v1 is running", "status": "active"}`))
		})

		r.Route("/sections", func(rt chi.Router) {
			rt.Get("/", sectionHandler.GetAll())
			rt.Get("/{id}", sectionHandler.GetByID())
			rt.Post("/", sectionHandler.Create())
			rt.Patch("/{id}", sectionHandler.Update())
			rt.Delete("/{id}", sectionHandler.DeleteByID())
		})
	})
}
