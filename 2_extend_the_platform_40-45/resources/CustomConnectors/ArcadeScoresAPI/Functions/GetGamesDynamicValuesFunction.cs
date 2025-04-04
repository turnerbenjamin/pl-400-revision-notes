using System.Net;
using ArcadeScoresAPI.Model;
using Microsoft.AspNetCore.Http;
using Microsoft.Azure.Functions.Worker;
using Microsoft.Azure.Functions.Worker.Extensions.Sql;
using Microsoft.Azure.Functions.Worker.Http;

namespace ArcadeScoresAPI.Functions
{
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
