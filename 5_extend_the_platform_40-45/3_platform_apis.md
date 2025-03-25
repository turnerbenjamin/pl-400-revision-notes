# Use Platform APIs

## Perform Operations with Dataverse Web API

### What we Need to Know

- Available at an OData v4 RESTful Endpoint
  - Use with any programming language that
    - supports HTTP requests, and
    - authentication with OAuth 2.0
    - newer service, preferred generally as it is not tied to dotnet
    - faster and lighter than Organisation Service

## Perform Operations with Organization Service

- Dotnet SDK
  - Typed class generators for tables classes
  - Preferred method for plugins
  - Instantiated and available without authentication in plugins
- Know that helper methods are wrappers around Execute and we can use Execute
for more control
- Associate/Disassociate used with relationships

### What we Need to Know

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
