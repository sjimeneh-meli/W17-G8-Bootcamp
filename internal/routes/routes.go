package routes

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/handlers"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/models"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/repositories"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/services"
	"github.com/sajimenezher_meli/meli-frescos-8/pkg/loader"
	"net/http"
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
	})

	router.Route("/sellers", func(r chi.Router) {

		// sellers
		sellerStorage := loader.NewJSONStorage[models.Seller](fmt.Sprintf("%s/%s", "docs/database", "sellers.json"))
		sellerRepo := repositories.NewJSONSellerRepository(sellerStorage)
		sellerService := services.NewJSONSellerService(sellerRepo)
		sellerHandler := handlers.NewSellerHandler(sellerService)

		r.Get("/", sellerHandler.GetAll)
	})
}
