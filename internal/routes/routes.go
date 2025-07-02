package routes

import (
	"fmt"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/handlers"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/models"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/repositories"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/services"
	"github.com/sajimenezher_meli/meli-frescos-8/pkg/loader"
)

func SetupRoutes(router *chi.Mux) {

	router.Route("/api/v1/warehouse", func(r chi.Router) {

		warehouseStorage := loader.NewJSONStorage[models.Warehouse](fmt.Sprintf("%s/%s", os.Getenv("folder_database"), "warehouse.json"))
		repository := repositories.NewWarehouseRepository(*warehouseStorage)
		service := services.NewWarehouseService(repository)
		handler := handlers.NewWarehouseHandler(service)

		r.Get("/{id}", handler.GetById)
		r.Get("/", handler.GetAll)
		r.Post("/", handler.Create)
	})
}
