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
