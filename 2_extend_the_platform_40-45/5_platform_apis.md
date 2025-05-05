# Use Platform APIs

Platform APIs relates to:

- Web API
- Organisation Service

The organisation service is covered in detail in the plug-ins document as it is
the preferred choice in that context. This document will focus primarily on the
Web API.

## Web API

Web API is an OData v4 service. Using this open standard helps to improve
cross-platform compatibility.

When choosing between Organisation Service and Web API, we should always go with
Web API unless:

- We are writing dotnet code, AND
- The code is written in the context of a plugin or Windows client

### Authentication and Authorisation

Microsoft Entra ID is used for authentication. Once authenticated, access to the
resource can be controlled with RBAC. So, there is a separation of concerns:

- OAuth: Authentication
- RBAC: Authorisation

To authenticate an application, the app should be registered with Entra ID by
creating an App registration in Azure.

#### Authentication Methods

Once an app registration has been created, additional configurations depend on
the method of authentication. There are two approaches we can take:

- App only access
- Delegated access

##### App Only Access

When using App only access, a client secret will need to be generated for the
app registration. This is used to authenticate the app.

Once authenticated, the app will act as itself and, therefore, requires the
appropriate permissions. In the context of Dataverse, we need to:

- Create an app user in make.powerapps for the App Registration in Azure
- Grant necessary permissions to the App user

Note, an app user with the required permissions will need to be set-up in all
environments that the app will be used.

##### Delegated Access

The approach is different for delegated access. In this instance, authorization
is based on the interactive user's security roles. The App Registration needs to
be configured specifically for this authentication method.

For delegated access, authentication will happen through OAuth 2.0 flows.
Depending on the application you may use interactive and/or non-interactive
flows. If an interactive authentication flow is used, it will be necessary to
configure a redirect URI on the App Registration. Entra ID will redirect the
user to this URL after authentication along with any necessary tokens.

Regardless of the OAuth 2.0 flows used in the application, it will also be
necessary to grant the App Registration the Dynamics Crm user_impersonation
permission.

##### Microsoft Authentication Libraries (MSAL)

There is a variety of platform specific Microsoft Authentication Libraries which
support user sign-in and access to protected Web APIs. Fundamentally, these are
OAuth 2.0 and OpenID connect libraries, the benefits of using these are that:

- They abstract the protocol level details
- They are regularly updated with security updates
- They cache and refresh tokens automatically
- They support any MS Identity including personal and work/school accounts

The [Go Odata demo](./demos/web_api_go_odata_demo.md) looks at using an MSAL
library to authenticate using both app only access and delegated access.

### Odata Queries

Web API is an OData 4.0 service, this offers a standardised way to query and
work with data. This [postman demo](./demos/web_api_postman_odata_demo.md)
demonstrates use of this api.

#### FetchXml

The Web API also supports fetchXML; this is not part of the oData standard but
can be used to perform more complex queries. Odata has the $apply query
parameter which may be used to perform grouping and aggregation, however, there
are various limitations to the Dataverse implementation of $apply:

- Distinct count not supported with count column
- Grouping by parts of a date always uses UTC
- No support for per query limits

FetchXml is used throughout dataverse, for instance, it has already been used in
the plug-in demo for [retrieve multiple](./demos/plug_ins_retreive_multiple.md).
When we use advance find in a model-driven app, we can also export the query to
XML. This can be added as the value to a fetchXml query parameter.

Note, the fetchXml query parameter is not prefixed with a $ as it is a
non-standard feature.

#### Call Actions from the Web Api

We can run actions, e.g. custom api actions and inbuilt dataverse actions, from
the Web Api. This is demonstrated
[here](./demos/web_api_and_custom_actions_demo.md).

### User Impersonation

If it is necessary to run logic for another user, then in Dataverse, the logic
will apply all roles and object-based security based on the user being
impersonated.

This is achieved by providing a CallerObjectId in the message header to indicate
that the message should run as that user. The value of the parameter will be the
impersonated user's Entra Id which can be queried using the MS Entra Graph API.

There are system tables used to track this behaviour, for instance createdby
and createdonbehalf columns.

This is demonstrated [here](./demos/web_api_user_impersonation_demo.md)

## Optimise for Performance, Concurrency, Transactions and Bulk Operations

This is a more general section with information relevant to both Organisation
Service and Web API.

### Work with Concurrency

Power Apps is a multi-threaded and multi-user system. There is a need,
therefore, to prevent race conditions. The information below relates only to:

