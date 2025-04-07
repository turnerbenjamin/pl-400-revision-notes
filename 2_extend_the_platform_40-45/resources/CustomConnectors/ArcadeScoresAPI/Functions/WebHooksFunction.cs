using ArcadeScoresAPI.Model;
using ArcadeScoresAPI.Service;
using Microsoft.AspNetCore.Http;
using Microsoft.AspNetCore.Mvc;
using Microsoft.Azure.Functions.Worker;
using Microsoft.Azure.Functions.Worker.Http;
using Microsoft.Extensions.Logging;
using Newtonsoft.Json;

namespace ArcadeScoresAPI.Functions;

public class WebHoodFunction
{
    /// <summary>
    /// Subscribes to a webhook by accepting a subscription request and storing
    /// it in the database.
    /// </summary>
    /// <param name="req">
    /// The HTTP request data, triggered by an HTTP POST request. The request
    /// body must contain a valid JSON representation of a WebHookSubscription
    /// object.
    /// </param>
    /// <param name="context">
    /// The function execution context, used for logging and other runtime
    /// operations.
    /// </param>
    /// <returns>
    /// An IActionResult indicating the result of the operation:
    /// - 200 (OK) if the subscription is successfully created.
    /// - 400 (Bad Request) if the request body is empty or invalid.
    /// - 500 (Internal Server Error) if an unexpected error occurs.
    /// </returns>
    [Function(nameof(SubscribeToWebhook))]
    public static async Task<IActionResult> SubscribeToWebhook(
        [HttpTrigger(AuthorizationLevel.Function, "post", Route = "webhooks/subscribe")]
            HttpRequestData req,
        FunctionContext context
    )
    {
        var logger = context.GetLogger(nameof(SubscribeToWebhook));
        try
        {
            logger.LogInformation($"BODY: {req.Body}");
            string requestBody = await new StreamReader(req.Body).ReadToEndAsync();
            if (string.IsNullOrWhiteSpace(requestBody))
            {
                logger.LogError("Request body is empty.");
                return new BadRequestObjectResult("Request body cannot be empty.");
            }

            var subscription = JsonConvert.DeserializeObject<WebHookSubscription>(requestBody);

            if (subscription == null || string.IsNullOrEmpty(subscription.Url))
            {
                return new BadRequestObjectResult(
                    "Invalid subscription request. A valid URL is required."
                );
            }
            logger.LogInformation($"ADDING: {subscription.Url}");
            subscription.Id = Guid.NewGuid();
            subscription.CreatedAt = DateTime.Now;
            subscription.IsActive = true;

            await WebHookSubscriptionService.Post(subscription);

            return new OkObjectResult($"Webhook subscription created for URL: {subscription.Url}");
        }
        catch (Exception ex)
        {
            logger.LogError(ex.ToString());
            return new ObjectResult("An error occurred while processing the subscription.")
            {
                StatusCode = StatusCodes.Status500InternalServerError,
            };
        }
    }
}
