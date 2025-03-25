# Apply Business Logic in Model-Driven Apps Using Client Scripting

## Introduction

We can write JS web resources to use with form event handlers, command bar
commands, HTML web resources and IFrames.

## Build JS Code Targeting Client API Object Model

### What is the Client API Object Model

The client API Object model consists for 4 root objects:

- Execution Context
- Form Context
- Grid Context
- Xrm

#### Execution Context

The execution context is the event context in which the code executes. When
defining an event handler in a MDA, we can pass this as the first argument to
our handler function.

This object contains a number of methods and properties but the most useful are:

- getFormContext: Form or form item depending on where method called;
- getEventSource: Reference to the object the event occurred on

Event source can be used to write generic handlers which run logic based on the
event source passed at runtime. This allows a single handler to be used for
multiple controls.

#### Form Context

Form context is a reference to a form or form item against which the script is
executed. Form items include quick view controls and a row in an editable grid.

Form context is used in favour of the now deprecated Xrm.Page object.

```js
const formContext = executionContext.getFormContext();
```

##### FormContext.Data

This object contains methods that may be used to work on the form data. It
includes:

- attributes: collection of non-table data
- entity: Methods to work with the form record
  - save
  - add/remove handlers for the save and post save events
  - get entity name/id/reference/primaryAttributeValue
  - get is dirty/valid
  - attributes: collection of table data
- process: Get stages and steps
  - get/set active stage
  - add/remove handlers for various process events

There are also methods on the Data object to:

- getIsDirty/Valid
- save and refresh the form
- add/remove onLoad handlers

Note that data.attributes does not relate to the fields in the form.
data.entity.attributes is the collection of form fields.

[Demonstration code](./resources/scripting/form_context_data_demo.js)

##### FormContext.UI

Properties and methods to work with the form ui.

- Controls: Controls on the form
- FormSelector: Query forms available to the user
- Navigation: Contains all nav items on the page
- Process: Work with the business process flow
- QuickForms: Collection of quick forms on the form
- Tabs: Work with form tabs

There are methods including:

- setFormNotification
- add/remove on load handler
- get viewPort Width/Height
- getFormType (create, update, readonly)
- close

[Demonstration code](./resources/scripting/form_context_ui_demo.js)

##### Form Context Methods

Form Context contains two important shortcut methods:

- getAttribute: Alias for formContext.data.entity.attributes.get
- getControl: Alias for formContext.ui.controls.get

```js
formContext.getAttribute(attributeLogicalName);
formContext.getControl(attributeLogicalName);
```

Note, if we do not pass an argument to these functions we will get an array of
all available attributes/controls.

##### Composite Fields

Generally, we use the name of a control or attribute to access it. The exception
is with composite fields such as the address composite field. We need to use
a special syntax:

```js
function getCompositeItemLogicalName(compositeLogicalName, itemLogicalName) {
  return `${compositeLogicalName}_compositionLinkControl_${itemLogicalName}`;
}
```

#### Grid Context

We can access grid context from the form context. We need to provide the name of
the grid/sub-grid that the context relates to:

```js
const gridContext = formContext.getControl(gridName);
```

There are various methods to:

- Add/remove event handlers
- refresh the grid/ribbon
- Work with rows and cells, including adding notifications to a cell
- Get and set the view
- Work with GridAttributes and Controls

Note that we use getGridAttribute and getGridControl

There is no example in the resources, however the WebApi example is intended for
use on a main grid control. PrimaryControl in this context will be GridContext.

#### Context in Async Functions

NOTE: When using context objects in async functions, be aware that the object
may become out of date. The same warning applies to set timeout/interval
callbacks.

#### XRM

The Xrm object is globally available in the code without need for the execution
context. It contains methods a variety of APIS:

- App: This contains methods to work with sidePanes and to set/clear app-level
notifications. App/global notifications are useful, e.g. in a command bar where
outside the context of a form.
- Device: Allows, with consent, access to device content such as image, video,
audio, location. We can also use getBarcodeValue and pickFile
- Encoding: Methods to encode strings to and from html/xml
- Navigation: Methods to navigate to forms, files and urls. Also includes
methods to open dialogues
- Panel: Contains just one method to load a side pane, replaced by App.sidePanes
- Utility: Various methods such as show/close progressIndicator and get allowed
status transitions. We can also use getGlobalContext to access metadata such as
client, device, user settings and roles.
- WebApi: This contains methods to use WebApi to create and manage records and
execute Web API actions and functions

See the resources section for usage examples focussed on Navigation and WebApi.

##### Deprecation Notes

- Xrm.Page: Deprecated in favour of formContext

#### Client-Side Script Dependencies

We can manage dependencies from the relevant solution.

We cannot delete a given item while other items depend on that item.

When we export solutions, we should try to ensure that all dependencies are
included.

We can create dependencies between Javascript library web resources. Where a
dependency exists, both scripts will be downloaded. However, we cannot control
the order in which the resources load.

## Determine Event Handler Registration Approach

There are two ways to add event handlers:

First we can access the properties of a form or form element and add a handler
from there. Note that we should ensure that we enable the pass execution context
as the first argument option. Note that this requires that we register an js
web resource and add it to the relevant form.

Second, we can add handlers programmatically, for example:

```js
formContext.data.entity.addOnPostSave(contactFormPostSaveHandler);
```

or:

```js
function addOnChangeHandler(formContext, attributeName, handler) {
  const attribute = formContext.getAttribute(attributeName);
  attribute.addOnChange(handler);
}
```

