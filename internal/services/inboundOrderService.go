package services

import (
	"context"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/models"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/repositories"
)

type InboundOrdersServiceI interface {
	GetAllInboundOrdersReports(ctx context.Context) ([]models.InboundOrderReport, error)
	GetInboundOrdersReportByEmployeeId(ctx context.Context, id int) (models.InboundOrderReport, error)
}

type InboundOrdersService struct {
	InboundOrderRepository repositories.InboundOrderRepositoryI
}

func GetInboundOrdersService(inboundOrderRepository repositories.InboundOrderRepositoryI, employeeRepository repositories.EmployeeRepositoryI) InboundOrdersServiceI {
	return &InboundOrdersService{
		InboundOrderRepository: inboundOrderRepository,
	}
}

func (s *InboundOrdersService) GetAllInboundOrdersReports(ctx context.Context) ([]models.InboundOrderReport, error) {
	return s.InboundOrderRepository.GetAllInboundOrdersReports(ctx)
}

func (s *InboundOrdersService) GetInboundOrdersReportByEmployeeId(ctx context.Context, id int) (models.InboundOrderReport, error) {
	return s.InboundOrderRepository.GetInboundOrdersReportByEmployeeId(ctx, id)
}
