package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/sajimenezher_meli/meli-frescos-8/internal/error_message"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/models"
)

var buyerRepositoryInstance BuyerRepositoryI

// GetNewBuyerMySQLRepository - Creates and returns a new instance of MySqlBuyerRepository using singleton pattern
// GetNewBuyerMySQLRepository - Crea y retorna una nueva instancia de MySqlBuyerRepository usando patrón singleton
func GetNewBuyerMySQLRepository(db *sql.DB) BuyerRepositoryI {
	if buyerRepositoryInstance != nil {
		return buyerRepositoryInstance
	}

	buyerRepositoryInstance = &MySqlBuyerRepository{
		Db: db,
	}
	return buyerRepositoryInstance
}

// BuyerRepositoryI - Interface defining the contract for buyer repository operations
// BuyerRepositoryI - Interfaz que define el contrato para las operaciones del repositorio de compradores
type BuyerRepositoryI interface {
	// GetAll - Retrieves all buyers from the database and returns them as a map with buyer ID as key
	// GetAll - Obtiene todos los compradores de la base de datos y los retorna como un mapa con el ID del comprador como clave
	GetAll(ctx context.Context) (map[int]models.Buyer, error)

	// GetById - Retrieves a specific buyer by their ID from the database
	// GetById - Obtiene un comprador específico por su ID de la base de datos
	GetById(ctx context.Context, id int) (models.Buyer, error)

	// DeleteById - Removes a buyer from the database by their ID
	// DeleteById - Elimina un comprador de la base de datos por su ID
	DeleteById(ctx context.Context, id int) error

	// Create - Inserts a new buyer into the database and returns the created buyer with its generated ID
	// Create - Inserta un nuevo comprador en la base de datos y retorna el comprador creado con su ID generado
	Create(ctx context.Context, buyer models.Buyer) (models.Buyer, error)

	// Update - Modifies an existing buyer in the database with support for partial updates
	// Update - Modifica un comprador existente en la base de datos con soporte para actualizaciones parciales
	Update(ctx context.Context, buyerId int, buyer models.Buyer) (models.Buyer, error)

	// GetCardNumberIds - Retrieves all card number IDs from the database for validation purposes
	// GetCardNumberIds - Obtiene todos los IDs de números de tarjeta de la base de datos para propósitos de validación
	GetCardNumberIds() ([]string, error)

	// ExistBuyerById - Checks if a buyer with the given ID exists in the database
	// ExistBuyerById - Verifica si un comprador con el ID dado existe en la base de datos
	ExistBuyerById(ctx context.Context, buyerId int) (bool, error)
}

// MySqlBuyerRepository - MySQL implementation of the BuyerRepositoryI interface
// MySqlBuyerRepository - Implementación MySQL de la interfaz BuyerRepositoryI
type MySqlBuyerRepository struct {
	Db *sql.DB // Database connection / Conexión a la base de datos
}

// GetAll - Retrieves all buyers from the MySQL database and returns them as a map with buyer ID as key
// GetAll - Obtiene todos los compradores de la base de datos MySQL y los retorna como un mapa con el ID del comprador como clave
func (r *MySqlBuyerRepository) GetAll(ctx context.Context) (map[int]models.Buyer, error) {
	buyers := make(map[int]models.Buyer)

	// SQL query to select all buyer fields / Consulta SQL para seleccionar todos los campos del comprador
	query := "select id, id_card_number, first_name, last_name from buyers"
	rows, err := r.Db.QueryContext(ctx, query)
	if err != nil {
		return buyers, fmt.Errorf("%w - %s", error_message.ErrInternalServerError, err.Error())
	}
	defer rows.Close()

	// Create temporary map to store buyers before returning / Crear mapa temporal para almacenar compradores antes de retornar
	tempBuyersMap := make(map[int]models.Buyer)
	// Iterate through all rows and map each buyer to the result map / Itera a través de todas las filas y mapea cada comprador al mapa de resultados
	for rows.Next() {
		buyer := models.Buyer{}
		err = rows.Scan(&buyer.Id, &buyer.CardNumberId, &buyer.FirstName, &buyer.LastName)
		if err != nil {
			return buyers, fmt.Errorf("%w - %s", error_message.ErrInternalServerError, err.Error())
		}
		tempBuyersMap[buyer.Id] = buyer
	}

	buyers = tempBuyersMap
	return buyers, nil
}

