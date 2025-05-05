# Configure Power Automate Cloud Flows

## Introduction

Power automate flows are an asynchronous declarative solution.

## Authentication

Dataverse connectors are set-up to authenticate with OAuth. They also support
authentication with a service principal:

- OAuth: Used when connection for an interactive user using an Entra ID
organisation account; e.g. user productivity flows
- Service Principal: For non-interactive users; e.g. for flows that support a
project and automation represents work being done for the application

## Triggers

Every flow needs a trigger, there are various triggers that we can use. For
instance:

- Recurrence: Flow runs at a specified frequency
- When a new email arrives: Outlook connector
- Manually trigger a flow

A simple example using recurrence can be found
[here](./demos/pa_recurrence.md)

### Dataverse Triggers

Dataverse has the following triggers:

- When a row is added, modified or deleted
- When an action is performed
- When a flow step is run from a business process flow
- When a row is selected

#### When a Row is Added, Modified or Deleted

We must specify:

- Change Type (e.g. Added, Added or Modified, Modified, Deleted)
- Table name
- Scope (e.g. Business unit, organisation, parent: Child business unit, User)

When multiple change types are selected, e.g. Added or Modified, we can query
the trigger body for the change type:

```exp
triggerBody()?['SdkMessage']
```

##### Select Columns

This applies when the change type includes "modified". It can be used to filter
modification events that will trigger the flow. For instance, if we define
select columns as firstname,lastname then only updates affecting one or both of
these columns will trigger the flow.

##### Filter Rows

With filter rows we define an Odata expression. The trigger will only run if the
expression returns true for the row.

```oData
firstname eq 'John'
contains(firstname,'John')
revenue lt 100000 and revenue gt 2000
(contains(name,'sample') or contains(name,'test')) and revenue gt 5000
primarycontactid/fullname eq 'Susanna (sample)'
```

This is mostly straight forward, just note the / syntax for accessing related
data.

##### Advanced Parameters

- Delay Until: Delays the flow until a specific time
- Run as: (Flow owner, modifying user, Row Owner)

Note, to use run as row owner or modifying user, the flow owner must have the
Act on Behalf of Another User privilege. The Delegate security role has this by
default.

## Actions

This section will focus on Dataverse actions

### Retrieve Actions

There are two commonly used retrieve actions:

- Get a row by id
- List rows

#### Retrieve Action Common Parameters

Both actions share the following parameters

##### Retrieve actions: Select Columns

This just defines a select query value, provide a comma-separated list of
columns to improve the efficiency of the query.

##### Retrieve Actions: Expand Query

We can use an OData style expand query to include data from related tables. For
example:

```exp
primarycontactid($select=contactid,fullname)
```

#### Get a Row by Id

This simply gets a row using a Guid. Generally this is used to get data for a
related entity.

The step will fail if the id is null or if you do not have permission to read
the data for the row

#### List Rows

If get row by id is equivalent to retrieve, list rows is equivalent to retrieve
multiple. We can specify a query using:

- Filter Rows: OData filter expression
- Fetch Xml Query: Fetch XML query expression

Note, at present aggregation queries are not supported when using the list rows
action with FetchXml

##### Retrieve more than 5,000 Rows

By default, list rows will not return more than 5,000 rows.

There are two options to get more than this:

First, in the settings for the list row action, we can enable pagination. It is
counter intuitive, but this will automatically retrieve all pages and return a
single set of results. We can retrieve up to 100,000 records with this means.
We can not use pagination with FetchXml expressions, but we can use OData
expressions.

Second, we can manually handle pagination using the nextLink parameter and a
skip token. After making an initial request, we would then initialise a skip
token and parse the value to skip from the nextLink parameter. Note that
nextLink will not be populated if pagination is turned on so this will need to
be disabled for the method to work.

### Create, Update, Delete and Associate

#### Add a New Row

We can select a table. The fields available to be set on that table are
automatically shown.

#### Update a Row and Upsert a Row

There are actions to both update and upsert a row. These are straight forward,
if you want to clear a field set the value to null.

