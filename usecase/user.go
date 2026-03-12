package usecase

import "github.com/example/clean-architecture/entity"

// UserRepository defines the interface for user data access
type UserRepository interface {
	GetByID(id int) (*entity.User, error)
	Create(user *entity.User) error
	Update(user *entity.User) error
	Delete(id int) error
	GetAll() ([]*entity.User, error)
}

// UserUsecase defines the interface for user business logic
type UserUsecase interface {
	GetUser(id int) (*entity.User, error)
	CreateUser(user *entity.User) error
	UpdateUser(user *entity.User) error
	DeleteUser(id int) error
	ListUsers() ([]*entity.User, error)
}
