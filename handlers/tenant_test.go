package handlers_test

import (
	"bytes"
	"encoding/json"
	"io"
	"naas/service"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"naas/domain"
	"naas/handlers"
	"naas/repositories"
)

func TestTenantHandler_CreateTenant(t *testing.T) {
	repo := repositories.NewTenantRepository()
	service := service.NewTenantService(repo)
	handler := handlers.NewTenantHandler(service)

	router := gin.Default()
	router.POST("/tenants", handler.CreateTenant)

	tenant := &domain.Tenant{
		ID:   "test-tenant",
		Name: "Test Tenant",
	}
	payload, err := json.Marshal(tenant)
	assert.NoError(t, err)

	req, err := http.NewRequest(http.MethodPost, "/tenants", bytes.NewBuffer(payload))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Equal(t, "application/json; charset=utf-8", w.Header().Get("Content-Type"))

	response := &domain.Tenant{}
	err = json.Unmarshal(w.Body.Bytes(), response)
	assert.NoError(t, err)
	assert.Equal(t, tenant, response)

	// Test creating the same tenant twice
	w = httptest.NewRecorder()
	req, err = http.NewRequest(http.MethodPost, "/tenants", bytes.NewBuffer(payload))
	assert.NoError(t, err)

	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusConflict, w.Code)
	assert.Equal(t, "application/json; charset=utf-8", w.Header().Get("Content-Type"))
	assert.Contains(t, w.Body.String(), "tenant already exists")
}

func TestTenantHandler_GetTenant(t *testing.T) {
	repo := repositories.NewTenantRepository()
	service := service.NewTenantService(repo)
	handler := handlers.NewTenantHandler(service)

	router := gin.Default()
	router.GET("/tenants/:id", handler.GetTenant)

	tenant := &domain.Tenant{
		ID:   "test-tenant",
		Name: "Test Tenant",
	}
	repo.CreateTenant(tenant)

	req, err := http.NewRequest(http.MethodGet, "/tenants/test-tenant", nil)
	assert.NoError(t, err)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "application/json; charset=utf-8", w.Header().Get("Content-Type"))

	response := &domain.Tenant{}
	err = json.Unmarshal(w.Body.Bytes(), response)
	assert.NoError(t, err)
	assert.Equal(t, tenant, response)

	// Test getting a non-existent tenant
	req, err = http.NewRequest(http.MethodGet, "/tenants/non-existent-tenant", nil)
	assert.NoError(t, err)

	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Equal(t, "application/json; charset=utf-8", w.Header().Get("Content-Type"))
	assert.Contains(t, w.Body.String(), "tenant not found")
}

func TestTenantHandler_ListTenants(t *testing.T) {
	repo := repositories.NewTenantRepository()
	service := service.NewTenantService(repo)
	handler := handlers.NewTenantHandler(service)

	// Define some mock tenants
	tenant1 := domain.Tenant{ID: "1", Name: "Tenant 1"}
	tenant2 := domain.Tenant{ID: "2", Name: "Tenant 2"}
	tenant3 := domain.Tenant{ID: "3", Name: "Tenant 3"}
	tenants := []domain.Tenant{tenant1, tenant2, tenant3}
	for _, tenant := range tenants {
		repo.CreateTenant(&tenant)
	}
	// Mock the ListTenants method of the repository to return the mock tenants
	req, err := http.NewRequest(http.MethodGet, "/tenants", nil)
	assert.NoError(t, err)

	w := httptest.NewRecorder()
	router := gin.Default()
	router.GET("/tenants", handler.ListTenants)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "application/json; charset=utf-8", w.Header().Get("Content-Type"))

	response, err := io.ReadAll(w.Body)
	assert.NoError(t, err)
	// Verify that the response matches the expected output
	expectedOutput := `[{"id":"1","name":"Tenant 1"},{"id":"2","name":"Tenant 2"},{"id":"3","name":"Tenant 3"}]`
	assert.Equal(t, expectedOutput, string(response))

}
