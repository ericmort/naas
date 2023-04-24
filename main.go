// main.go

package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"naas/handlers"
	"naas/repositories"
	"naas/service"
	"os"
)

func main() {
	// Initialize repositories
	tenantRepo := repositories.NewTenantRepository()
	namespaceRepo := repositories.NewNamespaceRepository()

	// Initialize services
	tenantService := service.NewTenantService(tenantRepo)
	namespaceService := service.NewNamespaceService(namespaceRepo)

	// Initialize handlers
	tenantHandler := handlers.NewTenantHandler(tenantService)
	namespaceHandler := handlers.NewNamespaceHandler(namespaceService)

	// Initialize Gin router
	router := gin.Default()

	// Configure CORS middleware
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true // Allow all origins for development
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Accept", "Authorization"}

	// Apply CORS middleware to your Gin instance
	router.Use(cors.New(config))

	// Define routes
	router.POST("/tenants", tenantHandler.CreateTenant)
	router.GET("/tenants", tenantHandler.ListTenants)
	router.GET("/tenants/:id", tenantHandler.GetTenant)
	router.POST("/namespaces/:tenantId", namespaceHandler.CreateNamespace)
	router.GET("/namespaces/all/:tenantId", namespaceHandler.GetAllNamespaces)
	router.GET("/namespaces/:tenantId/:name", namespaceHandler.GetNamespace)

	// Start server
	err := router.Run(":8082")
	if err != nil {
		os.Exit(1)
	}
}
