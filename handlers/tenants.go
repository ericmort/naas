// handlers/tenant.go

package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	. "naas/domain"
	. "naas/service"
)

type TenantHandler struct {
	service *TenantService
}

func NewTenantHandler(service *TenantService) *TenantHandler {
	return &TenantHandler{service: service}
}

func (h *TenantHandler) CreateTenant(c *gin.Context) {
	var tenant Tenant
	if err := c.ShouldBindJSON(&tenant); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.CreateTenant(&tenant); err != nil {
		if err.Error() == "tenant already exists" {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}

		return
	}

	c.JSON(http.StatusCreated, tenant)
}

func (h *TenantHandler) GetTenant(c *gin.Context) {
	id := c.Param("id")

	tenant, err := h.service.GetTenant(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, tenant)
}

func (h *TenantHandler) ListTenants(c *gin.Context) {
	tenants, err := h.service.ListTenants()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, tenants)
}
