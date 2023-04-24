package repositories_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"naas/domain"
	"naas/repositories"
)

func TestNamespaceRepository_CreateNamespace(t *testing.T) {
	repo := repositories.NewNamespaceRepository()

	namespace := &domain.Namespace{
		Name: "test-namespace",
	}

	err := repo.CreateNamespace("test-tenant", namespace)
	assert.NoError(t, err)

	err = repo.CreateNamespace("test-tenant", namespace)
	assert.EqualError(t, err, "namespace already exists")
}

func TestNamespaceRepository_GetAllNamespaces(t *testing.T) {
	repo := repositories.NewNamespaceRepository()

	namespaces := []domain.Namespace{
		{Name: "test-namespace-1"},
		{Name: "test-namespace-2"},
		{Name: "test-namespace-3"},
	}
	for _, ns := range namespaces {
		repo.CreateNamespace("test-tenant", &ns)
	}

	result, err := repo.GetAllNamespaces("test-tenant")
	assert.NoError(t, err)
	assert.ElementsMatch(t, namespaces, result)

	result, err = repo.GetAllNamespaces("non-existent-tenant")
	assert.Nil(t, result)
	assert.EqualError(t, err, "no namespaces found for tenant")
}

func TestNamespaceRepository_GetNamespace(t *testing.T) {
	repo := repositories.NewNamespaceRepository()

	namespace := &domain.Namespace{
		Name: "test-namespace",
	}
	repo.CreateNamespace("test-tenant", namespace)

	result, err := repo.GetNamespace("test-tenant", "test-namespace")
	assert.NoError(t, err)
	assert.Equal(t, namespace, result)

	result, err = repo.GetNamespace("test-tenant", "non-existent-namespace")
	assert.Nil(t, result)
	assert.EqualError(t, err, "namespace not found")

	result, err = repo.GetNamespace("non-existent-tenant", "test-namespace")
	assert.Nil(t, result)
	assert.EqualError(t, err, "namespace not found")
}
