# Dataverse Plugins

## Introduction

Plugins allow us to run dotnet code in response to Dataverse triggers. These are
efficient and powerful but they require developer skills and can detriment
performance if poorly written.

## Scaffolding a Plugin

We can use:

- pac plugin init: Includes a plug-in base class. Complex but easy to work with
- dotnet new classlib: Simple but more boilerplate required

Note the plugin base class in resources was not generated with pac. It is a
simple class used to share common boilerplate error handling logic between the
various demo plugins in the assembly.

Note that we must use .NET 4.6.2 to write plug-ins:

```xml
<TargetFramework>net462</TargetFramework>
```

So, no global usings, nullable types or file-scoped namespaces.

### Required Packages

We will need the Microsoft.CrmSdk.CoreAssemblies package. This can be downloaded
from Nuget.

This includes the official MS:

- Xrm sdk dll
- Crm sdk proxy dll

## Coding the Plugin

The plug-in must implement the IPlugin interface from the Xrm.Sdk. This
interface has a single method:

```cs
public void Execute(IServiceProvider serviceProvider) { }
```

### IServiceProvider

IServiceProvider also exposes a single method:

```cs
public interface IServiceProvider
  {
    object GetService(Type serviceType);
  }
```

We can use this to access various services, for instance:

```cs
serviceProvider.GetService(typeof(IPluginExecutionContext))
serviceProvider.GetService(typeof(IOrganizationServiceFactory))
serviceProvider.GetService(typeof(ITracingService))
```

GetService returns an Object or null. We need to cast the return value to the
appropriate type and perform null checks before using.

In the resources, this boilerplate code is contained in a base class:

[Plug-In Base Class](./resources/DemoPlugins/PluginBase.cs)

### IPluginExecutionContext

This API provides context about the environment the plugin runs in, the
execution pipeline and entity information.

We can use this to access

- input, output and shared parameters
- pre and post images
- business unit and user id
- organisation name and id
- message name
- primary entity name and id
- stage (i.e. pipeline stage)

Both IPluginExecutionContext and IWorkflowContext implement IExecutionContext.

#### Input Parameters

With regards to input parameters, there will be a Target key for most messages
this can be used to access the primary entity:

```cs
context.InputParameters["Target"]
```

The entity returned contains key value pairs with the key being the logical
column name. It will not include all properties, just those that will change. If
we need access to additional properties then we will need to use pre/post
images.

#### Output and Shared Parameters

We can also use:

- Output parameters: Return values from the plugin
- Shared Variables: Persist data between plugin steps

Note: We can only use output parameters in the postOperation Stage.
Note: Shared variables may be useful in certain situations but plugins should
generally be stateless

#### Pre/Post Images

We can use pre/post images to capture all information required on the target
entity in the execution context. This improves performance; without images we
would need to make an API call with the Organisation Service to get that data.

We can use the PRT tool to register an image on a given step. When registering
an image we define the parameters that we want to be included, the default is to
include all parameters.

When we register an image we provide a name. This name is used as a key in the
pre/post images collection:

```cs
context.PostEntityImages["ImageName"]
```

Note:

- PreEntityImages will not be available on create
- PostEntityImages will not be available on delete

Use of input parameters and pre-entity images is demonstrated in the following
demo:

[Plug-in parameters demo](./resources/DemoPlugins/PluginParametersDemo.cs)

### IOrganisationServiceFactory

#### Web API vs Organisation Service API

There are two mechanisms for interacting with the Dataverse:

- Web API
- Organisation Service API

The Web API is generally preferred as it implements OData v4 and may be used
with a variety of languages. The intention is for WebApi to replace the
Organisation Service. But this will be a gradual transition.

We should use Web API for code running on premises. However, it is not currently
easy to use this API for code running on the server such as plug-ins or workflow
assemblies. We should use the Organisation Service in this context.

Note the Organisation Service does (or at least will) utilise Web Api under the
hood.

Note: At this point WebApi requests should NEVER be used in plugins.
IOrganisationService methods allow for the transaction context to be passed
enabling the operation to make requests within the pipeline transaction.

#### Instantiating IOrganisationService

As noted above we can get an instance of IOrganisationServiceFactory from the
service provider. We can use this to create an instance of IOrganisationService:

```cs
var context = serviceProvider.GetService(typeof(IPluginExecutionContext)) 
  as IPluginExecutionContext;

var orgSvcFactory = 
  serviceProvider.GetService(typeof(IOrganizationServiceFactory)) 
    as IOrganizationServiceFactory;

var orgSvs = orgSvcFactory.CreateOrganizationService(context.UserId);
```

## Interacting with Dataverse

### Introduction to IOrganisation Service

As noted, we should use OrganisationService to interact with the Dataverse in
the context of plug-ins. Similar to Xrm.WebApi there are a number of shorthand
methods for performing simple CRUD operations:

```cs
orgSvc.Create(entityToCreate);
orgSvc.Retrieve(entityLogicalName, guid, columnSet);
orgSvc.Update(entityToUpdate);
orgSvc.Delete(entityLogicalName, guid);

orgSvc.RetrieveMultiple(query)
orgSvc.Associate(entityLogicalName, guid, relationship, relatedEntities)
orgSvc.Disassociate(entityLogicalName, guid, relationship, relatedEntities)
```

[Shorthand Method demos](./resources/DemoPlugins/OrgSvcShorthandMethods.cs)

There is also an Execute method which can be used to create more bespoke
requests.

```cs
var createRequest = new CreateRequest() { Target = accountToCreate };
var response = orgSvc.Execute(createRequest) as CreateResponse;
return response.id;
```

Usage of this method is demonstrated in the upsert demo:

[Upsert demo](./resources/DemoPlugins/Upsert.cs)

