package routes

import (
	"fmt"

	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/handlers"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/models"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/repositories"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/seeders"
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

	router.Route("/api/v1", func(r chi.Router) {
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"message": "API v1 is running", "status": "active"}`))
		})

		r.Route("/employee", func(rt chi.Router) {
			employeeLoader := loader.NewJSONStorage[models.Employee](fmt.Sprintf("%s/%s", os.Getenv("folder_database"), "employees.json"))
			employeeRepository, employeeLoaderErr := repositories.GetEmployeeRepository(employeeLoader)
			if employeeLoaderErr != nil {
				panic(employeeLoaderErr.Error())
			}
			employeeService := services.GetEmployeeService(employeeRepository)
			employeeValidation := validations.GetEmployeeValidation()
			employeeHandler := handlers.GetEmployeeHandler(employeeService, employeeValidation)

			rt.Get("/", employeeHandler.GetAll)
			rt.Get("/{id}", employeeHandler.GetById)
			rt.Post("/", employeeHandler.Create)
			rt.Patch("/{id}", employeeHandler.Update)
			rt.Delete("/{id}", employeeHandler.DeleteById)
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
			r.Patch("/{id}", handler.Update)
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

		r.Route("/sections", func(rt chi.Router) {
			rt.Get("/", sectionHandler.GetAll)
			rt.Get("/{id}", sectionHandler.GetByID)
			rt.Post("/", sectionHandler.Create)
			rt.Patch("/{id}", sectionHandler.Update)
			rt.Delete("/{id}", sectionHandler.DeleteByID)
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
