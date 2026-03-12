package entity

// Address represents the domain entity for an address
type Address struct {
	ID      int    `db:"id"`
	UserID  int    `db:"user_id"`
	Street  string `db:"street"`
	City    string `db:"city"`
	State   string `db:"state"`
	Country string `db:"country"`
	ZipCode string `db:"zip_code"`
}
