package services

import (
	"context"

	"github.com/sajimenezher_meli/meli-frescos-8/internal/error_message"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/models"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/repositories"
)

var inboundOrdersServiceInstance InboundOrdersServiceI

// GetInboundOrdersService - Creates and returns a new instance of InboundOrdersService with required repositories using singleton pattern
// GetInboundOrdersService - Crea y retorna una nueva instancia de InboundOrdersService con los repositorios requeridos usando patrón singleton
func GetInboundOrdersService(inboundOrderRepository repositories.InboundOrderRepositoryI, employeeRepository repositories.EmployeeRepositoryI) InboundOrdersServiceI {
	if inboundOrdersServiceInstance != nil {
		return inboundOrdersServiceInstance
	}
	inboundOrdersServiceInstance = &InboundOrdersService{
		InboundOrderRepository: inboundOrderRepository,
		EmployeeRepository:     employeeRepository,
	}
	return inboundOrdersServiceInstance
}

// InboundOrdersServiceI - Interface defining the contract for inbound order service operations with business logic
// InboundOrdersServiceI - Interfaz que define el contrato para las operaciones del servicio de órdenes de entrada con lógica de negocio
type InboundOrdersServiceI interface {
	// GetAllInboundOrdersReports - Retrieves inbound order reports for all employees
	// GetAllInboundOrdersReports - Obtiene reportes de órdenes de entrada para todos los empleados
	GetAllInboundOrdersReports(ctx context.Context) ([]models.InboundOrderReport, error)

	// GetInboundOrdersReportByEmployeeId - Retrieves an inbound order report for a specific employee by ID
	// GetInboundOrdersReportByEmployeeId - Obtiene un reporte de órdenes de entrada para un empleado específico por ID
	GetInboundOrdersReportByEmployeeId(ctx context.Context, id int) (models.InboundOrderReport, error)

	// Create - Creates a new inbound order with comprehensive business validation
	// Create - Crea una nueva orden de entrada con validación integral de negocio
	Create(ctx context.Context, order models.InboundOrder) (models.InboundOrder, error)
}

// InboundOrdersService - Implementation of InboundOrdersServiceI containing business logic for inbound order operations
// InboundOrdersService - Implementación de InboundOrdersServiceI que contiene la lógica de negocio para operaciones de órdenes de entrada
type InboundOrdersService struct {
	InboundOrderRepository repositories.InboundOrderRepositoryI // Repository dependency for inbound order data access / Dependencia del repositorio para acceso a datos de órdenes de entrada
	EmployeeRepository     repositories.EmployeeRepositoryI     // Repository dependency for employee validation / Dependencia del repositorio para validación de empleados
}

// GetAllInboundOrdersReports - Delegates retrieving all inbound order reports to the repository
// GetAllInboundOrdersReports - Delega la obtención de todos los reportes de órdenes de entrada al repositorio
func (s *InboundOrdersService) GetAllInboundOrdersReports(ctx context.Context) ([]models.InboundOrderReport, error) {
	return s.InboundOrderRepository.GetAllInboundOrdersReports(ctx)
}

// GetInboundOrdersReportByEmployeeId - Delegates retrieving an inbound order report by employee ID to the repository
// GetInboundOrdersReportByEmployeeId - Delega la obtención de un reporte de órdenes de entrada por ID de empleado al repositorio
func (s *InboundOrdersService) GetInboundOrdersReportByEmployeeId(ctx context.Context, id int) (models.InboundOrderReport, error) {
	return s.InboundOrderRepository.GetInboundOrdersReportByEmployeeId(ctx, id)
}

// Create - Creates a new inbound order with comprehensive business validation
// Create - Crea una nueva orden de entrada con validación integral de negocio
func (s *InboundOrdersService) Create(ctx context.Context, order models.InboundOrder) (models.InboundOrder, error) {
	// Business validation: Validate all required fields are provided
	// Validación de negocio: Validar que todos los campos requeridos estén proporcionados
	if order.OrderNumber == "" || order.EmployeeId == 0 || order.ProductBatchId == 0 || order.WarehouseId == 0 || order.OrderDate.IsZero() {
		return models.InboundOrder{}, error_message.ErrInvalidInput
	}

	// Business rule: Order number must be unique across all inbound orders
	// Regla de negocio: El número de orden debe ser único entre todas las órdenes de entrada
	exist, err := s.InboundOrderRepository.ExistsByOrderNumber(ctx, order.OrderNumber)
	if err != nil {
		return models.InboundOrder{}, err
	}
	if exist {
		return models.InboundOrder{}, error_message.ErrAlreadyExists
	}

	// Business validation: Verify that the referenced employee exists
	// Validación de negocio: Verificar que el empleado referenciado existe
	employeeExists, err := s.EmployeeRepository.ExistEmployeeById(ctx, order.EmployeeId)
	if err != nil {
		return models.InboundOrder{}, err
	}
	if !employeeExists {
		return models.InboundOrder{}, error_message.ErrDependencyNotFound
	}

	// If all validations pass, delegate to repository for persistence
	// Si todas las validaciones pasan, delegar al repositorio para la persistencia
	newOrder, err := s.InboundOrderRepository.Create(ctx, order)
	if err != nil {
		return models.InboundOrder{}, err
	}

	return newOrder, nil
}
