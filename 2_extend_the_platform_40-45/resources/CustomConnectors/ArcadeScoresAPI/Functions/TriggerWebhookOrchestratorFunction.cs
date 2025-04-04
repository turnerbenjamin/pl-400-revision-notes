using Azure.Messaging.ServiceBus;
using Microsoft.Azure.Functions.Worker;
using Microsoft.DurableTask.Client;
using Microsoft.Extensions.Logging;

namespace ArcadeScoresAPI.Functions
{
    public class TriggerWebhookOrchestratorFunction
    {
        private readonly ILogger<TriggerWebhookOrchestratorFunction> _logger;
        private readonly DurableTaskClient _durableTaskClient;

        public TriggerWebhookOrchestratorFunction(
            ILogger<TriggerWebhookOrchestratorFunction> logger,
            DurableTaskClient durableTaskClient
        )
        {
            _logger = logger;
            _durableTaskClient = durableTaskClient;
        }

        [Function(nameof(TriggerWebhookOrchestrator))]
        public async Task TriggerWebhookOrchestrator(
            [ServiceBusTrigger("webhooktriggerqueue", Connection = "ServiceBusConnectionString")]
                ServiceBusReceivedMessage message,
            ServiceBusMessageActions messageActions
        )
        {
            _logger.LogInformation("TRIGGER REACHED");
            await _durableTaskClient.ScheduleNewOrchestrationInstanceAsync(
                nameof(WebHookTriggerQueueOrchestrator.RunOrchestrator),
                message.Body.ToString()
            );
            await messageActions.CompleteMessageAsync(message);
        }
    }
}
