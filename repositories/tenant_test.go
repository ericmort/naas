// repositories/tenant_test.go

package repositories_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"naas/domain"
	"naas/repositories"
)

func TestTenantRepository_CreateTenant(t *testing.T) {
	repo := repositories.NewTenantRepository()

	tenant := &domain.Tenant{
		ID:   "test-tenant",
		Name: "Test Tenant",
	}

	err := repo.CreateTenant(tenant)
	assert.NoError(t, err)

	err = repo.CreateTenant(tenant)
	assert.EqualError(t, err, "tenant already exists")
}

func TestTenantRepository_GetTenant(t *testing.T) {
	repo := repositories.NewTenantRepository()

	tenant := &domain.Tenant{
		ID:   "test-tenant",
		Name: "Test Tenant",
	}
	repo.CreateTenant(tenant)

	result, err := repo.GetTenant("test-tenant")
	assert.NoError(t, err)
	assert.Equal(t, tenant, result)

	result, err = repo.GetTenant("non-existent-tenant")
	assert.Nil(t, result)
	assert.EqualError(t, err, "tenant not found")
}
