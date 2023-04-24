package handlers

import (
	. "naas/domain"
	. "naas/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type NamespaceHandler struct {
	service *NamespaceService
}

func NewNamespaceHandler(service *NamespaceService) *NamespaceHandler {
	return &NamespaceHandler{service: service}
}

func (h *NamespaceHandler) CreateNamespace(c *gin.Context) {
	tenantID := c.Param("tenantId")

	var namespace Namespace
	if err := c.ShouldBindJSON(&namespace); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.CreateNamespace(tenantID, &namespace); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, namespace)
}

func (h *NamespaceHandler) GetAllNamespaces(c *gin.Context) {
	tenantID := c.Param("tenantId")

	namespaces, err := h.service.GetAllNamespaces(tenantID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, namespaces)
}

func (h *NamespaceHandler) GetNamespace(c *gin.Context) {
	tenantID := c.Param("tenantId")
	name := c.Param("name")

	namespace, err := h.service.GetNamespace(tenantID, name)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, namespace)
}
