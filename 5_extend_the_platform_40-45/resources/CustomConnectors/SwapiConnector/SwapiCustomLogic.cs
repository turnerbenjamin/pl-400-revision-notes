using System.Net;
using Newtonsoft.Json;

// DELETE NAMESPACE ONCE UPLOADED
namespace SwapiConnector;

public class Script : ScriptBase
{
    private const string GET_ALL_PEOPLE_ID = "SwapiGetAllPeople";
    private const string GET_PERSON_BY_NAME_ID = "SwapiGetPerson";
    private const string GET_PLANET_BY_NAME_ID = "SwapiGetPlanet";

    private enum ResponseType
    {
        Single,
        Multiple,
    }

    private class Person
    {
        [JsonProperty("name")]
        public string? Name { get; set; }

        [JsonProperty("homeworld")]
        public string? HomeWorld { get; set; }
    }

    private class Planet
    {
        [JsonProperty("name")]
        public string? Name { get; set; }

        [JsonProperty("population")]
        public string? Population { get; set; }

        [JsonProperty("climate")]
        public string? Climate { get; set; }

        [JsonProperty("terrain")]
        public string? Terrain { get; set; }
    }

    private class SwapiResponse<T>
    {
        [JsonProperty("next")]
        public string? Next { get; set; }

        [JsonProperty("results")]
        public List<T>? Results { get; set; }
    }

    // For each custom connector, we can only define one code block. This code
    // block relates to multiple connectors so this initial method simply routes
    // requests
    public override async Task<HttpResponseMessage> ExecuteAsync()
    {
        var uri = GetUri(Context.Request.RequestUri);
        return Context.OperationId switch
        {
            GET_ALL_PEOPLE_ID => await GetAllPeople(uri),
            GET_PERSON_BY_NAME_ID => await GetPersonByName(uri),
            GET_PLANET_BY_NAME_ID => await GetPlanetByName(uri),
            _ => new HttpResponseMessage(HttpStatusCode.RequestedRangeNotSatisfiable),
        };
    }

    // This method maps the name parameter used in the custom connector to the
    // search parameter used by SWAPI
    private static Uri? GetUri(Uri? originalUri)
    {
        if (originalUri is null)
        {
            return null;
        }
        var uri = new UriBuilder(originalUri);
        var query = System.Web.HttpUtility.ParseQueryString(uri.Query);

        query["search"] = query["name"];
        query.Remove("name");
        uri.Query = query.ToString();
        return uri.Uri;
    }

    // Controller handling get all people from the SWAPI database
    private async Task<HttpResponseMessage> GetAllPeople(Uri? uri)
    {
        var people = await GetAll<Person>(uri);
        if (people is null || people.Count == 0)
        {
            return new HttpResponseMessage(HttpStatusCode.NoContent);
        }
        return GetSuccessResponse(people);
    }

    // Controller handling getting a single person by name from the SWAPI
    // database
    private async Task<HttpResponseMessage> GetPersonByName(Uri? uri)
    {
        var person = await GetOne<Person>(uri, ResponseType.Multiple);
        if (person is null)
        {
            return new HttpResponseMessage(HttpStatusCode.NotFound);
        }
        person.HomeWorld = await GetPlanetName(person);
        return GetSuccessResponse(person);
    }

    // Controller handling getting a single planet by name from the SWAPI
    // database
    private async Task<HttpResponseMessage> GetPlanetByName(Uri? uri)
    {
        var planet = await GetOne<Planet>(uri, ResponseType.Multiple);
        if (planet is null)
        {
            return new HttpResponseMessage(HttpStatusCode.NotFound);
        }
        return GetSuccessResponse(planet);
    }

    // The SWAPI list responses are paginated. This method iterates through all
    // pages and returns the complete results in a single response
    private async Task<List<T>> GetAll<T>(Uri? initialRequest)
        where T : class
    {
        List<T> parsedResults = [];
        Uri? next = initialRequest;

        while (next is not null)
        {
            var req = new HttpRequestMessage { Method = Context.Request.Method, RequestUri = next };
            var res = await MakeRequest<SwapiResponse<T>>(req);
            if (res is null || res.Results is null)
            {
                return parsedResults;
            }

            parsedResults.AddRange(res.Results);
            next = GetNextUri(res);
        }
        return parsedResults;
    }

    // Helper method, creates a uri from the next property on the SWAPI
    // response. Next contains a url or null.
    private Uri? GetNextUri<T>(SwapiResponse<T> res)
    {
        if (string.IsNullOrEmpty(res.Next) || Context?.Request?.RequestUri is null)
        {
            return null;
        }
        var baseUri = new Uri(Context.Request.RequestUri.GetLeftPart(UriPartial.Authority));
        return new Uri(baseUri, res.Next);
    }

    // Return a single response from a give request. Can be used with either
    // list or single response types.
    private async Task<T?> GetOne<T>(Uri? request, ResponseType responseType)
        where T : class
    {
        var req = new HttpRequestMessage { Method = Context.Request.Method, RequestUri = request };

        if (responseType == ResponseType.Single)
        {
            return await MakeRequest<T>(req);
        }

        var res = await MakeRequest<SwapiResponse<T>>(req);
        if (res?.Results is null || res.Results.Count == 0)
        {
            return null;
        }
        return res.Results[0];
    }

    // Helper method, makes a request and parses it to the relevant type
    private async Task<T?> MakeRequest<T>(HttpRequestMessage req)
        where T : class
    {
        var res = await Context.SendAsync(req, CancellationToken);
        if (!res.IsSuccessStatusCode)
        {
            return null;
        }
        var json = await res.Content.ReadAsStringAsync();
        return JsonConvert.DeserializeObject<T>(json);
    }

    // Used to populate the home world field of a person with the name of the
    // planet
    private async Task<string> GetPlanetName(Person person)
    {
        if (!Uri.TryCreate(person?.HomeWorld, UriKind.Absolute, out var homeWorldUri))
        {
            return "invalid home world URL";
        }

        var planet = await GetOne<Planet>(homeWorldUri, ResponseType.Single);
        if (string.IsNullOrEmpty(planet?.Name))
        {
            return homeWorldUri.ToString();
        }
        return planet.Name;
    }

    // Returns a success response with a json payload.
    private static HttpResponseMessage GetSuccessResponse(object payload)
    {
        var serializedPayload = JsonConvert.SerializeObject(payload);
        return new HttpResponseMessage(HttpStatusCode.OK)
        {
            Content = new StringContent(
                serializedPayload,
                System.Text.Encoding.UTF8,
                "application/json"
            ),
        };
    }
}
