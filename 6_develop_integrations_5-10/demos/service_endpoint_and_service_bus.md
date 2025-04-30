# Publishing Dataverse Data To Service Bus with Service Endpoint Registration

This demo will create a service endpoint to publish dataverse data to an Azure
Service Bus queue.

## 1. Creating a Service Bus Queue

First, we need a service bus queue:

- Create a service bus in portal.azure
- Once created add a queue

![create service bus](./screens/service_endpoint/1_create_service_bus.png)

To authenticate with the Service Bus, we use Shared Access Signatures (SAS).
Select shared access policies and copy the primary connection string for the
RootManageSharedAccessKey.

![shared access signature](./screens/service_endpoint/2_shared_access_signature.png)

## 2. Register Service Endpoint

Next, open PRT to create a service endpoint. This will prompt us to add the
connection string for the endpoint:

![enter connection string](./screens/service_endpoint/3_set_connection_string.png)

### Designation Type

There are various designation types:

- Queue: Post to a messaging queue in the cloud. No active listener is required
- One-Way: This requires an active listener. If no active listener the post will
fail after a number of retries
- Two-Way: As above but a string value can be returned from the listener to the
plug-in or custom workflow activity initiating the post
- REST: Like two-way but on a REST endpoint
- Topic: Like queue but active listeners can subscribe to receive messages from
the topic
- Event Hub: Relates to event hub integrations

### Queue Name

The connection string relates to the service bus. However a service bus can
contain multiple queues. Accordingly it is necessary to specify the queue to
post to.

This parameter is only present if we select queue as the designation type

### Message Format

We can select:

- JSON
- XML
- .NETBinary

Generally, we will use JSON as the message format

### User Information Sent

We can select either none or userid.

![service endpoint creation](./screens/service_endpoint/4_service_endpoint_creation.png)

## 3. Register a Step

![step registration](./screens/service_endpoint/5_register_step.png)

## 4. Test the Integration

To test the integration, we can create an account to trigger the step.

![account creation](./screens/service_endpoint/6_account_creation.png)

This will then add a message to the queue with the record json:

![message view](./screens/service_endpoint/7_message_view.png)
