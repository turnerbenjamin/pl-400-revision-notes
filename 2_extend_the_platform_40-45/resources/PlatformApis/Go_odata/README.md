# Go Dataverse OData Client

A terminal-based client application for interacting with Microsoft Dataverse
through the Web API. This application provides a command-line interface for
viewing and managing Dataverse entities such as Accounts and Contacts.

## Features

- Terminal-based user interface with keyboard navigation
- View, create, update, and delete Dataverse entities
- Pagination support for large result sets
- Search functionality to filter entities
- Support for both application-based and user-delegated authentication

## Prerequisites

- Go 1.24 or higher
- Microsoft Dataverse environment
- Application registration in Microsoft Entra ID

## Installation

1. Clone the repository:

```bash
git clone https://github.com/yourusername/go_odata.git
cd go_odata
```

2. Install dependencies:

```bash
go mod download
```

## Configuration

Create a .env file in the root directory with the following values:

```
CLIENT_ID=your-application-id
TENANT_ID=your-tenant-id
CLIENT_SECRET=your-client-secret
ENVIRONMENT_URL=https://yourorg.crm.dynamics.com/
API_PATH=api/data/v9.2/
AUTHORITY=https://login.microsoftonline.com/
```

### Authentication Setup

1. Register an application in the Microsoft Entra ID Admin Center
2. Add appropriate API permissions for Dataverse/Dynamics 365
3. Generate a client secret (if using application authentication)
4. Update the .env file with your application's details

## Usage

Run the application:

```bash
go run main.go
```

### Authentication Modes

When starting the application, you'll be prompted to choose an authentication
mode:

1. **Application** - Uses client credentials flow (requires CLIENT_SECRET)
2. **User** - Uses interactive browser-based authentication

### Navigation

- Use arrow keys (↑/↓) to navigate menus
- Press Enter to select options
- Follow on-screen prompts for creating/updating entities
- Use pagination controls to navigate through large result sets

## Architecture

The application is organized into the following packages:

- app - Application-level functionality and screens
- model - Data models for Dataverse entities
- msal - Authentication with Microsoft Authentication Library
- service - Service layer for API communication
- view - Terminal UI components
- constants - Application-wide constants and enumerations
- utilities - Helper functions
- request_builder - HTTP request construction
