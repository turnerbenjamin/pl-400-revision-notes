# Dataverse Plugins

## Introduction

Plugins allow us to run dotnet code in response to Dataverse triggers. These are
efficient and powerful but they require developer skills and can detriment
performance if poorly written.

We should consider alternatives such as:

- Power Automate Flows
- Custom Actions
- Calculated and Roll-Up Fields
- Azure Service Bus Integration and Webhooks

## Scaffolding a Plugin

We can use:

- pac plugin init: Includes a plug-in base class. Complex but easy to work with
- dotnet new classlib: Simple but more boilerplate required

Note the plugin base class in resources was not generated with pac. It is a
simple class used to share common boilerplate error handling logic between the
various demo plugins in the assembly.

If using VS use the class library template.

Note that we must use .NET 4.6.2 to write plug-ins

```xml
<TargetFramework>net462</TargetFramework>
```

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
appropriate type and perform null checks before using. This can be seen from
usage in the resources for this section.

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

We can use these overloads to identify records by their alternate keys:

```cs
// Update with shorthand
// *********************

var entityToFind = new Entity(entityLogicalName, ak, akv);
orgSvc.Update(entityToFind);

// Retrieve with execute (simple alternative key)
// **********************************************

var req = new RetrieveRequest(){
    Target = new EntityReference(entityLogicalName, ak, akv)
};
orgSvc.Execute(req);

// Retrieve with execute (compound alternative key)
// ************************************************

var req = new RetrieveRequest()
{
    Target = new EntityReference(
        entityLogicalName,
        new KeyAttributeCollection()
        {
            new KeyValuePair<string, object>(ak, akv),
            new KeyValuePair<string, object>(ak2, ak2v),
        }
    ),
};
orgSvc.Execute(req);
```

Note: A guid will always be present. However, keys may be deleted. This is a
disadvantage of alternate keys.

## Demonstrate Use of Different Event Execution Pipeline Stages

### What we Need to Know

- pre-validation
- pre-operation
- post-operation

## Implement Business Logic

### What we Need to Know

- Use for complex logic
  - custom workflow activities which may be reused in workflows
  - custom actions can create reusable messages that may be called from other
    workflows or web service endpoints
  - Service Bus integration and webhooks can be used to push data to external
  systems
- Consider alternatives
  - calculated fields
  - rollup fields
  - power automate

## Operations using Organisation Service

### What we Need to Know

- retrieve
- create
- update
- deleting
- operations on related entities
- custom actions
- optimisation with retrieve using column sets

## Optimise Plug-In Performance

### What we Need to Know

- Only get data needed
- Filter columns in steps
- Use pre and post images
- Register only where required
- Minimise speed (note 2-minute limit)
- Check transaction depth to avoid infinite loops
- Ensure comprehensive handling and catch IPluginExecutionFault

## Configure Dataverse Custom API Message

### What we Need to Know

- Allows us to build new messages in the pipeline
- Define the API, request parameters and response parameters
- They may be global or bound to a table
- May be called from code or power automate

- create with
  - Maker portal
  - PRT
  - Code

## Register Plugins Using the PRT

### What we Need to Know

- Understand how to register an assembly
- Register steps
- Filter attributes
- Stages and whether sync or async

## Develop a Plugin that Implements a Custom API

### What we Need to Know

- Add plug-in to perform logic when custom API called
- Register the assembly
- Rather than register a step, link plugin assembly to the custom API

## Configure Dataverse Business Events

### What we Need to Know

Historically plug-ins used for create and update of records. But we may have
multiple events at once, e.g. an invoice with various line items.

We can use business events, e.g. post invoice. There is a catalogue of events
from which we can choose.

We can then handle a single transaction with a single event which simplifies the
logic.

We need to look into this.
