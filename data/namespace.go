package data

import (
	"errors"
	"sync"

	"naas/domain"
)

var (
	ErrNamespaceNotFound = errors.New("namespace not found")
)

type NamespaceRepository interface {
	Create(namespace domain.Namespace) domain.Namespace
	Get(tenantID, name string) (domain.Namespace, error)
	Update(namespace domain.Namespace) (domain.Namespace, error)
	Delete(tenantID, name string) error
}

type InMemoryNamespaceRepository struct {
	m     map[string]map[string]domain.Namespace // tenantID -> namespaceName -> namespace
	mu    sync.RWMutex
	idGen func() string
}

func NewInMemoryNamespaceRepository(idGenerator func() string) *InMemoryNamespaceRepository {
	return &InMemoryNamespaceRepository{
		m:     make(map[string]map[string]domain.Namespace),
		idGen: idGenerator,
	}
}

func (r *InMemoryNamespaceRepository) Create(namespace domain.Namespace) domain.Namespace {
	namespace.ID = r.idGen()

	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.m[namespace.TenantID]; !ok {
		r.m[namespace.TenantID] = make(map[string]domain.Namespace)
	}
	r.m[namespace.TenantID][namespace.Name] = namespace

	return namespace
}

func (r *InMemoryNamespaceRepository) Get(tenantID, name string) (domain.Namespace, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if ns, ok := r.m[tenantID][name]; ok {
		return ns, nil
	}

	return domain.Namespace{}, ErrNamespaceNotFound
}

func (r *InMemoryNamespaceRepository) Update(namespace domain.Namespace) (domain.Namespace, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.m[namespace.TenantID][namespace.Name]; !ok {
		return domain.Namespace{}, ErrNamespaceNotFound
	}
	r.m[namespace.TenantID][namespace.Name] = namespace

	return namespace, nil
}

func (r *InMemoryNamespaceRepository) Delete(tenantID, name string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.m[tenantID][name]; !ok {
		return ErrNamespaceNotFound
	}
	delete(r.m[tenantID], name)

	return nil
}

func (r *InMemoryNamespaceRepository) GetForTenant(tenantID string) (map[string]domain.Namespace, error) {
	_ns := r.m[tenantID]
	if _ns == nil {
		return nil, errors.New("tenant ID not found")
	}
	return _ns, nil
}
