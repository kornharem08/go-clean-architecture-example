package address

import (
	"database/sql"
	"errors"

	"github.com/example/clean-architecture/entity"
)

type IRepository interface {
	GetByID(id int) (*entity.Address, error)
	GetByUserID(userID int) ([]*entity.Address, error)
	Create(address *entity.Address) error
	Update(address *entity.Address) error
	Delete(id int) error
	GetAll() ([]*entity.Address, error)
}

// Repository implements the IRepository interface for MSSQL
type Repository struct {
	db *sql.DB
}

// NewAddressRepository creates a new MSSQL address repository instance
func NewAddressRepository(db *sql.DB) IRepository {
	return &Repository{
		db: db,
	}
}

// GetByID retrieves an address by ID from MSSQL
func (r *Repository) GetByID(id int) (*entity.Address, error) {
	address := &entity.Address{}
	query := "SELECT id, user_id, street, city, state, country, zip_code FROM addresses WHERE id = @id"

	row := r.db.QueryRow(query, sql.Named("id", id))
	err := row.Scan(&address.ID, &address.UserID, &address.Street, &address.City, &address.State, &address.Country, &address.ZipCode)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("address not found")
		}
		return nil, err
	}

	return address, nil
}

// GetByUserID retrieves all addresses for a specific user
func (r *Repository) GetByUserID(userID int) ([]*entity.Address, error) {
	query := "SELECT id, user_id, street, city, state, country, zip_code FROM addresses WHERE user_id = @user_id ORDER BY id"

	rows, err := r.db.Query(query, sql.Named("user_id", userID))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	addresses := make([]*entity.Address, 0)

	for rows.Next() {
		address := &entity.Address{}
		if err := rows.Scan(&address.ID, &address.UserID, &address.Street, &address.City, &address.State, &address.Country, &address.ZipCode); err != nil {
			return nil, err
		}
		addresses = append(addresses, address)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return addresses, nil
}

// Create inserts a new address into MSSQL
func (r *Repository) Create(address *entity.Address) error {
	query := `
		INSERT INTO addresses (user_id, street, city, state, country, zip_code) 
		VALUES (@user_id, @street, @city, @state, @country, @zip_code)
	`

	result, err := r.db.Exec(
		query,
		sql.Named("user_id", address.UserID),
		sql.Named("street", address.Street),
		sql.Named("city", address.City),
		sql.Named("state", address.State),
		sql.Named("country", address.Country),
		sql.Named("zip_code", address.ZipCode),
	)

	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	address.ID = int(id)
	return nil
}

// Update modifies an existing address in MSSQL
func (r *Repository) Update(address *entity.Address) error {
	query := `
		UPDATE addresses 
		SET street = @street, city = @city, state = @state, country = @country, zip_code = @zip_code
		WHERE id = @id
	`

	result, err := r.db.Exec(
		query,
		sql.Named("street", address.Street),
		sql.Named("city", address.City),
		sql.Named("state", address.State),
		sql.Named("country", address.Country),
		sql.Named("zip_code", address.ZipCode),
		sql.Named("id", address.ID),
	)

	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("address not found")
	}

	return nil
}

// Delete removes an address from MSSQL
func (r *Repository) Delete(id int) error {
	query := "DELETE FROM addresses WHERE id = @id"

	result, err := r.db.Exec(query, sql.Named("id", id))
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("address not found")
	}

	return nil
}

// GetAll retrieves all addresses from MSSQL
func (r *Repository) GetAll() ([]*entity.Address, error) {
	query := "SELECT id, user_id, street, city, state, country, zip_code FROM addresses"

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	addresses := make([]*entity.Address, 0)

	for rows.Next() {
		address := &entity.Address{}
		if err := rows.Scan(&address.ID, &address.UserID, &address.Street, &address.City, &address.State, &address.Country, &address.ZipCode); err != nil {
			return nil, err
		}
		addresses = append(addresses, address)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return addresses, nil
}
