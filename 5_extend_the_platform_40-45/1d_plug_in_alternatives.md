# Plug-In Alternatives

## Declarative Alternatives

Plug-ins use imperative logic (how to do). They should be used only
when declarative processes (what to do) will not meet the requirements.

Declarative processes include:

- Business rules
- Flows
- Workflows

As a general rule, consider code the last resort. If a declarative solution is
suitable, use it.

## Length of Operations

Synchronous plug-ins and workflows, should be only be used for short running
actions. This is because:

- They are blocking
- They generally take place in a transaction and lock resources

We may use async plug-ins or workflows for medium length actions as they take
place outside of a transaction. However, note that actions will not be rolled
back in the event of failure.

## Specific Alternatives

### Plug-ins Vs Business Rules

Business rules allow us to write simple server-side logic. However, if the logic
becomes more complex, e.g. where a switch statement is preferable to a lengthy
if, else if, else sequence, we may determine that a plug-in is more appropriate.

### Plug-ins Vs Workflows

Workflows are a declarative workflow that can run either synchronously or
asynchronously.

The main limitations of workflows in relation to plugins are:

- Run only as a User, where as plugins can run as any licenced user
- They cannot access external data
- There are less triggers

On the other hand, unlike plugins, workflows can be triggered on demand.

[MS Learn](https://learn.microsoft.com/en-us/power-automate/bpf-add-on-demand-workflow)

### Power Automate Flows

Power Automate flows may be considered an alternative to asynchronous plug-ins.
Note, that they cannot run synchronously.

Power Automate flows have the benefit of not requiring development experience.
They can be less performant than plug-ins but the difference in performance is
becoming increasingly small.

### Custom Actions

### Calculated and Roll-Up Fields

### Azure Service Bus Integration

Plug-ins and workflows are not designed for batch processing. They should not
be used for long-running or high volume actions. We should offload such
processes to a separate service like an Azure worker role
