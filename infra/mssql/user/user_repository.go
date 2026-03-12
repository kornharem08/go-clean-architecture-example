package user

import (
	"database/sql"
	"errors"

	"github.com/example/clean-architecture/entity"
)

type IRepository interface {
	GetByID(id int) (*entity.User, error)
	Create(user *entity.User) error
	Update(user *entity.User) error
	Delete(id int) error
	GetAll() ([]*entity.User, error)
}

// Repository implements the IRepository interface for MSSQL
type Repository struct {
	db *sql.DB
}

// NewRepository creates a new MSSQL repository instance
func NewRepository(db *sql.DB) IRepository {
	return &Repository{
		db: db,
	}
}

// GetByID retrieves a user by ID from MSSQL
func (r *Repository) GetByID(id int) (*entity.User, error) {
	user := &entity.User{}
	query := "SELECT id, name, email, phone FROM users WHERE id = @id"

	row := r.db.QueryRow(query, sql.Named("id", id))
	err := row.Scan(&user.ID, &user.Name, &user.Email, &user.Phone)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return user, nil
}

// Create inserts a new user into MSSQL
func (r *Repository) Create(user *entity.User) error {
	query := `
		INSERT INTO users (name, email, phone) 
		VALUES (@name, @email, @phone)
	`

	result, err := r.db.Exec(
		query,
		sql.Named("name", user.Name),
		sql.Named("email", user.Email),
		sql.Named("phone", user.Phone),
	)

	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	user.ID = int(id)
	return nil
}

// Update modifies an existing user in MSSQL
func (r *Repository) Update(user *entity.User) error {
	query := `
		UPDATE users 
		SET name = @name, email = @email, phone = @phone 
		WHERE id = @id
	`

	result, err := r.db.Exec(
		query,
		sql.Named("name", user.Name),
		sql.Named("email", user.Email),
		sql.Named("phone", user.Phone),
		sql.Named("id", user.ID),
	)

	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("user not found")
	}

	return nil
}

// Delete removes a user from MSSQL
func (r *Repository) Delete(id int) error {
	query := "DELETE FROM users WHERE id = @id"

	result, err := r.db.Exec(query, sql.Named("id", id))
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("user not found")
	}

	return nil
}

// GetAll retrieves all users from MSSQL
func (r *Repository) GetAll() ([]*entity.User, error) {
	query := "SELECT id, name, email, phone FROM users"

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := make([]*entity.User, 0)

	for rows.Next() {
		user := &entity.User{}
		if err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.Phone); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}
