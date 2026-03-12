package address

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/example/clean-architecture/entity"
	"github.com/example/clean-architecture/internal/address/mocks"
)

func TestService_GetAddress(t *testing.T) {
	t.Run("invalid id", func(t *testing.T) {
		repo := mocks.NewAddressRepositoryMock(t)
		svc := NewService(repo)

		addr, err := svc.GetAddress(0)
		assert.Nil(t, addr)
		assert.EqualError(t, err, "invalid address ID")
	})

	t.Run("repo error", func(t *testing.T) {
		repo := mocks.NewAddressRepositoryMock(t)
		repo.On("GetByID", 1).Return((*entity.Address)(nil), errors.New("not found")).Once()

		svc := NewService(repo)
		addr, err := svc.GetAddress(1)
		assert.Nil(t, addr)
		assert.EqualError(t, err, "not found")
		repo.AssertExpectations(t)
	})

	t.Run("success", func(t *testing.T) {
		want := &entity.Address{ID: 1, UserID: 1, Street: "Main", City: "Bangkok", Country: "TH"}
		repo := mocks.NewAddressRepositoryMock(t)
		repo.On("GetByID", 1).Return(want, nil).Once()

		svc := NewService(repo)
		got, err := svc.GetAddress(1)
		assert.NoError(t, err)
		assert.Equal(t, want, got)
		repo.AssertExpectations(t)
	})
}

func TestService_GetAddressesByUser(t *testing.T) {
	t.Run("invalid user id", func(t *testing.T) {
		repo := mocks.NewAddressRepositoryMock(t)
		svc := NewService(repo)

		addresses, err := svc.GetAddressesByUser(0)
		assert.Nil(t, addresses)
		assert.EqualError(t, err, "invalid user ID")
	})

	t.Run("repo error", func(t *testing.T) {
		repo := mocks.NewAddressRepositoryMock(t)
		repo.On("GetByUserID", 1).Return(([]*entity.Address)(nil), errors.New("db down")).Once()

		svc := NewService(repo)
		addresses, err := svc.GetAddressesByUser(1)
		assert.Nil(t, addresses)
		assert.EqualError(t, err, "db down")
		repo.AssertExpectations(t)
	})

	t.Run("success", func(t *testing.T) {
		want := []*entity.Address{{ID: 1, UserID: 1, Street: "Main", City: "Bangkok", Country: "TH"}}
		repo := mocks.NewAddressRepositoryMock(t)
		repo.On("GetByUserID", 1).Return(want, nil).Once()

		svc := NewService(repo)
		got, err := svc.GetAddressesByUser(1)
		assert.NoError(t, err)
		assert.Equal(t, want, got)
		repo.AssertExpectations(t)
	})
}

func TestService_CreateAddress(t *testing.T) {
	t.Run("nil address", func(t *testing.T) {
		repo := mocks.NewAddressRepositoryMock(t)
		svc := NewService(repo)

		err := svc.CreateAddress(nil)
		assert.EqualError(t, err, "address cannot be nil")
	})

	t.Run("user id required", func(t *testing.T) {
		repo := mocks.NewAddressRepositoryMock(t)
		svc := NewService(repo)

		err := svc.CreateAddress(&entity.Address{Street: "Main", City: "Bangkok"})
		assert.EqualError(t, err, "user ID is required")
	})

	t.Run("street required", func(t *testing.T) {
		repo := mocks.NewAddressRepositoryMock(t)
		svc := NewService(repo)

		err := svc.CreateAddress(&entity.Address{UserID: 1, City: "Bangkok"})
		assert.EqualError(t, err, "street is required")
	})

	t.Run("city required", func(t *testing.T) {
		repo := mocks.NewAddressRepositoryMock(t)
		svc := NewService(repo)

		err := svc.CreateAddress(&entity.Address{UserID: 1, Street: "Main"})
		assert.EqualError(t, err, "city is required")
	})

	t.Run("repo error", func(t *testing.T) {
		repo := mocks.NewAddressRepositoryMock(t)
		repo.On("Create", mock.AnythingOfType("*entity.Address")).Return(errors.New("duplicate")).Once()

		svc := NewService(repo)
		err := svc.CreateAddress(&entity.Address{UserID: 1, Street: "Main", City: "Bangkok"})
		assert.EqualError(t, err, "duplicate")
		repo.AssertExpectations(t)
	})

	t.Run("success", func(t *testing.T) {
		repo := mocks.NewAddressRepositoryMock(t)
		repo.On("Create", mock.AnythingOfType("*entity.Address")).Return(nil).Once()

		svc := NewService(repo)
		err := svc.CreateAddress(&entity.Address{UserID: 1, Street: "Main", City: "Bangkok"})
		assert.NoError(t, err)
		repo.AssertExpectations(t)
	})
}

