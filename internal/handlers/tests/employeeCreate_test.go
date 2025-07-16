package tests

/*
import (
	"bytes"
	"encoding/json"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/handlers"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/handlers/requests"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/models"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/validations"
	"net/http"
	"net/http/httptest"
	"testing"
)

type mockEmployeeService struct{}

func (m *mockEmployeeService) Create(e *models.Employee) error {
	e.Id = 99
	return nil
}
func (m *mockEmployeeService) GetAll() []*models.Employee                        { return nil }
func (m *mockEmployeeService) GetById(id int) (*models.Employee, error)          { return nil, nil }
func (m *mockEmployeeService) DeleteById(id int) error                           { return nil }
func (m *mockEmployeeService) ExistsWhCardNumber(id int, cardNumber string) bool { return false }

func TestCreateEmployee_Success(t *testing.T) {
	handler := handlers.GetEmployeeHandler(&mockEmployeeService{}, validations.GetEmployeeValidation())

	requestPayload := requests.EmployeeRequest{
		CardNumberID: "EMP123",
		FirstName:    "Juan",
		LastName:     "Pérez",
		WarehouseID:  1,
	}

	body, _ := json.Marshal(requestPayload)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/employees", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()

	handler.Create().ServeHTTP(rr, req)

	if rr.Code != http.StatusCreated {
		t.Errorf("expected status 201 Created, got %d", rr.Code)
	}
}

func TestCreateEmployee_InvalidJSON(t *testing.T) {
	handler := handlers.GetEmployeeHandler(&mockEmployeeService{}, validations.GetEmployeeValidation())

	req := httptest.NewRequest(http.MethodPost, "/api/v1/employees", bytes.NewBuffer([]byte(`invalid-json`)))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler.Create().ServeHTTP(rr, req)

	if rr.Code != http.StatusExpectationFailed {
		t.Errorf("expected 417 ExpectationFailed, got %d", rr.Code)
	}
}
*/
