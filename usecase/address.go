package usecase

import "github.com/example/clean-architecture/entity"

// AddressRepository defines the interface for address data access
type AddressRepository interface {
	GetByID(id int) (*entity.Address, error)
	GetByUserID(userID int) ([]*entity.Address, error)
	Create(address *entity.Address) error
	Update(address *entity.Address) error
	Delete(id int) error
	GetAll() ([]*entity.Address, error)
}

// AddressUsecase defines the interface for address business logic
type AddressUsecase interface {
	GetAddress(id int) (*entity.Address, error)
	GetAddressesByUser(userID int) ([]*entity.Address, error)
	CreateAddress(address *entity.Address) error
	UpdateAddress(address *entity.Address) error
	DeleteAddress(id int) error
	ListAddresses() ([]*entity.Address, error)
}
