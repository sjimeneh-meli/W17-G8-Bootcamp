package tests

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/sajimenezher_meli/meli-frescos-8/internal/config"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/handlers"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/repositories"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/services"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/validations"
	"github.com/sajimenezher_meli/meli-frescos-8/pkg/database"
	"github.com/stretchr/testify/require"
)

func TestCreateOkSection(t *testing.T) {
	db := database.InitDB(&config.Config{
		Database: config.Database{
			DBUser:     "root",
			DBPassword: "test",
			DBHost:     "localhost",
			DBPort:     "3306",
			DBName:     "productos_frescos",
		},
	})
	defer db.Close()

	rp := repositories.GetSectionRepository(db)
	warehouseRp := repositories.NewWarehouseRepository(db)
	srv := services.GetSectionService(rp)
	warehouseSrv := services.NewWarehouseService(warehouseRp)
	vld := validations.GetSectionValidation()
	hdCreateFunc := handlers.GetSectionHandler(srv, warehouseSrv, vld).Create

	request := httptest.NewRequest(http.MethodPost, "/sections", strings.NewReader(`
		{
			"section_number": "G-01",
			"current_capacity": 2,
			"current_temperature": 3.43,
			"maximum_capacity": 2,
			"minimum_capacity": 2,
			"minimum_temperature": 2,
			"product_type_id": 1,
			"warehouse_id": 1
		}
	`))
	response := httptest.NewRecorder()
	hdCreateFunc(response, request)

	require.Equal(t, http.StatusCreated, response.Code)
	require.JSONEq(t, `{
		"data": {
			"id": 26,  
			"section_number": "G-01",
			"current_capacity": 2,
			"current_temperature": 3.43,
			"maximum_capacity": 2,
			"minimum_capacity": 2,
			"minimum_temperature": 2,
			"product_type_id": 2,
			"warehouse_id": 2
		}
	}`, response.Body.String())
}
