package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/sajimenezher_meli/meli-frescos-8/internal/error_message"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/models"
)

var purchaseOrderRepositoryInstance PurchaseOrderRepositoryI

// GetNewPurchaseOrderMySQLRepository - Creates and returns a new instance of MySqlPurchaseOrderRepository using singleton pattern
// GetNewPurchaseOrderMySQLRepository - Crea y retorna una nueva instancia de MySqlPurchaseOrderRepository usando patrón singleton
func GetNewPurchaseOrderMySQLRepository(db *sql.DB) PurchaseOrderRepositoryI {
	if purchaseOrderRepositoryInstance != nil {
		return purchaseOrderRepositoryInstance
	}

	purchaseOrderRepositoryInstance = &MySqlPurchaseOrderRepository{
		db: db,
	}
	return purchaseOrderRepositoryInstance
}

// PurchaseOrderRepositoryI - Interface defining the contract for purchase order repository operations
// PurchaseOrderRepositoryI - Interfaz que define el contrato para las operaciones del repositorio de órdenes de compra
type PurchaseOrderRepositoryI interface {
	// GetAll - Retrieves all purchase orders from the database and returns them as a map with order ID as key
	// GetAll - Obtiene todas las órdenes de compra de la base de datos y las retorna como un mapa con el ID de la orden como clave
	GetAll(ctx context.Context) (map[int]models.PurchaseOrder, error)

	// Create - Inserts a new purchase order into the database and returns the created order with its generated ID
	// Create - Inserta una nueva orden de compra en la base de datos y retorna la orden creada con su ID generado
	Create(ctx context.Context, order models.PurchaseOrder) (models.PurchaseOrder, error)

	// ExistPurchaseOrderByOrderNumber - Checks if a purchase order with the given order number already exists in the database
	// ExistPurchaseOrderByOrderNumber - Verifica si una orden de compra con el número de orden dado ya existe en la base de datos
	ExistPurchaseOrderByOrderNumber(ctx context.Context, orderNumber string) (bool, error)

	// GetPurchaseOrdersReportByBuyerId - Retrieves a purchase order report for a specific buyer ID showing buyer info and order count
	// GetPurchaseOrdersReportByBuyerId - Obtiene un reporte de órdenes de compra para un ID de comprador específico mostrando información del comprador y conteo de órdenes
	GetPurchaseOrdersReportByBuyerId(ctx context.Context, buyerId int) (models.PurchaseOrderReport, error)

	// GetAllPurchaseOrdersReports - Retrieves purchase order reports for all buyers showing buyer info and their order counts
	// GetAllPurchaseOrdersReports - Obtiene reportes de órdenes de compra para todos los compradores mostrando información del comprador y sus conteos de órdenes
	GetAllPurchaseOrdersReports(ctx context.Context) ([]models.PurchaseOrderReport, error)
}

// MySqlPurchaseOrderRepository - MySQL implementation of the PurchaseOrderRepositoryI interface
// MySqlPurchaseOrderRepository - Implementación MySQL de la interfaz PurchaseOrderRepositoryI
type MySqlPurchaseOrderRepository struct {
	db *sql.DB // Database connection / Conexión a la base de datos
}

// GetAll - Retrieves all purchase orders from the database and returns them as a map with order ID as key
// GetAll - Obtiene todas las órdenes de compra de la base de datos y las retorna como un mapa con el ID de la orden como clave
func (r *MySqlPurchaseOrderRepository) GetAll(ctx context.Context) (map[int]models.PurchaseOrder, error) {
	orders := make(map[int]models.PurchaseOrder)

	// SQL query to select all purchase order fields / Consulta SQL para seleccionar todos los campos de la orden de compra
	query := "select id, order_number, order_date, tracking_code, buyer_id, product_record_id from purchase_orders"

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return orders, fmt.Errorf("%w - %s", error_message.ErrInternalServerError, err.Error())
	}
	defer rows.Close()

	// Create temporary map to store orders before returning / Crear mapa temporal para almacenar órdenes antes de retornar
	tempOrdersMap := make(map[int]models.PurchaseOrder)

	// Iterate through all rows and map each order to the result map / Itera a través de todas las filas y mapea cada orden al mapa de resultados
	for rows.Next() {
		order := models.PurchaseOrder{}
		err := rows.Scan(&order.Id, &order.OrderNumber, &order.OrderDate, &order.TrackingCode, &order.BuyerId, &order.ProductRecordId)
		if err != nil {
			return orders, fmt.Errorf("%w - %s", error_message.ErrInternalServerError, err.Error())
		}

		tempOrdersMap[order.Id] = order
	}

	orders = tempOrdersMap
	return orders, nil
}

