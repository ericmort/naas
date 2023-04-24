package api

import (
	"github.com/gin-gonic/gin"
	"naas/data"
	"naas/domain"
	"net/http"
)

type NamespaceHandler struct {
	repo      data.NamespaceRepository
	tenantSvc domain.TenantService
}

func NewNamespaceHandler(repo data.NamespaceRepository, tenantSvc domain.TenantService) *NamespaceHandler {
	return &NamespaceHandler{
		repo:      repo,
		tenantSvc: tenantSvc,
	}
}

func (h *NamespaceHandler) Create(c *gin.Context) {
	tenantID := c.Param("tenant_id")

	// Check if tenant exists
	_, err := h.tenantSvc.Get(tenantID)
	if err != nil {
		if err == data.ErrNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "tenant not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	input := struct {
		Name string `json:"name" `
	}{}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := validateInput(input.Name); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	namespace := domain.Namespace{
		Name:     input.Name,
		TenantID: tenantID,
	}

	namespace = h.repo.Create(namespace)

	c.JSON(http.StatusCreated, namespace)
}

func (h *NamespaceHandler) Get(c *gin.Context) {
	id := c.Param("id")

	namespace, err := h.repo.Get(id)
	if err != nil {
		if err == data.ErrNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "namespace not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, namespace)
}

func (h *NamespaceHandler) Update(c *gin.Context) {
	id := c.Query("id")

	var input domain.Namespace
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	input.ID = id

	namespace, err := h.repo.Update(input)
	if err != nil {
		if err == data.ErrNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "namespace not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, namespace)
}

func (h *NamespaceHandler) Delete(c *gin.Context) {
	id := c.Query("id")

	if err := h.repo.Delete(id); err != nil {
		if err == data.ErrNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "namespace not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "namespace deleted"})
}

func (h *NamespaceHandler) SetupRoutes() *gin.Engine {
	router := gin.Default()

	router.POST("/tenants/:tenant_id/namespaces", h.Create)
	router.GET("/namespaces/:id", h.Get)
	router.PUT("/namespaces", h.Update)
	router.DELETE("/namespaces", h.Delete)

	return router
}
