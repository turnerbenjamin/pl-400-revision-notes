# Optimise and Troubleshoot Apps

## Optimising Canvas Apps

### Common Performance Issues

Most app performance problems are due to interactions with data sources:

- Too many refreshes: Power Apps handles refreshes automatically with form save
- Too many lookups: If used in a gallery label, lookup performed for each row
- Using wrong datasource: e.g. SQL for images rather than Azure Blob Storage
- Azure Blob Storage: This is a much better option for storing images
- Unoptimised media assets: Size and compress images
- Too many screens and components in a single app
- Outdated: Republish periodically to take advantage of MS updates
- Increasing row limit for non-delegable from 500 if users on devices with slow
connections

### App Load Times

A user's first impression of the app will be the time it takes to display the
first screen or provide other visual feedback such as progress indicators.
During startup, several steps need to be performed:

- Authentication: Authentication handled for all connectors
- GetMetadata: Data such as the version of Power Apps
- OnStart: Formulas in the App OnStart property are run
- Render screens: The initial screen is rendered

#### Evaluate App Load Times

In make.powerapps, we can see analytics by app by selecting ... -> Details
-> Analytics (preview). This includes metrics including time to first screen.

There is also the Power Apps Review Tool. This is an open source tool, packaged
as a solution, which can be imported into the Dataverse Environment. This can be
used to evaluate apps in the environment. The review tool can identify common
issues and check the App against best practices. Common rules the tool can check
related to startup times include

- Use of Concurrent
- Apps Settings Flags: Ensure delayed load and explicit column selection on
- Delegation: Ensure that ClearCollect and Filter operations are  delegable
- Asset Optimisation
- Data Loading Strategy: Move data calls from OnStart to OnVisible where
possible
- Cross-Screen Dependencies: Ensure repeating UI elements are encapsulated in
components, e.g. menus, headers, footers

### OnStart and OnVisible

Two common places to load data are the OnStart and OnVisible properties:

OnStart:    Blocking, app level property
OnVisible:  Non-Blocking, screen level property

#### Optimise App.OnStart

##### Navigate on App Startup

A common requirement on start is to select the first screen that the user will
view. The old way to achieve this was to use Navigate() in App.OnStart, however
this has now been retired and new apps cannot use this function in the OnStart
property.

The new approach is to us the App.StartScreen property. By default, this is
empty and will show the first screen in the tree view. If a formula is entered
and results in an error, the behaviour will be as though the property is empty.
If you want to handle errors, use the IfError() function to catch and handle,
e.g. by redirecting to an error screen.

Note that variables created in OnStart will not be available in StartScreen.

If you have an older app, update it to use the StartScreen Property. Once
Navigate has been removed from OnStart, update the settings to disable using
Navigate in OnStart.

##### Use OnVisible Instead

Review formulas in App.OnStart and consider moving these to OnVisible as this is
non-blocking.

##### Review App Settings

- If OnStart is not used, consider disabling the property from settings

### Data Loading Strategies

#### Direct Data Source Binding

When the items property of a gallery is set to a connector, of an expression is
used based on the table, then this is direct data source binding. Data will be
retrieved from the connector as the criteria on the filter changes.

The main advantage of this is that the Power Apps runtime will take control of
when to load and refresh data (although Refresh can be used to force a load).

Another advantage, is that if the connector supports delegation, you are not
limited by the Data row limit setting. In this instance, the gallery will load
an initial set of items and fetch more data as the user scrolls.

#### Preload Data into a Collection

Preloading into collection gives you control over the load of data in OnStart or
OnVisible. This is useful if the same data is used on multiple screens. Unlike
direct data source binding, it is possible to show a loading indicator when the
data is loaded.

Another advantage is that users can update and modify the collection locally.
The changes can then be committed with a single Patch call, or rolled back with
a call to ClearCollect.

The main issue with this approach is that Collect and ClearCollect are not
delegable so the Data row limit will apply.

For not tabular connectors, we can also preload into a variable rather than a
collection using Set.

#### Load from Persisted Storage

This is similar to the above strategy, however, data is persisted in local
storage. This is achieved using the SaveData and LoadData functions.

#### Load Data Concurrently

If you are caching multiple queries, wrap them in a Concurrent() call to load
data in parallel.

#### Load Data Once

If data is being cached in the OnVisible property, ensure to check that the
collection is empty before trying to retrieve data. OnVisible will run every
time that the screen is navigated to.

#### Other Data Loading Strategies

- Evaluate usefulness of any data shown when the App or a screen first renders
- Can filters be used to reduce the rows retrieved
- Remove LookUp invocations from gallery labels
- Cache static data or data with infrequent writes
- Show loading indicators to improve perceptions of load times
- If caching used, ensure data source does not have more rows than the limit for
non-delegable functions configured in settings
- Use delegable functions wherever possible
- Remove unused connectors to avoid redundant authentication