We can define up to 50 handler per event.

### Passing Data Between Handlers

We can use the setSharedVariable and getSharedVariable methods to pass a common
variable between event handlers.

We can also use getDepth to get the sequence that a handler is being executed
relative to other event handlers.

### Event Registration

There are two ways to register handlers for events:

First, we can use form properties to add handlers, for instance to the onLoad or
onSave event for a form. In this case we must explicitly enable the passing of
execution context.

Second, we can also register handlers for events with code in a handler.
Execution context will be passed to the handler automatically. Note that some
handlers may only be registered with code.

A common pattern is to register an OnLoad handler and register remaining
handlers in code with the OnLoad handler logic. This can simplify registration
on a given form and permit for dynamic handler registration.

Note, that there is a limit of 50 handlers per event.

#### Column Dependencies

We can configure column dependencies from the properties of a form/form element.
This will ensure that the column will always be available for the script logic
regardless of whether it is removed from the form.

#### Web Resource Dependencies

We can also configure dependencies on other web resources. This simplifies the
logic as the form need not explicitly load dependent resources. However, note
that resources are loaded in parallel.

## Create client Scripting Targeting Dataverse Web API

As noted above, we can access the WebApi using Xrm.WebApi methods. Check the
resources for an example using this API in the context of ribbon commands.

## Configure Commands and Buttons using Power FX and JS

### Command Bars

There are four command bars that we may edit:

- Main Grid: This is displayed for a full-page list of records
- Main Form: Shown when viewing an individual row in a main form
- Sub-Grid View: Displayed above a sub-grid in a form
- Associated View: Associated view is shown when a table is selected from the
related tab on a main form.

### Component Types

There are four component types:

- Dropdown: This is a list of commands
- Split-Button: As above, but with a primary command executed when the top-level
button is clicked
- Group: Groups commands within a dropdown/split
- Command: A button that executes a command

### Coding the Button

#### Parameters

We can specify the parameters passed to click handlers in the properties of a
command. We need to specify the type and, depending on the type, a value.

The most significant type is PrimaryControl. This is the equivalent of
formContext/gridContext.

#### Power FX/Javascript

We can use either PowerFx or JS web resources as click handlers. With JS the
code will be similar to client-scripting generally. We can pass PrimaryControl
to access the form/grid context as is relevant.

We may also write PowerFx functions, for instance as seen above in relation to
visibility rules.

We can use PowerFx for both Visible and OnSelect. JS may only be used directly
with OnSelect.

When we use Power Fx for the first time a component library will be created.
This is used to store the logic.

##### Power Fx and Context Data

We can access context data using PowerFx:

- Item: Record from the datasource, e.g. Self.Selected.Item
- AllItems: Table of Records, e.g. Self.Selected.AllItems
- State: Enum showing state, e.g. Edit, New or View
- Unsaved: Boolean indicating if Selected or SelectedItems have unsaved changes
- RecordInfo: Obtains info about a record, e.g. whether the user has an edit
permission

We can also access other Dataverse tables. However, currently Dataverse is the
only datasource that may be used with Command logic

### Classic Commands

Some commands cannot be customised with the visual command bar designer. These
are classic commands which are generally read-only.

The visibility of these commands is generally based on classic rules.

### Visibility with Classic Rules

We can define visibility and enable rules using a Power Fx formula. By default,
custom buttons will be hidden when an item is selected.

The classic designer used rules rather than PowerFx. To show/hide a field based
on data we would use a enable rule of type custom rule which would call a JS
function. If we wanted to show/hide based on privilege, we would need to use a
display rule of type EntityPrivilegeRule. This is because, display rules were
executed server-side and could check privileges while enable rules were executed
client-side and could call JS junctions in the browser.

We can now just use the PowerFx visible function. For instance:

- !IsBlank(Self.Selected.Item.Email) True if email is not blank
- RecordInfo(Self.Selected.Item, RecordInfo.EditPermission)

## Debugging

### Solution Checker

Solution checker will perform a static analysis of solutions based on best
practice. This can pick up problematic patterns in client-side scripts. It can
also identify issues with performance, stability, use of deprecated methods and
reliability. We should use this is part of an automated build process as part of
the Application Lifecycle Management process.

### Browser Tools

We can use browser debuggers such as the Google Chrome debugger. We can use the
search capability where applicable to find the file we are debugging.

### Alert and Logs

We can use window.alert or, better, console.log to record messages in the window
or console.

### Fiddler Auto-Responder

Editing resources under development can be time intensive due to the need to
keep rebublishing the files.

We can use tools like auto-responder in Telerik Fiddler to replace content of a
web resource with content from a local file rather that uploading and
republishing each time.

## Best Practices

### Define Unique Script Names

We need to avoid name conflicts with other libraries used in the same context.
We can do this by adding a prefix to function names or by using an IIFE to
define a namespace.

### Avoid Unsupported Methods

Avoid methods that are deprecated or not documented in the official
documentation. As noted, the solution checker can help identify issues here.

### Avoid using JQuery

Client-scripts and command bar commands do not support direct DOM manipulation.
Accordingly, we should avoid using JQuery and other dom manipulation libraries.

### Write Non-Blocking Code

Use asynchronous patterns to avoid blocking the UI for intensive activities.
Avoid blocking methods like dialogs and progress indicators where possible and
prefer notification methods.

### Write Code for Multiple Browsers

Test code on all browsers and form factors that the users use with Model-Driven
Apps.
