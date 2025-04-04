# Create Custom Connectors

## Intro

A connector is just a proxy/wrapper around a REST or SOAP API. It describes the
endpoints of the API and the data structures so that Power Apps may communicate
with that service.

Power Apps includes a large library of prebuilt connectors, however, if these do
not include the connection that we require we can create custom connectors.

Once created we can use the connector in:

- Power Automate
- Power Apps
- Logic Apps
- Copilot Studio

We can create a connector without writing a single line of code.

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
triggers it will be faster to use one of the other options.

The create from blank option just takes you through a wizard to create an Open
Api/Swagger file.

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

##### Triggers

In addition to actions, we can also create triggers. We can subscribe to these
triggers, e.g. with a Power Automate and take actions.

We can use one of two mechanisms for triggers:

- Webhooks: Listen for an event to occur on an endpoint
- Polling: Call service at a specified frequency to check for new data

The API will need to support polling or webhooks for these to be defined.

##### Policies

Policies can be used to modify behaviour at runtime, for instance based on the
language of the user.

We can create polities from the definition tab or using the paconn cli. A given
connector can have multiple policies, and they may be ordered to control the
order of execution.

Each policy can apply to one or more operations (actions/triggers)

###### Expressions

Policies use expression to specify where data is accessed from. Expressions are
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

## Configure Policy Templates to Modify Connector Behaviour at Runtime

### What we Need to Know

- A connector may group a number of API calls
  - Policy templates allow us to manage differences between calls
  - E.g. route to different endpoint, add info to header etc

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

## Extend Open API Definition for a Custom Connector

Custom connectors use OpenAPI (AKA Swagger) definitions to describe auth,
actions, triggers and parameters.

MS has defined a variety of extensions to the OpenAPI definition, which can be
identified by the x-ms- prefix. Many of the extensions can be applied through
the connection designer, for instance x-ms-visibility and x-ms-summary for
actions. Other extensions can only be specified using the Swagger editor.

### Configuring Extensions

There are four methods we may use to configure extensions:

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

##### x-ms-url-encoding

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

##### x-ms-dynamic-values & x-ms-dynamic-list

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

###### x-ms-dynamic-values

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

###### x-ms-dynamic-list

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

##### Dynamic Schema

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

### Entra ID

We can have an Entra ID secured custom connector. The architecture is similar to
a general custom connector, however, both the connector and app service will be
have an App Registration with Entra ID.

The benefit is that makers can apply appropriate permissions and then pass the
calling user's identity from an app in Power Apps/Power Automate/Logic Apps to
the underlying service. This helps to protect and authorise access to secured
resources using the identity of the caller.

The service may be an existing service registered in Entra ID or a custom
service, e.g. one implemented with Azure Functions or Microsoft Azure App
Service.

#### Securing a Connector with Entra ID

We need to create two App registrations to identify and protect the API service
and the connector.

We need to allow the registered app for the connector to make "on-behalf" calls
to the service's identity.

Next, we need to set up OAuth 2.0 in the connector with Azure Active Directory
as the provider. This will generate a redirect url which we can then add to the
connector app's registration.

If the service is set up with CORS, then you need to allow APIM domains,
generally azure-apim.net for CORS on the service.

While this is a more involved process, it is the recommended method for securing
connectors.

## Develop an Azure Function to be Used in a Custom Connector

We can create an API service using a HTTP Triggered Azure function and connect
with a Custom Connector.

TODO: Add a demo connecting to a HTTP Trigger created in the Azure Functions
section. There is no need to write a separate function for this.

## For Tomorrow

- continue from Policy Templates
- Set up a simple API with Azure Functions that: x
- uses custom policy templates and adjusts connection properties x
- Uses Entra ID for Auth
- Ideally has a trigger
- Ideally uses some API Extension methods
