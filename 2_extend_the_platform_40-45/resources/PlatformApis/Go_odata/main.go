// Package main provides the entry point for the Go OData client application.
// This application connects to Microsoft Dataverse using OAuth authentication
// and allows users to interact with Dataverse entities through the OData
// protocol.
package main

import (
	"log"
	"net/url"
	"os"

	goDotEnv "github.com/joho/godotenv"
	"github.com/turnerbenjamin/go_odata/app"
)

// Constants used throughout the application.
const (
	// unableToParseUrlMsgFmt is the error message format for URL parsing
	// failures.
	unableToParseUrlMsgFmt = "unable to parse url (%s)"
	// maxPageLimit defines the maximum number of records to retrieve per page.
	maxPageLimit = 5
)

// main initializes and runs the application.
// It loads configuration from environment variables, sets up URL paths for
// the API and authentication endpoints, and initializes the application with
// these settings.
func main() {

	// Load environment variables from .env file
	err := goDotEnv.Load()
	if err != nil {
		log.Fatal("unable to load .env file")
	}

	// Extract required configuration from environment variables
	clientID := os.Getenv("CLIENT_ID")
	tenantID := os.Getenv("TENANT_ID")
	clientSecret := os.Getenv("CLIENT_SECRET")
	apiPath := os.Getenv("API_PATH")

	// Parse the environment URL (Dataverse instance)
	environmentURL, err := url.Parse(os.Getenv("ENVIRONMENT_URL"))
	if err != nil {
		log.Fatalf(unableToParseUrlMsgFmt, os.Getenv("ENVIRONMENT_URL"))
	}

	// Create API base URL by combining environment URL with API path
	apiBaseURL := *environmentURL
	apiBaseURL.Path = apiPath

	// Parse the authority URL (Entra Id endpoint)
	authorityURL, err := url.Parse(os.Getenv("AUTHORITY"))
	if err != nil {
		log.Fatalf(unableToParseUrlMsgFmt, os.Getenv("AUTHORITY_URL"))
	}

	// Append tenant ID to authority path for tenant-specific authentication
	authorityURL.Path, err = url.JoinPath(authorityURL.Path, tenantID)
	if err != nil {
		log.Fatal(err.Error())
	}

	// Create application configuration
	c := app.AppConfig{
		ClientID:     clientID,
		TenantID:     tenantID,
		ClientSecret: clientSecret,
		ResourceURL:  environmentURL.String(),
		APIBaseURL:   apiBaseURL.String(),
		Authority:    authorityURL.String(),
		PageLimit:    maxPageLimit,
	}

	// Initialize the application
	a, err := app.NewApp(c)
	if err != nil {
		log.Fatal(err)
	}

	// Run the application
	err = a.Run()
	if err != nil {
		log.Fatalf("The application has experienced a fatal error: %s", err.Error())
	}
}
