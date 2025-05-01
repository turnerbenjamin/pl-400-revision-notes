# Create a Technical Design

## Intro

Power Platform is a low-code platform with 5 main components:

- Power Apps
- Power Automate
- Power BI
- Power Pages
- Copilot Studio

These components can be used individually but they are also deeply integrated.

The general approach is focussed on low-code to enable people with different
experience levels to contribute to development. Where there are capability gaps
or more complex requirements, the platform can be extended by professional
developers using the platform's extensibility model and native Azure
integration.

To develop efficiently on the platform it is important to understand the OOTB
features and how/when to extend these.

## Power Apps Applications

### Power Apps OOTB Functionality

We can develop canvas and model-driven applications to be consumed by users.

#### Model-Driven Apps

Model-driven applications are data-driven applications built on top of
Dataverse. We can quickly create an application to work with dataverse data. We
can also work with external data sources by embedding a canvas app with
connectors to those sources.

#### Canvas Apps

We can develop canvas applications for more control over the data sources and
services we use in an application. Canvas apps can be embedded into SharePoint,
Teams, Power BI and Dynamics 365 applications.

We can also implement logic using the low-code Power Fx language.

#### Business Rules

Business rules are a declarative solution which people of any experience level
can implement. We can use these rules to:

- Validate data and show error messages
- Set and clear column values
- Set column requirement levels
- Set column visibility
- Enable or disable columns
- Create business recommendations based on business intelligence.

These are demonstrated [here](./demos/business_rules.md)

Business rules are useful, but there is a long list of limitations:

- Actions either do not work, or do not work as expected, with canvas apps
- Not all actions work in editable grids and business rules do not work with
editable sub-grids at all
- Choices (multi-select), File and Language column types not supported
- It is not possible to access data in related tables
- Business rules run client-side, so validation  could be bypassed, e.g. by
using an editable grid
- They can quickly become unwieldy with more complex logic

The alternatives to business rules are:

- Plugins: For data validation logic
- Client Scripts: For display logic

Note that business rules run BEFORE the onLoad event. If an exam question
includes an interaction between client script handlers and values set by
business rules, look out for a gotcha question based on this fact.

### Power Apps Points of Extensibility

- PCF Controls
- Client Scripting
- Custom Connectors
- HTML web resources

## Dataverse

Dataverse is a cloud data store. This allows us to securely store and managed
data within tables. Dataverse is built on top of various technologies
including, Azure SQL, Cosmos DB and Azure storage.

There is a build in security model that allows access to data to be controlled
based on a user's roles (RBAC). We can define security at the table, row and
column level.

### Dataverse OOTB

#### Formula, Calculated and Roll-Up Columns

These are simple and performant ways to implement logic to created columns based
on calculations and aggregations. If any of these solutions does not meet
requirements, then a plug-in will generally be the most suitable alternative.

A demo of these capabilities can be found
[here](./demos/calculated_formula_and_rollup_columns.md)

##### Calculated Columns

These columns are used to perform calculations, we can specify conditions and
actions. In the action section we can define a calculation using data in the
table and related tables.

These are being depreciated in favour of formula columns:

##### Formula Columns

Similar to calculated columns, we can perform real time calculations on data in
the table and related tables. Unlike Calculated columns, formulas are expressed
in PowerFx.

Both calculated and formula columns have limitations, for instance:

- Neither can have cyclical references
- Sorting is disabled where there is reference to another table, a logical or
calculated column or Now/UtcNow is used
- There is a maximum depth of 5 and 10 respectively

Calculated and formula columns only run logic on retrieve. The values of the
columns will not be updated to reflect dirty data on a form. If we would like
the values to be updated dynamically then a client side script may be considered
as an alternative or a complimentary solution.

##### Rollup Columns

These are used to aggregate data from a related table. For instance, an account
table might have a 1:N relationship with an expenses table. We could use a
rollup column to display a total of all related expenses in the account column.

Unlike calculated and formula columns, these run on a schedule (default is every
12 hours)

#### Workflows

Workflows and Power Automate flows can be used to implement logic responding to
Dataverse messages. Both components are no-code methods of building such logic.

Workflows and plugins have different capabilities and limitations:

- Workflows can run either synchronously or asynchronously. Power Automate flows
can only run asynchronously
- Workflows cannot not be scheduled
- Workflows cannot access external data sources with connectors

