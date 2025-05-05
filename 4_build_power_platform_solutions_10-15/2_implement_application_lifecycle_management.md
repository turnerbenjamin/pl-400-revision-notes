# Application Lifecycle Management (ALM)

## Introduction

ALM contains three key areas relating to an Application's lifecycle:

- Plan and Track
- Develop
- Build and Test
- Deploy
- Operate
- Monitor and Learn
- Plan and Track
- ...

## Environments

Environments are used to store, manage and share business data, apps and
business processes. They are also containers to separate apps that may have
different roles, security roles or target audiences. Each environment may have
just one Dataverse Database.

### Environment Types

#### Sandbox Environments

These are any non-production environment of Dataverse. They are isolated from
production and may be used to develop and test application changes. They include
capabilities such as reset, delete and copy operations which would be harmful in
a production environment

#### Developer Environments

Single-user environment which cannot be used to run or share production apps.
The Developer plan gives access to premium functionality, Dataverse and Power
Automate for individual use. It is primarily intended for learning purposes.

#### Default Environment

A single default environment is created automatically for each tenant and shared
by all users in the tenant. The tenant identifies the customer which may have
one or more subscriptions and services associated with it.

#### Production Environments

Environment where apps are put into operation for intended use

### Environment Access

#### Development Environments

These should be used by app makers and developers. Users should not have access
to development environments.

Developers will require at least the Environment Maker role to create resources.

#### Testing Environments

These should be used by admins and those involved in testing. Makers, developers
and production app users should not have access and testers should have only the
privileges required for testing

#### Production

Admins and app users should have access to the production environment. Users
should only the permissions necessary to perform tasks on the apps they use. App
makers and developers should either have no access or only user-level
privileges.

#### Default

Every user in the tenant can create and edit apps in a Dataverse default
environment with a database. Environments should be created for specific
purposes with appropriate roles and privileges given to those who need them.

### Environment Strategy for ALM

To follow ALM principles, we need separate environments for development and
production. Generally, we will also have at lease one test environment. Some
organisations may have multiple test environments, e.g:

- User Acceptance Testing (UAT)
- Systems Integration Testing (SIT)
- Functional Acceptance Testing (FAT)

## Solutions

Solutions are an important component of ALM in the PowerPlatform. They enable us
to packaged components such as apps and tables. Solutions can then be
transported between environments.

A solution will not contain business data. It will contain:

- Schema
- User Interface
- Analytics
- Process & Code
- Templates
- Security
- Settings

### Managed and Unmanaged Solutions

There are two types of solutions:

- Unmanaged
- Managed

#### Unmanaged

These are used to design and build in a development environment. Unmanaged
solutions contain references to components rather than actual components.
The components themselves are stored in the default solution.

Accordingly, when an unmanaged solution is deleted, only the container is
removed. The components, will still be intact in the default solution.

When a solution is exported, a snapshot of the application behaviour of each
component is created and saved to an XML document. There are no restrictions on
what can be added, removed or modified.

A solution should only be unmanaged when actively under development.

#### Managed

A managed solution, is "sealed binary". It is not actually compiled binary, but
it is a useful way of thinking of the solution in ALM.

Components, cannot be added or removed. Managed solutions can also not be
exported.

We should use managed solutions when the solution is not being actively
customised. If we need to use a solution to satisfy a dependency for another
solution's development, the dependency should be imported as managed.

### Solution Components

Solution components are anything that we can add to a solution including apps,
flows and entities.

Some components are nested within other components, for example, a table will
contain forms, views, charts, columns, table relationships and business rules.
Each of these components requires a table to exist.

Choice columns are the only columns that can exist outside of a table.

Components allow us to track limitations on what may be customised by using
managed properties.

### Solution Layering

Layering is implemented at the component level. Managed and unmanaged solutions
within a Dataverse environment exist at two different levels:

#### Managed Layers

The system layer is at the base of the managed layers level. This will contain
the entities and components required for the platform to function.

Imported managed solutions will also exist in this layer. If multiple managed
solutions are installed, later layers will have a higher precedence. Conflicts
between solutions are handled as:

- Last one wins: General rule
- Merge Logic: Forms and sitemaps in Model-Driven Apps

If a managed solution is uninstalled, the managed solution below will take
effect. If all managed solutions are uninstalled, the default behaviour in the
system solution will apply.

#### Unmanaged Layer

Any customisations and imported unmanaged solutions will exist in the unmanaged
layer. All unmanaged solutions will share a single unmanaged layer.

