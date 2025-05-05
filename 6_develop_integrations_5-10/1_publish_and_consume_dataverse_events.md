# Publish and Consume Dataverse Events

## Dataverse Azure Solutions

Dataverse, as part of the Power Platform, can be easily integrated with Azure
solutions.

### Azure Logic Apps

Logic apps has a similar UI to Power Automate flows. In both cases it is
possible to use pre-built and custom connectors, including a Dataverse
connector.

Use to create workflows and orchestrate processes to connect both cloud and
on-premises services.

### Azure Service Bus

This is a cloud messaging as a service (MaaS) platform. Messages can be sent and
received from queues. There is also a publish-subscribe mechanism using the
Topics feature.

Use to connect on-premises and cloud-based applications and services with
messaging workflows.

### Azure API Management

This can be used to manage custom APIs built to use with Power Platform. The
logic in the APIs can work with internal data, external services or the
Dataverse API. We can export the API definition to create a Power Platform
custom connector

Use to publish an API for internal and external developers to use when
connecting to backend systems hosted anywhere.

### Event Grid

This is an event-driven, publish-subscribe framework allowing us to handle
various events. Unlike Event-Hubs and Service Bus, Dataverse has no out-of-the-
box integrations with Event Grid. However, it should be considered if event-
driven integrations are needed.

Use to connect supported Azure and third party solutions with a fully managed
event-routing service using a publish-subscribe model.

### Event Hubs

Event hubs is the MS version of Apache Kafka. This provides a real-time data
ingestion service that supports millions of events per second. It is good for
large data streams that need to be ingested in real time, e.g. application
telemetry data or IoT data. Typically, in the context, of dataverse, this is
generally used as an analytics solution rather than an integration solution. We
can publish events to an event hub with Dynamics.

## Expose Dataverse Data to Azure Service Bus

### Register a Service Bus Endpoint

Registering a service endpoint is a simple process. We can register the endpoint
in PRT by providing a service bus connection string. This is used to
authenticate using a Shared Access Signature (SAS) key passed with the
connection string

### Designation Types

When registering a service endpoint, there are a number of designation types
which may be selected. These are "Service Bus Contracts".

#### Queue

With a queue, there is no requirement for an active listener to send an event.
We can consume the messages at anytime using a destructive or non-destructive
read:

- Destructive: Removes the message from the queue after reading
- Non-destructive: Read without removing the message from the queue

#### Topic

This is similar to a queue. However, topics support a publish-subscribe model,
so listeners can subscribe to receive messages for a given topic. This is useful
if multiple consumers needed for a given message

#### One-Way

With a one-way contact, there must be an active event listener to consume the
message posted. If no active listener is available then the post will fail.
Dataverse will repost periodically until the job is eventually cancelled with a
status of failed

#### Two-Way

This is similar to a one-way contact, however, it is possible to return a string
value from the listener. For instance, if a custom Azure-aware plug-in posts the
message, then the data returned can be consumed by the plug-in.

For example, we might use this to retrieve the ID of a row created in an
external system to then store this in the Dataverse environment

#### REST

This is like a two-way contract except that we publish to a REST endpoint

### Posting Data to Service Bus

To post to the service bus, we can register a step directly on the service bus
endpoint. Note that this endpoint only supports asynchronous steps.

The process is demonstrated [here](./demos/service_endpoint_and_service_bus.md).

#### Using IServiceEndpointNotificationService

We can also use IServiceEndpointNotificationService to call a service endpoint
using it's guid and passing in the context. I was unable to access
IServiceEndpointNotificationService from the serviceProvider when testing this
interface.

[ms docs]<https://learn.microsoft.com/en-us/power-apps/developer/data-platform/write-custom-azure-aware-plugin>

## Publish events with Webhooks

### Webhooks vs Azure Service Bus

Azure Service Bus should be used when:

- High-scale async processing/queuing is required
- Decouples applications and protects from temporary peak
- Multiple subscribers may need to consume an event

Webhooks should be used when:

- Synchronous processing against an external system
- External operation needs to occur immediately
- Entire transaction should fail unless external service returns successfully
- Third party Web API endpoint already exists to be used for integration
purposes
- Shared Access Signatures are not preferred or feasible

### Webhook Authentication Options

- HttpHeader
- WebhookKey: Passed in the query string with code as the key
- HttpQueryString: Pass one or more key value pairs in query string

### Using Azure Functions with Webhooks

We can use Azure functions with webhooks to perform synchronous operations in a
similar manner to plugins. This can reduce load on Dataverse, however, since
Azure Functions do not run in the Dataverse event pipeline:

- Data needs to be updated in the most high-performing manner, e.g.
autoformatting a string before posting to Dataverse
- Any data operations in the Azure function will not roll back

[webhook Demo](./demos/webhook_demo.md)
