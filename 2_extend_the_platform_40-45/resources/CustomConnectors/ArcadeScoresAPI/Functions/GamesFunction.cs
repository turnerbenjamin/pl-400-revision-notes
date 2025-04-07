using System.Net;
using ArcadeScoresAPI.Model;
using ArcadeScoresAPI.Service;
using ArcadeScoresAPI.Utils;
using Azure.Messaging.ServiceBus;
using Microsoft.AspNetCore.Http;
using Microsoft.Azure.Functions.Worker;
using Microsoft.Azure.Functions.Worker.Extensions.Sql;
using Microsoft.Azure.Functions.Worker.Http;
using Microsoft.Data.SqlClient;
using Microsoft.Extensions.Logging;
using Newtonsoft.Json;

namespace ArcadeScoresAPI.Functions
{
    public class GamesFunction
    {
        /// <summary>
        /// Retrieves a list of game types from the database.
        /// </summary>
        /// <param name="req">
        /// The HTTP request data, triggered by an HTTP GET request.
        /// </param>
        /// <param name="gameTypes">
        /// A collection of game types retrieved from the database using the SQL
        /// input binding.
        /// </param>
        /// <returns>
        /// An HTTP response containing the list of game types in JSON format
        /// with a status code of 200 (OK).
        /// </returns>
        [Function(nameof(GetGameTypes))]
        public static async Task<HttpResponseData> GetGameTypes(
            [HttpTrigger(AuthorizationLevel.Function, "get", Route = "games")] HttpRequestData req,
            [SqlInput(
                commandText: "SELECT Id, label, doStoreDuration FROM dbo.GameTypes",
                connectionStringSetting: "SqlConnectionString"
            )]
                IEnumerable<GameType> gameTypes
        )
        {
            var response = req.CreateResponse(HttpStatusCode.OK);
            await response.WriteAsJsonAsync(gameTypes);
            return response;
        }

        /// <summary>
        /// Creates a new game type and stores it in the database.
        /// </summary>
        /// <param name="req">
        /// The HTTP request data, triggered by an HTTP POST request. The
        /// request body must contain a valid JSON representation of a GameType
        /// object.
        /// </param>
        /// <param name="context">
        /// The function execution context, used for logging and other runtime
        /// operations.
        /// </param>
        /// <returns>
        /// An HTTP response indicating the result of the operation:
        /// - 201 (Created) if the game type is successfully created.
        /// - 400 (Bad Request) if the request body is invalid.
        /// - 409 (Conflict) if a game type with the same unique identifier
        ///       already exists.
        /// - 500 (Internal Server Error) if an unexpected error occurs.
        /// </returns>
        [Function(nameof(PostGameType))]
        public static async Task<HttpResponseData> PostGameType(
            [HttpTrigger(AuthorizationLevel.Function, "post", Route = "games")] HttpRequestData req,
            FunctionContext context
        )
        {
            var logger = context.GetLogger(nameof(PostGameType));
            try
            {
                var gameType = await ParseGameType(req);
                gameType.Id = Guid.NewGuid();

                await GameTypesService.Post(gameType);
                await TriggerWebhooks(gameType, logger);
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
            catch (SqlException ex) when (ex.Number == 2627) // Handle UNIQUE constraint violation
            {
                return await Utilities.BuildResponse(
                    req,
                    HttpStatusCode.Conflict,
                    "A game with the same unique identifier already exists."
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
        /// Sends a message to a Service Bus queue to trigger webhooks for the
        /// specified game type.
        /// </summary>
        /// <param name="gameType">
        /// The game type object containing details to be sent as the message
        /// body.
        /// </param>
        /// <param name="logger">
        /// The logger instance used to log information or errors during the
        /// operation.
        /// </param>
        /// <returns>
        /// A task representing the asynchronous operation. Logs an error if the
        /// message fails to send.
        /// </returns>
        /// <exception cref="Exception">
        /// Thrown if the Service Bus connection string is not found in the
        /// environment variables.
        /// </exception>
        private static async Task TriggerWebhooks(GameType gameType, ILogger logger)
        {
            string? serviceBusConnectionString =
                Environment.GetEnvironmentVariable("SERVICEBUSCONNSTR_ServiceBusConnectionString")
                ?? throw new Exception("ServiceBusConnectionString is null");
            string queueName = "webhooktriggerqueue";

            var messageBody = JsonConvert.SerializeObject(gameType);

            await using var client = new ServiceBusClient(serviceBusConnectionString);
            ServiceBusSender sender = client.CreateSender(queueName);

            var message = new ServiceBusMessage(messageBody);
            try
            {
                await sender.SendMessageAsync(message);
            }
            catch (Exception ex)
            {
                logger.LogError($"Failed to send message to Service Bus: {ex.Message}");
            }
        }

        /// <summary>
        /// Parses the HTTP request body to extract a GameType object.
        /// </summary>
        /// <param name="req">
        /// The HTTP request data containing the JSON representation of a
        /// GameType object.
        /// </param>
        /// <returns>
        /// A GameType object parsed from the request body if valid.
        /// </returns>
        /// <exception cref="ArgumentException">
        /// Thrown if the request body is null, empty, or does not contain a
        /// valid GameType object.
        /// </exception>
        private static async Task<GameType> ParseGameType(HttpRequestData req)
        {
            string requestBody = await new StreamReader(req.Body).ReadToEndAsync();
            GameType? gameType = JsonConvert.DeserializeObject<GameType>(requestBody);

            if (gameType is null || !gameType.IsValid())
            {
                throw new ArgumentException("Invalid game type argument");
            }

            return gameType;
        }
    }
}
