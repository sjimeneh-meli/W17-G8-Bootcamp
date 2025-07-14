package routes

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/container"
)

func SetupRoutes(c *container.Container) *chi.Mux {

	router := chi.NewRouter()

	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	router.Route("/api/v1", func(r chi.Router) {

		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"message": "API v1 is running", "status": "active"}`))
		})

		r.Route("/employee", func(rt chi.Router) {

			rt.Get("/", c.EmployeeHandler.GetAll)
			rt.Get("/{id}", c.EmployeeHandler.GetById)
			rt.Post("/", c.EmployeeHandler.Create)
			rt.Patch("/{id}", c.EmployeeHandler.Update)
			rt.Delete("/{id}", c.EmployeeHandler.DeleteById)
		})

		r.Route("/buyers", func(r chi.Router) {

			r.Get("/", c.BuyerHandler.GetAll())
			r.Get("/{id}", c.BuyerHandler.GetById())
			r.Delete("/{id}", c.BuyerHandler.DeleteById())
			r.Post("/", c.BuyerHandler.PostBuyer())
			r.Patch("/{id}", c.BuyerHandler.PatchBuyer())
			r.Get("/reportPurchaseOrders", c.PurchaseOrderHandler.GetPurchaseOrdersReport())
		})

		r.Route("/warehouse", func(r chi.Router) {

			r.Get("/{id}", c.WarehouseHandler.GetById)
			r.Get("/", c.WarehouseHandler.GetAll)
			r.Post("/", c.WarehouseHandler.Create)
			r.Patch("/{id}", c.WarehouseHandler.Update)
			r.Delete("/{id}", c.WarehouseHandler.Delete)
		})

		r.Route("/sellers", func(r chi.Router) {

			r.Get("/", c.SellerHandler.GetAll)
			r.Get("/{id}", c.SellerHandler.GetById)
			r.Post("/", c.SellerHandler.Save)
			r.Patch("/{id}", c.SellerHandler.Update)
			r.Delete("/{id}", c.SellerHandler.Delete)
		})

		r.Route("/sections", func(rt chi.Router) {
			rt.Get("/", c.SectionHandler.GetAll)
			rt.Get("/{id}", c.SectionHandler.GetByID)
			rt.Post("/", c.SectionHandler.Create)
			rt.Patch("/{id}", c.SectionHandler.Update)
			rt.Delete("/{id}", c.SectionHandler.DeleteByID)
		})

		r.Route("/products", func(r chi.Router) {

			r.Get("/", c.ProductHandler.GetAll)
			r.Get("/{id}", c.ProductHandler.GetById)
			r.Post("/", c.ProductHandler.Save)
			r.Patch("/{id}", c.ProductHandler.Update)
			r.Delete("/{id}", c.ProductHandler.DeleteById)
		})

		r.Route("/purchaseOrders", func(r chi.Router) {
			r.Get("/", c.PurchaseOrderHandler.GetAll())
			r.Post("/", c.PurchaseOrderHandler.PostPurchaseOrder())
		})
	})

	return router
}
