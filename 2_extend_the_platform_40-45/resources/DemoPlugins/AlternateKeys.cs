using System;
using System.Collections.Generic;
using DemoPlugins.Model;
using Microsoft.Xrm.Sdk;
using Microsoft.Xrm.Sdk.Messages;
using Microsoft.Xrm.Sdk.Query;

namespace DemoPlugins
{
    public class AlternateKeys : PluginBase
    {
        private readonly string emailValue = "albert@altkey.co.uk";
        private readonly string firstNameValue = "albert";
        private readonly string faxValue = "111222333";

        private readonly ColumnSet columnSet = new ColumnSet(Contact.Fields.FullName);

        // In this demo alternate keys are used to identify a contact record.
        // Two keys have been created:
        //    - Simple: email
        //    - Compound: first name and fax
        public override void ExecuteBusinessLogic(IServiceProvider serviceProvider)
        {
            var c1 = RetrieveBySimpleKey().ToEntity<Contact>();
            var c2 = RetrieveByCompoundKey().ToEntity<Contact>();
            tracingService.Trace($"Simple key: {c1.FullName}, {c1.Id}");
            tracingService.Trace($"Compound key: {c2.FullName}, {c2.Id}");
        }

        // Fetch by the email key
        private Entity RetrieveBySimpleKey()
        {
            var target = new EntityReference(
                Contact.EntityLogicalName,
                Contact.Fields.EMailAddress1,
                emailValue
            );
            return ExecuteRetrieve(target);
        }

        // Fetch by first name/fax compound key
        private Entity RetrieveByCompoundKey()
        {
            var keyValues = new KeyAttributeCollection
            {
                new KeyValuePair<string, object>(Contact.Fields.FirstName, firstNameValue),
                new KeyValuePair<string, object>(Contact.Fields.Fax, faxValue),
            };
            var target = new EntityReference(Contact.EntityLogicalName, keyValues);
            return ExecuteRetrieve(target);
        }

        // Helper method, retrieve the target entity
        private Entity ExecuteRetrieve(EntityReference target)
        {
            var req = new RetrieveRequest() { Target = target, ColumnSet = columnSet };
            var res = organisationService.Execute(req) as RetrieveResponse;
            return res.Entity;
        }
    }
}
