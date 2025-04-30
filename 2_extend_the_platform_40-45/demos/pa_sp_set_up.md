# Power Automate Flow Service Principle Owner

This is a quick demo showing how to make a service principle owned flow.

First, an app registration is created in Azure:

![app reg](./screens/pa/sp/1_sp_creation.png)

Next, the app is added as an app user to the Dataverse environment in
admin.powerapps:

![app user](./screens/pa/sp/2_ap_creation.png)

Finally, we need to either:

- Change the primary owner to the app user
- Add the app user as a co-owner

To change the primary user, the flow must be solution aware. After adding the
flow to a solution:

![primary user](./screens/pa/sp/3_update_primary_owner.png)
