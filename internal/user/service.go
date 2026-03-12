package user

import (
	"errors"

	"github.com/example/clean-architecture/entity"
	"github.com/example/clean-architecture/infra/mssql/address"
	"github.com/example/clean-architecture/infra/mssql/user"
)

type IService interface {
	GetUser(id int) (*entity.User, error)
	CreateUser(user *entity.User) error
	UpdateUser(user *entity.User) error
	DeleteUser(id int) error
	ListUsers() ([]*entity.User, error)
}

// Service implements the UserUsecase interface
type Service struct {
	repo        user.IRepository
	addressRepo address.IRepository // Added for address management
}

// NewService creates a new UserService instance
func NewService(repo user.IRepository) IService {
	return &Service{
		repo: repo,
	}
}

// NewServiceWithAddress creates a new UserService instance with address repository
func NewServiceWithAddress(repo user.IRepository, addressRepo address.IRepository) IService {
	return &Service{
		repo:        repo,
		addressRepo: addressRepo,
	}
}

// GetUser retrieves a user by ID
func (s *Service) GetUser(id int) (*entity.User, error) {
	if id <= 0 {
		return nil, errors.New("invalid user ID")
	}

	return s.repo.GetByID(id)
}

// CreateUser creates a new user
func (s *Service) CreateUser(user *entity.User) error {
	if user == nil {
		return errors.New("user cannot be nil")
	}

	if user.Name == "" {
		return errors.New("user name is required")
	}

	if user.Email == "" {
		return errors.New("user email is required")
	}

	return s.repo.Create(user)
}

// UpdateUser updates an existing user
func (s *Service) UpdateUser(user *entity.User) error {
	if user == nil {
		return errors.New("user cannot be nil")
	}

	if user.ID <= 0 {
		return errors.New("invalid user ID")
	}

	return s.repo.Update(user)
}

// DeleteUser deletes a user by ID
func (s *Service) DeleteUser(id int) error {
	if id <= 0 {
		return errors.New("invalid user ID")
	}

	return s.repo.Delete(id)
}

// ListUsers retrieves all users
func (s *Service) ListUsers() ([]*entity.User, error) {
	return s.repo.GetAll()
}
