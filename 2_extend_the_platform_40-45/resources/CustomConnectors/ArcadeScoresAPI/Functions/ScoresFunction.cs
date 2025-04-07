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
    /// <summary>
    /// Retrieves a list of game scores for a specific game from the database.
    /// </summary>
    /// <param name="req">
    /// The HTTP request data, triggered by an HTTP GET request. The route
    /// parameter `{game}` specifies the game for which scores are retrieved.
    /// </param>
    /// <param name="gameScores">
    /// A collection of game scores retrieved from the database using the SQL
    /// input binding.
    /// </param>
    /// <param name="context">
    /// The function execution context, used for logging and other runtime
    /// operations.
    /// </param>
    /// <returns>
    /// An HTTP response containing the list of game scores in JSON format with
    /// a status code of 200 (OK). If an error occurs, returns a 500 (Internal
    /// Server Error) response with an error message.
    /// </returns>
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

    /// <summary>
    /// Creates a new game score and stores it in the database.
    /// </summary>
    /// <param name="req">
    /// The HTTP request data, triggered by an HTTP POST request. The request
    /// body must contain a valid JSON representation of a GameScore object.
    /// </param>
    /// <param name="context">
    /// The function execution context, used for logging and other runtime
    /// operations.
    /// </param>
    /// <returns>
    /// An HTTP response indicating the result of the operation:
    /// - 201 (Created) if the game score is successfully created.
    /// - 400 (Bad Request) if the request body is invalid or missing required
    ///       properties.
    /// - 409 (Conflict) if the gamer tag exceeds the allowed length.
    /// - 500 (Internal Server Error) if an unexpected error occurs.
    /// </returns>
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
        catch (SqlException ex) when (ex.Number == 8152)
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

    /// <summary>
    /// Parses the HTTP request body to extract a GameScore object.
    /// </summary>
    /// <param name="req">
    /// The HTTP request data containing the JSON representation of a GameScore
    /// object.
    /// </param>
    /// <returns>
    /// A GameScore object parsed from the request body if valid.
    /// </returns>
    /// <exception cref="ArgumentException">
    /// Thrown if the request body is null, empty, or does not contain a valid
    /// GameScore object.
    /// </exception>
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
