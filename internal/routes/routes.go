package routes

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/handlers"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/repositories"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/services"
)

func SetupRoutes(router *chi.Mux) {
	buyerRepository, err := repositories.GetJsonBuyerRepository()
	if err != nil {
		panic(err.Error())
	}
	buyerService := services.GetBuyerService(buyerRepository)
	buyerHandler := handlers.GetBuyerHandler(buyerService)

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
			r.Get("/", buyerHandler.GetAll())
			r.Get("/{id}", buyerHandler.GetById())
			r.Delete("/{id}", buyerHandler.DeleteById())
			r.Post("/", buyerHandler.PostBuyer())
			r.Patch("/{id}", buyerHandler.PatchBuyer())
		})
	})
}
