# Plug-In Performance

## Introduction

It is important to optimise plug-ins. Slow plug-ins can lead to:

- Slow response times, e.g. when submitting a form
- Generic SQL errors, generally caused by a SQL timeout
- Deadlocks
- Slow throughput in batch loading

Plug-ins are a powerful tool, but a badly designed plug-in could severely impact
the system. To mitigate this, the Dataverse Platform imposes constraints to
reduce the impact a user can have on the system. We need to design for these
constraints, i.e:

- Optimise to try and avoid them
- Implement error handling to recover when a constraint is triggered

Dataverse is not designed for long-running or batch processing. In these
situations separate services should be used which drive shorter transactional
requests to Dataverse. For instance:

- Using a flow
- Hosting MS SQL Server Integration Services to drive requests

## Dataverse Platform Constraints

- Plug-ins timeout after 2-minutes, they should not be used for long-running
operations
- Requests to SQL server also generally timeout after 2 minutes
- Service Protection Limits

Running into these constraints can indicate a problem with a plug-in or a series
of plug-ins.

### Service Protection Limits

There are a number of service protection limits:

- a user cannot make more than 6000 requests within a 5-minute period
- a user cannot make requests with a total execution time of 20 mins within a
5-minute period
- A user cannot make more than 51 concurrent requests

Service protection limits do not directly apply to plug-ins or custom workflow
activities. Requests made by a plug-in triggered by an application will not
count to service protection API limits. However, the computation time of these
operations are added to the initial request and this is part of the service
protection limits.

#### Retry Operations

!! TO BE ADDED

## Transactions and Locks

The protect the integrity of the database, operations are preformed within a
transaction. Each request will place locks on resources which will not be
released until the transaction has been committed or aborted. For instance:

- Create: Write lock against the record
- Update: Write lock against the record

### Transactions and Plugins

Whether or not a plug-in is within a transaction depends on the stage:

- pre-validation: Before the transaction is started
- pre-operation: Inside the transaction
- post-operation (sync): Inside the transaction
- post-operation (async): Outside the transaction

A plug-in registered in the pre-validation stage will be within a transaction if
the pipeline event was triggered by a message within an existing transaction.

With async post-operation plug-ins, no transaction is created. Each message is
acted on independently.

### Transactions and Workflows

- Sync workflows: Within the transaction
- Async workflows: Outside the transaction

Custom workflow activities act within the parent context.

### Transactions and Custom Actions

Custom actions may create their own transactions if:

- Enable rollback is set
- There is no existing transaction

### Transactions and Web Service Requests

With plug-ins, it is the execution context that maintains the execution context.
External requests with web services will create a pipeline with transaction
handling but the transaction will close once the response is returned. In this
scenario, there are two messages that may be used to perform multiple actions
within a single request:

#### ExecuteMultiple

Execute multiple is used to pass multiple independent requests. No transaction
context is held since the requests are independent

#### ExecuteTransaction

Here related actions can be processed in a single transaction. This should be
used with care as it has the potential to create blocking issues in the
platform.

### Concurrency and Transactions

When there are multiple concurrent requests, there is a higher chance of
collisions on locks.

#### Async Operations

Async workflows or plug-ins are not processed serially from a queue. Multiple
activities are processed by dataverse in parallel within a given service
instance and across service instances. Each async service instance retrieves
jobs in batches of around 20 depending on configuration and load.

Multiple asynchronous activities from a single event are likely to be processed
in parallel and, if they operate on the same resource lock conflict risk is
high.

## General Considerations

### Only Execute when Necessary

As noted above, we can specify filtering attributes in the step definition. This
can be used to filter the messages triggering plugin execution.

### Only Fetch Required Data

We should only fetch the data we require. For instance:

- Use column sets to limit data returned by the organisation service
- Specify only required data in pre and post images

### Use Pre and Post Images

Pre and post images should be used where appropriate to avoid making a separate
call to the organisation service to retrieve data.

### Use Depth to Avoid Infinite Loops

We can check the transaction depth with:

```cs
context.Depth
```

This can be used to detect infinite loops. For instance, we may throw an error
if depth exceeds a certain value so that problematic implementations may be
fixed.

## Performance Analysis

admin.powerplatform has an analytics page with a plug-ins tab. We can use this
to view:

- Pass Rate
- Executions
- Most active plug-ins
- Average execution time
- Top plug-ins by failure
