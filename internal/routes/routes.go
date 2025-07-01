package routes

import (
	"fmt"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/handlers"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/models"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/repositories"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/seeders"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/services"
	"github.com/sajimenezher_meli/meli-frescos-8/pkg/loader"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
)

func SetupRoutes(router *chi.Mux) {
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

		r.Route("/products", func(r chi.Router) {
			storage := loader.NewJSONStorage[models.Product](fmt.Sprintf("%s/%s", os.Getenv("folder_database"), "products.json"))
			repository := repositories.NewProductRepository(*storage)
			service := services.NewProductService(repository)
			productHandler := handlers.NewProductHandler(service)

			ps := seeders.NewSeeder(service)
			ps.LoadAllData()

			r.Get("/", productHandler.GetAll)
			r.Get("/{id}", productHandler.GetById)
			r.Post("/", productHandler.Save)
			r.Patch("/{id}", productHandler.Update)
			r.Delete("/{id}", productHandler.DeleteById)
		})
	})
}
