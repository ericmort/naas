package main

import (
	"github.com/spf13/cobra"
	"log"
	"naas/api"
	"naas/data"
	"os"
	"strconv"
)

var port int

var rootCmd = &cobra.Command{
	Use:   "myapp",
	Short: "MyApp is a CLI tool for managing tenants",
	Long:  `MyApp is a CLI tool for managing tenants in a system`,
}

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Starts the API server",
	Run: func(cmd *cobra.Command, args []string) {
		// Call the function that starts the API server
		// Get the port to run the server on from the "port" flag

		startAPIServer()
	},
}

func startAPIServer() {
	repo := data.NewInMemoryTenantRepository(data.DefaultGenerateID)
	handler := api.NewTenantHandler(repo)

	router := handler.SetupRoutes()

	// Start the server on the specified port
	if err := router.Run(":" + strconv.Itoa(port)); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}

func init() {
	// Add the "serve" command to the root command
	rootCmd.AddCommand(serveCmd)

	// Set any flags for the "serve" command
	serveCmd.Flags().String("port", "8080", "Port to run the API server on")
}

func main() {
	_port, err := serveCmd.Flags().GetString("port")
	if err != nil {
		log.Fatalf("Error getting port: %v", err)
	}

	port, err = strconv.Atoi(_port)
	if err := rootCmd.Execute(); err != nil {
		log.Fatalf("Error executing command: %v", err)
		os.Exit(1)
	}
}
