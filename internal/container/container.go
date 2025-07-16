// internal/container/container.go
package container

import (
	"database/sql"
	"fmt"

	"github.com/sajimenezher_meli/meli-frescos-8/internal/handlers"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/repositories"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/services"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/validations"
)

type Container struct {
	EmployeeHandler      handlers.EmployeeHandlerI
	BuyerHandler         handlers.BuyerHandlerI
	WarehouseHandler     *handlers.WarehouseHandler
	SellerHandler        *handlers.SellerHandler
	SectionHandler       handlers.SectionHandlerI
	ProductBatchHandler  handlers.ProductBatchHandlerI
	ProductHandler       *handlers.ProductHandler
	PurchaseOrderHandler handlers.PurchaseOrderHandlerI
	ProductRecordHandler handlers.ProductRecordHandlerI
	LocalityHandler      *handlers.LocalityHandler
	CarryHandler         *handlers.CarryHandler
	InboundOrderHandler  handlers.InboundOrderHandlerI
	StorageDB            *sql.DB
}

// Strategy para manejo de errores
type ErrorHandler interface {
	Execute(tasks []Task) error
}

type Task struct {
	Name string
	Func func() error
}

// Implementación concreta que maneja los errores
type InitializationErrorHandler struct{}

func (h InitializationErrorHandler) Execute(tasks []Task) error {
	for _, task := range tasks {
		if err := task.Func(); err != nil {
			return fmt.Errorf("failed to initialize %s: %w", task.Name, err)
		}
	}
	return nil
}

func NewContainer(storeDB *sql.DB) (*Container, error) {
	container := &Container{
		StorageDB: storeDB,
	}
	errorHandler := InitializationErrorHandler{}

	tasks := []Task{
		{"employee handler", container.initializeEmployeeHandler},
		{"buyer handler", container.initializeBuyerHandler},
		{"warehouse handler", container.initializeWarehouseHandler},
		{"seller handler", container.initializeSellerHandler},
		{"section handler", container.initializeSectionHandler},
		{"product handler", container.initializeProductHandler},
		{"productBatch handler", container.initializeProductBatchHandler},
		{"purchase order handler", container.initializePurchaseOrderHandler},
		{"product record handler", container.initializeProductRecordHandler},
		{"locality handler", container.initializeLocalityHandler},
		{"carry handler", container.initializeCarryHandler},
		{"inbound order handler", container.initializeInboundOrderHandler},
	}

	if err := errorHandler.Execute(tasks); err != nil {
		return nil, err
	}

	return container, nil
}

func (c *Container) initializeEmployeeHandler() error {
	employeeRepository := repositories.GetNewEmployeeMySQLRepository(c.StorageDB)
	employeeService := services.GetEmployeeService(employeeRepository)
	employeeValidation := validations.GetEmployeeValidation()
	c.EmployeeHandler = handlers.GetEmployeeHandler(employeeService, employeeValidation)
	return nil
}

func (c *Container) initializeBuyerHandler() error {
	buyerRepository := repositories.GetNewBuyerMySQLRepository(c.StorageDB)
	buyerService := services.GetBuyerService(buyerRepository)
	c.BuyerHandler = handlers.GetBuyerHandler(buyerService)
	return nil
}

func (c *Container) initializeWarehouseHandler() error {
	repository := repositories.NewWarehouseRepository(c.StorageDB)
	service := services.NewWarehouseService(repository)
	c.WarehouseHandler = handlers.NewWarehouseHandler(service)
	return nil
}

func (c *Container) initializeSellerHandler() error {
	sellerRepo := repositories.NewSQLSellerRepository(c.StorageDB)
	sellerService := services.NewJSONSellerService(sellerRepo)
	c.SellerHandler = handlers.NewSellerHandler(sellerService)
	return nil
}
func (c *Container) initializeLocalityHandler() error {
	localityRepo := repositories.NewSQLLocalityRepository(c.StorageDB)
	localityService := services.NewSQLLocalityService(localityRepo)
	c.LocalityHandler = handlers.NewLocalityHandler(localityService)
	return nil
}

func (c *Container) initializeSectionHandler() error {
	sectionRepository := repositories.GetSectionRepository(c.StorageDB)
	sectionService := services.GetSectionService(sectionRepository)
	sectionValidation := validations.GetSectionValidation()
	c.SectionHandler = handlers.GetSectionHandler(sectionService, sectionValidation)
	return nil
}

func (c *Container) initializeProductHandler() error {
	productDB := c.StorageDB
	productRepository := repositories.NewProductRepository(productDB)
	productService := services.NewProductService(productRepository)
	c.ProductHandler = handlers.NewProductHandler(productService)
	return nil
}

func (c *Container) initializeProductRecordHandler() error {
	productRecordDb := c.StorageDB
	productRecordRepository := repositories.NewProductRecordRepository(productRecordDb)

	productRepository := repositories.NewProductRepository(productRecordDb)
	productService := services.NewProductService(productRepository)

	productRecordService := services.NewProductRecordService(productRecordRepository, productService)
	productRecordHandler := handlers.NewProductRecordHandler(productRecordService)

	c.ProductRecordHandler = productRecordHandler

	return nil
}

func (c *Container) initializeProductBatchHandler() error {
	sectionRepository := repositories.GetSectionRepository(c.StorageDB)
	sectionService := services.GetSectionService(sectionRepository)

	productBatchRepository := repositories.GetProductBatchRepository(c.StorageDB)
	productBatchService := services.GetProductBatchService(productBatchRepository)
	productBatchValidation := validations.GetProductBatchValidation()
	c.ProductBatchHandler = handlers.GetProductBatchHandler(productBatchService, sectionService, *productBatchValidation)
	return nil
}

func (c *Container) initializePurchaseOrderHandler() error {
	purchaseOrderRepository := repositories.GetNewPurchaseOrderMySQLRepository(c.StorageDB)
	buyerRepository := repositories.GetNewBuyerMySQLRepository(c.StorageDB)
	productRecordsRepository := repositories.NewProductRecordRepository(c.StorageDB)

	purchaseOrderService := services.GetPurchaseOrderService(purchaseOrderRepository, buyerRepository, productRecordsRepository)
	c.PurchaseOrderHandler = handlers.GetPurchaseOrderHandler(purchaseOrderService)
	return nil
}
func (c *Container) initializeCarryHandler() error {
	carryRepo := repositories.NewCarryRepository(c.StorageDB)
	localityRepo := repositories.NewSQLLocalityRepository(c.StorageDB)
	carryService := services.NewCarryService(carryRepo, localityRepo)
	c.CarryHandler = handlers.NewCarryHandler(carryService)
	return nil
}
func (c *Container) initializeInboundOrderHandler() error {
	inboundOrderRepository := repositories.GetNewInboundOrderMySQLRepository(c.StorageDB)
	employeeRepository := repositories.GetNewEmployeeMySQLRepository(c.StorageDB)
	inboundOrderService := services.GetInboundOrdersService(inboundOrderRepository, employeeRepository)
	c.InboundOrderHandler = handlers.GetInboundOrderHandler(inboundOrderService)
	return nil
}
