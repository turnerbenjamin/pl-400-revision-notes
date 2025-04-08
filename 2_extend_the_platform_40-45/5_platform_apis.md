# Use Platform APIs

This section relates to the two platform APIS:

- Web APi
- Organisation Service

The organisation service is covered in detail in the plug-ins document as it is
the preferred choice in that context. This document will focus primarily on the
Web API.

Web API is an OData v4 service. Using this open standard helps to improve
cross-platform compatibility. There is no specific Dataverse tooling for the
WebApi, instead, we can use language specific OData libraries. We can write
queries using both either syntax and FetchXml.

We should use use the WebApi unless:

- We are writing dotnet code, AND
- The code is written for a plugin or the project is a Windows client

## Authentication and Authorisation

Microsoft use Microsoft Entra Id for authentication. Once authenticated, access
to the resource can be controlled with RBAC. So, there is a separation of
concerns:

- OAuth: Authentication
- RBAC: Authorisation

To authenticate an application, the app must be registered with Entra Id by
creating an App registration.

### App Registration Account Types

When creating an app registration there are four options for supported account
types:

- Accounts in this org directory: Target audience is internal
- Accounts in any org directory: Target is business and educational customers or
to support multi-tenancy
- Accounts in any org directory + personal MS accounts
- personal MS accounts only

### Authentication

There are two methods of authentication for an application accessing the
Dataverse:

- Acting on behalf of an interactive user
- Acting directly with S2S authentication

#### Acting on Behalf of an Interactive User

To act on behalf of an interactive user, we need to set-up the API permissions
on the application registration. Specifically, the user_impersonation delegated
permission is required for Dynamics CRM.

#### Acting Directly with S2S Authentication

If acting directly with S2S authentication, an application user will need to be
set up in each Dataverse environment through admin.powerplatform. The app user
may then be associated with an Entra Id app or managed identity. The Dataverse
security roles for the application user can be configured in power platform.

#### Microsoft Authentication Libraries (MSAL)

There is a variety of platform specific Microsoft Authentication Libraries which
support user sign-in and access to protected Web Apis. Fundamentally, these are
OAuth 2.0 and OpenId connect libraries, the benefits of using these are that:

- They abstract the protocol level details
- They are regularly updated with security updates
- They cache and refresh tokens automatically
- Supports any Microsoft Identity including personal and work/school accounts

The general pattern is:

- Try to acquire a token silently (i.e. with no user input)
- If this fails, acquire a token by asking the user to input credentials

In OAuth, there are two client types:

- Confidential: Clients capable of maintaining the confidentiality of their
credentials, e.g. a client on a secure server with restricted access
- Public: Clients that cannot maintain confidentiality of their credentials,
e.g. clients on a device installed by the resource owner

## ODataQueries

## FetchXmlQueries

## Perform Operations with Dataverse Web API

### What we Need to Know

- Available at an OData v4 RESTful Endpoint
  - Use with any programming language that
    - supports HTTP requests, and
    - authentication with OAuth 2.0
    - newer service, preferred generally as it is not tied to dotnet
    - faster and lighter than Organisation Service

## Implement API Retry Policies

### What we Need to Know

- Service Protection Limits
  - Evaluated over 5 min rolling window
  - Combination of requests, execution time and number of concurrent requests
  - 429 error thrown with a retryAfterDuration parameter sent
  - Build 429 error handling into the code

- API limits
  - Evaluated over 24 hours
  - Based on licence
  - Can purchase more calls

## Optimise for Performance, Concurrency, Transactions and Bulk Operations

### What we Need to Know

- Concurrency
  - Set concurrency behaviors
- Transactions
  - Handle up to 1000 requests in a transaction with roll back if one fails
- Bulk Operations
  - The same but without rollback, we can set continue on error to false. Also
  1,000 requests
- Likely questions on which to use, options to set and limitations

## Perform Authentication by Using OAuth

### What we Need to Know

- We generally create an Application User for integrations and a user in Power
platform for the service principle with a relevant security role
- MSAL library available for various languages
