# Design Solution Components

## Intro

Questions here will be on design decisions when implementing components. For
instance, the appropriate pipeline stage for a plug-in. The information to make
these decisions is, for the most part, contained in the documents for those
components.

This document mostly just points to content elsewhere in this repo. However,
there is some additional content.

## Design Power Apps reusable components

- [PCF controls](../1_extend_user_experience_10_15/2_power_apps_component_framework.md)
- [Client scripting](../1_extend_user_experience_10_15/1_client_scripting.md)

### Canvas App Components

We can create components within a canvas app. This is a good way to separate
some of the functionality. The process is very simple and not worth covering
here.

However, if we want to reuse components across applications, then we will need
to use a component library. This is the recommended way to reuse components
across apps, an app will maintain dependencies on the component used and will
be alerted if updates become available.

Importing components from one canvas app to another has been retired.

Creating a component in a component library is the same as when creating a
component in an app, there are a few differences though:

- We can enable "Allow customisation" in the components properties. This is
enabled by default. When the component is edited a local copy is created and
the association to the component library is removed.
- To add to an app, we need to import the component into the app from the
component library. If it cannot be seen in the list, ensure that it has been
shared with the app maker.
- Any updates to the component must be published to be used in existing apps
- Component libraries cannot be deleted when it is referenced by a canvas app
the dependency/ies must be removed first

## Design custom connectors

- [custom connectors](../2_extend_the_platform_40-45/4_custom_connectors.md)

## Design Dataverse code components

- [Plug-ins](../2_extend_the_platform_40-45/1_plug_ins.md)
- [Custom Apis](../2_extend_the_platform_40-45/2_custom_apis.md)

## Design automations including Power Automate cloud flows

- [Power Automate Flows](../2_extend_the_platform_40-45/6_power_automate_flows.md)

## Design inbound and outbound integrations using Dataverse and Azure

### Inbound

- [Dataverse Web Api](../2_extend_the_platform_40-45/5_platform_apis.md)
- [Alternate keys and upsert](../2_extend_the_platform_40-45/1_plug_ins.md)
- [Custom apis](../2_extend_the_platform_40-45/2_custom_apis.md)
- [Azure functions](../2_extend_the_platform_40-45/3_azure_functions.md)

#### Inbound Patterns

##### Processing an Incoming Message

Let's say we have an external application that needs to notify Power Apps to
create a new contact. A pattern we might use here would be:

- External application adds a message to a service bus queue
- An Azure function listens for a message and use the .Net sdk to perform the
operation using managed identities for authentication to remove need for secrets

The related service principle we need the appropriate permissions in Dataverse.

We may also want to implement a plugin to validate the contacts created and
alternate keys to help prevent duplicates.

An alternate method for implementing this pattern would be to use a Power
Automate flow with a Dataverse Connector and a Service Bus Connector.

We could also use a http trigger for the Azure Function rather than a service
bus queue, but we would lose the automatic retry and load balancing
functionality provided by service bus

##### Bulk Import of Records

Let's say that we have an excel file containing data that we need to bulk import
into the Dataverse. One solution would be to:

- Import input file into a Azure Sql database as a staging table
- Use a logic app with a Dataverse connector to read rows from the table and
write to Dataverse
- Azure Data Factory can be used to manage the pipeline

##### Embed External Data in Model-Driven Apps

Before importing data into dataverse we should consider whether we can expose
external data in the application in real time. We show external data using:

- PCF Controls
- Embedded Power BI reports

We may also consider defining virtual tables.

### Outbound

- [Event publishing](../6_develop_integrations_5-10/1_publish_and_consume_dataverse_events.md)
- [Change Tracking](../6_develop_integrations_5-10/2_implement_data_synchronisation_with_dataverse.md)

#### Event Model

Dataverse has an event model to integrate with other systems. The event model
can trigger:

- Plug-ins
- Classic workflows and custom workflow activities
- Power Automate cloud flows
- Messages to Azure Service Bus and Event Hub
- Webhooks

##### Plug-ins

Plugins and custom workflow activities, registered in sandbox mode, can access
the network with http/https to call external services. There are some
limitations here:

- A named web address must be used not an IP
- There is no provision for prompting a user for credentials or saving them
- Connections must be allowed from Power Platform and Dynamics 365 services

Service protection limits are not applied directly to plugins or custom workflow
limits, but computation time will be added to the request that triggered them.

##### Azure Relay

This can be used for one way, two way and REST contracts.

##### Azure Event Hub

Big data streaming platform and data ingestion service. It can receive and
process millions of events

##### Azure Service Bus

Can be used to post and listen for messages. Useful when there is a likely hood
that a listener is not available or there is a high volume of messages. We can
use:

- Queues: Deliver message to one or more competing consumers
- Topics: Like queues but one or more listeners can receive messages for a given
topic

We can post plugin context to the service bus without code using a service
endpoint with a step. We can also use code to call the service bus endpoint, we
might use this to add data to the shared variables in the context.

Security is managed using Shared Access Signatures (SAS)

##### Webhook

Webhooks are a HTTP pattern for connecting web apis and services with a publish-
subscribe models. Webhooks enable both synchronous and asynchronous steps. As
wit Service Endpoint we can register a webhook and step with no code or use a
plugin/custom workflow to trigger the webhook.

Webhooks will timeout after 60 seconds.

##### Batch with Azure Synapse Link for Dataverse

We can use this service to continuously export data from Dataverse to Azure
Data Lake Storage and Azure Synapse Analytics. It is designed for enterprise
big data analytics. Any actions or metadata changes where change tracking
enabled will automatically pushed.

The service is integrated with the maker portal so it is simple to configure.
