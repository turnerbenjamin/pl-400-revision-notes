# Create Custom Connectors

## Create Open API Definition for Existing REST API

### What we Need to Know

- They are wrappers around existing APIS
  - Allows use with Power Apps and Power Automate
  - Just Create an Open API Definition (Swagger File)
  - We can use the wizard
  - Import a definition
  - For Azure Information Protection Management there is a button to create a
  custom connector
  - Import from GitHub directly
  
## Implement Authentication for Custom Connectors

### What we Need to Know

- Authentication
  - No auth
  - Basic
  - OAuth 2.0 (Generic or prebuilt configs)
  - API Key (Particularly with Azure Functions)

## Configure Policy Templates to Modify Connector Behaviour at Runtime

### What we Need to Know

- A connector may group a number of API calls
  - Policy templates allow us to manage differences between calls
  - E.g. route to different endpoint, add info to header etc

## Import Definitions from Existing APIS

### What we Need to Know

## Create a Custom Connector for an Azure Service

### What we Need to Know

- Create from portal.azure

## Develop an Azure Function to be Used in a Custom Connector

### What we Need to Know

- Azure functions commonly used with Custom Connectors if HTTP trigger used
- Build with wizard or use APIM to create the connector

## Extend Open API Definition for a Custom Connector

### What we Need to Know

## Develop Code for a Custom Connector to Transform Data

### What we Need to Know

- Use the wizard to write code to call the api and then edit data to return a
transformed response
