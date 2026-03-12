package user

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/example/clean-architecture/entity"
	"github.com/example/clean-architecture/internal/user/mocks"
)

func TestService_GetUser(t *testing.T) {
	t.Run("invalid id", func(t *testing.T) {
		repo := mocks.NewUserRepositoryMock(t)
		svc := NewService(repo)

		user, err := svc.GetUser(0)
		assert.Nil(t, user)
		assert.EqualError(t, err, "invalid user ID")
	})

	t.Run("repo error", func(t *testing.T) {
		repo := mocks.NewUserRepositoryMock(t)
		repo.On("GetByID", 1).Return((*entity.User)(nil), errors.New("not found")).Once()

		svc := NewService(repo)
		user, err := svc.GetUser(1)
		assert.Nil(t, user)
		assert.EqualError(t, err, "not found")
		repo.AssertExpectations(t)
	})

	t.Run("success", func(t *testing.T) {
		want := &entity.User{ID: 1, Name: "Alice", Email: "alice@example.com", Phone: "081"}
		repo := mocks.NewUserRepositoryMock(t)
		repo.On("GetByID", 1).Return(want, nil).Once()

		svc := NewService(repo)
		got, err := svc.GetUser(1)
		assert.NoError(t, err)
		assert.Equal(t, want, got)
		repo.AssertExpectations(t)
	})
}

func TestService_CreateUser(t *testing.T) {
	t.Run("nil user", func(t *testing.T) {
		repo := mocks.NewUserRepositoryMock(t)
		svc := NewService(repo)

		err := svc.CreateUser(nil)
		assert.EqualError(t, err, "user cannot be nil")
	})

	t.Run("name required", func(t *testing.T) {
		repo := mocks.NewUserRepositoryMock(t)
		svc := NewService(repo)

		err := svc.CreateUser(&entity.User{Email: "a@a.com"})
		assert.EqualError(t, err, "user name is required")
	})

	t.Run("email required", func(t *testing.T) {
		repo := mocks.NewUserRepositoryMock(t)
		svc := NewService(repo)

		err := svc.CreateUser(&entity.User{Name: "A"})
		assert.EqualError(t, err, "user email is required")
	})

	t.Run("repo error", func(t *testing.T) {
		repo := mocks.NewUserRepositoryMock(t)
		repo.On("Create", mock.AnythingOfType("*entity.User")).Return(errors.New("duplicate email")).Once()

		svc := NewService(repo)
		err := svc.CreateUser(&entity.User{Name: "A", Email: "a@a.com"})
		assert.EqualError(t, err, "duplicate email")
		repo.AssertExpectations(t)
	})

	t.Run("success", func(t *testing.T) {
		repo := mocks.NewUserRepositoryMock(t)
		repo.On("Create", mock.AnythingOfType("*entity.User")).Return(nil).Once()

		svc := NewService(repo)
		err := svc.CreateUser(&entity.User{Name: "A", Email: "a@a.com"})
		assert.NoError(t, err)
		repo.AssertExpectations(t)
	})
}

func TestService_UpdateUser(t *testing.T) {
	t.Run("nil user", func(t *testing.T) {
		repo := mocks.NewUserRepositoryMock(t)
		svc := NewService(repo)

		err := svc.UpdateUser(nil)
		assert.EqualError(t, err, "user cannot be nil")
	})

	t.Run("invalid id", func(t *testing.T) {
		repo := mocks.NewUserRepositoryMock(t)
		svc := NewService(repo)

		err := svc.UpdateUser(&entity.User{ID: 0, Name: "A", Email: "a@a.com"})
		assert.EqualError(t, err, "invalid user ID")
	})

	t.Run("repo error", func(t *testing.T) {
		repo := mocks.NewUserRepositoryMock(t)
		repo.On("Update", mock.AnythingOfType("*entity.User")).Return(errors.New("not found")).Once()

		svc := NewService(repo)
		err := svc.UpdateUser(&entity.User{ID: 1, Name: "A", Email: "a@a.com"})
		assert.EqualError(t, err, "not found")
		repo.AssertExpectations(t)
	})

	t.Run("success", func(t *testing.T) {
		repo := mocks.NewUserRepositoryMock(t)
		repo.On("Update", mock.AnythingOfType("*entity.User")).Return(nil).Once()

		svc := NewService(repo)
		err := svc.UpdateUser(&entity.User{ID: 1, Name: "A", Email: "a@a.com"})
		assert.NoError(t, err)
		repo.AssertExpectations(t)
	})
}

func TestService_DeleteUser(t *testing.T) {
	t.Run("invalid id", func(t *testing.T) {
		repo := mocks.NewUserRepositoryMock(t)
		svc := NewService(repo)

		err := svc.DeleteUser(0)
		assert.EqualError(t, err, "invalid user ID")
	})

	t.Run("repo error", func(t *testing.T) {
		repo := mocks.NewUserRepositoryMock(t)
		repo.On("Delete", 1).Return(errors.New("not found")).Once()

		svc := NewService(repo)
		err := svc.DeleteUser(1)
		assert.EqualError(t, err, "not found")
		repo.AssertExpectations(t)
	})

	t.Run("success", func(t *testing.T) {
		repo := mocks.NewUserRepositoryMock(t)
		repo.On("Delete", 1).Return(nil).Once()

		svc := NewService(repo)
		err := svc.DeleteUser(1)
		assert.NoError(t, err)
		repo.AssertExpectations(t)
	})
}

func TestService_ListUsers(t *testing.T) {
	t.Run("repo error", func(t *testing.T) {
		repo := mocks.NewUserRepositoryMock(t)
		repo.On("GetAll").Return(([]*entity.User)(nil), errors.New("db down")).Once()

		svc := NewService(repo)
		users, err := svc.ListUsers()
		assert.Nil(t, users)
		assert.EqualError(t, err, "db down")
		repo.AssertExpectations(t)
	})

	t.Run("success", func(t *testing.T) {
		want := []*entity.User{{ID: 1, Name: "A", Email: "a@a.com"}}
		repo := mocks.NewUserRepositoryMock(t)
		repo.On("GetAll").Return(want, nil).Once()

		svc := NewService(repo)
		got, err := svc.ListUsers()
		assert.NoError(t, err)
		assert.Equal(t, want, got)
		repo.AssertExpectations(t)
	})
}
