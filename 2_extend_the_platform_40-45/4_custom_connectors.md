# Create Custom Connectors

## Intro

Connectors allow the Power Platform to communicate with external APIs, for
instance, an API built with Azure Function Apps. A connector is just a
proxy/wrapper around a REST or SOAP API. It describes the endpoints of the API
and the data structures it receives/returns so that Power Platform may
communicate with that service.

There is a substantial library of prebuilt connectors, however, if these do not
include the service or functionality we require then a custom creator may be
created. Once created we can use the connector in:

- Power Automate
- Power Apps
- Logic Apps
- Copilot Studio

We can build a connector without writing a single line of code, however there
are advanced features of connectors that may require the API to have supporting
capabilities. We may also need to interact with the OpenAPI definition directly
writing JSON or YAML, as some advanced features cannot be customised through the
custom API designer.

## Creating a Connector

To create a custom connector navigate to make.powerautomate. From the left pane
select More -> Discover All -> Custom Connectors.

There are various options to quickly configure a custom connector:

- Create from Azure Service
- Import an OpenAPI file
- Import an OpenAPI from URL
- Import a Postman collection
- Import from GitHub

In this section, we will use Create from Blank. If you have multiple actions and
triggers it will be faster to use one of the other options. When we select
Create from Blank, the custom API designer will be opened. This provides a UI
to build an OpenAI/Swagger definition from scratch. If we need to interact with
the definition directly, we can switch to the Swagger editor.

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

For OAuth 2.0 the set-up is more involved. There is a quick path for a number of
common OAuth services.

### Creating a Connector: Definition Tab

This tab allows us to create:

- Actions
- Triggers
- References
- Policies

#### Actions

Actions are essentially endpoints we can we can call on demand using the
connector.

In the general section, we provide a summary,description and operation id for
the action. We can also specify the visibility:

- None: Displays normally
- Advanced: Hidden under a show advanced menu
- Internal: Hidden from users
- Important: Always shown to users first

Under the hood, this option sets the x-ms-visibility header to the appropriate
value.

There is then a request and response section, these allow us to import from a
sample.

##### Triggers

Triggers allow us to trigger a workflow when an event occurs in an external
service, for example, when a row is added in the Dataverse.

We can use one of two mechanisms for triggers:

- Webhooks: Listen for an event to occur on an endpoint
- Polling: Call service at a specified frequency to check for new data

In both instances, the API needs to provide the relevant capability. For
instance, with webhooks, the API will need to provide an endpoint that will
receive and store callback urls which may then be used to inform Power Platform
when a relevant event has occurred.

##### Policies

Policies can be used to modify behaviour at runtime, for example based on the
language of the user.

We can create polities from the definition tab or using the paconn cli. A given
connector can have multiple policies, and they may be ordered to control the
order of execution.

Each policy can apply to one or more operations.

###### Expressions

Policies use expressions to specify where data is accessed from. Expressions are
prefixed with an @ symbol.

- Wrapping an expression in curly braces will cast a numeric value to a string
automatically
- String literals must be defined with single quotation marks

```exp
@{connectionParameters('HostPortNumber')}   //Cast number to string
@body().invoices[0].invoiceId               //dot and bracket notation
```

###### DataSources

- headers
- body()
- queryParameters
- connectionParameters: Values entered when connection set-up, e.g. auth or host
values
- {paramName}: Action parameters

###### Connection Parameters

When we set up a connector, an apiProperties.json file is generated with the
connection parameters, for instance, an API key parameter. We can manually edit
this file to add additional parameters which may be used in the policies.

###### Templates

- Set header
- Set query parameter
- Set Property (preview): Property in body of request of response
- Set host url: modify host
- Route request: modify path

There are also a set of templates for data conversion, these are currently all
in preview:

- Convert an array to an object
- Convert an object to an array
- Convert delimited string into an array of objects

We can select whether the policy should run on the request or the response. Note
if having transformed the request/response the relevant templates should be
updated.

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

Note:

- Don't include any usings before checking they will be available at runtime.
There is a limited set of available namespaces we may use.
- Only one code block may be defined regardless of the number of actions.
However, the code block may be applied to multiple operations

### Creating a Connector: Test Tab

The final tab is the test tab. We can use this to run actions and validate that
the responses are in the correct shape.

## Extend Open API Definition for a Custom Connector

Custom connectors use OpenAPI (AKA Swagger) definitions to describe auth,
actions, triggers and parameters.

MS has defined a variety of extensions to the OpenAPI definition, which can be
identified by the x-ms-prefix. Many of the extensions can be applied through
the connection designer, for instance x-ms-visibility and x-ms-summary for
actions. Other extensions can only be specified using the Swagger editor.

### Configuring Extensions

There are three methods we may use to configure extensions:

- Import from an OpenAPI file that contains ms extensions
- Use either the Power Platform CLI (pac) or Power Platform Connectors CLI
(paconn) to download the API file, edit locally and push the updated file to the
environment
- The swagger editor may be used in the connector designer, this is generally
the fastest method to add extensions

Note, the swagger editor uses YAML. If we use a CLI table then data is
downloaded and uploaded in JSON.

### Specific Microsoft Extensions

#### x-ms-capabilities

There are two options we may configure for capabilities:

##### Chunk Transfer

The connector runtime limits message content to a maximum size. To handle
messages that are larger than this limit, chunking may be used. To use chunking:

- The API must support chunking
- The custom connector must enable chunkTransfer capability on the action
- The maker using the connector action must enable chunk transfer in the flow
step

In terms of the OpenAPI definition, we would add the following to an
x-ms-capabilities extension at the operation level:

```yaml
{chunkTransfer: true}
```

##### Test Connection

