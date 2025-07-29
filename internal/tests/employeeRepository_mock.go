package tests

import (
	"context"

	"github.com/sajimenezher_meli/meli-frescos-8/internal/models"
	"github.com/stretchr/testify/mock"
)

// EmployeeRepositoryMock - Mock del repositorio Employee para tests del service
// Permite simular respuestas de base de datos y testear lógica de negocio aisladamente
type EmployeeRepositoryMock struct {
	mock.Mock
}

// GetNewEmployeeRepositoryMock - Factory para crear nueva instancia del mock
func GetNewEmployeeRepositoryMock() *EmployeeRepositoryMock {
	return &EmployeeRepositoryMock{}
}

// GetAll - Mock para obtener todos los empleados de base de datos
func (r *EmployeeRepositoryMock) GetAll(ctx context.Context) (map[int]models.Employee, error) {
	args := r.Called(ctx)
	return args.Get(0).(map[int]models.Employee), args.Error(1)
}

// GetById - Mock para obtener empleado por ID de base de datos
func (r *EmployeeRepositoryMock) GetById(ctx context.Context, id int) (models.Employee, error) {
	args := r.Called(ctx, id)
	return args.Get(0).(models.Employee), args.Error(1)
}

// DeleteById - Mock para eliminar empleado por ID de base de datos
func (r *EmployeeRepositoryMock) DeleteById(ctx context.Context, id int) error {
	args := r.Called(ctx, id)
	return args.Error(0)
}

// Create - Mock para insertar nuevo empleado en base de datos
func (r *EmployeeRepositoryMock) Create(ctx context.Context, employee models.Employee) (models.Employee, error) {
	args := r.Called(ctx, employee)
	return args.Get(0).(models.Employee), args.Error(1)
}

// Update - Mock para actualizar empleado en base de datos
func (r *EmployeeRepositoryMock) Update(ctx context.Context, employeeId int, employee models.Employee) (models.Employee, error) {
	args := r.Called(ctx, employeeId, employee)
	return args.Get(0).(models.Employee), args.Error(1)
}

// GetCardNumberIds - Mock para obtener lista de card numbers existentes
func (r *EmployeeRepositoryMock) GetCardNumberIds() ([]string, error) {
	args := r.Called()
	return args.Get(0).([]string), args.Error(1)
}

// ExistEmployeeById - Mock para verificar existencia de empleado por ID
func (r *EmployeeRepositoryMock) ExistEmployeeById(ctx context.Context, employeeId int) (bool, error) {
	args := r.Called(ctx, employeeId)
	return args.Bool(0), args.Error(1)
}
