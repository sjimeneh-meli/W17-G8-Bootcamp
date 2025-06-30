package repositories

import (
	"fmt"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/models"
	"github.com/sajimenezher_meli/meli-frescos-8/pkg/loader"
	"os"
)

type EmployeeRepositoryI interface {
	GetAll() []*models.Employee
}

type employeeRepository struct {
	storage []*models.Employee
}

func (r *employeeRepository) GetAll() []*models.Employee {
	return r.storage
}

func GetEmployeeRepository() EmployeeRepositoryI {
	jsonLoader := loader.NewJSONStorage[models.Employee](fmt.Sprintf("%s/%s", os.Getenv("folder_database"), "employees.json"))

	// Leer todos los datos
	storage, err := jsonLoader.ReadAll()
	if err != nil {
		fmt.Println("error loading employees:", err)
	}

	return &employeeRepository{
		storage: jsonLoader.MapToSlice(storage),
	}
}
