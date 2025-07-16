package services

import (
	"context"

	"github.com/sajimenezher_meli/meli-frescos-8/internal/error_message"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/models"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/repositories"
)

var inboundOrdersServiceInstance InboundOrdersServiceI

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

type InboundOrdersServiceI interface {
	GetAllInboundOrdersReports(ctx context.Context) ([]models.InboundOrderReport, error)
	GetInboundOrdersReportByEmployeeId(ctx context.Context, id int) (models.InboundOrderReport, error)

	Create(ctx context.Context, order models.InboundOrder) (models.InboundOrder, error)
}

type InboundOrdersService struct {
	InboundOrderRepository repositories.InboundOrderRepositoryI
	EmployeeRepository     repositories.EmployeeRepositoryI
}

func (s *InboundOrdersService) GetAllInboundOrdersReports(ctx context.Context) ([]models.InboundOrderReport, error) {
	return s.InboundOrderRepository.GetAllInboundOrdersReports(ctx)
}

func (s *InboundOrdersService) GetInboundOrdersReportByEmployeeId(ctx context.Context, id int) (models.InboundOrderReport, error) {
	return s.InboundOrderRepository.GetInboundOrdersReportByEmployeeId(ctx, id)
}
func (s *InboundOrdersService) Create(ctx context.Context, order models.InboundOrder) (models.InboundOrder, error) {
	if order.OrderNumber == "" || order.EmployeeId == 0 || order.ProductBatchId == 0 || order.WarehouseId == 0 || order.OrderDate.IsZero() {
		return models.InboundOrder{}, error_message.ErrInvalidInput
	}

	exist, err := s.InboundOrderRepository.ExistsByOrderNumber(ctx, order.OrderNumber)
	if err != nil {
		return models.InboundOrder{}, err
	}
	if exist {
		return models.InboundOrder{}, error_message.ErrAlreadyExists
	}

	employeeExists, err := s.EmployeeRepository.ExistEmployeeById(ctx, order.EmployeeId)
	if err != nil {
		return models.InboundOrder{}, err
	}
	if !employeeExists {
		return models.InboundOrder{}, error_message.ErrDependencyNotFound
	}

	newOrder, err := s.InboundOrderRepository.Create(ctx, order)
	if err != nil {
		return models.InboundOrder{}, err
	}

	return newOrder, nil
}
