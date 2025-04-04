# Microsoft Entra ID

## Intro

### Active Directory Domain Services

This has a directory structure. There were organisational unit directories that
would contain user and device objects. These objects would have attributes such
as first and last name. This was used to manage users and devices on a Windows
server. This was used in the 90s and 00s.

### Entra ID

If Active Directory is a house, Entra ID is a hotel. You can have tenancies
within the hotel and use the services. You pay a subscription, as long as this
is paid you can stay within the tenant. Entra ID provides Identity as a service.

Microsoft takes responsibility for the directory, however, we are responsible
for our tenant, users and devices.

Once authenticated, you have access to the various SAAS offerings. Identity is
the key to everything.

### Central Identity System

Entra ID is a Centralised Identity Management System. It is used to store and
managed credentials to provide authentication and authorisation capabilities:

- Credentials are verified when stored
- Management is by a single authority (admin or admin group)
- Used for identity and access management

## entra.microsoft.com

We can manage tenants, users etc.

Every tenant has a unique tenant id, this is how we refer to a tenant, e.g. with
scripts.

### Users

We can create users from the users tab. Note, we cannot licence users from Entra
Id, so this would need to be done from the 365 portal.

An important concept is Role-Based Access Control (RBAC). We can assign security
roles to security roles to manage permissions. We should follow the principle of
least privilege when assigning roles.

There are member users and external users. We can create external users
in entra id by inviting them with an email. External users can be either guests
or members.

## App Registrations

App registrations can be created with permissions. For credentials we can use:

- Certificates
- Client Secrets
- Federated Credentials

These expire and this is something that needs to be managed. Within an instance
of an App registration we can manage things like users and provisioning. The
object ID will be the actor in terms of access control.

## Managed Identities

With managed identities, microsoft manage the app registration for us. We do not
need to work with certificates and secrets.

On many Azure resources, there is an identity tab. There are two types, system
assigned and user assigned.

### System Assigned

Once created we will get the Object Id for the service principle. Behind the
scenes, an App registration is created and the certificate expiry status will be
managed by Microsoft.

The identity is assigned to the resource, so it is tied specifically to the
resource and cannot be used by any other resource. If the resource is deleted,
the service principle will also be deleted.

### User Assigned Managed Identity

System assigned managed identities are:

- Tied to a given app (limited if we want to use a single identity for multiple
apps)

With user assigned managed identities, we can assign multiple apps to the same
identity.
