
# Custom APIs and Custom Messages

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

MS recommend avoiding this method for creating custom business events. This is
because:

- As with other workflows they can be disabled in the UI
- It is less intuitive to prevent developers from registering synchronous steps
against the message
- We cannot define execute privileges

Instead, we should use:

## Custom APIs

Custom APIs are an alternative to custom process actions. These may be created
without a plug-in to pass data about an event, but generally a plugin is used to
perform some business logic and return a result.

These can be used to extend the OOTB messaging to extend the Dataverse Web Api.
For instance, we may have some complex logic to find a customer. We can simplify
this for callers by creating a findCustomer message which implements the logic
and returns the results.

## Custom API Tables

The metadata for custom apis is stored in the following Tables:

- Custom API
- Custom API Request Parameter
- Custom API Response Parameter

## Create a Custom API - Overview

There are various ways to create a custom API:

- Register with PRT
- Power apps maker portal (Solution -> New -> More -> Other -> Custom API)
- Code

Note, the custom API and any input/output parameters will be customisable by
default. It is recommended that we change this to false so that it cannot be
modified when exported as a managed solution. Find the API and parameters in the
solution, select the ... menu and edit the managed properties.

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
operation. Since only async steps may be registered, any logic added by other
developers will always be outside of the transaction.

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
invoked. Functions must return some data else it will not appear in the metadata
service document and a 404 error will be returned when the function is invoked.

We cannot use functions when the API is enabled for workflow. The Dataverse
Connector currently only enables actions.

Actions use POST:

- Params are passed in the body which gets around the url length limit
- They can be enabled for workflow
- There is no requirement to return data.

### Making APIs Private

We can make Custom APIs private, this will keep the API out of the metadata
service document. While this will flag that the API should not be used, it does
not prevent other developers using the message if they have awareness of it.

We can add some security by adding an Execute Privilege property.

### Input/Output Parameters

There's not much to say here, we can define input and output parameters. Just
remember that:

- A function must have at least one output parameter
- Input parameters for a function are passed as a query string which can be
limiting

When using a plug-in, we can access both the input and output parameters from
the execution context.

## Custom API Business Logic

We can attach business logic to a custom API by:

- Attaching a plug-in to the custom API
- Registering a plugin step against the custom API message (if enabled)
- Listening for the message with a Power Automate Flow (actions only)

Generally, a plug-in is attached to the custom API. However note that:

- The profiler cannot be used for debugging as the profiler is attached to a
step
- Secure and unsecure configurations cannot be set

In both cases the workaround is to choose the second option, register the plugin
against the custom api message.

## Triggering a Custom API Message

Again, there are multiple ways to trigger a custom API message.

- Use a bound or unbound action in a power automate flow depending on the
binding type of the API
- Trigger from a plugin
- Trigger from WebApi

```cs
var req = new OrganizationRequest(customApiName)
```

A demonstration of Custom APIs can be found [here](./demos/custom_api_demo.md).
This demo demonstrates the use of custom APIs from a technical perspective but
I am not sure it represents a good use case in that it is just a wrapper around
an external API rather than an extension of the Dataverse API.
