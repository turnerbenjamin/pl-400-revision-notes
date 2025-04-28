# Implement_Power_Apps_Improvements

## Variables

There are three types of variable types in Power Apps

### Global Variables

These are defined using the Set function. These can be accessed and modified
from anywhere in the application.

These can be useful in improving performance and allowing a more declarative
syntax elsewhere in the app. For instance, writing User().FullName throughout
the app is:

- Ugly
- Will slow down the application with frequent invocations of this function

A better approach is to call this function once in OnStart and store the result
in a global variable.

### Context Variables

These are defined using the UpdateContext function. These can only be accessed
and modified from the screen where the context variable is created. Unlike Set,
UpdateContext is variadic so we can define multiple context variables at once.

Since context variables are scoped to a screen, they can be used to limit
namespace pollution.

### Collections

Collections are created using the Collect or ClearCollect functions. These are
a special variable type used to store a table of data. Like global variables,
collections are available throughout the application.

Collections are most commonly used to optimise performance by reducing calls to
the same table in a datasource. We can query data once and store it locally in
a collection. This can be useful if we will be making frequent reads and
infrequent writes to the data source.

Collections, like all variables are temporary, if we wish to persist changes
then we will need to write to a data source before closing the application. This
can be done using the Patch function.

Collect is not delegable, so by default only the first 500 records from the
datasource will be retrieved and stored.

Collections may also be used in isolation from a datasource, for example, to
provide values for a drop-down menu.

## PowerFx Formulas

### Using Patch to Create and Edit Records

Patch, can be used to both create and edit records. We provide the:

- datasource,
- record to patch or Defaults(datasource)
- An object literal with any updates

### Using Remove, RemoveIf and Clear to Delete Records

- Remove: Removes a given record from a datasource
- RemoveIf: Removes any records that meet the condition
- Clear: Removes all records in a collection

To delete all records from a datasource we cannot use Clear. We would instead
use RemoveIf(datasource, true)

## Working With Choice Columns

Choice columns can be useful in providing an enumerated list of valid values
which cannot be edited by App Users. These are particularly useful as a
mechanism for filtering table rows.

Note, if you want to clear an existing stored value in a choice column using the
patch function then set the value to Blank().

### Local vs Global Choice Options

Available choice values may be created as either a local or global list:

- Local: Only used in the column and table where created
- Global: List can be used in multiple columns and tables

When a new choice column is created, the default value for the Sync with global
choice option is Yes. This is recommended and should be kept unless you are
confident that the values apply only to that single column in that table.

### Choice Values

When data for a choice only the numeric value (single option) or comma separated
numeric options (multiple options) is stored.

### Choice Vs Lookup

A common data modelling decision is between using a choice column and a lookup
column or between multiple choices and a many-many relationship.

- Choices can be modified by maker only, whereas with tables normal security
applies
- Choices are stored as a number or comma-separated numbers, while lookups are
stored as table references
- Choices cannot be made inactive or retired, whereas with tables the state row
can be used to filter out inactive options
- Choices are treated as solution components with full ALM support whereas
lookups are treated as reference data
- Choices only have label, value and colour and only label may be used in
formulas. With Tables, additional columns may be added
- Choices have localisation built-in whereas we must handle this ourselves with
lookups
- With choices there is no build-in support of dependent choice columns. With
lookups it is easy to model and implement dependent columns

Once a column has been created the data type cannot be changed. So we should be
careful when creating a new choice/lookup column.

## Dataverse Table Relationships

To build a good user experience in a canvas app, we should abstract unnecessary
complexities of the data model and allow efficient navigation of the tables,
this can be achieved using relationships.

### One to Many

With a one to many relationship, rows in a child table, will point to records in
a parent table.

By default, the parent or primary table will be a single table. However,
Dataverse also supports polymorphic lookups which allows a lookup field to point
to a row in multiple tables. An example of this is the Customer lookup which can
point to either an Account or a Contact.

When working with polymorphic lookups, the IsType and AsType PowerFx functions
can be useful to check and cast the actual record type.