### Solution Dependencies

Because of the way that managed solutions are layered, some managed solutions
may be dependent on solution components in other managed solutions. This may
be used to build modular solutions.

The system will track dependencies between solutions, if you try to install a
solution that requires a base solution that is not installed, then you will not
be able to install the solution until the base is installed. Conversely, the
base solution cannot be deleted while a solution depends on it.

### Solution Component Dependencies

The solutions framework will automatically track dependencies for solution
components. To maintain the integrity of the system:

- A component cannot be deleted while another components depends on it
- There will be a warning if the solution contains any missing components that
may cause a failure on import if the dependencies are missing in the target
environment.

There are three types of component dependency:

- Solution internal: Internal dependencies are those managed by Dataverse. They
exist when a component cannot exist without another
- Published: These are created when two solution components are related to each
other and then published. To remove this dependency, the association must be
removed and the tables republished
- Unpublished: These apply to the unpublished version of a publishable solution
component that is being updated. One published it will be a published dependency

Internal dependencies may lead to published dependencies. For example, columns
in a table are dependent on the table. If a table is deleted then all columns
will be deleted to along with any relationships. If we have a lookup field on
a table form and then delete the primary table in the relationship, the deletion
can not be completed until the lookup column has been removed from the related
form and the form has been published.

### Solutions and Environment Variables

We can add an environment variable to a solution with a name, description and
data type. We can also provide a default and a current value. The current value
will override the default in the environment.

This can be useful in ALM as we can use a variable, e.g. as the host for a
connector and override this value as necessary in different environments.

### Solution Segmentation

Segmentation allows for granular development. Here, we would avoid selecting the
Add All Assets option.

### Patch Solutions

A Patch solution contains only changes for a parent-managed solution. These may
be used when making small updates. When imported, they will be layered on top of
the parent solution. However, the use of patch solutions is not recommended.

### Importing and Exporting Solutions

Generally, we will export solutions as a managed solution. When exporting a
solution we need to define a version number. Version numbers are in the format:

MAJOR.MINOR.BUILD.REVISION

Exporting a solution creates a zip file which may then be imported into another
environment.

When exporting we can run a solution checker to check for issues before export.

### Updates and Upgrades to Managed Solutions

#### Updates

This allows us to update a managed solution, the process is optimised. We cannot
delete components with an update.

To leverage optimisations, we cannot use the overwrite customisations option

Only changes to a solution will be imported

#### Upgrades

With an upgrade, we can upgrade immediately or stage the upgrade to allow
additional actions prior to the upgrade.

With an upgrade, all existing patches will be deleted from the base layer. This
is an inherently slower process, but does allow for the removal of unused
components.

When an upgrade is staged, the upgraded solution will be layered on top of the
base layer. When the upgrade is applied, the layers are flattened and a new base
solution will be created.

### Solutions and Source Control

Tools like Azure DevOps and GitHub can be used for Source Control. This is an
important aspect of ALM as it is the ultimate source of truth for a project.
With source control, dev environments are replaceable as it is easy to rebuild
from source control.

There are two paths that may be taken:

First, the unmanaged solution can be exported and included, unpacked, in the
source control system. Te build process will then import that packed solution
as unmanaged into a temporary build environment (sandbox). Finally, the solution
will be exported as managed and stored as a build artifact in the source control
system.

Second, the solution can be exported as unmanaged and managed. Both can then be
placed in the source control system. This does not require a build environment
but does require maintaining two copies of all components in the solution.

Note, if you simply export a solution into a file and add it to source control,
it will be clear that something has changed but not what has changed. Automation
can be used to unpack the file, create separate files for each part and add
these to source control. This enables us to see a detailed log of all changes.

### Automate Solution Management

We can automate a range of aspects of ALM using tools such as Azure Pipelines
and GitHub Actions. Both have pre-build Power Platform task and action support.
We can automate:

- Creation of new dev environments and installing solutions from source control
- Taking changes from dev environments and updating source control
- Running solution checker to identify quality problems
- Provisioning and de-provisioning environments
- Running automated tests such as Power Apps Test Studio tests
- Building managed solutions from source control and deploying to downstream
environments like test and production

Generally, a DevOps team will build the automations.

## Power Platform Build Tools

These are a toolbox designed for managing apps on Power Platform. Common
operations include:

- Power Platform Checker: Static analysis of solutions to look for problems
- Export solutions
- Import Solutions
- Unpack solution
- Pack Solution
- Set Solution Version
- Create, Delete and Copy Environments
