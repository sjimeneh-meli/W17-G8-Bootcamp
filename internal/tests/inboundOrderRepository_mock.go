package tests

import (
	"context"

	"github.com/sajimenezher_meli/meli-frescos-8/internal/models"
	"github.com/stretchr/testify/mock"
)

// InboundOrderRepositoryMock - Mock del repositorio InboundOrder para tests del service
// Permite simular respuestas de base de datos y testear lógica de negocio aisladamente
type InboundOrderRepositoryMock struct {
	mock.Mock
}

// GetNewInboundOrderRepositoryMock - Factory para crear nueva instancia del mock
func GetNewInboundOrderRepositoryMock() *InboundOrderRepositoryMock {
	return &InboundOrderRepositoryMock{}
}

// GetAllInboundOrdersReports - Mock para obtener reportes de todas las órdenes de entrada
func (r *InboundOrderRepositoryMock) GetAllInboundOrdersReports(ctx context.Context) ([]models.InboundOrderReport, error) {
	args := r.Called(ctx)
	return args.Get(0).([]models.InboundOrderReport), args.Error(1)
}

// GetInboundOrdersReportByEmployeeId - Mock para obtener reporte de órdenes por ID de empleado
func (r *InboundOrderRepositoryMock) GetInboundOrdersReportByEmployeeId(ctx context.Context, employeeId int) (models.InboundOrderReport, error) {
	args := r.Called(ctx, employeeId)
	return args.Get(0).(models.InboundOrderReport), args.Error(1)
}

// Create - Mock para crear nueva orden de entrada
func (r *InboundOrderRepositoryMock) Create(ctx context.Context, inbound models.InboundOrder) (models.InboundOrder, error) {
	args := r.Called(ctx, inbound)
	return args.Get(0).(models.InboundOrder), args.Error(1)
}

// ExistsByOrderNumber - Mock para verificar si existe orden por número de orden
func (r *InboundOrderRepositoryMock) ExistsByOrderNumber(ctx context.Context, orderNumber string) (bool, error) {
	args := r.Called(ctx, orderNumber)
	return args.Bool(0), args.Error(1)
}
