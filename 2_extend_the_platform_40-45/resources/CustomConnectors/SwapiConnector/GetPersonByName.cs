using System.Net;
using Newtonsoft.Json;

// DELETE NAMESPACE ONCE UPLOADED
namespace SwapiConnector;

// CHANGE CLASS NAME TO SCRIPT ONCE UPLOADED
public class GetPersonByName : ScriptBase
{
    //Used to deserialize people objects in response.results
    public class Person
    {
        [JsonProperty("name")]
        public string? Name { get; set; }
    }

    // Used to deserialize response
    public class PersonResponse
    {
        [JsonProperty("next")]
        public string? Next { get; set; }

        [JsonProperty("results")]
        public List<Person>? Results { get; set; }
    }

    // Returns person or 404 based on the Name parameter
    public override async Task<HttpResponseMessage> ExecuteAsync()
    {
        var uri = TranslateUri(Context.Request.RequestUri);
        var person = await GetPerson(uri);
        if (person is null)
        {
            return new HttpResponseMessage(HttpStatusCode.NoContent);
        }
        return GetSuccessResponse(person);
    }

    // SWAPI uses search to filter records. The connector uses Name to indicate
    // that a single result will be returned for a given name. This method just
    // changes the query string key and returns an updated uri
    private static Uri? TranslateUri(Uri? originalUri)
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

    // Returns the person found or null
    private async Task<Person?> GetPerson(Uri? uri)
    {
        var req = new HttpRequestMessage { Method = Context.Request.Method, RequestUri = uri };
        return await MakeRequest(req);
    }

    // Makes request and parses the results
    private async Task<Person?> MakeRequest(HttpRequestMessage req)
    {
        var res = await Context.SendAsync(req, CancellationToken);
        if (!res.IsSuccessStatusCode)
        {
            return null;
        }
        return await ParseResults(res);
    }

    // Uses Newtonsoft.Json to parse results. The MS docs provide a list of
    // supported namespaces
    private static async Task<Person?> ParseResults(HttpResponseMessage res)
    {
        var json = await res.Content.ReadAsStringAsync();
        var people = JsonConvert.DeserializeObject<PersonResponse>(json);

        if (people?.Results is null || people.Results.Count == 0)
        {
            return null;
        }
        return people.Results[0];
    }

    // Returns a success response containing the serialised person.
    private static HttpResponseMessage GetSuccessResponse(Person person)
    {
        var serializedPerson = JsonConvert.SerializeObject(person);
        return new HttpResponseMessage(HttpStatusCode.OK)
        {
            Content = new StringContent(
                serializedPerson,
                System.Text.Encoding.UTF8,
                "application/json"
            ),
        };
    }
}
