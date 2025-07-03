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
	"github.com/sajimenezher_meli/meli-frescos-8/pkg/loader"
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

		r.Route("/buyers", func(r chi.Router) {
			buyerLoader := loader.NewJSONStorage[models.Buyer](fmt.Sprintf("%s/%s", os.Getenv("folder_database"), "buyers.json"))
			buyerRepository, err := repositories.GetNewBuyerRepository(buyerLoader)
			if err != nil {
				panic(err.Error())
			}
			buyerService := services.GetBuyerService(buyerRepository)
			buyerHandler := handlers.GetBuyerHandler(buyerService)

			r.Get("/", buyerHandler.GetAll())
			r.Get("/{id}", buyerHandler.GetById())
			r.Delete("/{id}", buyerHandler.DeleteById())
			r.Post("/", buyerHandler.PostBuyer())
			r.Patch("/{id}", buyerHandler.PatchBuyer())
		})

		r.Route("/warehouse", func(r chi.Router) {
			warehouseStorage := loader.NewJSONStorage[models.Warehouse](fmt.Sprintf("%s/%s", os.Getenv("folder_database"), "warehouse.json"))
			repository := repositories.NewWarehouseRepository(*warehouseStorage)
			service := services.NewWarehouseService(repository)
			handler := handlers.NewWarehouseHandler(service)
			r.Get("/{id}", handler.GetById)
			r.Get("/", handler.GetAll)
			r.Post("/", handler.Create)
			r.Put("/{id}", handler.Update)
			r.Delete("/{id}", handler.Delete)
		})
		r.Route("/sellers", func(r chi.Router) {
			sellerStorage := loader.NewJSONStorage[models.Seller](fmt.Sprintf("%s/%s", "docs/database", "sellers.json"))
			sellerRepo := repositories.NewJSONSellerRepository(sellerStorage)
			sellerService := services.NewJSONSellerService(sellerRepo)
			sellerHandler := handlers.NewSellerHandler(sellerService)

			r.Get("/", sellerHandler.GetAll)
			r.Get("/{id}", sellerHandler.GetById)
			r.Post("/", sellerHandler.Save)
			r.Patch("/{id}", sellerHandler.Update)
			r.Delete("/{id}", sellerHandler.Delete)
		})
	})
}
