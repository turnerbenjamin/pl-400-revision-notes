using System;
using Microsoft.Xrm.Sdk;

namespace CustomApi
{
    // This plugin is registered to run when the address1_country field of a
    // contact is updated. It calls the custom api, get country data, and
    // records the response using the tracing service.
    // Note that the implicit input, Target, must be an entity reference
    public class CustomApiPluginTrigger : PluginBase
    {
        public override void ExecuteBusinessLogic(IServiceProvider serviceProvider)
        {
            var target = GetContextParameterOrDefault<Entity>(context.InputParameters, "Target");

            var req = new OrganizationRequest("cr950_get_country_data")
            {
                ["Target"] = target.ToEntityReference(),
                ["cr950_default_country"] = "France",
            };
            var res = organisationService.Execute(req);

            var capital = GetContextParameterOrDefault<string>(res.Results, "cr950_capital");
            var population = GetContextParameterOrDefault<int>(res.Results, "cr950_population");

            tracingService.Trace($"Capital: {capital}, Population: {population / 1000000} million");
        }
    }
}
