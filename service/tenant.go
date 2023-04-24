package service

import (
	. "naas/domain"
	. "naas/repositories"
)

type TenantService struct {
	repo *TenantRepository
}

func NewTenantService(repo *TenantRepository) *TenantService {
	return &TenantService{repo: repo}
}

func (s *TenantService) CreateTenant(tenant *Tenant) error {
	return s.repo.CreateTenant(tenant)
}

func (s *TenantService) GetTenant(id string) (*Tenant, error) {
	return s.repo.GetTenant(id)
}

func (s *TenantService) ListTenants() ([]Tenant, error) {
	return s.repo.ListTenants()
}
