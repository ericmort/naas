package api

import (
	"errors"
	"naas/service"
	"net/http"
	"strings"
	"unicode"

	"github.com/gin-gonic/gin"
	"naas/data"
	"naas/domain"
)

type TenantHandler struct {
	service service.TenantService
}

func NewTenantHandler(s service.TenantService) *TenantHandler {
	return &TenantHandler{
		service: s,
	}
}

func (h *TenantHandler) Create(c *gin.Context) {
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
	tenant, err := h.service.CreateTenant(input.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, tenant)
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

func validateInput(name string) error {
	if len(strings.TrimSpace(name)) == 0 {
		return errors.New("name cannot be empty")
	}

	// Validate that Name is a string
	if !isValidVariableName(name) {
		return errors.New("name must be a valid variable name and cannot contain a string")
	}

	return nil
}

func (h *TenantHandler) Get(c *gin.Context) {
	id := c.Param("id")

	tenant, err := h.repo.Get(id)
	if err != nil {
		if err == data.ErrNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "tenant not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, tenant)
}

func (h *TenantHandler) Update(c *gin.Context) {
	id := c.Query("id")

	var input domain.Tenant
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	input.ID = id

	tenant, err := h.repo.Update(input)
	if err != nil {
		if err == data.ErrNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "tenant not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, tenant)
}

func (h *TenantHandler) Delete(c *gin.Context) {
	id := c.Query("id")

	if err := h.repo.Delete(id); err != nil {
		if err == data.ErrNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "tenant not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "tenant deleted"})
}

func (h *TenantHandler) SetupRoutes() *gin.Engine {
	router := gin.Default()

	router.POST("/tenants", h.Create)
	router.GET("/tenants/:id", h.Get)
	router.PUT("/tenants", h.Update)
	router.DELETE("/tenants", h.Delete)

	return router
}
