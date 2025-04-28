# Azure Functions

## Introduction

Azure Function Apps can be used to generate serverless applications. Azure
Functions supports a variety of triggers bindings. We can use a variety of
runtimes when defining a function app including:

- .Net
- Node.js
- Python
- Java
- PowerShell Core

## Triggers

Azure Functions run in response to triggers, each function must be associated
with a single trigger. Many triggers will have associated data which will
generally be available as the payload of the function:

- Timer: Runs on a schedule defined with a Cron Expression
- Queue Storage: Runs when a message added to specified storage queue
- Service Bus Queue: Ditto but with a service bus queue
- Bus Topic Queue: Runs when a message is added to a specified Service Bus Topic
queue
- Blob Storage: When a blob is added to a specified container
- Event Hub Trigger: When an event hub receives a new event
- IoT Hub: When an IoT Hub receives a new message

### Cron Expressions

A Cron Expression specifies a frequency in:

- Seconds
- Minutes
- Hours
- Days
- Months
- Days of Week (Monday is 1, Sunday is 0 or 7)

"*" Means any
"/n" declares an interval, e.g. 0 0 /2 *** will run every two hours, every day,
every month on any day of the week.
"n-n" specifies a range
"n,n,n" specifies a list

## Process Long-Running Operations with Azure Functions

Azure functions can be used for long-running operations:

- 5-10 minutes: Consumption plan
- 30-unlimited: Premium plan

This stands in contrast to Powerplatform Plug-ins which will timeout after 2
minutes.

### Durable Functions

If we are defining a long-running operation, we should use Durable Functions.
This is an open source extension to Azure functions and allows us to define
complex orchestration flows in code.

#### Durable Functions Theory

As a general rule:

- Functions MUST be stateless
- Functions SHOULD not call other functions
- Functions SHOULD have a single responsibility

Without the Durable Functions Extension, we may the following process:

- F1 is triggered and runs
- F1 stores a result in storage or a queue
- F2 is triggered by a blob being added to the storage or a message in a queue

Such a process can meet the three rules above, but the workflow may be difficult
to understand. Durable Functions, essentially allow us to define these workflows
in code. Behind the scenes, there will still be queues and table storage, but
Azure Durable Functions will handle all of this behind the scenes.

The extension also contains a toolbox to help with common tasks like timers and
retries.

#### Durable Function Components

A basic durable function will contain three core components:

- Orchestrator Function: This contains the workflow logic. It is used to
dispatch activities in the workflow
- Activity Functions: The basic worker functions in the workflow. Triggered or
scheduled by the orchestrator
- Client Function: This is used to trigger an orchestration function. It can
also get information about the status of the orchestration.

#### Durable Function Patterns

##### Function Chaining

With function chaining, the orchestrator function will execute activity
functions sequentially. This enables output from earlier functions to be used as
input to subsequent functions.

To prevent double-billing, the orchestrator may be put to sleep while awaiting
one or more activity functions. This is handled automatically by the Durable
Functions extension.

##### Fan Out/Fan In

With this pattern, activity functions are run in parallel. We can then wait for
all tasks to complete in the orchestrator function.

Without durable functions, this would involve a lot of state management in
Azure.

##### Async HTTP APIs

With this pattern, an endpoint will trigger a durable function and the client is
redirected it an endpoint which may be used to poll the status of the action.
The Durable Functions runtime automatically exposes inbuilt APIS for polling.

##### Monitor

This is the reverse of the above. The Durable function will monitor a polling
endpoint for status changes and return when either the status is complete or a
timeout is reached.

##### Human Interaction

Similar to the above, we may trigger some human interaction and the await a
response.

##### Aggregator

An aggregator function can be designed to read data from various sources and
send it to be processed.

## Authenticate to Power Platform with Identity-Based Connection

Azure functions use the applications settings functionality of Azure App Service
to securely store secrets and tokens when connecting to other services. App
Settings in Azure are stored encrypted and accessed at runtime by the app as
environment variables.

Some connections are configured to use an identity instead of a secret. Support
for this depends on the extension using the connection.

Note that the identities must have the required permissions to perform actions
this is done using RBAC.

## Authenticate to Power Platform with Managed Identities

### What we Need to Know

- We can use App registrations
- We can create managed identities in Azure
  - Uses client ids and secrets
