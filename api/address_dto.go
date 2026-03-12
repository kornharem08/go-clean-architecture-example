package api

import "github.com/example/clean-architecture/entity"

// CreateAddressRequest is the DTO for creating an address
type CreateAddressRequest struct {
	UserID  int    `json:"user_id" binding:"required"`
	Street  string `json:"street" binding:"required"`
	City    string `json:"city" binding:"required"`
	State   string `json:"state"`
	Country string `json:"country" binding:"required"`
	ZipCode string `json:"zip_code"`
}

// UpdateAddressRequest is the DTO for updating an address
type UpdateAddressRequest struct {
	Street  string `json:"street" binding:"required"`
	City    string `json:"city" binding:"required"`
	State   string `json:"state"`
	Country string `json:"country" binding:"required"`
	ZipCode string `json:"zip_code"`
}

// AddressResponse is the DTO for returning address data
type AddressResponse struct {
	ID      int    `json:"id"`
	UserID  int    `json:"user_id"`
	Street  string `json:"street"`
	City    string `json:"city"`
	State   string `json:"state"`
	Country string `json:"country"`
	ZipCode string `json:"zip_code"`
}

// UserWithAddressesResponse returns user data with associated addresses
type UserWithAddressesResponse struct {
	ID        int                `json:"id"`
	Name      string             `json:"name"`
	Email     string             `json:"email"`
	Phone     string             `json:"phone"`
	Addresses []*AddressResponse `json:"addresses"`
}

// NewAddressResponse creates a new AddressResponse from an Address entity
func NewAddressResponse(address *entity.Address) *AddressResponse {
	return &AddressResponse{
		ID:      address.ID,
		UserID:  address.UserID,
		Street:  address.Street,
		City:    address.City,
		State:   address.State,
		Country: address.Country,
		ZipCode: address.ZipCode,
	}
}

// NewUserWithAddressesResponse creates a UserWithAddressesResponse
func NewUserWithAddressesResponse(user *entity.User, addresses []*entity.Address) *UserWithAddressesResponse {
	addressResponses := make([]*AddressResponse, 0, len(addresses))
	for _, addr := range addresses {
		addressResponses = append(addressResponses, NewAddressResponse(addr))
	}

	return &UserWithAddressesResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Phone:     user.Phone,
		Addresses: addressResponses,
	}
}
