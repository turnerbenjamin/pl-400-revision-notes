using Microsoft.AspNetCore.Http;
using Microsoft.AspNetCore.Mvc;
using Microsoft.Azure.Functions.Worker;
using Microsoft.Extensions.Logging;

namespace AzureWorker.Functions
{
    public class AzureWorkerFunc
    {
        private readonly ILogger<AzureWorkerFunc> _logger;

        public AzureWorkerFunc(ILogger<AzureWorkerFunc> logger)
        {
            _logger = logger;
        }

        [Function("AzureWorkerFunc")]
        public IActionResult Run(
            [HttpTrigger(AuthorizationLevel.Function, "get", "post")] HttpRequest req
        )
        {
            string queryParams = "";
            foreach (var q in req.Query)
            {
                queryParams += $"Key: {q.Key} Value: {q.Value}\n";
            }

            string requestHeader = "";
            foreach (var h in req.Headers)
            {
                requestHeader += $"Key: {h.Key} Value: {h.Value}\n";
            }
            _logger.LogInformation("Query Parameters:\n" + queryParams);
            _logger.LogInformation("Request Header: \n" + requestHeader);
            _logger.LogInformation("Request Body:\n" + req.BodyReader.ToString());

            return new OkObjectResult("All good");
        }
    }
}
