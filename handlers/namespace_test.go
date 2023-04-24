package handlers_test

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"naas/domain"
	"naas/handlers"
	"naas/repositories"
	"naas/service"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNamespaceHandler_CreateNamespace(t *testing.T) {
	repo := repositories.NewNamespaceRepository()
	service := service.NewNamespaceService(repo)
	handler := handlers.NewNamespaceHandler(service)

	namespace := &domain.Namespace{
		Name: "test-namespace",
	}

	body, err := json.Marshal(namespace)
	assert.NoError(t, err)

	req, err := http.NewRequest(http.MethodPost, "/namespaces/test-tenant", bytes.NewBuffer(body))
	assert.NoError(t, err)

	rec := httptest.NewRecorder()

	router := gin.Default()
	router.POST("/namespaces/:tenantId", handler.CreateNamespace)

	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusCreated, rec.Code)
	assert.Equal(t, "application/json; charset=utf-8", rec.Header().Get("Content-Type"))

	var result domain.Namespace
	err = json.Unmarshal(rec.Body.Bytes(), &result)
	assert.NoError(t, err)

	assert.Equal(t, namespace.Name, result.Name)
}

func TestNamespaceHandler_GetAllNamespaces(t *testing.T) {
	repo := repositories.NewNamespaceRepository()
	service := service.NewNamespaceService(repo)
	handler := handlers.NewNamespaceHandler(service)

	namespaces := []domain.Namespace{
		{Name: "test-namespace-1"},
		{Name: "test-namespace-2"},
		{Name: "test-namespace-3"},
	}
	for _, ns := range namespaces {
		repo.CreateNamespace("test-tenant", &ns)
	}

	req, err := http.NewRequest(http.MethodGet, "/namespaces/all/test-tenant", nil)
	assert.NoError(t, err)

	rec := httptest.NewRecorder()

	router := gin.Default()
	router.GET("/namespaces/all/:tenantId", handler.GetAllNamespaces)

	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "application/json; charset=utf-8", rec.Header().Get("Content-Type"))

	var result []domain.Namespace
	err = json.Unmarshal(rec.Body.Bytes(), &result)
	assert.NoError(t, err)

	assert.ElementsMatch(t, namespaces, result)
}

func TestNamespaceHandler_GetNamespace(t *testing.T) {
	repo := repositories.NewNamespaceRepository()
	service := service.NewNamespaceService(repo)
	handler := handlers.NewNamespaceHandler(service)

	namespace := &domain.Namespace{
		Name: "test-namespace",
	}
	repo.CreateNamespace("test-tenant", namespace)

	req, err := http.NewRequest(http.MethodGet, "/namespaces/test-tenant/test-namespace", nil)
	assert.NoError(t, err)

	rec := httptest.NewRecorder()

	router := gin.Default()
	router.GET("/namespaces/:tenantId/:name", handler.GetNamespace)

	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "application/json; charset=utf-8", rec.Header().Get("Content-Type"))

	var result domain.Namespace
	err = json.Unmarshal(rec.Body.Bytes(), &result)
	assert.NoError(t, err)

	assert.Equal(t, namespace, &result)
}