Imperative alternatives include:

- Plugins
- Azure Functions

In both cases we can run the functions on demand, using a custom API (plugins)
or a webhook (Azure functions). Generally, plugins will be used to extend
dataverse, but we may consider Azure functions:

- For batch processing or long-running high volume actions
- To reduce the load on dataverse
- For scheduled tasks

Note that Azure functions will never take place within a transaction so we lose
the benefit of rollback here.

##### Long-Running Tasks

To expand on the point about long running actions, we should only use
synchronous plug-ins and workflows for short running tasks:

- They are blocking
- They generally take place in a transaction and lock resources

Async plug-ins may be used for medium length actions as they take place outside
of a transaction, however they will timeout after 2 minutes.

Azure functions, particularly durable functions are better suited for long
running tasks. Azure functions timeout, by default after 5 minutes but this can
be increased to 10 minutes. With the durable functions extension we can have
tasks that run for much longer than this by building an orchestration involving
multiple Azure functions.

### Dataverse Points of Extensibility

- Virtual Tables
- Elastic Tables
- Plug-Ins
- Webhooks and service endpoints
- Custom APIS

## Power Automate

Power automate allows use to automate tasks and orchestrate activities. With
connectors, we can integrate with external services. The main points of
extensibility are:

- Using custom connectors
- Using workflow definition functions to build complex expressions

## Extending Power Platform with Azure

A number of extension points allow integration with the cloud services available
through Azure.

### Azure Functions

Azure functions may be used to offload logic outside of Dataverse. We can use
Custom Connectors to integrate Dataverse with these functions. For Azure
functions with a HTTP trigger we can integrate using:

- Custom connectors
- Calling the endpoint in code, e.g. a plug-in or client script

### API Management (APIM)

APIM allows us to manage APIS across clouds and on premises. We can also use
APIM to export API definitions directly to Power Platform. When exported a
custom connector will automatically be created.

### Service Bus

This is a Messaging as a Service (Maas) framework for real-time, asynchronous
messaging. This can be a valuable feature to integrate with external systems in
a serverless, distributed fashion.

Dataverse can be set-up to publish events to Azure service bus queues and topics
automatically using service endpoints or in code.

### Event Grid

Event grid is a managed service for managing event routing from any source to
any destination. We can use this to route events between Power Platform and
other Azure services like Azure Functions

### Logic Apps

Logic apps are a cloud service we can use to automate and orchestrate tasks. The
Power Automate service is built on top of Logic Apps and includes integration
with Power Apps and Dataverse.

We can use Logic Apps to compliment Power Automate where Power Automate does not
meet requirements, for example, to use the Enterprise integration pack or SOAP
connectors.

### AI Services

Azure AI services is a collection of AI and cognitive APIs to add AI
functionality to apps. Power Platform has an AI Builder service which is a low
code option for some of the APIS. For more complex requirements we can use
Azure AI services directly.

## Extend Power Platform with Developer Tools

Official tooling is provided through:

- Power Platform CLI (pac)
- Nuget Package Manager

There are also community tools such as XRM toolbox

### Power Platform CLI (pac)

This is a very useful tool that enables us to:

- launch other tools such as PRT and CMT
- Initialise components such as plugins and PCF controls
- Work with solutions and packages
- Deploy packages, solutions and PCF controls
- Work with custom connectors
- etc

There is also the package deployer and solution packager tools, the
functionality of these tools can be accessed with pac.

A demonstration, using the pac cli and Configuration Migration Tool can be found
[here](./demos/pac_tool_demo.md)

## Security

### Authentication and Authorisation

Authentication is controlled using Entra ID. Authorisation is controlled with
licences, security roles and the sharing of apps and flows.

Creation of apps and flows is controlled by security roles.

The ability to see and use apps is controlled by sharing:

- Canvas Apps: Share directly with user or Entra group and subject to Dataverse
security roles
- Model-Driven Apps: Shared via Dataverse security roles

To manage security settings in Dataverse you MUST be a system administrator.

### Business Units

Business units are a building block of security modelling. We can use these
to manage users and the data that they can access. Every environment will have
a single root business unit. We can construct a hierarchy of business units
beneath this.

