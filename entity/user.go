package entity

// User represents the domain entity for a user
type User struct {
	ID    int    `db:"id"`
	Name  string `db:"name"`
	Email string `db:"email"`
	Phone string `db:"phone"`
}
