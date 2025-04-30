# Power Automate Flow Recurrence

This is a simple powerautomate flow based on the following MS learn activity:

<https://learn.microsoft.com/en-us/training/modules/get-started-flows/6-flow-scheduled-flows>

First create a table in excel:

![create table](./screens/pa/recurrence/1_create_excel_data.png)

Next create a flow with a recurrence trigger. I have set the trigger to run
once every minute (the frequency cannot be less than 60 seconds).

![recurrence](./screens/pa/recurrence/2_recurrence.png)

There is then an action to read rows from the table. The connector uses dynamic
values to give a dropdown of available tables. There is also a dynamic schema
for the imported table so we can access the fields without deserializing the
body

![read rows](./screens/pa/recurrence/3_list_rows.png)

Finally, the notifications are sent:

![notifications](./screens/pa/recurrence/4_notifications.png)
