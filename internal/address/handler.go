package address

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/example/clean-architecture/api"
	"github.com/example/clean-architecture/entity"
)

// Handler handles HTTP requests for address operations
type Handler struct {
	service IService
}

// NewHandler creates a new address handler
func NewHandler(service IService) *Handler {
	return &Handler{
		service: service,
	}
}

// GetAddress handles GET /addresses/:id
func (h *Handler) GetAddress(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid address ID"})
		return
	}

	address, err := h.service.GetAddress(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, api.NewAddressResponse(address))
}

// GetAddressesByUser handles GET /users/:user_id/addresses
func (h *Handler) GetAddressesByUser(c *gin.Context) {
	userID, err := strconv.Atoi(c.Param("user_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	addresses, err := h.service.GetAddressesByUser(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := make([]*api.AddressResponse, 0, len(addresses))
	for _, addr := range addresses {
		response = append(response, api.NewAddressResponse(addr))
	}

	c.JSON(http.StatusOK, gin.H{"data": response})
}

// CreateAddress handles POST /addresses
func (h *Handler) CreateAddress(c *gin.Context) {
	var req api.CreateAddressRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	address := &entity.Address{
		UserID:  req.UserID,
		Street:  req.Street,
		City:    req.City,
		State:   req.State,
		Country: req.Country,
		ZipCode: req.ZipCode,
	}

	if err := h.service.CreateAddress(address); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, api.NewAddressResponse(address))
}

// UpdateAddress handles PUT /addresses/:id
func (h *Handler) UpdateAddress(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid address ID"})
		return
	}

	var req api.UpdateAddressRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	address := &entity.Address{
		ID:      id,
		Street:  req.Street,
		City:    req.City,
		State:   req.State,
		Country: req.Country,
		ZipCode: req.ZipCode,
	}

	if err := h.service.UpdateAddress(address); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, api.NewAddressResponse(address))
}

// DeleteAddress handles DELETE /addresses/:id
func (h *Handler) DeleteAddress(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid address ID"})
		return
	}

	if err := h.service.DeleteAddress(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "address deleted successfully"})
}

// ListAddresses handles GET /addresses
func (h *Handler) ListAddresses(c *gin.Context) {
	addresses, err := h.service.ListAddresses()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := make([]*api.AddressResponse, 0, len(addresses))
	for _, addr := range addresses {
		response = append(response, api.NewAddressResponse(addr))
	}

	c.JSON(http.StatusOK, gin.H{"data": response})
}
