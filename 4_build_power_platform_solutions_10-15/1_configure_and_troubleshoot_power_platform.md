# Configure and Troubleshoot Power Platform

## Complete Testing and Performance Checks in a Canvas App

The content for this requirement has been covered
[here](../5_implement_power_apps_improvements_10-15/2_optimise_and_troubleshoot_apps.md)

## Troubleshoot Operational Security Issues Found in Testing

This primarily involves an understanding of the security set-set up in the
environment. This is covered
[here](../3_create_a_technical_design_10-15/1_design_technical_architecture.md)
We need to understand authentication, authorisation and data loss prevention
policies.

### Threat Detection and Monitoring

Monitoring is about gathering information on events that have occurred to gain
an awareness of suspicious activities and use these to predict future security
incidents. The purpose of security monitoring is threat detection.

#### Dataverse Auditing

Dataverse auditing is designed to meet internal auditing, compliance, security
and governance policies. Auditing logs changes made to customer records in an
environment with a Dataverse Database and access through an app or SDK in the
environment.

Auditing is supported on all custom and must customisable tables and columns.
Audit logs are stored in Dataverse and consume log storage capacity. They can
be viewed from the Audit History tab for a single record and in the Audit
summary view for all audited operations in a given environment.

##### Configure Auditing

Auditing can be configured at three levels:

- Environment
- Table
- Column

First, we must turn on auditing at the environment level. To log changes in a
table auditing must be turned on at both the table and the column level.

To enable user access auditing or activity logging, we must enable these at the
environment level. We can also enable auditing for security roles.

To configure auditing we need the System Administrator or Customiser role.

We can define a retention policy to configure the time that logs are kept in an
environment. Changing this will not affect logs created before the change as
each log is individually timestamped with a retention period. We can also
manually delete logs if we need to free up log capacity.

#### Microsoft Sentinel

We can use the Sentinel Solution with Power Platform to detect suspicious
activities.

#### Microsoft Purview Activity Logging

Power Apps, Power Automate, Connectors, Data Loss Prevention and Power Platform
administrative activity logging are tracked and viewed from the Microsoft
Purview compliance portal.

## Configure Dataverse Security Roles to support Code Components

This is all covered
[here](../3_create_a_technical_design_10-15/1_design_technical_architecture.md).
Just remember the principle of least privilege
