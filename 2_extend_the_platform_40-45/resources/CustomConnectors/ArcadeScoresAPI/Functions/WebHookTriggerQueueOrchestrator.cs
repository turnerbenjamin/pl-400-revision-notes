using ArcadeScoresAPI.Model;
using Microsoft.Azure.Functions.Worker;
using Microsoft.Azure.Functions.Worker.Extensions.Sql;
using Microsoft.DurableTask;
using Microsoft.Extensions.Logging;
using Newtonsoft.Json;

namespace ArcadeScoresAPI.Functions
{
    public static class WebHookTriggerQueueOrchestrator
    {
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
            logger.LogInformation("IN ORCHESTRATOR");

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