// GetById - Retrieves a specific buyer by their ID from the MySQL database
// GetById - Obtiene un comprador específico por su ID de la base de datos MySQL
func (r *MySqlBuyerRepository) GetById(ctx context.Context, id int) (models.Buyer, error) {
	buyer := models.Buyer{}

	// SQL query to select buyer by specific ID / Consulta SQL para seleccionar comprador por ID específico
	query := "select id, id_card_number, first_name, last_name from buyers where id = ?"
	row := r.Db.QueryRowContext(ctx, query, id)
	err := row.Err()
	if err != nil {
		return buyer, fmt.Errorf("%w - %s", error_message.ErrInternalServerError, err.Error())
	}

	err = row.Scan(&buyer.Id, &buyer.CardNumberId, &buyer.FirstName, &buyer.LastName)
	if err != nil {
		// Handle case when no buyer is found / Maneja el caso cuando no se encuentra ningún comprador
		if errors.Is(err, sql.ErrNoRows) {
			return models.Buyer{}, fmt.Errorf("%w. %s %d %s", error_message.ErrNotFound, "Buyer with Id", id, "doesn't exists.")
		}
		return models.Buyer{}, fmt.Errorf("%w - %s", error_message.ErrInternalServerError, err.Error())
	}

	return buyer, nil
}

// DeleteById - Removes a buyer from the MySQL database by their ID
// DeleteById - Elimina un comprador de la base de datos MySQL por su ID
func (r *MySqlBuyerRepository) DeleteById(ctx context.Context, id int) error {
	// SQL query to delete buyer by ID / Consulta SQL para eliminar comprador por ID
	query := "delete from buyers where id = ?"

	result, err := r.Db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("%w. %s", error_message.ErrInternalServerError, err)
	}

	// Check if any rows were affected to confirm deletion / Verifica si alguna fila fue afectada para confirmar la eliminación
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("%w. %s", error_message.ErrInternalServerError, err)
	}

	// If no rows affected, buyer doesn't exist / Si ninguna fila fue afectada, el comprador no existe
	if rowsAffected == 0 {
		return fmt.Errorf("%w. %s %d %s", error_message.ErrNotFound, "Buyer with Id", id, "doesn't exists.")
	}

	return nil
}

// Create - Inserts a new buyer into the MySQL database and returns the created buyer with its generated ID
// Create - Inserta un nuevo comprador en la base de datos MySQL y retorna el comprador creado con su ID generado
func (r *MySqlBuyerRepository) Create(ctx context.Context, buyer models.Buyer) (models.Buyer, error) {
	// SQL query to insert new buyer / Consulta SQL para insertar nuevo comprador
	query := `insert into buyers (id_card_number, first_name, last_name) values (?, ?, ?)`

	result, err := r.Db.ExecContext(ctx, query, buyer.CardNumberId, buyer.FirstName, buyer.LastName)

	if err != nil {
		return models.Buyer{}, fmt.Errorf("%w - %s", error_message.ErrInternalServerError, err.Error())
	}

	// Get the auto-generated ID from the database / Obtiene el ID autogenerado de la base de datos
	lastId, err := result.LastInsertId()
	if err != nil {
		return models.Buyer{}, fmt.Errorf("%w - %s", error_message.ErrInternalServerError, err.Error())
	}

	buyer.Id = int(lastId)
	return buyer, nil
}