func TestService_UpdateAddress(t *testing.T) {
	t.Run("nil address", func(t *testing.T) {
		repo := mocks.NewAddressRepositoryMock(t)
		svc := NewService(repo)

		err := svc.UpdateAddress(nil)
		assert.EqualError(t, err, "address cannot be nil")
	})

	t.Run("invalid id", func(t *testing.T) {
		repo := mocks.NewAddressRepositoryMock(t)
		svc := NewService(repo)

		err := svc.UpdateAddress(&entity.Address{ID: 0, UserID: 1, Street: "Main", City: "Bangkok"})
		assert.EqualError(t, err, "invalid address ID")
	})

	t.Run("repo error", func(t *testing.T) {
		repo := mocks.NewAddressRepositoryMock(t)
		repo.On("Update", mock.AnythingOfType("*entity.Address")).Return(errors.New("not found")).Once()

		svc := NewService(repo)
		err := svc.UpdateAddress(&entity.Address{ID: 1, UserID: 1, Street: "Main", City: "Bangkok"})
		assert.EqualError(t, err, "not found")
		repo.AssertExpectations(t)
	})

	t.Run("success", func(t *testing.T) {
		repo := mocks.NewAddressRepositoryMock(t)
		repo.On("Update", mock.AnythingOfType("*entity.Address")).Return(nil).Once()

		svc := NewService(repo)
		err := svc.UpdateAddress(&entity.Address{ID: 1, UserID: 1, Street: "Main", City: "Bangkok"})
		assert.NoError(t, err)
		repo.AssertExpectations(t)
	})
}

func TestService_DeleteAddress(t *testing.T) {
	t.Run("invalid id", func(t *testing.T) {
		repo := mocks.NewAddressRepositoryMock(t)
		svc := NewService(repo)

		err := svc.DeleteAddress(0)
		assert.EqualError(t, err, "invalid address ID")
	})

	t.Run("repo error", func(t *testing.T) {
		repo := mocks.NewAddressRepositoryMock(t)
		repo.On("Delete", 1).Return(errors.New("not found")).Once()

		svc := NewService(repo)
		err := svc.DeleteAddress(1)
		assert.EqualError(t, err, "not found")
		repo.AssertExpectations(t)
	})

	t.Run("success", func(t *testing.T) {
		repo := mocks.NewAddressRepositoryMock(t)
		repo.On("Delete", 1).Return(nil).Once()

		svc := NewService(repo)
		err := svc.DeleteAddress(1)
		assert.NoError(t, err)
		repo.AssertExpectations(t)
	})
}

func TestService_ListAddresses(t *testing.T) {
	t.Run("repo error", func(t *testing.T) {
		repo := mocks.NewAddressRepositoryMock(t)
		repo.On("GetAll").Return(([]*entity.Address)(nil), errors.New("db down")).Once()

		svc := NewService(repo)
		addresses, err := svc.ListAddresses()
		assert.Nil(t, addresses)
		assert.EqualError(t, err, "db down")
		repo.AssertExpectations(t)
	})

	t.Run("success", func(t *testing.T) {
		want := []*entity.Address{{ID: 1, UserID: 1, Street: "Main", City: "Bangkok", Country: "TH"}}
		repo := mocks.NewAddressRepositoryMock(t)
		repo.On("GetAll").Return(want, nil).Once()

		svc := NewService(repo)
		got, err := svc.ListAddresses()
		assert.NoError(t, err)
		assert.Equal(t, want, got)
		repo.AssertExpectations(t)
	})
}
