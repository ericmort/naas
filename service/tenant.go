package service

import (
	"errors"
	"unicode"

	"naas/data"
	"naas/domain"
)

type TenantService struct {
	repo data.TenantRepository
}

func NewTenantService(repo data.TenantRepository) *TenantService {
	return &TenantService{
		repo: repo,
	}
}

func (s *TenantService) CreateTenant(name string) (domain.Tenant, error) {
	if err := validateInput(name); err != nil {
		return domain.Tenant{}, err
	}

	tenant := domain.Tenant{
		Name: name,
	}

	return s.repo.Create(tenant), nil
}

func (s *TenantService) GetTenant(id string) (domain.Tenant, error) {
	return s.repo.Get(id)
}

func (s *TenantService) UpdateTenant(tenant domain.Tenant) (domain.Tenant, error) {
	if err := validateInput(tenant.Name); err != nil {
		return domain.Tenant{}, err
	}

	return s.repo.Update(tenant)
}

func (s *TenantService) DeleteTenant(id string) error {
	return s.repo.Delete(id)
}

func validateInput(name string) error {
	if len(name) == 0 {
		return errors.New("name cannot be empty")
	}

	// Validate that Name is a string
	if !isValidVariableName(name) {
		return errors.New("name must be a valid variable name and cannot contain a string")
	}

	return nil
}

func isValidVariableName(s string) bool {
	if len(s) == 0 || !unicode.IsLetter(rune(s[0])) {
		return false
	}

	for _, c := range s[1:] {
		if !unicode.IsLetter(c) && !unicode.IsDigit(c) {
			return false
		}
	}
	if s == "true" || s == "false" {
		return false
	}

	return true
}
