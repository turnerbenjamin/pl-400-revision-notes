using System.Net;
using ArcadeScoresAPI.Model;
using Microsoft.AspNetCore.Http;
using Microsoft.Azure.Functions.Worker;
using Microsoft.Azure.Functions.Worker.Extensions.Sql;
using Microsoft.Azure.Functions.Worker.Http;

namespace ArcadeScoresAPI.Functions
{
    /// <summary>
    /// Retrieves a list of dynamic values for game types from the database.
    /// </summary>
    /// <param name="req">
    /// The HTTP request data, triggered by an HTTP GET request.
    /// </param>
    /// <param name="gameTypes">
    /// A collection of game types retrieved from the database using the SQL
    /// input binding.
    /// Each game type includes the Id, label, and doStoreDuration properties.
    /// </param>
    /// <returns>
    /// An HTTP response containing a list of dynamic values in JSON format with
    /// a status code of 200 (OK). Each dynamic value includes a label and a
    /// corresponding ID.
    /// </returns>
    public class GetGamesDynamicValuesFunction
    {
        [Function(nameof(GetGamesDynamicValues))]
        public static async Task<HttpResponseData> GetGamesDynamicValues(
            [HttpTrigger(AuthorizationLevel.Function, "get", Route = "games/dynamic-values")]
                HttpRequestData req,
            [SqlInput(
                commandText: "SELECT Id, label, doStoreDuration FROM dbo.GameTypes",
                connectionStringSetting: "SqlConnectionString"
            )]
                IEnumerable<GameType> gameTypes
        )
        {
            var dynamicValues = new List<DynamicValue>();
            foreach (var gameType in gameTypes)
            {
                dynamicValues.Add(new DynamicValue(gameType.Label, gameType.Id.ToString()));
            }
            var response = req.CreateResponse(HttpStatusCode.OK);
            await response.WriteAsJsonAsync(dynamicValues);
            return response;
        }
    }
}
