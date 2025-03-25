using System;
using DemoPlugins.Model;
using Microsoft.Xrm.Sdk;
using Microsoft.Xrm.Sdk.Messages;

namespace DemoPlugins
{
    public class Upsert : PluginBase
    {
        // Simple demonstration of upsert and execute. Intended to run on
        // account creation
        public override void ExecuteBusinessLogic(IServiceProvider serviceProvider)
        {
            var newAccount = GetTargetEntity().ToEntity<Account>();
            UpsertDemo(newAccount);
        }

        // Upsert an account record with name Account Tracker and the city field
        // containing the most recently created account name
        //
        // Note, there is no shorthand method for an upsert so execute is used
        // Note, alternate key used to identify record by name
        private void UpsertDemo(Account createdAccount)
        {
            var accountTracker = new Entity(
                Account.EntityLogicalName,
                Account.Fields.Name,
                "Account Tracker"
            );
            accountTracker[Account.Fields.Address1_City] = createdAccount.Name;

            var req = new UpsertRequest() { Target = accountTracker };
            var res = organisationService.Execute(req) as UpsertResponse;
            tracingService.Trace($"Upsert: {res.Target.Id}");
        }

        private Entity GetTargetEntity()
        {
            return GetContextParameterOrDefault<Entity>(context.InputParameters, "Target")
                ?? throw new ArgumentException("Unable to access Target from plug-in context");
        }
    }
}
