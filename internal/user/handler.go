package user

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/example/clean-architecture/api"
	"github.com/example/clean-architecture/entity"
	"github.com/example/clean-architecture/internal/address"
)

// Handler handles HTTP requests for user operations
type Handler struct {
	userUsecase    IService
	addressUsecase address.IService
}

// NewHandlerWithUsecases creates a user handler with multiple injected usecases
func NewHandlerWithUsecases(userUsecase IService, addressUsecase address.IService) *Handler {
	return &Handler{
		userUsecase:    userUsecase,
		addressUsecase: addressUsecase,
	}
}

// GetUser handles GET /users/:id
func (h *Handler) GetUser(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	user, err := h.userUsecase.GetUser(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	// If address usecase is available, fetch and include addresses
	if h.addressUsecase != nil {
		addresses, _ := h.addressUsecase.GetAddressesByUser(id)
		if addresses == nil {
			addresses = make([]*entity.Address, 0)
		}
		c.JSON(http.StatusOK, api.NewUserWithAddressesResponse(user, addresses))
		return
	}

	c.JSON(http.StatusOK, api.NewUserResponse(user))
}

// CreateUser handles POST /users
func (h *Handler) CreateUser(c *gin.Context) {
	var req api.CreateUserRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user := &entity.User{
		Name:  req.Name,
		Email: req.Email,
		Phone: req.Phone,
	}

	if err := h.userUsecase.CreateUser(user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, api.NewUserResponse(user))
}

// UpdateUser handles PUT /users/:id
func (h *Handler) UpdateUser(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	var req api.UpdateUserRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user := &entity.User{
		ID:    id,
		Name:  req.Name,
		Email: req.Email,
		Phone: req.Phone,
	}

	if err := h.userUsecase.UpdateUser(user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, api.NewUserResponse(user))
}

// DeleteUser handles DELETE /users/:id
func (h *Handler) DeleteUser(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	if err := h.userUsecase.DeleteUser(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "user deleted successfully"})
}

// ListUsers handles GET /users
func (h *Handler) ListUsers(c *gin.Context) {
	users, err := h.userUsecase.ListUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := make([]*api.UserResponse, 0, len(users))
	for _, user := range users {
		response = append(response, api.NewUserResponse(user))
	}

	c.JSON(http.StatusOK, gin.H{"data": response})
}
