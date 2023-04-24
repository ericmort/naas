we are going to implement an REST  API with Golang and the Gin web framework. We have two
resources; "tenant" and "namespace". A tenant can have many namespaces, but one namespace
belongs to only one tenant. We have the following endpoints:

POST /tenants (create tenant)
GET /tenants/:id (retrieve tenant with :id)
POST /namespaces/:tenantId (create namespace under tenant with ID :tenantId
GET /namespaces/all/:tenantId (retrieve all namespaces for the tenant with ID :tenantId)
GET /namespaces/:tenantId/:name (retrieve namespace under tenant with id :tenantid with name :name)

We want a three layered architecture: for example, for the "tenant" resource:
- TenantHandler handles the API part
- TenantService handles the business logic, and is used by the TenantHandler
- TenantRepository handles the data logic, and is used by the TenantService 

The same principles work for the "namespace" resource.

The TenantRepository and NamespaceRepository will be implemented using an in-memory map structure.

For the TenantRepository the in-memory data structure should be map[string]domain.Tenant.
In other words, the key is a :tenant ID and the value is the corresponding tenant object.

For the NamespaceREpository the in-memory data structure should be map[string]map[string]domain.Namespace
In other words, the key is the tenant ID and the value is a map of "name" -> "domain.Namespace"

TenantHandler.Create(
    call tenantService.Create
)

TenantService.Create(
    call tenantRepo.Create()
)

TenantRepository.Create(
    add domain.Tenant to map[string]domain.Tenant 
)


NamespaceHandler.Create(
    call namespaceService.Create
)
NamespaceService.Create(
    call namespaceRepo.Create
)
NamespaceRepository.Create(
    add domain.Namespace to map["tenant-id:namespace-name]domain.Namespace
)



TenantHandler.ListNamespaces(
    tenantService.ListNamespaces()
)

TenantService.ListNamespaces(
    tenantRepo.ListNamespaces()
)

TenantRepo.ListNamespaces(
    return map["tenantID:namespace-name]
)

NamespaceHandler.Get(tenantId, namespaceName
    namespaceService.Get("tenantId")
)

NamespaceService.Get(tenantId, namespaceName
    namespaceRepo.Get(tenantId)
)

NamespaceRepo.Get(tenantId, namespaceName)
    return map["tenantid:namespace-name]
)