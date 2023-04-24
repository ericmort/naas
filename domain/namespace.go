package domain

type Namespace struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	TenantID  string `json:"tenantId"`
	CreatedAt int64  `json:"createdAt"`
	UpdatedAt int64  `json:"updatedAt"`
}

func NewNamespace(name, tenantID string) *Namespace {
	return &Namespace{
		Name:     name,
		TenantID: tenantID,
	}
}
