package address

import (
	"errors"

	"github.com/example/clean-architecture/entity"
	addrRepo "github.com/example/clean-architecture/infra/mssql/address"
)

type IService interface {
	GetAddress(id int) (*entity.Address, error)
	GetAddressesByUser(userID int) ([]*entity.Address, error)
	CreateAddress(address *entity.Address) error
	UpdateAddress(address *entity.Address) error
	DeleteAddress(id int) error
	ListAddresses() ([]*entity.Address, error)
}

// Service implements the IService interface
type Service struct {
	repo addrRepo.IRepository
}

// NewService creates a new AddressService instance
func NewService(repo addrRepo.IRepository) IService {
	return &Service{
		repo: repo,
	}
}

// GetAddress retrieves an address by ID
func (s *Service) GetAddress(id int) (*entity.Address, error) {
	if id <= 0 {
		return nil, errors.New("invalid address ID")
	}

	return s.repo.GetByID(id)
}

// GetAddressesByUser retrieves all addresses for a specific user
func (s *Service) GetAddressesByUser(userID int) ([]*entity.Address, error) {
	if userID <= 0 {
		return nil, errors.New("invalid user ID")
	}

	return s.repo.GetByUserID(userID)
}

// CreateAddress creates a new address
func (s *Service) CreateAddress(address *entity.Address) error {
	if address == nil {
		return errors.New("address cannot be nil")
	}

	if address.UserID <= 0 {
		return errors.New("user ID is required")
	}

	if address.Street == "" {
		return errors.New("street is required")
	}

	if address.City == "" {
		return errors.New("city is required")
	}

	return s.repo.Create(address)
}

// UpdateAddress updates an existing address
func (s *Service) UpdateAddress(address *entity.Address) error {
	if address == nil {
		return errors.New("address cannot be nil")
	}

	if address.ID <= 0 {
		return errors.New("invalid address ID")
	}

	return s.repo.Update(address)
}

// DeleteAddress deletes an address by ID
func (s *Service) DeleteAddress(id int) error {
	if id <= 0 {
		return errors.New("invalid address ID")
	}

	return s.repo.Delete(id)
}

// ListAddresses retrieves all addresses
func (s *Service) ListAddresses() ([]*entity.Address, error) {
	return s.repo.GetAll()
}