We can assign security roles, and column level security profiles to business
units and set row level security based on business units.

For each business unit, an associated team is created. We can associate this
team with an MS security group to help manage user administration and roles
assignment.

Each record has an Owning business unit column determining which business unit
owns the record. This defaults to the user's business unit when a record is
created.

#### Modernised Business Units/Matrix Organisation

In the traditional business unit hierarchy, users are restricted to their
respective business units. With modernised business units, users can navigate
and collaborate across business units.

The term matrix organisation may be used here to contrast with the traditional
hierarchical approach.

A user may still belong to a single business unit, but they can have security
roles from multiple business units. The user or admin, is also enabled to set
the business unit ownership of a record.

To use this feature it needs to be enabled in the environment settings in the
feature tab. The control is Record ownership across business units. Once enabled
there is a drop down to select the business unit when assigning security roles.

When a user is moved between business units, there are two helper settings:

- DoNotRemoveRolesOnChangeBusinessUnit: Default false
- AlwaysMoveRecordToOwnerBusinessUnit: Default is true

### Teams

Teams are another security building block. Teams are owned by a business unit,
however, they can include members from other business units.  There are three
team types:

- Owning teams: Can own records which gives any team member direct access to
that record.
- Access Teams: These are used with record sharing
- Entra group teams: Owner teams corresponding to an Entra security or office
group type

### System and Application Users

When a system is provisioned, a collection of special system and application
users will be created:

- System users created for integration and support
- Application users created for setup and configuration management

We should NOT delete or modify these users.

We  can add additional Application users to perform back-end services. Their
data access is managed by the security role that is assigned.

### Security Groups

If we apply a security group to an environment, then security roles in the
environment can only be applied to users in that group.

### Table/Record Ownership

There are two types of record ownership:

- organisation owned
- User/Team owned

This is defined when the table is created and cannot be changed.

With records that are organisation owned, the only access level choices are
whether the user can perform and operation or not. For records that are user/
team owned, access level choices are tiered allowing for row level security.
This is covered in more detail below.

### Role-Based Access Control (RBAC)

Dataverse uses a role-based security model to control access to the database
and tables in an environment. Security roles are used to grant a set of
permissions to users and teams.

Dataverse has a common RBAC layer built in at the datasource. Every request is
made in the context of a person and only data they are allowed to see will be
returned.

Security roles are used to group a collection of privileges. These roles can be
assigned to users, teams and business units. Privilege grands are accumulative
with the greatest access prevailing. Role assignments should be additive, we
could not give broad access and then use role assignments to diminish access.

The process of determining access is:

1. Privilege check: Does the user have the required privilege on the table
2. Access check: Does the user have access rights to perform the action

#### Access Rights

There are four ways in which a user can gain access rights to perform an action
on a given record:

##### Access Rights: Ownership

A user can have access to a record because they own the record or belong to a
team that owns the record. In this instance, any access level is sufficient to
access the record regardless of the business unit the record belongs to

##### Role Access

A user may have access to perform an action based on a record due to their
security roles:

- If the record belongs to the user or a team they belong to the user can access
with the user-level access privilege
- If the record belongs to the same business unit as the user or a team they
belong to they will need to have a role with at least Business Unit-level access
privilege
- If the record belongs to a business unit that is a descendant of the business
unit the user or a team they belong to is a member of, they will need to have a
role with at least Parent: Child business unites access privilege
- If the record belongs to a business unit that is not a descendant of the user/
teams business unit then the organisation level access privilege is required.

Configuring permissions on this basis allows for row-level security.

The access-levels, from least to most access then are:

- User
- Business unit
- Parent: Child business units
- Organisation

For roles assigned to teams with Basic-level access privilege, the role's
inheritance configuration is also relevant. If the team has Member's privilege
inheritance set to Team privileges only, they will only be able to use the
privilege for records owned by the team. If we set this property to User
Privileges then the user will be granted the privilege directly.

##### Shared Access

Individual records may be shared with a user or team. This is a useful way to
managed exceptions to the security model. However, this should be an exception
as it is a less performant way of controlling access.

When shared, the user can access both the record and related records. But in
both instances, they must still first pass the privilege check.

