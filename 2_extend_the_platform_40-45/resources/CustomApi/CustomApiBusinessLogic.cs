using System;
using CustomApi.DTO;
using Microsoft.Xrm.Sdk;
using Microsoft.Xrm.Sdk.Query;

namespace CustomApi
{
    // This plugin is registered against the custom api, get country data, it
    // implements the business logic for this message
    public class CustomApiBusinessLogic : PluginBase
    {
        public override void ExecuteBusinessLogic(IServiceProvider serviceProvider)
        {
            var targetRef = GetContextParameterOrDefault<EntityReference>(
                context.InputParameters,
                "Target"
            );
            var defaultCountry = GetContextParameterOrDefault<string>(
                context.InputParameters,
                "cr950_default_country"
            );

            var countryName = GetTargetCountryOrDefault(targetRef, defaultCountry);
            var country = GetCountryDetails(countryName);

            var capital = country?.Capital?[0] ?? "unknown";
            var population = country?.Population ?? -1;

            context.OutputParameters["cr950_capital"] = capital;
            context.OutputParameters["cr950_population"] = population;
        }

        // Call the rest countries api to extract details about the country
        private Country GetCountryDetails(string countryName)
        {
            using (var client = new System.Net.WebClient())
            {
                try
                {
                    var url = $"https://restcountries.com/v3.1/name/{countryName}";
                    var res = client.DownloadString(url);
                    var countries = Newtonsoft.Json.JsonConvert.DeserializeObject<Country[]>(res);
                    return countries?.Length > 0 ? countries[0] : null;
                }
                catch (Exception ex)
                {
                    tracingService.Trace($"Error fetching country details: {ex.Message}");
                    return null;
                }
            }
        }

        // This function tries to access the Target entity's country data and
        // returns the name of the country. If the country cannot be accessed
        // then the default is returned.
        //
        // Target in the context is an EntityReference. I assume that this is
        // due to the limits of sending parameters by query string. Accordingly,
        // it is necessary to make a call to the org service here.
        private string GetTargetCountryOrDefault(EntityReference tRef, string d)
        {
            var t = organisationService.Retrieve(
                tRef.LogicalName,
                tRef.Id,
                new ColumnSet("address1_country")
            );
            if (t.TryGetAttributeValue("address1_country", out string c))
            {
                return c;
            }
            return d;
        }
    }
}
