package repositories

import (
	"database/sql"

	"github.com/sajimenezher_meli/meli-frescos-8/internal/error_message"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/models"
)

var sellerRepositoryInstance SellerRepository

// NewSQLSellerRepository - Creates and returns a new instance of SQLSellerRepository using singleton pattern
// NewSQLSellerRepository - Crea y retorna una nueva instancia de SQLSellerRepository usando patrón singleton
func NewSQLSellerRepository(db *sql.DB) SellerRepository {
	if sellerRepositoryInstance != nil {
		return sellerRepositoryInstance
	}

	sellerRepositoryInstance = &SQLSellerRepository{db: db}
	return sellerRepositoryInstance
}

// SellerRepository - Interface defining the contract for seller repository operations
// SellerRepository - Interfaz que define el contrato para las operaciones del repositorio de vendedores
type SellerRepository interface {
	// GetAll - Retrieves all sellers from the database
	// GetAll - Obtiene todos los vendedores de la base de datos
	GetAll() ([]models.Seller, error)

	// Save - Creates a new seller in the database with validation checks for CID uniqueness and locality existence
	// Save - Crea un nuevo vendedor en la base de datos con validaciones de unicidad de CID y existencia de localidad
	Save(seller models.Seller) ([]models.Seller, error)

	// Update - Modifies an existing seller in the database with partial update support
	// Update - Modifica un vendedor existente en la base de datos con soporte para actualizaciones parciales
	Update(id int, seller models.Seller) ([]models.Seller, error)

	// Delete - Removes a seller from the database by their ID
	// Delete - Elimina un vendedor de la base de datos por su ID
	Delete(id int) error
}

// SQLSellerRepository - SQL implementation of the SellerRepository interface
// SQLSellerRepository - Implementación SQL de la interfaz SellerRepository
type SQLSellerRepository struct {
	db *sql.DB // Database connection / Conexión a la base de datos
}

// GetAll - Retrieves all sellers from the database and returns them as a slice
// GetAll - Obtiene todos los vendedores de la base de datos y los retorna como un slice
func (r *SQLSellerRepository) GetAll() ([]models.Seller, error) {
	// Execute query to select all seller fields / Ejecutar consulta para seleccionar todos los campos del vendedor
	rows, err := r.db.Query("SELECT id, cid, company_name, address, telephone, locality_id FROM sellers")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Iterate through all rows and scan each seller into the results slice
	// Itera a través de todas las filas y escanea cada vendedor en el slice de resultados
	var sellers []models.Seller
	for rows.Next() {
		var s models.Seller
		if err := rows.Scan(&s.Id, &s.CID, &s.CompanyName, &s.Address, &s.Telephone, &s.LocalityID); err != nil {
			return nil, err
		}
		sellers = append(sellers, s)
	}
	return sellers, nil
}

// Save - Creates a new seller in the database with validation for CID uniqueness and locality existence
// Save - Crea un nuevo vendedor en la base de datos con validación de unicidad de CID y existencia de localidad
func (r *SQLSellerRepository) Save(seller models.Seller) ([]models.Seller, error) {
	// Validate that no other seller exists with the same CID / Validar que no exista otro vendedor con el mismo CID
	var exists bool
	err := r.db.QueryRow("SELECT EXISTS(SELECT 1 FROM sellers WHERE cid = ?)", seller.CID).Scan(&exists)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, error_message.ErrAlreadyExists
	}

	// Validate that the locality exists before creating the seller / Validar que la localidad existe antes de crear el vendedor
	var existsLocality bool
	err = r.db.QueryRow("SELECT EXISTS(SELECT 1 FROM localities WHERE id = ?)", seller.LocalityID).Scan(&existsLocality)
	if err != nil {
		return nil, err
	}
	if !existsLocality {
		return nil, error_message.ErrDependencyNotFound
	}

	// Execute insert statement with seller data / Ejecutar declaración de inserción con datos del vendedor
	res, err := r.db.Exec("INSERT INTO sellers (cid, company_name, address, telephone, locality_id) VALUES (?, ?, ?, ?, ?)",
		seller.CID, seller.CompanyName, seller.Address, seller.Telephone, seller.LocalityID)
	if err != nil {
		return nil, err
	}

	// Get the auto-generated ID and assign it to the seller / Obtener el ID autogenerado y asignarlo al vendedor
	lastID, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}
	seller.Id = int(lastID)
	return []models.Seller{seller}, nil
}

// Update - Modifies an existing seller in the database with support for partial updates
// Update - Modifica un vendedor existente en la base de datos con soporte para actualizaciones parciales
func (r *SQLSellerRepository) Update(id int, seller models.Seller) ([]models.Seller, error) {
	// Retrieve existing seller data to perform partial update / Obtener datos del vendedor existente para realizar actualización parcial
	var existing models.Seller
	err := r.db.QueryRow("SELECT id, cid, company_name, address, telephone, locality_id FROM sellers WHERE id = ?", id).
		Scan(&existing.Id, &existing.CID, &existing.CompanyName, &existing.Address, &existing.Telephone, &existing.LocalityID)

	// Handle case when seller is not found / Manejar el caso cuando el vendedor no se encuentra
	if err == sql.ErrNoRows {
		return nil, error_message.ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	// Update only the provided fields (partial update logic) / Actualizar solo los campos proporcionados (lógica de actualización parcial)
	if seller.CID != "" {
		existing.CID = seller.CID
	}
	if seller.CompanyName != "" {
		existing.CompanyName = seller.CompanyName
	}
	if seller.Address != "" {
		existing.Address = seller.Address
	}
	if seller.Telephone != "" {
		existing.Telephone = seller.Telephone
	}

	if seller.LocalityID != 0 {
		existing.LocalityID = seller.LocalityID
	}

	// Execute update statement with merged data / Ejecutar declaración de actualización con datos combinados
	_, err = r.db.Exec("UPDATE sellers SET cid = ?, company_name = ?, address = ?, telephone = ?, locality_id = ? WHERE id = ?",
		existing.CID, existing.CompanyName, existing.Address, existing.Telephone, existing.LocalityID, id)

	if err != nil {
		return nil, err
	}

	return []models.Seller{existing}, nil
}

// Delete - Removes a seller from the database by their ID
// Delete - Elimina un vendedor de la base de datos por su ID
func (r *SQLSellerRepository) Delete(id int) error {
	// Execute delete statement for the specified seller ID / Ejecutar declaración de eliminación para el ID del vendedor especificado
	res, err := r.db.Exec("DELETE FROM sellers WHERE id = ?", id)
	if err != nil {
		return err
	}

	// Check if any rows were affected to confirm deletion / Verificar si alguna fila fue afectada para confirmar la eliminación
	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	// If no rows affected, seller doesn't exist / Si ninguna fila fue afectada, el vendedor no existe
	if count == 0 {
		return error_message.ErrNotFound
	}
	return nil
}
