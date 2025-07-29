package tests

import (
	"context"

	"github.com/sajimenezher_meli/meli-frescos-8/internal/models"
	"github.com/stretchr/testify/mock"
)

// EmployeeServiceMock - Mock del servicio Employee para tests del handler
// Permite controlar respuestas del servicio y testear handler de forma aislada
type EmployeeServiceMock struct {
	mock.Mock
}

// GetNewEmployeeServiceMock - Factory para crear nueva instancia del mock
func GetNewEmployeeServiceMock() *EmployeeServiceMock {
	return &EmployeeServiceMock{}
}

// GetAll - Mock para obtener todos los empleados
func (r *EmployeeServiceMock) GetAll(ctx context.Context) (map[int]models.Employee, error) {
	args := r.Called(ctx)
	return args.Get(0).(map[int]models.Employee), args.Error(1)
}

// GetById - Mock para obtener empleado por ID
func (r *EmployeeServiceMock) GetById(ctx context.Context, id int) (models.Employee, error) {
	args := r.Called(ctx, id)
	return args.Get(0).(models.Employee), args.Error(1)
}

// DeleteById - Mock para eliminar empleado por ID
func (r *EmployeeServiceMock) DeleteById(ctx context.Context, id int) error {
	args := r.Called(ctx, id)
	return args.Error(0)
}

// Create - Mock para crear nuevo empleado
func (r *EmployeeServiceMock) Create(ctx context.Context, employee models.Employee) (models.Employee, error) {
	args := r.Called(ctx, employee)
	return args.Get(0).(models.Employee), args.Error(1)
}

// Update - Mock para actualizar empleado existente
func (r *EmployeeServiceMock) Update(ctx context.Context, employeeId int, employee models.Employee) (models.Employee, error) {
	args := r.Called(ctx, employeeId, employee)
	return args.Get(0).(models.Employee), args.Error(1)
}
