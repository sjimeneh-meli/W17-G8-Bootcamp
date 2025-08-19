package tests

import (
	"context"

	"github.com/sajimenezher_meli/meli-frescos-8/internal/models"
	"github.com/stretchr/testify/mock"
)

// InboundOrderServiceMock - Mock del servicio InboundOrder para tests del handler
// Permite controlar respuestas del servicio y testear handler de forma aislada
type InboundOrderServiceMock struct {
	mock.Mock
}

// GetNewInboundOrderServiceMock - Factory para crear nueva instancia del mock
func GetNewInboundOrderServiceMock() *InboundOrderServiceMock {
	return &InboundOrderServiceMock{}
}

// GetAllInboundOrdersReports - Mock para obtener reportes de todas las órdenes de entrada
func (s *InboundOrderServiceMock) GetAllInboundOrdersReports(ctx context.Context) ([]models.InboundOrderReport, error) {
	args := s.Called(ctx)
	return args.Get(0).([]models.InboundOrderReport), args.Error(1)
}

// GetInboundOrdersReportByEmployeeId - Mock para obtener reporte de órdenes por ID de empleado
func (s *InboundOrderServiceMock) GetInboundOrdersReportByEmployeeId(ctx context.Context, id int) (models.InboundOrderReport, error) {
	args := s.Called(ctx, id)
	return args.Get(0).(models.InboundOrderReport), args.Error(1)
}

// Create - Mock para crear nueva orden de entrada
func (s *InboundOrderServiceMock) Create(ctx context.Context, order models.InboundOrder) (models.InboundOrder, error) {
	args := s.Called(ctx, order)
	return args.Get(0).(models.InboundOrder), args.Error(1)
}
