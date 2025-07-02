package routes

import (
	"fmt"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/handlers"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/models"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/repositories"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/services"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/validations"
	"github.com/sajimenezher_meli/meli-frescos-8/pkg/loader"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
)

func SetupRoutes(router *chi.Mux) {
	employeeLoader := loader.NewJSONStorage[models.Employee](fmt.Sprintf("%s/%s", os.Getenv("folder_database"), "employee.json"))
	employeeRepository, employeeLoaderErr := repositories.GetEmployeeRepository(employeeLoader)
	if employeeLoaderErr != nil {
		panic(employeeLoaderErr.Error())
	}
	employeeService := services.GetEmployeeService(employeeRepository)
	employeeValidation := validations.GetEmployeeValidation()
	employeeHandler := handlers.GetEmployeeHandler(employeeService, employeeValidation)

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

		r.Route("/employee", func(rt chi.Router) {
			rt.Get("/", employeeHandler.GetAll())
			rt.Get("/{id}", employeeHandler.GetById())
			rt.Post("/", employeeHandler.Create())
			rt.Patch("/{id}", employeeHandler.PatchEmployee())
			rt.Delete("/{id}", employeeHandler.DeleteById())
		})
	})
}
