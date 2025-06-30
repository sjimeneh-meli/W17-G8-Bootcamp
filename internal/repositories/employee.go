package repositories

import (
	"fmt"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/models"
	"github.com/sajimenezher_meli/meli-frescos-8/pkg/loader"
	"os"
)

func GetEmployeeRepository() EmployeeRepositoryI {
	jsonLoader := loader.NewJSONStorage[models.Employee](fmt.Sprintf("%s/%s", os.Getenv("folder_database"), "employees.json"))
	storage, err := jsonLoader.ReadAll()
	if err != nil {
		fmt.Println("error loading employees:", err)
	}
	return &employeeRepository{
		storage: jsonLoader.MapToSlice(storage),
	}
}

type EmployeeRepositoryI interface {
	GetAll() []*models.Employee
	Add(e *models.Employee) (*models.Employee, error)
}

type employeeRepository struct {
	storage []*models.Employee
}

func (r *employeeRepository) GetAll() []*models.Employee {
	return r.storage
}
func (r *employeeRepository) Add(e *models.Employee) (*models.Employee, error) {
	// Asignar nuevo ID
	maxID := 0
	for _, emp := range r.storage {
		if emp.ID > maxID {
			maxID = emp.ID
		}
	}
	e.ID = maxID + 1

	// Agregar al slice
	r.storage = append(r.storage, e)
	return e, nil
}
