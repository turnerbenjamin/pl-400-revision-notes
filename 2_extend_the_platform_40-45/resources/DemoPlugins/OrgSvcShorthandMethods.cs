using System;
using System.Collections.Generic;
using DemoPlugins.Model;
using Microsoft.Xrm.Sdk;
using Microsoft.Xrm.Sdk.Messages;
using Microsoft.Xrm.Sdk.Query;

namespace DemoPlugins
{
    public class OrgSvsShorthandMethods : PluginBase
    {
        /*
        Demonstrates usage of the organisation service shorthand methods.

        Retrieval of the services has been offloaded to the base class to keep
        the logic relevant to the demonstration.

        Early bound classes used to access logical names of entities and fields.
        */
        public override void ExecuteBusinessLogic(IServiceProvider serviceProvider)
        {
            var accountId = CreateSampleAccount();
            RetrieveSampleAccount(accountId);

            UpdateSampleAccount(accountId);

            var contactId = CreateSampleContact();
            AssociateSampleAccountWithContact(accountId, contactId);

            RetrieveSampleAccount(accountId);
            DeleteSampleEntity(Contact.EntityLogicalName, contactId);
            DeleteSampleEntity(Account.EntityLogicalName, accountId);

            RetrieveAllSampleAccounts();
        }

        // Create as sample contact to be associated with the sample account
        private Guid CreateSampleContact()
        {
            var contactToCreate = new Entity(Contact.EntityLogicalName);
            contactToCreate[Contact.Fields.FirstName] = "Sample";
            contactToCreate[Contact.Fields.LastName] = "Contact";

            var contactId = organisationService.Create(contactToCreate);
            tracingService.Trace($"Contact created with id {contactId}");
            return contactId;
        }

        // Create a sample account for the to demonstrate create. Note that we
        // need to specify the entity logical name in the constructor else we
        // will get a runtime error.
        private Guid CreateSampleAccount()
        {
            var accountToCreate = new Entity(Account.EntityLogicalName);
            accountToCreate[Account.Fields.Name] = "Sample Account";

            var accountId = organisationService.Create(accountToCreate);
            tracingService.Trace($"Account created with id {accountId}");
            return accountId;
        }

        // Demonstrate retrieval of an account by id. Alternative keys will be
        // looked at in a separate plugin
        private void RetrieveSampleAccount(Guid entityId)
        {
            var columnSet = new ColumnSet(Account.Fields.Name, Account.Fields.PrimaryContactId);

            var result = organisationService.Retrieve(
                Account.EntityLogicalName,
                entityId,
                columnSet
            );

            var a = result.ToEntity<Account>();
            tracingService.Trace(
                $"Retrieved account - NAME: {a.Name}, "
                    + $"Primary Contact Id: {a.PrimaryContactId?.Id.ToString() ?? "NONE"}"
            );

            var r = new RetrieveRequest();
        }

        // Update the name of the sample account. Here an overload of entity is
        // used to identify the relevant record by its Guid
        private void UpdateSampleAccount(Guid accountId)
        {
            var updatedAccount = new Entity(Account.EntityLogicalName, accountId);
            updatedAccount[Account.Fields.Name] = "Updated Account Name";
            organisationService.Update(updatedAccount);
        }

        // Associates a sample contact with a contact. I am not demonstrating
        // disassociate but it has the same signature.
        private void AssociateSampleAccountWithContact(Guid accountId, Guid contactId)
        {
            organisationService.Associate(
                Account.EntityLogicalName,
                accountId,
                new Relationship(Account.Fields.account_primary_contact),
                new EntityReferenceCollection()
                {
                    new EntityReference(Contact.EntityLogicalName, contactId),
                }
            );
        }

        // Demonstrate delete by deleting the sample account and contact
        private void DeleteSampleEntity(string entityLogicalName, Guid entityId)
        {
            organisationService.Delete(entityLogicalName, entityId);
        }

        // Demonstrate retrieve all. I am using QueryExpression here as it the
        // simplest implementation for the functionality. QueryByAttribute can
        // not be used with Condition Operators. A separate plugin will demo
        // the other options for retrieve all.
        //
        // The code just finds all accounts with sample in the name - all sample
        // data accounts have a (sample) suffix
        private void RetrieveAllSampleAccounts()
        {
            var query = new QueryExpression(Account.EntityLogicalName)
            {
                ColumnSet = new ColumnSet(Account.Fields.Name),
                PageInfo = new PagingInfo() { Count = 50, PageNumber = 1 },
            };
            query.Criteria.AddCondition(Account.Fields.Name, ConditionOperator.Like, "%sample%");

            var entities = new List<Entity>();
            EntityCollection results;
            do
            {
                results = organisationService.RetrieveMultiple(query);
                entities.AddRange(results.Entities);
                query.PageInfo.PageNumber++;
            } while (results.MoreRecords);

            tracingService.Trace("Sample accounts:");
            foreach (Entity entity in entities)
            {
                var a = entity.ToEntity<Account>();
                tracingService.Trace(a.Name);
            }
        }
    }
}
