using ArcadeScoresAPI.Model;
using Microsoft.Azure.Functions.Worker;
using Microsoft.Azure.Functions.Worker.Extensions.Sql;
using Microsoft.DurableTask;
using Microsoft.Extensions.Logging;

namespace ArcadeScoresAPI.Functions
{
    public static class WebHookTriggerQueueOrchestrator
    {
        /// <summary>
        /// Executes the Durable Task orchestration to trigger webhooks for
        /// active subscriptions.
        /// </summary>
        /// <param name="context">
        /// The orchestration context, used to manage the orchestration flow and
        /// call activity functions.
        /// </param>
        /// <param name="webHookSubs">
        /// A collection of webhook subscriptions retrieved from the database
        /// using the SQL input binding.
        /// </param>
        /// <param name="executionContext">
        /// The function execution context, used for logging and other runtime
        /// operations.
        /// </param>
        /// <returns>
        /// A task representing the asynchronous operation. Triggers webhooks
        /// for all active subscriptions
        /// by calling the `TriggerWebhook` activity function.
        /// </returns>
        [Function(nameof(RunOrchestrator))]
        public static async Task RunOrchestrator(
            [OrchestrationTrigger] TaskOrchestrationContext context,
            [SqlInput(
                commandText: "SELECT Url, IsActive FROM dbo.WebHookSubscriptions",
                connectionStringSetting: "SqlConnectionString"
            )]
                IEnumerable<WebHookSubscription> webHookSubs,
            FunctionContext executionContext
        )
        {
            var logger = executionContext.GetLogger(nameof(RunOrchestrator));

            var tasks = new List<Task>();
            foreach (var webHookSub in webHookSubs)
            {
                if (!webHookSub.IsActive)
                {
                    continue;
                }
                var message = context.GetInput<string>();
                if (string.IsNullOrWhiteSpace(message))
                {
                    logger.LogError("Message not found");
                    continue;
                }

                var workerInput = new WebHookTriggerWorkerInput(webHookSub.Url, message);
                tasks.Add(context.CallActivityAsync(nameof(TriggerWebhook), workerInput));
            }

            await Task.WhenAll(tasks);
        }

        /// <summary>
        /// Sends a POST request to a webhook URL with the specified message
        /// payload.
        /// </summary>
        /// <param name="input">
        /// The input containing the webhook URL and the message payload to be
        /// sent.
        /// </param>
        /// <param name="executionContext">
        /// The function execution context, used for logging and other runtime
        /// operations.
        /// </param>
        /// <returns>
        /// A task representing the asynchronous operation. Logs the response
        /// status code or any errors encountered during the operation.
        /// </returns>
        [Function(nameof(TriggerWebhook))]
        public static async Task TriggerWebhook(
            [ActivityTrigger] WebHookTriggerWorkerInput input,
            FunctionContext executionContext
        )
        {
            ILogger logger = executionContext.GetLogger(nameof(TriggerWebhook));
            try
            {
                using var httpClient = new HttpClient();
                var content = new StringContent(
                    input.Message,
                    System.Text.Encoding.UTF8,
                    "application/json"
                );

                var res = await httpClient.PostAsync(input.Url, content);
                logger.LogInformation(
                    $"Response to webhook callback url ({input.Url}): {res.StatusCode}"
                );
            }
            catch (Exception ex)
            {
                logger.LogError(ex.ToString());
            }
        }
    }
}
