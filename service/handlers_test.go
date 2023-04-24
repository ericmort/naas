package service_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"naas/domain"
	"naas/repositories"
	"naas/service"
)

func TestTenantService_CreateTenant(t *testing.T) {
	repo := repositories.NewTenantRepository()
	service := service.NewTenantService(repo)

	tenant := &domain.Tenant{
		ID:   "test-tenant",
		Name: "Test Tenant",
	}

	err := service.CreateTenant(tenant)
	assert.NoError(t, err)

	err = service.CreateTenant(tenant)
	assert.EqualError(t, err, "tenant already exists")
}

func TestTenantService_GetTenant(t *testing.T) {
	repo := repositories.NewTenantRepository()
	service := service.NewTenantService(repo)

	tenant := &domain.Tenant{
		ID:   "test-tenant",
		Name: "Test Tenant",
	}
	err := repo.CreateTenant(tenant)
	assert.NoError(t, err)

	result, err := service.GetTenant("test-tenant")
	assert.NoError(t, err)
	assert.Equal(t, tenant, result)

	result, err = service.GetTenant("non-existent-tenant")
	assert.Nil(t, result)
	assert.EqualError(t, err, "tenant not found")
}

func TestNamespaceService_CreateNamespace(t *testing.T) {
	repo := repositories.NewNamespaceRepository()
	service := service.NewNamespaceService(repo)

	namespace := &domain.Namespace{
		Name: "test-namespace",
	}

	err := service.CreateNamespace("test-tenant", namespace)
	assert.NoError(t, err)

	err = service.CreateNamespace("test-tenant", namespace)
	assert.EqualError(t, err, "namespace already exists")
}

func TestNamespaceService_GetAllNamespaces(t *testing.T) {
	repo := repositories.NewNamespaceRepository()
	service := service.NewNamespaceService(repo)

	namespaces := []domain.Namespace{
		{Name: "test-namespace-1"},
		{Name: "test-namespace-2"},
		{Name: "test-namespace-3"},
	}
	for _, ns := range namespaces {
		err := repo.CreateNamespace("test-tenant", &ns)
		assert.NoError(t, err)
	}

	result, err := service.GetAllNamespaces("test-tenant")
	assert.NoError(t, err)
	assert.ElementsMatch(t, namespaces, result)

	result, err = service.GetAllNamespaces("non-existent-tenant")
	assert.Nil(t, result)
	assert.EqualError(t, err, "no namespaces found for tenant")
}

func TestNamespaceService_GetNamespace(t *testing.T) {
	repo := repositories.NewNamespaceRepository()
	service := service.NewNamespaceService(repo)

	namespace := &domain.Namespace{
		Name: "test-namespace",
	}
	err := repo.CreateNamespace("test-tenant", namespace)
	assert.NoError(t, err)

	result, err := service.GetNamespace("test-tenant", "test-namespace")
	assert.NoError(t, err)
	assert.Equal(t, namespace, result)

	result, err = service.GetNamespace("test-tenant", "non-existent-namespace")
	assert.Nil(t, result)
	assert.EqualError(t, err, "namespace not found")

	result, err = service.GetNamespace("non-existent-tenant", "test-namespace")
	assert.Nil(t, result)
	assert.EqualError(t, err, "namespace not found")
}
