package routes

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/handlers"
)

func SetupRoutes(router *chi.Mux) {
	sectionHandler := handlers.GetSectionHandler()

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