### Other Ways to Improve Performance

- Reduce the number of components on each screen to speed-up rendering
- Review settings: e.g. Debug published App, if enabled, can damage performance

### Preview and Experimental Features

Preview features are well tested and close to being released. Experimental
features, are less stable and may disappear at any time. These may include
performance features but it is risky to include such features in production
projects.

## Testing and Troubleshooting Canvas Apps

We can use various techniques to test and troubleshoot apps:

- Timers to get a sense of function execution time
- Labels to display the values in variables
- Using different browsers and devices
- Using monitor to inspect network activity and duration of network calls

### Monitor

Azure Monitor can be launched from Power Apps studio to troubleshoot problems
and improve the quality of apps. We can use this to identify:

- Errors in using connections
- Extensive data being sent/received
- Slow connector responses
- Duplicate data actions
- Unexpected user control interactions

Monitor, when activated, will show a stream of events including:

- User interaction with controls
- Network calls

We can also use Trace() to log custom messages to the monitor.

Captured data will include necessary context, for instance, when looking at a
line for patchRow you will be able to see the formula, request and response.

When using monitor, start the app from a cold start rather than manually
rerunning OnStart or OnVisible. This ensures that the data is not influenced by
data caching. To run from a cold start:

- Enable Debug Published App in settings
- Save and publish the app
- Launch monitor

#### Collaborate with Monitor

There are three options to collaborate with monitor:

##### Download/Upload

This is the most simple method. A user can download a monitor session's events
to a file and share this. The collaborator can then upload this file to view
the stream.

##### Invite a User

Invite can be used to invite a user in the organisation to access monitor and
view events generated when previewing the app in Power Apps Studio. Invited
users are able to see the events in monitor, but they will not be able to see
what you are doing in the app.

Since the invited user does not run the app, there is no need to share the app
or app permissions with them.

Invite a user is not supported for deployed apps, it works when the app is run
in preview from Power Apps Studio

##### Connect User

This can be used to monitor usage of a published app. To use this feature, the
Debug published app setting must be enabled and the app must then be
republished.

Note that the Debug published app option will detriment performance of the
application.

As with invite, a link will be generated which can then be shared with the user
that you want to run the app. Since the collaborator will be running the app,
the app will need to be shared with them and they will need any necessary
permissions to run the app.

#### Tips when Using Monitor

First, use meaningful control names, these are included in captured events and
well named control names can help when trying to understand the stream of
events.

Second, use filter to find similar events. This can help to identify patterns or
to pinpoint areas to investigate when using monitor reactively to address a
specific problem. If there are slow actions, we may filter by the duration
column to find events with a duration greater than a certain value.

Third, often errors are revealed when an error message is presented to a user.
Monitor can be used to obtain more detailed information about specific errors.
This is particularly useful when error handling logic is used to recover and
hide the error from users. Lines with errors will include a red circle indicator
to increase their visibility.

Forth, monitor can be used to identify and fix delegation issues. Often, App
Checker will pick up these issues, but monitor can flag these events and provide
additional detail. Often, delegation problems can be solved by rewriting the
formula to remove dynamic calculations.

### Application Insights

Application Insights is a feature of MS Azure Monitor. Once an App has been
deployed, insights can be used to gain visibility into the app's performance.
Insights can tell us:

- Count of people using the app
- Screens used most often
- Time spent by users on a screen
- Which screens are slow

From Azure Portal, we can view prebuild visualisations of the telemetry data. We
can also use Power Bi to create custom visualisations of this data.

#### Enable Application Insights

In portal.azure:

- Create a new resource
- Select Application Insights
- Create the resource instance
- Once deployed, go to the resource and copy the connection string
- In make.powerapps, edit the app and select the app from the tree view
- Add the connection string to the app's properties

Note, formally an instrumentation key was used here, but insights has now
migrated away from this.

#### View Application Insights Data

Once set up we can access data using:

- Visualisations: For instance, users report
- Workbooks: Prebuild and custom workbooks for insights into availability,
performance, usage etc.
- Logs: Query the raw event data

We can use cohorts to define a set of users, events or operations. We can also
set up alerts based on metrics, for instance, if average page load time exceeds
a given value.

There is an Azure Monitor connector which can be used to build automated
workflows using data from Applications Insights.

#### Trace Logging

We can use the Trace function to send custom events to Monitor/Application
Insights. The syntax is

```PowerFx
Trace(message, trace_severity, custom_record)
```

We can query traces using properties in the custom record, this can be a
powerful tool. Note that we should keep the data sent in trace to a minimum and
avoid sending sensitive data which may cause compliance issues, e.g. personal
data.
