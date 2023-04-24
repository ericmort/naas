package service

import (
	. "naas/domain"
	. "naas/repositories"
)

type NamespaceService struct {
	repo *NamespaceRepository
}

func NewNamespaceService(repo *NamespaceRepository) *NamespaceService {
	return &NamespaceService{repo: repo}
}

func (s *NamespaceService) CreateNamespace(tenantID string, namespace *Namespace) error {
	return s.repo.CreateNamespace(tenantID, namespace)
}

func (s *NamespaceService) GetAllNamespaces(tenantID string) ([]Namespace, error) {
	return s.repo.GetAllNamespaces(tenantID)
}

func (s *NamespaceService) GetNamespace(tenantID string, name string) (*Namespace, error) {
	return s.repo.GetNamespace(tenantID, name)
}
