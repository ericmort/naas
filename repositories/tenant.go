package repositories

import (
	"errors"
	"sync"

	"naas/domain"
)

type TenantRepository struct {
	mtx     sync.RWMutex
	tenants map[string]domain.Tenant
}

func NewTenantRepository() *TenantRepository {
	return &TenantRepository{tenants: make(map[string]domain.Tenant)}
}

func (r *TenantRepository) CreateTenant(tenant *domain.Tenant) error {
	r.mtx.Lock()
	defer r.mtx.Unlock()

	if _, ok := r.tenants[tenant.ID]; ok {
		return errors.New("tenant already exists")
	}

	r.tenants[tenant.ID] = *tenant
	return nil
}

func (r *TenantRepository) GetTenant(id string) (*domain.Tenant, error) {
	r.mtx.RLock()
	defer r.mtx.RUnlock()

	if tenant, ok := r.tenants[id]; ok {
		return &tenant, nil
	}

	return nil, errors.New("tenant not found")
}

func (r *TenantRepository) ListTenants() ([]domain.Tenant, error) {
	tenants := make([]domain.Tenant, 0, len(r.tenants))
	for _, tenant := range r.tenants {
		tenants = append(tenants, tenant)
	}
	return tenants, nil
}