// Update - Modifies an existing buyer in the MySQL database with support for partial updates
// Update - Modifica un comprador existente en la base de datos MySQL con soporte para actualizaciones parciales
func (r *MySqlBuyerRepository) Update(ctx context.Context, buyerId int, buyer models.Buyer) (models.Buyer, error) {
	updates := []string{}
	values := []interface{}{}

	// Build dynamic UPDATE query based on provided fields / Construye consulta UPDATE dinámica basada en campos proporcionados
	if buyer.FirstName != "" {
		updates = append(updates, "first_name = ?")
		values = append(values, buyer.FirstName)
	}
	if buyer.LastName != "" {
		updates = append(updates, "last_name = ?")
		values = append(values, buyer.LastName)
	}
	if buyer.CardNumberId != "" {
		updates = append(updates, "id_card_number = ?")
		values = append(values, buyer.CardNumberId)
	}

	// Execute dynamic UPDATE query / Ejecuta consulta UPDATE dinámica
	query := "UPDATE buyers SET " + strings.Join(updates, ", ") + " WHERE id = ?"
	values = append(values, buyerId)

	result, err := r.Db.ExecContext(ctx, query, values...)
	if err != nil {
		return models.Buyer{}, fmt.Errorf("%w - %s", error_message.ErrInternalServerError, err.Error())
	}

	// Check if any rows were affected to confirm update / Verifica si alguna fila fue afectada para confirmar la actualización
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return models.Buyer{}, fmt.Errorf("%w. %s", error_message.ErrInternalServerError, err)
	}

	// If no rows affected, buyer doesn't exist / Si ninguna fila fue afectada, el comprador no existe
	if rowsAffected == 0 {
		return models.Buyer{}, fmt.Errorf("%w. %s %d %s", error_message.ErrNotFound, "Buyer with Id", buyerId, "doesn't exists.")
	}

	// Retrieve and return the updated buyer / Obtiene y retorna el comprador actualizado
	updatedUser, err := r.GetById(ctx, buyerId)
	if err != nil {
		return models.Buyer{}, fmt.Errorf("%w - %s", error_message.ErrInternalServerError, err.Error())
	}
	return updatedUser, nil
}

// GetCardNumberIds - Retrieves all card number IDs from the MySQL database for validation purposes
// GetCardNumberIds - Obtiene todos los IDs de números de tarjeta de la base de datos MySQL para propósitos de validación
func (r *MySqlBuyerRepository) GetCardNumberIds() ([]string, error) {
	cardNumberIds := []string{}

	// SQL query to select all card number IDs / Consulta SQL para seleccionar todos los IDs de números de tarjeta
	query := "select id_card_number from buyers"
	rows, err := r.Db.Query(query)
	if err != nil {
		return []string{}, fmt.Errorf("%w - %s", error_message.ErrInternalServerError, err.Error())
	}
	defer rows.Close()

	// Iterate through all rows and collect card number IDs / Itera a través de todas las filas y recolecta los IDs de números de tarjeta
	for rows.Next() {
		cardNumberId := ""
		err = rows.Scan(&cardNumberId)
		if err != nil {
			return []string{}, fmt.Errorf("%w - %s", error_message.ErrInternalServerError, err.Error())
		}

		cardNumberIds = append(cardNumberIds, cardNumberId)
	}

	return cardNumberIds, nil
}

// ExistBuyerById - Checks if a buyer with the given ID exists in the MySQL database
// ExistBuyerById - Verifica si un comprador con el ID dado existe en la base de datos MySQL
func (r *MySqlBuyerRepository) ExistBuyerById(ctx context.Context, buyerId int) (bool, error) {
	// Simple query to check buyer existence using LIMIT 1 for efficiency / Consulta simple para verificar existencia del comprador usando LIMIT 1 por eficiencia
	query := "SELECT 1 FROM buyers WHERE id = ? LIMIT 1"

	var exists int64
	err := r.Db.QueryRowContext(ctx, query, buyerId).Scan(&exists)

	if err != nil {
		// If no rows found, buyer doesn't exist (not an error) / Si no se encuentran filas, el comprador no existe (no es un error)
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, fmt.Errorf("%w - %s", error_message.ErrInternalServerError, err.Error())
	}
	return true, nil
}
