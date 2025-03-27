
# Plug-Ins and Custom Messages

## Introduction

There are two ways to create custom messages in Power Apps:

- Custom process actions
- Custom Apis

## Custom Process Actions

Custom process actions are are declarative, but more limited alternative to
custom APIs. We can create them from:

Solutions -> New -> Automation -> Process -> Action

We can define a name for the action and also any input/output arguments.

A custom process action is a workflow, we can define steps and actions to
execute business logic when the message is called.

Alternatively, a common pattern is to create a custom process action without any
steps and to register a plug-in step against the message to perform the logic.
Once the action has been created we are able to access the message from the PRT.

## Custom APIs

Custom APIs are an alternative to custom process actions. These may be created
without a plug-in to pass data about an event, but generally a plugin is used to
perform some business logic and return a result.

### Create a Custom API

There are various ways to create a custom API:

- Register with PRT
- Power apps maker portal (Solution -> New -> More -> Other -> Custom API)
- Code

Note, the custom api and any input/output parameters will be customisable by
default. It is recommended that we change this to false so that it cannot be
modified when exported as a managed solution. Find the Api and params in the
solution, select the kebab menu and edit the managed properties.

### Custom API Tables

The metadata for custom apis is stored in the following Tables:

- Custom API
- Custom API Request Parameter
- Custom API Response Parameter

## Creating a Custom API with PRT

This is straight forward, select Register -> Custom API. We can set various
properties when defining the API:

### Step Types

There are three step types:

#### None

This should be used when the custom API provides functionality that should not
be customisable. Other developers are unable to register any more steps against
the message.

#### Async Only

We can use this when we want to allow developers to detect when the operation
occurs but not permit them to cancel or customise the behaviour of the
operation.

This is recommended when using the business event pattern where a business event
creates a trigger in Power Automate you can use when the event occurs

#### Sync and Async

Most messages will use this option. Other developers are able to register
synchronous steps that can modify and cancel the operation.

### Binding Types

Binding in OData, associates an operation with a specific table. There are three
options here:

- Global: Unbound
- Entity: Bound to a record in a table
- EntityCollection: Applies to a collection of records in a table

When entity is used, a Target input parameter is created automatically.

### Functions vs Actions

The IsFunction property determines whether the custom Api is a function or an
action.

In OData a function is an operation called with GET that returns data without
making any changes. Parameters are passed in the url when the function is
invoked.

Functions must return some data else it will not appear in the metadata service
document and a 404 error will be returned when the function is invoked.

We cannot use functions when the API is enabled for workflow. The Dataverse
Connector currently only enables actions.

Note, there is a limit on the length of a url we can send.

Actions use POST:

- Params are passed in the body which gets around the url length limit
- They can be enabled for workflow
- There is no requirement to return data.

### Making APIs Private

We can make Custom APIs private, this will keep the API out of the metadata
service document. While this will flag that the API should not be used it does
not prevent other developers using the message if they have awareness of it.

We can add some security by adding an Execute Privilege property.

### Input/Output Parameters

There's not much to say here, we can define input and output parameters. Just
remember that:

- A function must have at least on output parameter
- Input parameters for a function are passed as a query string which can be
limiting

## Implement Business Logic for a Custom API with a Plug-in

## Invoking the Custom Message

### Invoking with a Plug-In

### Invoking with Power Automate

### Invoking with Web Api

## Develop a Plugin that Implements a Custom API

### What we Need to Know (Custom API Plugin)

- Add plug-in to perform logic when custom API called
- Register the assembly
- Rather than register a step, link plugin assembly to the custom API

## Configure Dataverse Business Events

### What we Need to Know (Events)

Historically plug-ins used for create and update of records. But we may have
multiple events at once, e.g. an invoice with various line items.

We can use business events, e.g. post invoice. There is a catalogue of events
from which we can choose.

We can then handle a single transaction with a single event which simplifies the
logic.

We need to look into this.
