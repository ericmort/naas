package service

import (
	"naas/data"
	"naas/domain"
)

type NamespaceService struct {
	repo data.NamespaceRepository
}

func NewNamespaceService(repo data.NamespaceRepository) *NamespaceService {
	return &NamespaceService{
		repo: repo,
	}
}

func (s *NamespaceService) CreateNamespace(name string, tenantID string) (domain.Namespace, error) {
	if err := validateInput(name); err != nil {
		return domain.Namespace{}, err
	}

	namespace := domain.Namespace{
		Name:     name,
		TenantID: tenantID,
	}

	return s.repo.Create(namespace), nil
}

func (s *NamespaceService) GetNamespace(tenantID, namespace string) (domain.Namespace, error) {
	return s.repo.Get(tenantID, namespace)
}

func (s *NamespaceService) UpdateNamespace(namespace domain.Namespace) (domain.Namespace, error) {
	if err := validateInput(namespace.Name); err != nil {
		return domain.Namespace{}, err
	}

	return s.repo.Update(namespace)
}

func (s *NamespaceService) DeleteNamespace(id string) error {
	return s.repo.Delete(id)
}

func (s *NamespaceService) GetNamespacesForTenant(tenantID string) ([]domain.Namespace, error) {
	return s.repo.GetForTenant(tenantID)
}