Note that update will function like upsert, if we provide an existing GUID the
existing account is created. If we use a new GUID then a new record will be
created.

#### Relate Data

There is a relate rows action which can be used to define a relationship. This
is the only method of defining many to many relationships.

For n:1 relationships, we can simply add a guid to the lookup column to create
the relationship.

### Changeset Request

The Perform a changeset request action can be used to perform dataverse
operations in a transaction with rollback if any of the operations fail. This
can be used with:

- Add a new row
- Delete a row
- Update a row

This action is not supported in the new designer, so this will need to be
disabled if you want to use this action.

### Returning Data

There are two actions commonly used to return data from a Power Automate Flow:

#### Respond to a Power App or Flow (Standard)

This action is specifically designed to respond to a Power App or Flow. The
benefits of this action are that it is:

- Free (standard)
- Simple to use

The main limitation with this action is that there are limited types that may
be returned:

- Text
- Yes/No
- File
- Email
- Number
- Date

There is no way to return a collection of values, other than perhaps serialising
and returning as text and then deserializing in Power Apps.

#### Response (Premium)

This is a general use action. It is:

- More complex
- More powerful
- Premium

We provide a response as JSON and generate a schema.

The use of both actions is demonstrated [here](./demos/pa_response.md)

## Flow Permissions

### Sharing Flows (Co-Owners)

With shared flows:

- Multiple people can share and manage a flow
- If the creator of a shared flow leaves the organisation, other owners can
still run it
- All owners of a shared flow can views its history, managed properties, edit
it, delete it and add/remove owners
- Access the connections in the flow

To create a shared flow, select the share command for the flow and add users or
groups as owners. All owners will be added to the co-owners flow.

Once a flow has been shared it will not appear in the cloud flows tab but in the
shared with me tab.

### Run Only Permissions

From the flow page we can select edit run only users to give permission to run
but not edit the flow. Unlike co-owners, there is no default access to the
connections to the flow we must specify, for each connection, whether:

- Share the connection in the flow
- Require the user to provide their own connection

### Service Principle Owned Flows

We can set-up a service principle to own and run Power Automate flows. A Service
principle is an identity representing an application or service that can own and
manage resources within Azure and Power Platform.

We should use this feature for:

- Mission critical flows: Remove risk of flow owner leaving or having a licence
unassigned to run premium connectors
- Where devOps pipelines used to deploy flows

Since service principles cannot have licences assigned to them, a per flow
licence is needed if premium connectors are used.

This is a simple process demonstrated [here](./demos/pa_sp_set_up.md)

## Connections

A connector is a wrapper around an API, when a connection is created and there
is authentication set up, the credentials must be provided to establish the
connection.

Solution aware canvas apps, and operations in a solution aware flow are bound to
a connection reference rather than a connection. When a solution is imported
into a target environment, a connection will need to be created for all
connection references.

Flows created outside of a solution use connections directly. We can update a
flow to used connection references by:

- Exporting and then importing the solution
- Selecting the Remove connections so connection references can be added action
recommended by the flow checker

## Troubleshoot Flows

We can see a summary of all cloud flow activity from make.powerautomate by
accessing More -> Cloud flow activity from the side bar.

### Authentication Failures

Authentication errors will generally involve 401 or 403 errors. Generally, we
can resolve these by fixing the connection.

### Action Configuration Issues

These issues will generally be indicated by 404 or 404 errors, we should
investigate the error message to determine how to fix the flow

### Temporary Issues

If an error code of 500 or 502 appears the issue is likely to be temporary or
transient. In this case, resubmit the flow to try and run it again.

### Pricing Plan Issues

Some issues with flows may be due to the pricing plan. For instance, we may
have run out of data. Investigate the plan and the associated usage limits.
Similarly, the plan may have invocation limits, for instance, every 15 minutes
on the free plan.

### General Flow Limitations

Each account can have up to:

- 600 flows
- 50 custom connectors
- 20 connections per API and 100 connections in total

In addition, some connectors, e.g. C, implement connection throttling.

## Expressions in Power Automate

