package user

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/example/clean-architecture/entity"
	addr "github.com/example/clean-architecture/internal/address"
	addrmocks "github.com/example/clean-architecture/internal/address/mocks"
	usermocks "github.com/example/clean-architecture/internal/user/mocks"
)

func init() {
	gin.SetMode(gin.TestMode)
}

// ─── Helpers ─────────────────────────────────────────────────────────────────

func newUserRouter(userSvc IService, addrSvc addr.IService) *gin.Engine {
	r := gin.New()
	h := NewHandlerWithUsecases(userSvc, addrSvc)
	r.GET("/users", h.ListUsers)
	r.POST("/users", h.CreateUser)
	r.GET("/users/:id", h.GetUser)
	r.PUT("/users/:id", h.UpdateUser)
	r.DELETE("/users/:id", h.DeleteUser)
	return r
}

func toJSON(v any) *bytes.Buffer {
	b, _ := json.Marshal(v)
	return bytes.NewBuffer(b)
}

// ─── GetUser ─────────────────────────────────────────────────────────────────

func TestHandlerGetUser(t *testing.T) {
	t.Run("Success_WithAddresses", func(t *testing.T) {
		userSvc := usermocks.NewUserServiceMock(t)
		addrSvc := addrmocks.NewAddressServiceMock(t)
		userSvc.On("GetUser", 1).Return(&entity.User{ID: 1, Name: "John", Email: "john@example.com", Phone: "111"}, nil).Once()
		addrSvc.On("GetAddressesByUser", 1).Return([]*entity.Address{{ID: 1, UserID: 1, Street: "123 Main", City: "Bangkok", Country: "TH"}}, nil).Once()

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/users/1", nil)
		newUserRouter(userSvc, addrSvc).ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var body map[string]any
		json.Unmarshal(w.Body.Bytes(), &body)
		assert.Equal(t, float64(1), body["id"])
		assert.Equal(t, "John", body["name"])
		addresses := body["addresses"].([]any)
		assert.Len(t, addresses, 1)
		userSvc.AssertExpectations(t)
		addrSvc.AssertExpectations(t)
	})

	t.Run("Success_AddressServiceError_StillReturnsUser", func(t *testing.T) {
		userSvc := usermocks.NewUserServiceMock(t)
		addrSvc := addrmocks.NewAddressServiceMock(t)
		userSvc.On("GetUser", 1).Return(&entity.User{ID: 1, Name: "John", Email: "john@example.com"}, nil).Once()
		addrSvc.On("GetAddressesByUser", 1).Return(([]*entity.Address)(nil), errors.New("address db down")).Once()

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/users/1", nil)
		newUserRouter(userSvc, addrSvc).ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var body map[string]any
		json.Unmarshal(w.Body.Bytes(), &body)
		assert.Equal(t, "John", body["name"])
		userSvc.AssertExpectations(t)
		addrSvc.AssertExpectations(t)
	})

	t.Run("InvalidID", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/users/abc", nil)
		newUserRouter(usermocks.NewUserServiceMock(t), addrmocks.NewAddressServiceMock(t)).ServeHTTP(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("NotFound", func(t *testing.T) {
		userSvc := usermocks.NewUserServiceMock(t)
		addrSvc := addrmocks.NewAddressServiceMock(t)
		userSvc.On("GetUser", 99).Return((*entity.User)(nil), errors.New("user not found")).Once()
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/users/99", nil)
		newUserRouter(userSvc, addrSvc).ServeHTTP(w, req)
		assert.Equal(t, http.StatusNotFound, w.Code)
		userSvc.AssertExpectations(t)
	})
}

// ─── CreateUser ──────────────────────────────────────────────────────────────

func TestHandlerCreateUser(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		userSvc := usermocks.NewUserServiceMock(t)
		addrSvc := addrmocks.NewAddressServiceMock(t)
		userSvc.On("CreateUser", mock.AnythingOfType("*entity.User")).Return(nil).Once()
		body := toJSON(map[string]string{"name": "Alice", "email": "alice@example.com", "phone": "999"})
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/users", body)
		req.Header.Set("Content-Type", "application/json")
		newUserRouter(userSvc, addrSvc).ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
		var resp map[string]any
		json.Unmarshal(w.Body.Bytes(), &resp)
		assert.Equal(t, "Alice", resp["name"])
		assert.Equal(t, "alice@example.com", resp["email"])
		userSvc.AssertExpectations(t)
	})

	t.Run("InvalidBody", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/users", bytes.NewBufferString("not-json"))
		req.Header.Set("Content-Type", "application/json")
		newUserRouter(usermocks.NewUserServiceMock(t), addrmocks.NewAddressServiceMock(t)).ServeHTTP(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("MissingRequiredField", func(t *testing.T) {
		// missing email (required)
		body := toJSON(map[string]string{"name": "Alice"})
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/users", body)
		req.Header.Set("Content-Type", "application/json")
		newUserRouter(usermocks.NewUserServiceMock(t), addrmocks.NewAddressServiceMock(t)).ServeHTTP(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("ServiceError", func(t *testing.T) {
		userSvc := usermocks.NewUserServiceMock(t)
		addrSvc := addrmocks.NewAddressServiceMock(t)
		userSvc.On("CreateUser", mock.AnythingOfType("*entity.User")).Return(errors.New("duplicate email")).Once()
		body := toJSON(map[string]string{"name": "Alice", "email": "alice@example.com"})
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/users", body)
		req.Header.Set("Content-Type", "application/json")
		newUserRouter(userSvc, addrSvc).ServeHTTP(w, req)
		assert.Equal(t, http.StatusInternalServerError, w.Code)
		userSvc.AssertExpectations(t)
	})
}

// ─── UpdateUser ──────────────────────────────────────────────────────────────

func TestHandlerUpdateUser(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		userSvc := usermocks.NewUserServiceMock(t)
		addrSvc := addrmocks.NewAddressServiceMock(t)
		userSvc.On("UpdateUser", mock.AnythingOfType("*entity.User")).Return(nil).Once()
		body := toJSON(map[string]string{"name": "Bob", "email": "bob@example.com", "phone": "555"})
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPut, "/users/1", body)
		req.Header.Set("Content-Type", "application/json")
		newUserRouter(userSvc, addrSvc).ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var resp map[string]any
		json.Unmarshal(w.Body.Bytes(), &resp)
		assert.Equal(t, "Bob", resp["name"])
		userSvc.AssertExpectations(t)
	})

	t.Run("InvalidID", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPut, "/users/xyz", toJSON(map[string]string{"name": "Bob", "email": "bob@example.com"}))
		req.Header.Set("Content-Type", "application/json")
		newUserRouter(usermocks.NewUserServiceMock(t), addrmocks.NewAddressServiceMock(t)).ServeHTTP(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("MissingRequiredField", func(t *testing.T) {
		body := toJSON(map[string]string{"name": "Bob"}) // missing email
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPut, "/users/1", body)
		req.Header.Set("Content-Type", "application/json")
		newUserRouter(usermocks.NewUserServiceMock(t), addrmocks.NewAddressServiceMock(t)).ServeHTTP(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("ServiceError", func(t *testing.T) {
		userSvc := usermocks.NewUserServiceMock(t)
		addrSvc := addrmocks.NewAddressServiceMock(t)
		userSvc.On("UpdateUser", mock.AnythingOfType("*entity.User")).Return(errors.New("user not found")).Once()
		body := toJSON(map[string]string{"name": "Bob", "email": "bob@example.com"})
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPut, "/users/1", body)
		req.Header.Set("Content-Type", "application/json")
		newUserRouter(userSvc, addrSvc).ServeHTTP(w, req)
		assert.Equal(t, http.StatusInternalServerError, w.Code)
		userSvc.AssertExpectations(t)
	})
}

// ─── DeleteUser ──────────────────────────────────────────────────────────────

func TestHandlerDeleteUser(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		userSvc := usermocks.NewUserServiceMock(t)
		addrSvc := addrmocks.NewAddressServiceMock(t)
		userSvc.On("DeleteUser", 1).Return(nil).Once()
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodDelete, "/users/1", nil)
		newUserRouter(userSvc, addrSvc).ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var resp map[string]any
		json.Unmarshal(w.Body.Bytes(), &resp)
		assert.Equal(t, "user deleted successfully", resp["message"])
		userSvc.AssertExpectations(t)
	})

	t.Run("InvalidID", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodDelete, "/users/abc", nil)
		newUserRouter(usermocks.NewUserServiceMock(t), addrmocks.NewAddressServiceMock(t)).ServeHTTP(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("ServiceError", func(t *testing.T) {
		userSvc := usermocks.NewUserServiceMock(t)
		addrSvc := addrmocks.NewAddressServiceMock(t)
		userSvc.On("DeleteUser", 1).Return(errors.New("user not found")).Once()
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodDelete, "/users/1", nil)
		newUserRouter(userSvc, addrSvc).ServeHTTP(w, req)
		assert.Equal(t, http.StatusInternalServerError, w.Code)
		userSvc.AssertExpectations(t)
	})
}

// ─── ListUsers ───────────────────────────────────────────────────────────────

func TestHandlerListUsers(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		userSvc := usermocks.NewUserServiceMock(t)
		addrSvc := addrmocks.NewAddressServiceMock(t)
		userSvc.On("ListUsers").Return([]*entity.User{{ID: 1, Name: "Alice", Email: "alice@example.com"}, {ID: 2, Name: "Bob", Email: "bob@example.com"}}, nil).Once()
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/users", nil)
		newUserRouter(userSvc, addrSvc).ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var resp map[string]any
		json.Unmarshal(w.Body.Bytes(), &resp)
		data := resp["data"].([]any)
		assert.Len(t, data, 2)
		userSvc.AssertExpectations(t)
	})

	t.Run("EmptyList", func(t *testing.T) {
		userSvc := usermocks.NewUserServiceMock(t)
		addrSvc := addrmocks.NewAddressServiceMock(t)
		userSvc.On("ListUsers").Return([]*entity.User{}, nil).Once()
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/users", nil)
		newUserRouter(userSvc, addrSvc).ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var resp map[string]any
		json.Unmarshal(w.Body.Bytes(), &resp)
		data := resp["data"].([]any)
		assert.Len(t, data, 0)
		userSvc.AssertExpectations(t)
	})

	t.Run("ServiceError", func(t *testing.T) {
		userSvc := usermocks.NewUserServiceMock(t)
		addrSvc := addrmocks.NewAddressServiceMock(t)
		userSvc.On("ListUsers").Return(([]*entity.User)(nil), errors.New("db error")).Once()
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/users", nil)
		newUserRouter(userSvc, addrSvc).ServeHTTP(w, req)
		assert.Equal(t, http.StatusInternalServerError, w.Code)
		userSvc.AssertExpectations(t)
	})
}
