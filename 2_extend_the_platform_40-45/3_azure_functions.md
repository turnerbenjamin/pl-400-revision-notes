# Azure Functions

## Introduction

Power Platform can integrate natively with a variety of Azure Services. One of
these services is Azure Function Apps. We can use this service to generate
serverless functions. Azure Functions supports a variety of triggers and
bindings. We can also select from a range of runtimes including:

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
- Queue Storage: Runs when a message is added to specified storage queue
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

## Building an Azure Function

The first step is to create a function app in Azure Portal. Once the app has
been created we can develop the code locally. If using VS Code, it is useful to
install the following extensions:

- Azure Functions
- Azure Tools
- C# Extension

Once the extensions have been installed, search the command palette for Install
or Update Azure Functions Core Tools.

### Scaffolding the Function Project

Search the command palette for Create New Project and work through the wizard.
We will need to select the directory, runtime and template for the project. If
using a .NET runtime, select an isolated runtime; support for the in process
model is ending.

Once the project has been scaffolded, it will include two json files:

- host.json: Allows us to configure the host
- local.settings.json: Allows us to configure things like secrets for local
development

### Running the Function Project

To run the function we can use the func start command. This is a simple process
for the HTTP and Webhook triggers. However, for other triggers we will need some
additional set-up, e.g. setting up Values.AzureWebJobsStorage key to a valid
Azure storage account connection string.

### Publishing the Function Project

We can publish the function using the command palette; search for the Deploy to
Function App command and follow the wizard.

## Process Long-Running Operations with Azure Functions

Azure functions can be used for long-running operations:

- 5-10 minutes: Consumption plan
- 30-unlimited: Premium plan

This stands in contrast to Power Platform Plug-ins which will timeout after 2
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

Without the Durable Functions Extension, we may implement the following process:

- F1 is triggered and runs
- F1 stores a result in storage or a queue
- F2 is triggered by a blob being added to the storage or a message in a queue

Such a process can meet the three rules above, but the workflow may be difficult
to understand. Durable Functions, essentially allow us to define these workflows
in code. Behind the scenes, there will still be queues and table storage, but
Azure Durable Functions will handle this for us.

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

With some connections, we will need to use tokens and secrets to authenticate.
With Azure functions, we can use the application settings to securely store
these credentials. They will be encrypted and accessed by the app at runtime
as environment variables. For local development these credentials can be stored
in local.settings.json and added to .gitignore.

With some connections, including to Power Platform with the Web API, we can
authenticate with an Entra ID identity rather than use secrets. The process is
simple, we can create an app registration to act as a service principle. We can
then create an App User in the Power Platform environment linked to the service
principle. Finally we can grant the app user the necessary permissions.

This process is demonstrated in a
[demo](./demos/web_api_go_odata_demo.md)
for the platform APIs document in the context of an external console
application.

With Azure functions we can also use managed identities. When we create a
managed identity for a function app, an application registration is created
which we can use as described in the platform API demo. The main difference here
is that azure will manage the certificates and secrets for the app registration
on our behalf so we do not need to keep renewing these manually when they
expire.

## Demonstrations

I have not created any specific demonstrations of Azure functions in this
document, however, they are used to demonstrate other concepts:

- [custom connectors](./demos/custom_connectors_arcade_scores_demo.md)
- [webhooks](../6_develop_integrations_5-10/demos/webhook_demo.md)

The custom connectors demo, covered in the next document, contains:

- Standard and durable functions
- HTTP and service bus queue triggers
- Use of SQL bindings

There is a write-up to explain the creation of this Azure Function App and a few
of its methods [here](./demos/azure_functions_demo.md)