A more advanced sharing concept is access teams. This provides for the
auth-creation of a team and sharing record access with the team is based on an
Access Team Template which is a table of permissions to be applied. This is a
more performant method of sharing because the team does not own records and
security roles are not assigned to the team. Access is granted because the *
record is shared with the team and the user is a member.

##### Hierarchy Access

This can be used when hierarchy security management is enabled in the
organisation and for the table if the user is a manager. The manager would have
access to the record if:

- They have a security role that has the access level Business Unit or
Parent:Child business units, AND
  - The record is owned by a direct report, OR
  - A direct report is a member of the owner team, OR
  - The record was shared to perform the required action with a direct report
  - The record was shared to perform the required action with a team a direct
  report belongs to

#### Security Roles

Security roles define access to tables in the environment. For each table we can
specify permissions for:

- Create: Create a record
- Read: Read a record
- Write: Update a record
- Delete: Delete a record
- Append: Append record to the table
- Append To: Append record to another table
- Assign: Used to give ownership of a record to another user
- Share: Ability to share a record with a user or team

When a record is shared we can define the access level to that shared record,
i.e. which of the above permissions they have.

We can assign roles on a user, business unit or team basis.

##### Miscellaneous Privileges

We can also define miscellaneous privileges in a security role, e.g.:

- Act on Behalf of Another user
- Activate Business Process Flows/Business Rules
- Bulk edit and Bulk delete
- Run Flows

##### Pre-defined Security Roles

- System Admin: Full permission to customise and administer the environment
including creating, modifying and assigning security roles. Can view all data
- System Customiser: Full permission to customise the environment. Can view all
custom data in the environment but has only user level access for Account,
Contact and Activity tables
- Basic user: For OOTB entities only, can run an app in the environment and
perform tasks on records they own.

#### Form Security

In the settings for a form, we can define the security roles required to access
the form. This can be useful to show different UIs for a form, but it will not
secure the underlying data. To control access to fields use column-level
security.

##### Column Level Security

We can use column security to control access to a given column. We can enable
column level security for a given column in a table. Once enabled it will only
be accessible to those with permissions through a column security profile.

We can create a column security profile in an environment and add the column and
the level of access.

We can add users, teams and business units to the profile.

The permissions we can configure for a column security profile are:

- Read (allowed/not allowed)
- Read unmask (all records, one record, not allowed)
- update (allowed/not allowed)
- create (allowed/not allowed)

Some columns cannot be secured this way:

- Virtual table columns
- Lookup columns
- Formula columns
- Primary name columns
- System columns

##### Column-Level Security: Calculated columns

When a calculated column included a secured column, access to data may be given
unintentionally, so both should be secured.

Similarly, composite columns inherit data from multiple columns, to secure this
we would need to secure all columns that belong to it.

###### Column-Level Security: Masking Rules

Masking rules can be used to hide sensitive information, we can created masking
rules in a solution. We use a regexp to define the masking rule and a character
to be used for the masking.

When a column is enabled for column security we can add a masking rule to be
applied to the column.

These can only be used for text and number fields.

Note, if the permission for masking is set to one record then the user can read
masked values but must retrieve the unmasked values one record at a time.

#### Security and Solutions

We can add security roles and column security profiles to solutions. This
enables the transport of these entities across environments.

## Data Loss Prevention Policies

Data Loss Prevention is an important aspect of data security and compliance. We
can create policies to reduce the risk of users unintentionally exposing
organisational data.

Policies can be defined at either a tenant or environment level, we need the
Power Platform Administrator/System Administrator roles respectively to work
with these policies. The scopes are:

- Add all environments
- Excluded certain environments (all but those excluded)
- Add multiple environments

The policy defines, for prebuilt and custom connectors, whether each connector
is:

- Business
- Non-Business
- Blocked

By default, all connectors are non-business. We can change the default group
when defining a policy.

Blocked policies cannot be used and business connectors cannot be used in the
same flow/app with non-business connectors.

Some connectors are not blockable, this includes, Teams, Excel, SharePoint,
OneDrive, PowerBi and Outlook. It includes most MS enterprise standard
connectors and the dataverse connectors.

If an app or flow violates an active DLP then the maker will be unable to save
the app/flow. If a policy is activated which causes a previously saved app or
flow to be in violation of the policy, the app/flow will not run.

I can't find confirmation, but I understand that we can have multiple policies
and if any of these are violated the app/flow will not work.
