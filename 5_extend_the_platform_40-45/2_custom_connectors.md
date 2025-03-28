# Create Custom Connectors

## Intro

Power Apps includes a large library of prebuilt connectors, however, if these do
not include a connection that we require we can create custom connectors.

A custom connector is just a wrapper around a REST or SOAP API. It describes the
endpoints of the API and the data structures used.

Once created we can use the connector in:

- Power Automate
- Power Apps
- Logic Apps
- Copilot Studio

## Creating a Connector

To create a custom connector navigate to make.powerautomate. From the left pane
select More -> Discover All -> Custom Connectors.

There are various options to quickly configure a custom connector:

- Create from Azure Service
- Import an OpenAPI file
- Import an OpenAPI from URL
- Import a Postman collection
- Import from GitHub

In this section, we will use Create from Blank.

### Creating a Connector: General Tab

The general tab contains high level information for the connector. There are
two broad categories:

1. Display information: Icon, icon background colour and description
2. Root connection info: Connect via on-premises gateway, http/s, host and base
url

### Creating a Connector: Security Tab

There are 4 options for security:

- No auth
- Basic auth: username and password
- Api Key: key and location (header/query)
- OAuth 2.0

For basic auth and api keys we specify a label and variable name for the
credentials rather than enter them directly.

For OAuth 2.0 the set-up is more involved, however, there is a large set of
identity providers which, when selected, offer a more streamlined configuration.

### Creating a Connector: Definition Tab

This tab allows us to create:

- Actions
- Triggers
- References
- Policies

#### Actions

Actions are essentially endpoints we can we can call on demand using the
connector.

##### Actions: General

In the general section, we provide a summary and description for the action. We
also need to provide a unique operation id and specify the visibility:

- None: Displays normally
- Advanced: Hidden under a show advanced menu
- Internal: Hidden from users
- Important: Always shown to users first

Under the hood, this option sets the x-ms-visibility header to the appropriate
value.

##### Actions: Request

The request contains:

- HTTP Verb
- Endpoint url
- Headers
- Query Parameters
- Body

We can import from a sample.

##### Actions: Response

We can define the response from a sample copying in the headers the JSON
returned.

### Creating a Connector: AI Plugin Tab

To use a connector as an AI connector it must be certified. This section relates
to this certification. This is entirely optional.

### Creating a Connector: Code Tab

This section is also optional, however, it can be useful if we need to perform
any data transformations or add additional functionality.

The code must be:

- Written in C#
- Have a maximum execution time of five seconds
- Have a file size no larger than 1mb

The resources contain a ludicrous example of custom code. Although, it is likely
unrealistic, some key takeaways from the exercise were:

- Don't include any usings before checking they will be available at runtime.
There is a limited set of available namespaces we may use.
- Only one code block may be defined regardless of the number of actions. It is
possible to enable multiple actions to use the custom code but you may need to
add some routing in the class if each action has custom logic.

### Creating a Connector: Test Tab

The final tab is the test tab.

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
