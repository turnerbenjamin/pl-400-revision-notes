using System;
using System.Collections.Generic;
using DemoPlugins.Model;
using Microsoft.Xrm.Sdk;
using Microsoft.Xrm.Sdk.Query;

namespace DemoPlugins
{
    public class RetrieveMultiple : PluginBase
    {
        // Retrieve Multiple on the org service accepts an instance of plug-in
        // base. There are three commonly used implementations, each explored
        // here
        public override void ExecuteBusinessLogic(IServiceProvider serviceProvider)
        {
            DemonstrateQueryByAttribute();
            DemonstrateQueryExpression();
            DemonstrateFetchExpression();
        }

        // The most simple implementation. We can define a column set, page info
        // and an order but filtering is very limited
        private void DemonstrateQueryByAttribute()
        {
            tracingService.Trace($"Query By Attribute\n");
            var query = new QueryByAttribute(Account.EntityLogicalName)
            {
                ColumnSet = new ColumnSet(Account.Fields.Name),
                PageInfo = new PagingInfo() { Count = 50, PageNumber = 1 },
            };

            query.AddOrder(Account.Fields.Name, OrderType.Descending);
            query.AddAttributeValue(Account.Fields.Address1_City, "Redmond");

            var entities = GetResults(query);

            ForEachEntity(
                entities,
                (e) =>
                {
                    var a = e.ToEntity<Account>();
                    tracingService.Trace(a.Name);
                }
            );
        }

        // Query expression is a more powerful implementation. In this example
        // a join is made with the account table to filter contacts linked to
        // an account containing "alpine" in the name. An additional filter
        // checks that the contact fullname contains "sample"
        private void DemonstrateQueryExpression()
        {
            tracingService.Trace($"\nQuery Expression\n");
            var query = new QueryExpression(Contact.EntityLogicalName)
            {
                ColumnSet = new ColumnSet(Contact.Fields.FullName),
                PageInfo = new PagingInfo() { Count = 50, PageNumber = 1 },
            };

            query.AddOrder(Contact.Fields.FullName, OrderType.Descending);

            // Create an inner join with the account table
            var accountJoin = query.AddLink(
                Account.EntityLogicalName,
                Contact.Fields.AccountId,
                Account.Fields.AccountId,
                JoinOperator.Inner
            );

            // Customise the join with a column set and filter
            accountJoin.EntityAlias = Account.EntityLogicalName;
            accountJoin.Columns = new ColumnSet(Account.Fields.Name);
            accountJoin.LinkCriteria.AddCondition(
                Account.Fields.Name,
                ConditionOperator.Like,
                "%Alpine%"
            );

            // Standard filter on the contact table
            query.Criteria.AddCondition(
                Contact.Fields.FullName,
                ConditionOperator.Like,
                "%sample%"
            );

            var entities = GetResults(query);

            ForEachEntity(
                entities,
                (e) =>
                {
                    var c = e.ToEntity<Contact>();
                    var accountName =
                        e.GetAttributeValue<AliasedValue>(
                            $"{accountJoin.EntityAlias}.{Account.Fields.Name}"
                        )?.Value as string;
                    tracingService.Trace($"{c.FullName} - {accountName}");
                }
            );
        }

        // The most powerful/ugly implementation. This expression groups
        // accounts by the city and aggregates the count of accounts related
        // to that city
        private void DemonstrateFetchExpression()
        {
            tracingService.Trace($"\nFetch Expression\n");
            var queryString =
                @"<fetch aggregate='true'>
            <entity name='account'>
               <attribute name='address1_city' alias='Count' aggregate='count' />
               <attribute name='address1_city' alias='City' groupby='true' />
               <order alias='City' />
            </entity>
        </fetch>";

            var query = new FetchExpression(queryString);
            var entities = GetResults(query);
            ForEachEntity(
                entities,
                (e) =>
                {
                    var city = e.GetAttributeValue<AliasedValue>("City")?.Value ?? "Undefined";
                    var count = e.GetAttributeValue<AliasedValue>("Count")?.Value;

                    tracingService.Trace($"City: {city} Count: {count}");
                }
            );
        }

        private List<Entity> GetResults(QueryBase query)
        {
            var entities = new List<Entity>();
            EntityCollection result;
            do
            {
                result = organisationService.RetrieveMultiple(query);
                entities.AddRange(result.Entities);
            } while (result.MoreRecords);
            return entities;
        }

        private void ForEachEntity(List<Entity> entities, Action<Entity> action)
        {
            foreach (Entity entity in entities)
            {
                action(entity);
            }
        }
    }
}
