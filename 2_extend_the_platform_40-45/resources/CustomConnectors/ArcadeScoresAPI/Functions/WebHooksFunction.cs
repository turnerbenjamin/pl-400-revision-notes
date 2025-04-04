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
