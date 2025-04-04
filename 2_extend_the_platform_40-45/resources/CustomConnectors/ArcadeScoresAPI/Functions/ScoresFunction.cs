using System.Net;
using ArcadeScoresAPI.Model;
using ArcadeScoresAPI.Service;
using ArcadeScoresAPI.Utils;
using Microsoft.AspNetCore.Http;
using Microsoft.Azure.Functions.Worker;
using Microsoft.Azure.Functions.Worker.Extensions.Sql;
using Microsoft.Azure.Functions.Worker.Http;
using Microsoft.Data.SqlClient;
using Microsoft.Extensions.Logging;
using Newtonsoft.Json;

namespace ArcadeScoresAPI.Functions;

public class ScoresFunction
{
    [Function(nameof(GetGameScores))]
    public static async Task<HttpResponseData> GetGameScores(
        [HttpTrigger(AuthorizationLevel.Function, "get", Route = "scores/{game}")]
            HttpRequestData req,
        [SqlInput(
            commandText: "SELECT Id, game, gamerTag, teamName, score, "
                + "CONVERT(VARCHAR(8), duration, 108) AS duration "
                + "FROM dbo.GameScores "
                + "WHERE game = @game",
            commandType: System.Data.CommandType.Text,
            parameters: "@game={game}",
            connectionStringSetting: "SqlConnectionString"
        )]
            IEnumerable<GameScore> gameScores,
        FunctionContext context
    )
    {
        var logger = context.GetLogger(nameof(gameScores));
        try
        {
            var response = req.CreateResponse(HttpStatusCode.OK);
            await response.WriteAsJsonAsync(gameScores);
            return response;
        }
        catch (Exception ex)
        {
            logger.LogError(ex.ToString());
            var response = req.CreateResponse(HttpStatusCode.InternalServerError);
            await response.WriteAsJsonAsync(new { error = "An unexpected error occurred." });
            return response;
        }
    }

    [Function(nameof(PostGameScore))]
    public static async Task<HttpResponseData> PostGameScore(
        [HttpTrigger(AuthorizationLevel.Function, "post", Route = "scores")] HttpRequestData req,
        FunctionContext context
    )
    {
        var logger = context.GetLogger(nameof(PostGameScore));
        try
        {
            var gameScore = await ParseGameScore(req);
            gameScore.Id = Guid.NewGuid();

            await GameScoresService.Post(gameScore);

            return await Utilities.BuildResponse(
                req,
                HttpStatusCode.Created,
                "Game successfully Created"
            );
        }
        catch (ArgumentException)
        {
            return await Utilities.BuildResponse(
                req,
                HttpStatusCode.BadRequest,
                "Request body is invalid. Game Type object required with a label property"
            );
        }
        catch (SqlException ex) when (ex.Number == 8152) // Handle tag too long
        {
            return await Utilities.BuildResponse(
                req,
                HttpStatusCode.Conflict,
                "The gamer tag can not be longer than 3 characters"
            );
        }
        catch (Exception ex)
        {
            logger.LogError(ex.ToString());
            return await Utilities.BuildResponse(
                req,
                HttpStatusCode.InternalServerError,
                "The server has experienced an unexpected error"
            );
        }
    }

    private static async Task<GameScore> ParseGameScore(HttpRequestData req)
    {
        string requestBody = await new StreamReader(req.Body).ReadToEndAsync();
        GameScore? gameType = JsonConvert.DeserializeObject<GameScore>(requestBody);

        if (gameType is null || !gameType.IsValid())
        {
            throw new ArgumentException("Invalid game type argument");
        }

        return gameType;
    }
}
