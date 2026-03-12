package address

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
	mocks "github.com/example/clean-architecture/internal/address/mocks"
)

func init() {
	gin.SetMode(gin.TestMode)
}

// ─── Helpers ─────────────────────────────────────────────────────────────────

func newAddressRouter(svc IService) *gin.Engine {
	r := gin.New()
	h := NewHandler(svc)
	r.GET("/addresses", h.ListAddresses)
	r.POST("/addresses", h.CreateAddress)
	r.GET("/addresses/:id", h.GetAddress)
	r.PUT("/addresses/:id", h.UpdateAddress)
	r.DELETE("/addresses/:id", h.DeleteAddress)
	r.GET("/users/:user_id/addresses", h.GetAddressesByUser)
	return r
}

func addrToJSON(v any) *bytes.Buffer {
	b, _ := json.Marshal(v)
	return bytes.NewBuffer(b)
}

// ─── GetAddress ──────────────────────────────────────────────────────────────

func TestHandlerGetAddress(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		svc := mocks.NewAddressServiceMock(t)
		svc.On("GetAddress", 1).Return(&entity.Address{ID: 1, UserID: 1, Street: "123 Main", City: "Bangkok", Country: "TH"}, nil).Once()
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/addresses/1", nil)
		newAddressRouter(svc).ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var resp map[string]any
		json.Unmarshal(w.Body.Bytes(), &resp)
		assert.Equal(t, float64(1), resp["id"])
		assert.Equal(t, "123 Main", resp["street"])
		assert.Equal(t, "Bangkok", resp["city"])
		svc.AssertExpectations(t)
	})

	t.Run("InvalidID", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/addresses/abc", nil)
		newAddressRouter(mocks.NewAddressServiceMock(t)).ServeHTTP(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("NotFound", func(t *testing.T) {
		svc := mocks.NewAddressServiceMock(t)
		svc.On("GetAddress", 99).Return((*entity.Address)(nil), errors.New("address not found")).Once()
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/addresses/99", nil)
		newAddressRouter(svc).ServeHTTP(w, req)
		assert.Equal(t, http.StatusNotFound, w.Code)
		svc.AssertExpectations(t)
	})
}

// ─── GetAddressesByUser ──────────────────────────────────────────────────────

func TestHandlerGetAddressesByUser(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		svc := mocks.NewAddressServiceMock(t)
		svc.On("GetAddressesByUser", 5).Return([]*entity.Address{{ID: 1, UserID: 5, Street: "A", City: "Bangkok", Country: "TH"}, {ID: 2, UserID: 5, Street: "B", City: "Phuket", Country: "TH"}}, nil).Once()
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/users/5/addresses", nil)
		newAddressRouter(svc).ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var resp map[string]any
		json.Unmarshal(w.Body.Bytes(), &resp)
		data := resp["data"].([]any)
		assert.Len(t, data, 2)
		svc.AssertExpectations(t)
	})

	t.Run("InvalidUserID", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/users/abc/addresses", nil)
		newAddressRouter(mocks.NewAddressServiceMock(t)).ServeHTTP(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("ServiceError", func(t *testing.T) {
		svc := mocks.NewAddressServiceMock(t)
		svc.On("GetAddressesByUser", 1).Return(([]*entity.Address)(nil), errors.New("db error")).Once()
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/users/1/addresses", nil)
		newAddressRouter(svc).ServeHTTP(w, req)
		assert.Equal(t, http.StatusInternalServerError, w.Code)
		svc.AssertExpectations(t)
	})
}

// ─── CreateAddress ───────────────────────────────────────────────────────────

func TestHandlerCreateAddress(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		svc := mocks.NewAddressServiceMock(t)
		svc.On("CreateAddress", mock.AnythingOfType("*entity.Address")).Return(nil).Once()
		body := addrToJSON(map[string]any{
			"user_id": 1, "street": "123 Main St", "city": "Bangkok", "country": "TH",
		})
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/addresses", body)
		req.Header.Set("Content-Type", "application/json")
		newAddressRouter(svc).ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
		var resp map[string]any
		json.Unmarshal(w.Body.Bytes(), &resp)
		assert.Equal(t, "123 Main St", resp["street"])
		assert.Equal(t, "Bangkok", resp["city"])
		svc.AssertExpectations(t)
	})

	t.Run("InvalidBody", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/addresses", bytes.NewBufferString("not-json"))
		req.Header.Set("Content-Type", "application/json")
		newAddressRouter(mocks.NewAddressServiceMock(t)).ServeHTTP(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("MissingRequiredField", func(t *testing.T) {
		// missing country (required) and street
		body := addrToJSON(map[string]any{"user_id": 1, "city": "Bangkok"})
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/addresses", body)
		req.Header.Set("Content-Type", "application/json")
		newAddressRouter(mocks.NewAddressServiceMock(t)).ServeHTTP(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("ServiceError", func(t *testing.T) {
		svc := mocks.NewAddressServiceMock(t)
		svc.On("CreateAddress", mock.AnythingOfType("*entity.Address")).Return(errors.New("insert failed")).Once()
		body := addrToJSON(map[string]any{
			"user_id": 1, "street": "123 Main", "city": "Bangkok", "country": "TH",
		})
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/addresses", body)
		req.Header.Set("Content-Type", "application/json")
		newAddressRouter(svc).ServeHTTP(w, req)
		assert.Equal(t, http.StatusInternalServerError, w.Code)
		svc.AssertExpectations(t)
	})
}

// ─── UpdateAddress ───────────────────────────────────────────────────────────

func TestHandlerUpdateAddress(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		svc := mocks.NewAddressServiceMock(t)
		svc.On("UpdateAddress", mock.AnythingOfType("*entity.Address")).Return(nil).Once()
		body := addrToJSON(map[string]any{
			"street": "New Street", "city": "Bangkok", "country": "TH",
		})
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPut, "/addresses/2", body)
		req.Header.Set("Content-Type", "application/json")
		newAddressRouter(svc).ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var resp map[string]any
		json.Unmarshal(w.Body.Bytes(), &resp)
		assert.Equal(t, "New Street", resp["street"])
		svc.AssertExpectations(t)
	})

	t.Run("InvalidID", func(t *testing.T) {
		w := httptest.NewRecorder()
		body := addrToJSON(map[string]any{"street": "X", "city": "Y", "country": "Z"})
		req, _ := http.NewRequest(http.MethodPut, "/addresses/xyz", body)
		req.Header.Set("Content-Type", "application/json")
		newAddressRouter(mocks.NewAddressServiceMock(t)).ServeHTTP(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("MissingRequiredField", func(t *testing.T) {
		body := addrToJSON(map[string]any{"street": "X"}) // missing city & country
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPut, "/addresses/1", body)
		req.Header.Set("Content-Type", "application/json")
		newAddressRouter(mocks.NewAddressServiceMock(t)).ServeHTTP(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("ServiceError", func(t *testing.T) {
		svc := mocks.NewAddressServiceMock(t)
		svc.On("UpdateAddress", mock.AnythingOfType("*entity.Address")).Return(errors.New("address not found")).Once()
		body := addrToJSON(map[string]any{"street": "X", "city": "Y", "country": "Z"})
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPut, "/addresses/1", body)
		req.Header.Set("Content-Type", "application/json")
		newAddressRouter(svc).ServeHTTP(w, req)
		assert.Equal(t, http.StatusInternalServerError, w.Code)
		svc.AssertExpectations(t)
	})
}

// ─── DeleteAddress ───────────────────────────────────────────────────────────

func TestHandlerDeleteAddress(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		svc := mocks.NewAddressServiceMock(t)
		svc.On("DeleteAddress", 3).Return(nil).Once()
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodDelete, "/addresses/3", nil)
		newAddressRouter(svc).ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var resp map[string]any
		json.Unmarshal(w.Body.Bytes(), &resp)
		assert.Equal(t, "address deleted successfully", resp["message"])
		svc.AssertExpectations(t)
	})

	t.Run("InvalidID", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodDelete, "/addresses/abc", nil)
		newAddressRouter(mocks.NewAddressServiceMock(t)).ServeHTTP(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("ServiceError", func(t *testing.T) {
		svc := mocks.NewAddressServiceMock(t)
		svc.On("DeleteAddress", 1).Return(errors.New("address not found")).Once()
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodDelete, "/addresses/1", nil)
		newAddressRouter(svc).ServeHTTP(w, req)
		assert.Equal(t, http.StatusInternalServerError, w.Code)
		svc.AssertExpectations(t)
	})
}

// ─── ListAddresses ───────────────────────────────────────────────────────────

func TestHandlerListAddresses(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		svc := mocks.NewAddressServiceMock(t)
		svc.On("ListAddresses").Return([]*entity.Address{{ID: 1, UserID: 1, Street: "A", City: "Bangkok", Country: "TH"}, {ID: 2, UserID: 2, Street: "B", City: "Phuket", Country: "TH"}}, nil).Once()
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/addresses", nil)
		newAddressRouter(svc).ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var resp map[string]any
		json.Unmarshal(w.Body.Bytes(), &resp)
		data := resp["data"].([]any)
		assert.Len(t, data, 2)
		svc.AssertExpectations(t)
	})

	t.Run("EmptyList", func(t *testing.T) {
		svc := mocks.NewAddressServiceMock(t)
		svc.On("ListAddresses").Return([]*entity.Address{}, nil).Once()
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/addresses", nil)
		newAddressRouter(svc).ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var resp map[string]any
		json.Unmarshal(w.Body.Bytes(), &resp)
		data := resp["data"].([]any)
		assert.Len(t, data, 0)
		svc.AssertExpectations(t)
	})

	t.Run("ServiceError", func(t *testing.T) {
		svc := mocks.NewAddressServiceMock(t)
		svc.On("ListAddresses").Return(([]*entity.Address)(nil), errors.New("db error")).Once()
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/addresses", nil)
		newAddressRouter(svc).ServeHTTP(w, req)
		assert.Equal(t, http.StatusInternalServerError, w.Code)
		svc.AssertExpectations(t)
	})
}