If defining a relationship from the primary/parent table use one-to-many. If
defining from a child table use many-to-one.

#### Relationship Behaviours

With one-to-many relationships, relationship behaviours can be defined to
determine what should happen if the primary table row is deleted, assigned,
shared, unshared or re-parented.

The default behaviour is reference, this will remove the link between the two
tables if the primary row is deleted.

If we set the relationship to parental, then any child records will be deleted
when the parent row is deleted. This can be used to ensure that there are no
orphaned child rows.

We can also define custom relationship behaviour.

#### One-to-Many Relationships with PowerFx

Many desks can be associated with one location. With this relationship we can
use dot notation, rather than filter/lookup to access the desks at a location
and the location of a given desk respectively:

```PowerFx
LocationRecord.Desks
DeskRecord.Location.Address
```

### Many to Many Relationships

With many-to-many relationships, a hidden (relationship or intersect table) is
created. This will be used to map relationships.

We can access related records with dot notation as with one-to-many
relationships.

```PowerFx
DeskFeatureRecord.Desks
DeskRecord.'Desk Features'
```

We need to be careful when accessing relationships for many records, in the
example, for every desk, the following formula creates a label of features:

```PowerFx
Concat(ThisItem.'Desk Features', Name, ", ")
```

In Monitor, for each row a call is made to retrieve the desk features.

#### Multi-Select

If we wanted to show desks for a range of features, then this would involve some
complex logic. PowerFx does not have a simple way to get an intersection of two
tables.

## Delegation

Delegation helps Power Apps to work efficiently with data sources by reducing
the amount of data transferred. With delegation, whenever possible, Power Apps
will delegate data processing to the datasource. This includes tasks like
filtering, searching and sorting.

Whether the data can be delegated depends on both the function used and the data
source.

### When Delegation isn't Available

As noted, for certain functions and data sources, delegation will not be
available. An example is the Search function against the SharePoint data source.
In this instance, Power Apps will need to request all records and then filter on
the frontend.

By default, only the first 500 records will be requested. This can be increased
to a maximum of 2,000 records. If the SharePoint source contained 5,000 records
then 3,000 records would not be processed or displayed.

We should consider the availability of delegation whenever selecting a
datasource and the functions that we can use.

### Dataverse Delegation Availability

- Filter: Yes
- Lookup: Yes
- Sort: Yes, other than for option sets
- SortByColumns: Yes other than for option sets
- Equality (=, <>): Yes
- Comparison (<, <=, >, >=) Yes, other than for option sets
- Logical(And, Or, Not): Yes
- StartsWith: yes for text
- IsBlank: Yes, other than for option sets
- Sum, Min,Max, Avg: Yes for numbers, but not DateTime
- First,FirstN,Last,LastN: Not delegable
- Choices: Not delegable
- Collect and Clear collect: Not delegable
- CountIf, RemoveIf, UpdateIf: Not delegable
- GroupBy, Ungroup: Not delegable

Note that numbers with arithmetic expressions are not delegable, e.g.

```PowerFx
col + 10 > 100
```

Note, that the aggregate functions, sum, min, max and average are limited to a
collection of 50,000 records.

#### Partially Supported Delegable Functions

Table shaping functions are partially delegable:

- Add Columns
- Drop Columns
- Show Columns
- Rename Columns

Formulas in the arguments may be delegated, however, the output is subject to
the non delegation record limit.

### SharePoint and SQL Delegation Availability

When using SharePoint or SQL, check the documentation to see supported delegable
functions.

### Delegation Warnings

If a function is not delegable, a warning will be shown to indicate this.

## Implement Complex Power FX Formulas

### What we Need to Know

The focus is on more complex Power FX formulas:

- Set values
- Format output
- Perform calculations
- Control visibility and display mode of controls
- Patch
- Filter Data
- Search Data

## Build reusable Component Libraries

### What we Need to Know

- Discover and search components
- Building components in a library
- Updating components in a library

## Use Power Automate Flows from a Canvas App

### What we Need to Know

- Calling flows from a canvas app
- Storing and using data returned