### Retrieve Multiple

Retrieve multiple is more complex than the other shorthand queries. The query
parameter is of type QueryBase; this is an abstract class with a number of
concrete implementations:

- QueryByAttribute: Simple implementation, we can filter by attributes, order
results and use column sets
- Query Expression: As above but we can create more advanced filters and create
joins
- Fetch Expression: Queries are passed as xml strings. Gross and powerful, we
can perform aggregation and grouping with this implementation

[Retrieve multiple demo](./resources/DemoPlugins/RetrieveMultiple.cs)

### Upsert Request

An upsert request will:

- Update a record if it is present in the database
- Insert the record if it is not present in the database

This improves performance as we can perform the operation in a single trip.

There is no short hand method for an upsert so we must use the Execute method:

```cs
var upsertRequest = new UpsertRequest() { Target = entityToUpsert };
var response = orgSvc.Execute(upsertRequest) as UpsertResponse;
```

[Upsert demo](./resources/DemoPlugins/Upsert.cs)

### Entity and Keys

The methods above generally use a guid to identify a specific record. This is
either passed as a separate argument or stored as a parameter within entity.
However, we can also use both simple and compound alternative keys.

The entity constructor has 5 overloads:

```cs
Entity()
Entity(entityLogicalName)
Entity(entityLogicalName, guid)
Entity(entityLogicalName, keyName, keyValue) // Simple alternate keys (AKs)
Entity(entityLogicalName, keyAttributeCollection) // Compound and simple AKs
```

There is also an EntityReference class which has the same overload signatures.
This is used when we need to reference an Entity and do not require access to
its attributes.

We can use these overloads to identify records by their alternate keys.

Note: A guid will always be present. However, keys may be deleted. This is a
disadvantage of alternate keys.

[Alternate keys demo](./resources/DemoPlugins/AlternateKeys.cs)

### Error Handling

We should throw InvalidPluginExecutionExceptions from Plugins. In the resources
other exceptions are thrown from the concrete demo classes but these are nested
in an InvalidPluginExecutionException by the abstract base class.

Error messages will be displayed to users in Power Apps clients.

Avoid writing HTML in the message as the raw html will be displayed to users.

### Work with Concurrency

Power Apps is a multi-threaded and multi-user system. There is a need, therefore
to prevent race conditions. The information below relates only to:

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

#### Always Overwrite

The alternative to optimistic concurrency, is alwaysOverwrite. This will
overwrite the record regardless of the version number.

#### Default

The default behaviour will be always overwrite except:

- Where optimistic concurrency is enabled on the table, AND
- The source is WebService as opposed to a plug-in/custom workflow activity

#### Which Should You Use?

Assuming both behaviours are available, optimistic concurrency adds some
overhead but prioritises data integrity.

Always overwrite is faster but will be at the cost of data integrity where we do
not want the latest update to always take precedence.

## Register Plugins Using the PRT

The recommended method for deploying plugins is to use the Plugin Registration
Tool. This can be found in Nuget.

Alternatively, if you have the Power Apps CLI installed run:

```console
pac tool prt
```

### Strong Name Key

We can use a strong name key to sign the assembly. In visual studio this can be
done from the properties of the project.

We can also create a key by running:

```console
sn -k PACKAGE_NAME.snk
```

We then need to update the csproj file:

```xml
<PropertyGroup>
    <SignAssembly>true</SignAssembly>
    <AssemblyOriginatorKeyFile>VerySimplePlugin.snk</AssemblyOriginatorKeyFile>
</PropertyGroup>
```

### Register an Assembly

Once connected to the environment in PRT, we can register an assembly.
Registering is simple, provide the dll and select the plugin classes to
register.

### Register Steps

Once the assembly has been registered we need to define when the plugin will be
executed. This is achieved by registering one or more steps for the plugin.

The main fields to note are:

- Message
- Primary Entity
- Filtering Attributes

Note that the filtering attributes defaults to all attributes. This should
generally be edited to relevant attributes to avoid unnecessary runs.

#### Event Pipeline Stages

There are three pipeline stages:

- PreValidation: Before the transaction and any security checks (SYNC)
- PreOperation: Within the transaction but prior to write (SYNC)
- PostOperation: After the transaction (SYNC | ASYNC)

We should use the pre-validation stage for any validation logic and if we want
to conditionally cancel the transaction. We can cancel an operation by throwing
an InvalidPluginExecutionException within the plugin.

Pre-operation should be used if we want to change and of the table values before
it is saved.

Post-operation can be used to modify and message properties before the response
is returned. Be careful about entering an infinite loop.

NOTE: We can rollback in the pre operation stage but this does take time. If we
are looking to cancel a transaction use pre-validation.

NOTE: We cannot roll back in the post-operation stage. We should also be
cautious when updating values as this will trigger an update message.

## Tracing Service

Use the tracing service to get visibility into plugin runs. The tracing service
can be accessed from the service provider:

```cs
serviceProvider.getService(typeof(ITracingService)) as ITracingService
```

The Trace method may be used to add a message to the tracing service. Any
exceptions thrown will also be recorded in the tracing service.

Note:

- By default, a bulk execution job deletes trace logs every 24 hours
- Logging is asynchronous
- Trace logs for plugins need to be enabled in the advanced settings

The shorthand methods demo uses the tracing service extensively:

[Shorthand Methods Demo](./resources/DemoPlugins/OrgSvcShorthandMethods.cs)

## Debugging Plugins

We can use the profiler to debug plugins. This is very straight forward.

- Install and start the profiler on a step with the PRT tool
- Trigger the step manually to capture a profile
- Select debug in the PRT and select the relevant profile and assembly
- In VS attach debugger to the PRT process and add breakpoints
- Select start execution in the PRT
