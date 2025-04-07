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

        /// <summary>
        /// Initializes a new instance of the <see
        /// cref="TriggerWebhookOrchestratorFunction"/> class.
        /// </summary>
        /// <param name="logger">
        /// The logger instance used to log information and errors during the
        /// function's execution.
        /// </param>
        /// <param name="durableTaskClient">
        /// The Durable Task client used to schedule and manage orchestration
        /// instances.
        /// </param>
        public TriggerWebhookOrchestratorFunction(
            ILogger<TriggerWebhookOrchestratorFunction> logger,
            DurableTaskClient durableTaskClient
        )
        {
            _logger = logger;
            _durableTaskClient = durableTaskClient;
        }

        /// <summary>
        /// Triggers a Durable Task orchestration instance based on a message
        /// received from a Service Bus queue.
        /// </summary>
        /// <param name="message">
        /// The message received from the Service Bus queue.
        /// </param>
        /// <param name="messageActions">
        /// The Service Bus message actions used to complete or abandon the
        /// message after processing.
        /// </param>
        /// <returns>
        /// A task representing the asynchronous operation. Completes the
        /// message in the queue after successfully scheduling the
        /// orchestration.
        /// </returns>
        [Function(nameof(TriggerWebhookOrchestrator))]
        public async Task TriggerWebhookOrchestrator(
            [ServiceBusTrigger("webhooktriggerqueue", Connection = "ServiceBusConnectionString")]
                ServiceBusReceivedMessage message,
            ServiceBusMessageActions messageActions
        )
        {
            try
            {
                await _durableTaskClient.ScheduleNewOrchestrationInstanceAsync(
                    nameof(WebHookTriggerQueueOrchestrator.RunOrchestrator),
                    message.Body.ToString()
                );
                await messageActions.CompleteMessageAsync(message);
            }
            catch (Exception ex)
            {
                _logger.LogError(ex.Message);
                await messageActions.AbandonMessageAsync(message);
            }
        }
    }
}