For more complex operations, e.g. calculations and data transformation, we will
need to use expressions. FLows in Power Automate run on top of Azure Logic Apps
and both use the same functions.

### Writing Expressions

When we click on an input, an icon will show to open the expression editor. We
can type expressions here and can use dynamic content within these expressions.

Expressions can become quite long, a common practice is to add a note to the
step and paste the expression to improve the visibility.

### Defining Text

String literals are defined with single quotes, e.g.:

```exp
concat('Today is', outputs('compose'))
```

### Function Types

There are 10 different function types including:

#### Math Functions

- rand(1,10): Random number from 1 to 10
- add(12,13)

#### Referencing Functions

These functions are used to work with the outputs of actions and triggers. Often
these functions are written for you when you use dynamic content, but we can
hover over the content to see the underlying expression:

```exp
triggerOutputs()?['headers']?['x-ms-user-name-encoded']
```

#### Workflow Functions

These are used to retrieve information about the flow

- workflow().run.id

#### Manipulation functions

These are various functions in this category, one useful function is:

- coalesce(null, 'Power Automate', 'Power Apps'): Returns first non-null value

### Complex Expressions

A complex expression is one that combines multiple functions:

```exp
addDays(triggerBody()?['date'], triggerBody()?['number'])
```

## Manage Sensitive Input and Output Parameters

In the settings for a trigger/action we can enable secure inputs and secure
outputs. When enabled, the inputs/outputs will not be shown in the logs for flow
runs.

### Using Azure Key Vault

Marking inputs and outputs as secure will remove them from the logs. However, if
we edit the flow hardcoded values will still be visible. In addition, if these
values expire then we will need to update them everywhere that they are used.

To resolve this we can use Azure Key Vault to store the values and access these
with the Azure Key Vault connector. We can create a connection as an interactive
user or, preferably, by using a service principle.

The use of Azure key vault authenticated as both an interactive user and a
service principle is demonstrated [here](./demos/pa_key_vault.md).

## Error Handling in Power Automate Flows

### Run After Settings

Actions and controls have a Run after section in their settings. By default,
when we add an action to a flow this is set to run after:

- The previous action/control
- If that action/control ran successfully

If we want to add error handling to a flow, we should add one or more parallel
branches to the success branch to handle:

- Has timed out
- Is Skipped
- Has failed

We can select one or more of these conditions to handle all errors the same or
provide more granular error handling as required.

### Scopes

We can use a scope control to group related actions. A common policy is to
define try, catch and finally scopes to handle errors. Try would represent the
success path.

Catch, set to run after try if there are any failures, would contain error
handling logic.

We could also add a finally scope which would run after success or failure.

### Retry Policies

We can use retry policies to handle transient failures, e.g. due to network
issues or service unavailability.

There are four options:

- None
- Default
- Exponential Interval
- Fixed Interval

The default is an exponential interval set to retry 5 times. Time periods are
set using the ISO-8601 "Period Time" format, e.g.

- PT2H: 2 hours
- PT30M: 30 minutes
- PT5S: 5 seconds
- P3Y6M4DT12H30M5S: 3 years, 6 months, 4 days, 12 hours, 30 mins and 5 seconds

### Terminate Action

The terminate action can be used to end the flow and set a status:

- Succeeded
- Failed
- Cancelled

This is used to explicitly end a flow due to an error.

### Analytics

Each flow has an analytics tab which can provide a view of:

- Action requests
- Usage
- Errors

The errors dashboard provides a break down of errors by error status code.

As with Power Apps, we can also set up Application Insights to work with Power
Automate for detailed analytics. This also allows for alerts to be set-up to
increase awareness of errors.

## Child Flows

We can create child flows to break-up and reuse functionality. Creating a child
flow is easily, we need to ensure that:

- The trigger is Manually trigger a flow
- The flow should have a Response or Respond to Power App or Flow action to
return data

From a parent flow, we can use the Run a Child Flow action to run a given child
flow.

The process is very simple, however, the child flow and parent flow MUST be in
the same solution.