// Create - Inserts a new purchase order into the database and returns the created order with its generated ID
// Create - Inserta una nueva orden de compra en la base de datos y retorna la orden creada con su ID generado
func (r *MySqlPurchaseOrderRepository) Create(ctx context.Context, order models.PurchaseOrder) (models.PurchaseOrder, error) {
	// SQL query to insert new purchase order / Consulta SQL para insertar nueva orden de compra
	query := `insert into purchase_orders (order_number, order_date, tracking_code, buyer_id, product_record_id)
	values (?, ?, ?, ?, ?)`

	result, err := r.db.ExecContext(ctx, query, order.OrderNumber, order.OrderDate, order.TrackingCode, order.BuyerId, order.ProductRecordId)
	if err != nil {
		return models.PurchaseOrder{}, fmt.Errorf("%w - %s", error_message.ErrInternalServerError, err.Error())
	}

	// Get the auto-generated ID from the database / Obtiene el ID autogenerado de la base de datos
	lastId, err := result.LastInsertId()
	if err != nil {
		return models.PurchaseOrder{}, fmt.Errorf("%w - %s", error_message.ErrInternalServerError, err.Error())
	}

	order.Id = int(lastId)
	return order, nil
}

// GetPurchaseOrdersReportByBuyerId - Retrieves a purchase order report for a specific buyer ID showing buyer info and order count
// GetPurchaseOrdersReportByBuyerId - Obtiene un reporte de órdenes de compra para un ID de comprador específico mostrando información del comprador y conteo de órdenes
func (r *MySqlPurchaseOrderRepository) GetPurchaseOrdersReportByBuyerId(ctx context.Context, buyerId int) (models.PurchaseOrderReport, error) {
	report := models.PurchaseOrderReport{}

	// Complex SQL query using INNER JOIN to get buyer info and count their purchase orders / Consulta SQL compleja usando INNER JOIN para obtener información del comprador y contar sus órdenes de compra
	query := `select b.id, b.id_card_number, b.first_name, b.last_name, count(po.id) as "purchase_orders_count"
from productos_frescos.buyers b
inner join productos_frescos.purchase_orders po on po.buyer_id = b.id
where b.id = ?
group by b.id`
	row := r.db.QueryRowContext(ctx, query, buyerId)

	err := row.Scan(&report.Id, &report.IdCardNumber, &report.FirstName, &report.LastName, &report.PurchaseOrderCount)
	if err != nil {
		// Handle case when no buyer is found / Maneja el caso cuando no se encuentra ningún comprador
		if errors.Is(err, sql.ErrNoRows) {
			return models.PurchaseOrderReport{}, fmt.Errorf("%w. %s %d %s", error_message.ErrNotFound, "Buyer with Id", buyerId, "doesn't exists.")
		}
		return models.PurchaseOrderReport{}, fmt.Errorf("%w - %s", error_message.ErrInternalServerError, err.Error())
	}

	return report, nil
}

// GetAllPurchaseOrdersReports - Retrieves purchase order reports for all buyers showing buyer info and their order counts
// GetAllPurchaseOrdersReports - Obtiene reportes de órdenes de compra para todos los compradores mostrando información del comprador y sus conteos de órdenes
func (r *MySqlPurchaseOrderRepository) GetAllPurchaseOrdersReports(ctx context.Context) ([]models.PurchaseOrderReport, error) {
	reports := []models.PurchaseOrderReport{}

	// Complex SQL query using INNER JOIN to get all buyers info and count their purchase orders / Consulta SQL compleja usando INNER JOIN para obtener información de todos los compradores y contar sus órdenes de compra
	query := `select b.id, b.id_card_number, b.first_name, b.last_name, count(po.id) as "purchase_orders_count"
from productos_frescos.buyers b
inner join productos_frescos.purchase_orders po on po.buyer_id = b.id
group by b.id
order by b.id`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return []models.PurchaseOrderReport{}, fmt.Errorf("%w - %s", error_message.ErrInternalServerError, err.Error())
	}
	defer rows.Close()

	// Iterate through all rows and scan each report into the results slice / Itera a través de todas las filas y escanea cada reporte en el slice de resultados
	for rows.Next() {
		report := models.PurchaseOrderReport{}
		err := rows.Scan(&report.Id, &report.IdCardNumber, &report.FirstName, &report.LastName, &report.PurchaseOrderCount)
		if err != nil {
			return []models.PurchaseOrderReport{}, fmt.Errorf("%w - %s", error_message.ErrInternalServerError, err.Error())
		}

		reports = append(reports, report)
	}
	return reports, nil
}

// ExistPurchaseOrderByOrderNumber - Checks if a purchase order with the given order number already exists in the database
// ExistPurchaseOrderByOrderNumber - Verifica si una orden de compra con el número de orden dado ya existe en la base de datos
func (r *MySqlPurchaseOrderRepository) ExistPurchaseOrderByOrderNumber(ctx context.Context, orderNumber string) (bool, error) {
	// Simple query to check purchase order existence using LIMIT 1 for efficiency / Consulta simple para verificar existencia de la orden de compra usando LIMIT 1 por eficiencia
	query := "SELECT 1 FROM purchase_orders WHERE order_number = ? LIMIT 1 "

	var exists int64
	err := r.db.QueryRowContext(ctx, query, orderNumber).Scan(&exists)

	if err != nil {
		// If no rows found, order number doesn't exist (not an error) / Si no se encuentran filas, el número de orden no existe (no es un error)
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, fmt.Errorf("error al verificar la existencia del order number: %w", err)
	}
	return true, nil
}
