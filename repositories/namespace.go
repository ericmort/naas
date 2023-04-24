// repositories/namespace.go

package repositories

import (
	"errors"
	. "naas/domain"
	"sync"
)

type NamespaceRepository struct {
	mtx        sync.RWMutex
	namespaces map[string]map[string]Namespace
}

func NewNamespaceRepository() *NamespaceRepository {
	return &NamespaceRepository{namespaces: make(map[string]map[string]Namespace)}
}

func (r *NamespaceRepository) CreateNamespace(tenantID string, namespace *Namespace) error {
	r.mtx.Lock()
	defer r.mtx.Unlock()

	if _, ok := r.namespaces[tenantID]; !ok {
		r.namespaces[tenantID] = make(map[string]Namespace)
	}

	if _, ok := r.namespaces[tenantID][namespace.Name]; ok {
		return errors.New("namespace already exists")
	}

	r.namespaces[tenantID][namespace.Name] = *namespace
	return nil
}

func (r *NamespaceRepository) GetAllNamespaces(tenantID string) ([]Namespace, error) {
	r.mtx.RLock()
	defer r.mtx.RUnlock()

	if namespaces, ok := r.namespaces[tenantID]; ok {
		result := make([]Namespace, 0, len(namespaces))
		for _, ns := range namespaces {
			result = append(result, ns)
		}
		return result, nil
	}

	return nil, errors.New("no namespaces found for tenant")
}

func (r *NamespaceRepository) GetNamespace(tenantID string, name string) (*Namespace, error) {
	r.mtx.RLock()
	defer r.mtx.RUnlock()

	if namespaces, ok := r.namespaces[tenantID]; ok {
		if namespace, ok := namespaces[name]; ok {
			return &namespace, nil
		}
	}

	return nil, errors.New("namespace not found")
}
