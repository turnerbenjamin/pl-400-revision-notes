# Apply Business Logic in Model-Driven Apps Using Client Scripting

## Introduction

The extend user experience section is focussed on writing **client-side logic**
in JS/TS; it is concerned with two main technologies:

1. PCF components: These are a self-contained bundle of JS, HTML and CSS, that
define both the component UI and logic. These may be consumed by canvas and
model-driven apps (MDAs)
2. Client-side scripting with JS web resources: These resources are simple
scripts which create custom logic in MDAs by interacting with APIs made
available in this context.

This document will look at the second of these, client-side scripting with JS
resources. The primary focus will be on the APIs made available by MDAs and the
Xrm API (available in client-side scripts and PCF components). It will also look
at the contexts in which client-side scripts may operate, including as handlers
in forms and command bars.

## The Client API Object Model

### What is the Client API Object Model

Client-side scripts should not interact directly with the DOM. Instead, they
should implement logic by interacting with the Client API model. This model
includes 4 root objects:

- Execution Context
- Form Context
- Grid Context
- Xrm

### Execution Context

Client-side web resources may be used as event handlers in MDA forms. In this
context, the execution is event driven, for instance, a handler may be
registered against an onLoad, onChange or onSave event. When a JS web resource
is registered against an event in a MDA form, we may pass the execution context
as the first argument.

The execution context defines the event context in which the handler executes.
For example, we can call getEventSource to retrieve the event target. This may
be useful when writing generic handlers that can respond to a variety of events.

Generally, we primarily use the execution context to access the form/grid
context.

When a client-side web resource is used as a click handler for a
command bar button, we pass the primaryControl rather than execution context,
this gives direct access to the form/grid context depending on the command bar
type. Here, the event will always be the same and the execution is context,
rather than event, driven.

### Form Context

Form context is a reference to a form or form item against which the script is
executed. Form items include quick view controls and a row in an editable grid.

Form context is used in favour of the now deprecated Xrm.Page object.

```js
const formContext = executionContext.getFormContext();
```

#### FormContext.Data

This object contains various methods that can be used to work with the form
data. For instance, we can access and set the value of fields on the form, check
for the form status for dirty or invalid data and save the form. It is also
possible to interact with a business process where available, e.g. getting and
setting process stages and steps.

This object is demonstrated
[here](./demos/1_scripting_form_context_data_demo.md).

#### FormContext.UI

This object contains methods to work with the form UI. For instance, to access
and manipulate controls on the form, navigate to tabs and set field level
notifications on the form and its controls.

This object is demonstrated [here](./demos/2_scripting_form_context_ui_demo.md).

#### FormContext Convenience Methods

Form context defines two convenience methods at its root which give quick access
to two useful, deeply nested methods:

- getAttribute: Alias for formContext.data.entity.attributes.get
- getControl: Alias for formContext.ui.controls.get

##### Composite Fields

To access a given control or attribute, we generally provide the logical name
of the related dataverse column. The exception to this, is where the control or
field belongs to a composite field such as the address composite field. In this
instance we need to use a special syntax:

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

Grid context allows us to work with rows and cells in a grid, including the
ability to add notifications to a cell. We can also get and set views, add and
remove event listeners and refresh the grid/ribbon.

As with form context, we can access attributes and controls within the grid
using:

- getGridAttribute
- getGridControl

#### XRM

The Xrm object is globally available in the code without need for the execution
context. It contains methods a variety of APIS:

##### Xrm.App

This contains methods to work with sidePanes and to set/clear app-level
notifications. App/global notifications are useful, e.g. in a command bar where
outside the context of a form.

##### Xrm.Device

Allows, with consent, access to device content such as image, video, audio,
location. We can also use getBarcodeValue and pickFile

##### Xrm.Encoding

This contains methods to encode strings to and from html/xml

##### Xrm.Navigation

The navigation object contains methods to navigate to forms, files and urls.
Also includes methods to open dialogues boxes.

##### Xrm.Panel

Panel contains just one method to load a side pane, replaced by App.sidePanes

##### Xrm.Utility

Utility contains various methods such as show/close progressIndicator and get
allowed status transitions. We can also use getGlobalContext to access metadata
such as client, device, user settings and roles.

##### Xrm.WebApi

WebApi contains methods to use WebApi to create and manage records and execute
Web API actions and functions

##### Xrm Demo

A demonstrations focussed on the Navigation, Panel, App and WebApi methods can
be found [here](./demos/3_scripting_xrm_demos.md).

## Managing Dependencies

Web resource dependencies can be managed from the solution. We can set-up
dependencies between scripts so that one js resource cannot be deleted while
another js resource depends on it.

If a script is dependent on another, then both scripts will be downloaded.
However, the scripts will be downloaded in parallel; there is no way to control
the order in which resources load.

When we register a handler for a form, we can also set-up column dependencies to
help ensure that any columns the script relies on are present within the client
API columns. The column will be accessible even if the column control is removed
from the form.

## Registering Event Handlers

There are two ways to add event handlers:

First we can access the properties of a form or form element and add a handler
from there. Note that we should ensure that we enable the pass execution context
as the first argument option. Note that this requires that we register an js
web resource and add it to the relevant form.

Second, we can add handlers programmatically, for example:

```js
function addOnChangeHandler(formContext, attributeName, handler) {
  const attribute = formContext.getAttribute(attributeName);
  attribute.addOnChange(handler);
}
```

When handlers are added programmatically, the first argument to the handler will
always be the execution context.

A common pattern, is for the script to expose a single onLoad handler which will
then register additional handlers programmatically. This makes the resource
easier for makers to use and reduces the risk of errors when configuring the
resource in the form.

We can define up to 50 handler per event.

## Client-Side Logic in Command Bars

We can add client-side logic to custom command bar commands. For instance, the
[Xrm demos](./demos/3_scripting_xrm_demos.md) execute in the context of a
command bar.

### Command Bar Types

There are four command bars that we may edit:

- Main Grid: This is displayed for a full-page list of records
- Main Form: Shown when viewing an individual row in a main form
- Sub-Grid View: Displayed above a sub-grid in a form
- Associated View: Associated view is shown when a table is selected from the
related tab on a main form.

### Component Types

A command bar is built from four component types:

- Dropdown: This is a list of commands
- Split-Button: As above, but with a primary command executed when the top-level
button is clicked
- Group: Groups commands within a dropdown/split
- Command: A button that executes a command

### Building Command Bar Logic

#### Parameters

We can specify the parameters passed to click handlers in the properties of a
command. We need to specify the type and, depending on the type, a value.

The most significant type is PrimaryControl. This is the equivalent of
formContext/gridContext depending on the command bar type.

#### Power FX/Javascript

We can choose to build command bar logic using either PowerFx or JS web
resources. The choice between the two will depend largely on the complexity of
the logic.

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

#### Visibility with Classic Rules

We can define visibility and enable rules using a Power Fx formula. This is very
simple, just note that by default, custom buttons will be hidden when an item is
selected.

The classic designer used rules rather than PowerFx. The classic rules are much
less intuitive, for instance:

To show/hide a field based on data we would use an **enable** rule of type
**custom** rule which would call a JS function.

If we wanted to show/hide based on privilege, we would need to use a **display
rule** of type **EntityPrivilegeRule**.

This is because, display rules were executed server-side and could check
privileges while enable rules were executed client-side and could call JS
functions in the browser.

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

### Monitor

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