By default, connections are not verified as valid or not at the point of
creation. For instance, if an invalid host URL or invalid API key is provided
then the creation could be created but would fail at runtime.

To enable connection testing, a simple operation must be defined that returns a
200 success status code. This may be an existing action or an action created
specifically for this purpose.

This extension is defined at the connector level and should include the
operationId and any parameters to be passed:

```yaml
x-ms-capabilities:
  testConnection:
    operationId: TestConnection
    parameters: {}
```

#### x-ms-url-encoding

This extension relates to url parameters, e.g. planets/{id}. Since id will be a
parameter supplied by the user of the connector, they will need to be
url-encoded so that they may be included in the url.

By default, single encoding is used, however the api may use double encoding to
resolve potential ambiguities where symbols such as @,/,\ are used. To use
double encoding add the following at the parameter level:

```yaml
x-ms-url-encoding: double
```

This means that users can use the connector without having to encode path
parameters manually.

#### x-ms-dynamic-values & x-ms-dynamic-list

By default, action parameters are added as a simple text box. MS Learn provide a
good example here of an Invoice Type ID parameter with two options:

1. Purchase Order
2. Non-Purchase Order

By default, the user would have to know the available options, e.g. 1 or 2 and
what each mean.

We could use the enum attribute, this is not an extension to OpenAPI, e.g:

```yaml
enum: [1,2]
```

This would then display the two options that may be selected. However, the user
would still need to know which operations each of these ids map to.

##### x-ms-dynamic-values

This extension may be used to request a list of values from the API. The
benefits of this approach are that:

- We can configure a label that the user can see in addition to the value
- More maintainable as values need only be updated in the API itself
- The API can filter values based on context, e.g. security
- Parameters may be passed to the API to further filter the list
- Output may be dependent on other options

To use this approach, the API must provide an operation that returns an array
of values, and optionally a name property for each as a label.

We can then connect this with Swagger:

```yaml
x-ms-dynamic-values: {
  operationId: ListInvoiceTypes,
  value-path: typeId,
  value-title: name,
  value-connection: types,
}
```

##### x-ms-dynamic-list

This is an updated version of the above extension. For older, existing flows the
recommendation is to implement both. For newer flows we should just use
x-ms-dynamic-list.

```yaml
x-ms-dynamic-list: {
  operationId: ListInvoiceTypes,
  itemValuePath: typeId,
  itemTitlePath: name,
  itemsPath: types,
}
```

The benefit of this newer extension is that it can resolve ambiguity, e.g. where
the request has both a path and body parameter with the same name.

#### Dynamic Schema

By default the parameters for a given operation are static. However, we may want
these to be variable, e.g.:

- Parameters based on a category
- Parameters based on security
- Common parameters for an action type applied to multiple actions

We can achieve this with a dynamic schema. Where configured, the custom
connector runtime will call an operation to retrieve the schema for the action.
This can include parameters, display name and parameter descriptions.

As with dynamic-value and dynamic-list, there are two extensions available:

- x-ms-dynamic-schema: v1
- x-ms-dynamic-properties: v2

We should use properties for newer flows and both if supporting older flows.

```yaml
x-ms-dynamic-properties:
  operationId: GetSchema,
  parameters:
    typeId: {
      parameterReference: typeId
      }
```

So we pass the operation used to get the schema and we can also pass parameters
to the operation.

## Sharing Connectors

Once created in an environment, connectors can be exported and imported into
other environments.

We can also share custom definitions on GitHub so that other developers can
import them into their environments. If the connector is certified then it will
be available as a prebuild custom connector in all customer environments.

To certify a connector, you must either own the API or have explicit permission
to publish a connector from the API owner.

Connections cannot be within solutions, instead, connection references are
created and implicitly included when an app or flow using a custom connector is
added to a solution. The reference is pointer to a connection outside the
solution.

Note, by default, custom connectors are not available to other users within the
same environment. They must be shared to make them available to other users. If
the connector is within a solution access then there is no need to share.
Role-based security is used here and availability depends on the user's access
to the Connector and Connection Reference tables.

## Authentication

Custom connectors are underpinned by Azure API Management infrastructure. When a
connection is created, the APIM gateway stores credentials in a token store.

Once created there is no need to authenticate again; Power Apps will pass a
connection id to the custom connector which will then use this id to access the
credentials from the token store and pass them on to the API/On-premises
network.

Depending on how connectors and connections are shared, credentials will be
collected at different stages.

### Use a Service Principle for Authentication

A service principle represents the identity of an application rather than a
user. This allows us to authenticate and perform actions on Azure resources
without providing user credentials.

Role-Based Access Control (RBAC), can be used to assign the service principal
with specific roles in order to control access to Azure resources.

To connect to certain services, e.g. Microsoft Graph and Azure CLI with a custom
connector, we must use a service principal. Azure will copy authentication
details for the service principle within the connector security definition.

#### Authenticating with a Service Principal

At a high level:

- Create an application registration and copy the application and tenant id
- Add a secret for the application and copy the secret value
- In the connector security tab set-up oAuth with Azure AD and enable service
principal support. Enter credentials for the principal
- Save the connector and copy the redirect url from the security tab
- In the service principal app registration, use the auth section to configure
a web platform and copy in the redirect url from the custom connector

## Demos

Three demos have been created for custom connectors

- [SWAPI](./demos/custom_connectors_swapi_demo.md): Uses a pre-existing API and
custom code to transform the responses
- [Arcade Game Scores](./demos/custom_connectors_arcade_scores_demo.md): Uses a
custom API to explore the use of triggers and the dynamic list extension
- [Graph API](./demos/custom_connectors_service_principal_auth_demo.md): Follows
MS learn exercise to authenticate with a service principal
