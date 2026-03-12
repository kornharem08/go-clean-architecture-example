package address

import (
	"database/sql"
	"errors"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/example/clean-architecture/entity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newTestRepo(t *testing.T) (*Repository, sqlmock.Sqlmock, func()) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	repo := &Repository{db: db}
	cleanup := func() {
		mock.ExpectClose()
		require.NoError(t, db.Close())
		require.NoError(t, mock.ExpectationsWereMet())
	}

	return repo, mock, cleanup
}

func TestRepositoryGetByID(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo, mock, cleanup := newTestRepo(t)
		defer cleanup()

		query := regexp.QuoteMeta("SELECT id, user_id, street, city, state, country, zip_code FROM addresses WHERE id = @id")
		rows := sqlmock.NewRows([]string{"id", "user_id", "street", "city", "state", "country", "zip_code"}).
			AddRow(1, 10, "123 Main", "Bangkok", "BKK", "TH", "10110")

		mock.ExpectQuery(query).WillReturnRows(rows)

		got, err := repo.GetByID(1)
		require.NoError(t, err)
		assert.Equal(t, 1, got.ID)
		assert.Equal(t, 10, got.UserID)
		assert.Equal(t, "123 Main", got.Street)
	})

	t.Run("not found", func(t *testing.T) {
		repo, mock, cleanup := newTestRepo(t)
		defer cleanup()

		query := regexp.QuoteMeta("SELECT id, user_id, street, city, state, country, zip_code FROM addresses WHERE id = @id")
		mock.ExpectQuery(query).WillReturnError(sql.ErrNoRows)

		got, err := repo.GetByID(999)
		assert.Nil(t, got)
		assert.EqualError(t, err, "address not found")
	})
}

func TestRepositoryGetByUserID(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo, mock, cleanup := newTestRepo(t)
		defer cleanup()

		query := regexp.QuoteMeta("SELECT id, user_id, street, city, state, country, zip_code FROM addresses WHERE user_id = @user_id ORDER BY id")
		rows := sqlmock.NewRows([]string{"id", "user_id", "street", "city", "state", "country", "zip_code"}).
			AddRow(1, 10, "123 Main", "Bangkok", "BKK", "TH", "10110").
			AddRow(2, 10, "456 Oak", "Bangkok", "BKK", "TH", "10120")

		mock.ExpectQuery(query).WillReturnRows(rows)

		got, err := repo.GetByUserID(10)
		require.NoError(t, err)
		assert.Len(t, got, 2)
		assert.Equal(t, "456 Oak", got[1].Street)
	})

	t.Run("query error", func(t *testing.T) {
		repo, mock, cleanup := newTestRepo(t)
		defer cleanup()

		query := regexp.QuoteMeta("SELECT id, user_id, street, city, state, country, zip_code FROM addresses WHERE user_id = @user_id ORDER BY id")
		mock.ExpectQuery(query).WillReturnError(errors.New("db down"))

		got, err := repo.GetByUserID(10)
		assert.Nil(t, got)
		assert.EqualError(t, err, "db down")
	})
}

func TestRepositoryCreate(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo, mock, cleanup := newTestRepo(t)
		defer cleanup()

		query := regexp.QuoteMeta("INSERT INTO addresses (user_id, street, city, state, country, zip_code) \n\t\tVALUES (@user_id, @street, @city, @state, @country, @zip_code)")
		mock.ExpectExec(query).WillReturnResult(sqlmock.NewResult(5, 1))

		addr := &entity.Address{UserID: 10, Street: "123 Main", City: "Bangkok", State: "BKK", Country: "TH", ZipCode: "10110"}
		err := repo.Create(addr)
		require.NoError(t, err)
		assert.Equal(t, 5, addr.ID)
	})

	t.Run("exec error", func(t *testing.T) {
		repo, mock, cleanup := newTestRepo(t)
		defer cleanup()

		query := regexp.QuoteMeta("INSERT INTO addresses (user_id, street, city, state, country, zip_code) \n\t\tVALUES (@user_id, @street, @city, @state, @country, @zip_code)")
		mock.ExpectExec(query).WillReturnError(errors.New("insert failed"))

		addr := &entity.Address{UserID: 10, Street: "123 Main", City: "Bangkok", Country: "TH"}
		err := repo.Create(addr)
		assert.EqualError(t, err, "insert failed")
	})
}

func TestRepositoryUpdate(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo, mock, cleanup := newTestRepo(t)
		defer cleanup()

		query := regexp.QuoteMeta("UPDATE addresses \n\t\tSET street = @street, city = @city, state = @state, country = @country, zip_code = @zip_code\n\t\tWHERE id = @id")
		mock.ExpectExec(query).WillReturnResult(sqlmock.NewResult(0, 1))

		addr := &entity.Address{ID: 1, UserID: 10, Street: "123 Main", City: "Bangkok", State: "BKK", Country: "TH", ZipCode: "10110"}
		err := repo.Update(addr)
		assert.NoError(t, err)
	})

	t.Run("not found", func(t *testing.T) {
		repo, mock, cleanup := newTestRepo(t)
		defer cleanup()

		query := regexp.QuoteMeta("UPDATE addresses \n\t\tSET street = @street, city = @city, state = @state, country = @country, zip_code = @zip_code\n\t\tWHERE id = @id")
		mock.ExpectExec(query).WillReturnResult(sqlmock.NewResult(0, 0))

		addr := &entity.Address{ID: 999, UserID: 10, Street: "123 Main", City: "Bangkok", Country: "TH"}
		err := repo.Update(addr)
		assert.EqualError(t, err, "address not found")
	})
}

func TestRepositoryDelete(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo, mock, cleanup := newTestRepo(t)
		defer cleanup()

		query := regexp.QuoteMeta("DELETE FROM addresses WHERE id = @id")
		mock.ExpectExec(query).WillReturnResult(sqlmock.NewResult(0, 1))

		err := repo.Delete(1)
		assert.NoError(t, err)
	})

	t.Run("not found", func(t *testing.T) {
		repo, mock, cleanup := newTestRepo(t)
		defer cleanup()

		query := regexp.QuoteMeta("DELETE FROM addresses WHERE id = @id")
		mock.ExpectExec(query).WillReturnResult(sqlmock.NewResult(0, 0))

		err := repo.Delete(999)
		assert.EqualError(t, err, "address not found")
	})
}

func TestRepositoryGetAll(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo, mock, cleanup := newTestRepo(t)
		defer cleanup()

		query := regexp.QuoteMeta("SELECT id, user_id, street, city, state, country, zip_code FROM addresses")
		rows := sqlmock.NewRows([]string{"id", "user_id", "street", "city", "state", "country", "zip_code"}).
			AddRow(1, 10, "123 Main", "Bangkok", "BKK", "TH", "10110").
			AddRow(2, 11, "456 Oak", "Chiang Mai", "CM", "TH", "50000")

		mock.ExpectQuery(query).WillReturnRows(rows)

		got, err := repo.GetAll()
		require.NoError(t, err)
		assert.Len(t, got, 2)
		assert.Equal(t, "Chiang Mai", got[1].City)
	})

	t.Run("query error", func(t *testing.T) {
		repo, mock, cleanup := newTestRepo(t)
		defer cleanup()

		query := regexp.QuoteMeta("SELECT id, user_id, street, city, state, country, zip_code FROM addresses")
		mock.ExpectQuery(query).WillReturnError(errors.New("db down"))

		got, err := repo.GetAll()
		assert.Nil(t, got)
		assert.EqualError(t, err, "db down")
	})
}
