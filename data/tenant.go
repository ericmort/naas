// data/tenant_repository.go
package data

import (
	"errors"
	"github.com/google/uuid"
	"naas/domain"
	"sync"
)

var ErrNotFound = errors.New("not found")

type TenantRepository interface {
	Create(tenant domain.Tenant) domain.Tenant
	Get(id string) (domain.Tenant, error)
	Update(tenant domain.Tenant) (domain.Tenant, error)
	Delete(id string) error
}

type inMemoryTenantRepository struct {
	tenants     map[string]domain.Tenant
	mu          sync.RWMutex
	idGenerator func() string
}

func NewInMemoryTenantRepository(idGenerator func() string) TenantRepository {
	return &inMemoryTenantRepository{
		tenants:     make(map[string]domain.Tenant),
		idGenerator: idGenerator,
	}
}

// GenerateID generates a unique id for a tenant
func DefaultGenerateID() string {
	// TODO: implement your own id generation algorithm
	return uuid.New().String()
}

func (r *inMemoryTenantRepository) Create(tenant domain.Tenant) domain.Tenant {
	tenant.ID = r.idGenerator()

	r.mu.Lock()
	defer r.mu.Unlock()

	r.tenants[tenant.ID] = tenant

	return tenant
}

func (r *inMemoryTenantRepository) Get(id string) (domain.Tenant, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	tenant, ok := r.tenants[id]
	if !ok {
		return domain.Tenant{}, ErrNotFound
	}

	return tenant, nil
}

func (r *inMemoryTenantRepository) Update(tenant domain.Tenant) (domain.Tenant, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	_, ok := r.tenants[tenant.ID]
	if !ok {
		return domain.Tenant{}, ErrNotFound
	}

	r.tenants[tenant.ID] = tenant

	return tenant, nil
}

func (r *inMemoryTenantRepository) Delete(id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	_, ok := r.tenants[id]
	if !ok {
		return ErrNotFound
	}

	delete(r.tenants, id)

	return nil
}