- UpdateRequests
- DeleteRequests

#### Optimistic Concurrency

Power Apps supports optimistic concurrency for:

- All custom tables
- All out-of-the-box tables enabled for offline sync

This strategy allows processes to proceed without locks on the assumption that
no conflicts will occur. If a conflict is detected then the operation will be
rolled back.

We can enable optimistic concurrency when using Organisation Service by setting
the concurrency behaviour of a request to IfRowVersionMatches.

```cs
var request = new UpdateRequest(){
    Target = updatedEntity,
    ConcurrencyBehavior = ConcurrencyBehavior.IfRowVersionMatches,
};
```

When using the Web Api, we can use the If-Matches header and provide the record
version as a value. The record version can be found in the odata.etag property
included when records are returned with the API. A demonstration can be found
[here](./demos/web_api_optimistic_concurrency.md)

#### Always Overwrite

The alternative to optimistic concurrency, is alwaysOverwrite. This will
overwrite the record regardless of the version number.

#### Default

The default behaviour will be always overwrite except:

- Where optimistic concurrency is enabled on the table, AND
- The source is WebService as opposed to a plug-in/custom workflow activity

#### Which Should You Use?

Assuming both behaviours are available, optimistic concurrency adds some
overhead but prioritises data integrity. Always overwrite is faster but will be
at the cost of data integrity where we do not want the latest update to always
take precedence.

### Dataverse Platform Constraints

- Plug-ins timeout after 2-minutes, they should not be used for long-running
  operations
- Requests to SQL server also generally timeout after 2 minutes
- Service Protection Limits

Running into these constraints can indicate a problem with a plugin, series of
plug-in or web API client

#### Service Protection Limits

There are a number of service protection limits:

- a user cannot make more than 6000 requests within a 5-minute period
- a user cannot make requests with a total execution time of 20 mins within a
  5-minute period
- A user cannot make more than 51 concurrent requests

##### Service Protection Limits and Plug-Ins

Service protection limits do not directly apply to plug-ins or custom workflow
activities. Requests made by a plug-in triggered by an application will not
count to service protection API limits. However, the computation time of these
operations are added to the initial request and this is part of the service
protection limits.

##### Retry Policies

When a user exceeds service protection limits they will receive an:

- OrganizationServiceFault (Organisation Service)
- 429 Too Many Requests error (Web Api)

The OrganisationsServiceFault ErrorDetails collection will contain a Retry-After
key. Similarly, a 429 response will include a Retry-After header.

One implementation of IOrganisation Service, CrmServiceClient, will handle these
errors automatically. This class is in the Crm.Tooling.Connector library.

### Transactions and Locks

The protect the integrity of the database, operations are performed within a
transaction. Each request will place locks on resources which will not be
released until the transaction has been committed or aborted. For instance:

- Create: Write lock against the record
- Update: Write lock against the record

Transactions can include up to 1000 operations and can be rolled back if any one
operation within a transaction fails.

#### Transactions and Plugins

Whether or not a plug-in is within a transaction depends on the stage:

- pre-validation: Before the transaction is started
- pre-operation: Inside the transaction
- post-operation (sync): Inside the transaction
- post-operation (async): Outside the transaction

A plug-in registered in the pre-validation stage will be within a transaction if
the pipeline event was triggered by a message within an existing transaction.

With async post-operation plug-ins, no transaction is created. Each message is
acted on independently.

#### Transactions and Web Service Requests

With plug-ins, the execution context that maintains the transaction context.
External requests with web services will create a pipeline with transaction
handling but the transaction will close once the response is returned. In this
scenario, there are two methods that may be used to perform multiple actions
within a single request:

##### ExecuteMultiple

Execute multiple is used to pass multiple independent requests. No transaction
context is held since the requests are independent. Since the operations run
outside the context of a transaction, there is no rollback, however, we can
pass options to stop execution when an error is encountered.

##### ExecuteTransaction

Here related actions can be processed in a single transaction. This should be
used with care as it has the potential to create blocking issues in the
platform.

#### Concurrency and Transactions

When there are multiple concurrent requests, there is a higher chance of
collisions on locks.

##### Async Operations

Async workflows or plug-ins are not processed serially from a queue. Multiple
activities are processed by dataverse in parallel within a given service
instance and across service instances. Each async service instance retrieves
jobs in batches of around 20 depending on configuration and load.

Multiple asynchronous activities from a single event are likely to be processed
in parallel and, if they operate on the same resource, lock conflict risk is
high.
