package api_test

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"io"
	"naas/api"
	"naas/data"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"naas/domain"
)

func mockGenerateId() string {
	return "John Doe"
}

func TestCreateTenant(t *testing.T) {
	// Setup
	repo := data.NewInMemoryTenantRepository(mockGenerateId)
	h := api.NewTenantHandler(repo)

	router := gin.Default()
	router.POST("/tenants", h.Create)

	tests := []struct {
		name      string
		inputName string
		expected  interface{}
		status    int
	}{
		{"empty name", "", gin.H{"error": "name cannot be empty"}, http.StatusBadRequest},
		{"numeric name", "123", gin.H{"error": "name must be a valid variable name and cannot contain a string"}, http.StatusBadRequest},
		{"boolean name", "true", gin.H{"error": "name must be a valid variable name and cannot contain a string"}, http.StatusBadRequest},
		{"contains space", "a b", gin.H{"error": "name must be a valid variable name and cannot contain a string"}, http.StatusBadRequest},
		{"valid name", "JohnDoe", "JohnDoe", http.StatusCreated},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input := gin.H{"name": tt.inputName}

			// Perform request
			w := performRequest(router, "POST", "/tenants", input)

			// Verify response
			assert.Equal(t, tt.status, w.Code)

			if tt.status == http.StatusCreated {
				var tenant domain.Tenant
				err := json.Unmarshal(w.Body.Bytes(), &tenant)
				require.NoError(t, err)

				assert.NotEqual(t, tt.expected, tenant.ID)
				assert.Equal(t, tt.inputName, tenant.Name)
			} else {
				var response gin.H
				err := json.Unmarshal(w.Body.Bytes(), &response)
				require.NoError(t, err)

				assert.Equal(t, tt.expected, response)
			}
		})
	}
}

func performRequest(r http.Handler, method, path string, body interface{}) *httptest.ResponseRecorder {
	var reqBody io.Reader
	if body != nil {
		jsonBytes, err := json.Marshal(body)
		if err != nil {
			panic(err)
		}
		reqBody = bytes.NewBuffer(jsonBytes)
	}

	req, err := http.NewRequest(method, path, reqBody)
	if err != nil {
		panic(err)
	}

	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func TestGetTenantByID(t *testing.T) {
	// Define a custom ID generator function that always returns "test-id"
	idGenerator := func() string {
		return "test-id"
	}

	// Create a new in-memory tenant repository using the custom ID generator function
	repo := data.NewInMemoryTenantRepository(idGenerator)

	// Create a new mock tenant and add it to the repository
	tenant := repo.Create(domain.Tenant{Name: "Test Tenant"})

	// setup GET endpoint
	h := api.NewTenantHandler(repo)
	router := gin.Default()
	router.GET("/tenants/:id", h.Get)

	t.Run("Happy Path", func(t *testing.T) {

		// Create a new HTTP request
		req, err := http.NewRequest("GET", "/tenants/test-id", nil)
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}

		// Execute the request
		recorder := httptest.NewRecorder()
		router.ServeHTTP(recorder, req)

		// Check the response code and body
		if recorder.Code != http.StatusOK {
			t.Errorf("Expected status code %d but got %d", http.StatusOK, recorder.Code)
		}

		var response domain.Tenant
		err = json.Unmarshal(recorder.Body.Bytes(), &response)
		if err != nil {
			t.Fatalf("Failed to unmarshal response: %v", err)
		}

		if response.ID != "test-id" {
			t.Errorf("Expected tenant ID %s but got %s", "test-id", response.ID)
		}

		if response.Name != tenant.Name {
			t.Errorf("Expected tenant name %s but got %s", tenant.Name, response.Name)
		}
	})

	t.Run("ID Not Found", func(t *testing.T) {
		// Create a new HTTP request with a non-existent tenant ID
		req, err := http.NewRequest("GET", "/tenants/invalid-id", nil)
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}

		// Create a new router and execute the request
		router := h.SetupRoutes()
		recorder := httptest.NewRecorder()
		router.ServeHTTP(recorder, req)

		// Check the response code
		if recorder.Code != http.StatusNotFound {
			t.Errorf("Expected status code %d but got %d", http.StatusNotFound, recorder.Code)
		}
	})
}
