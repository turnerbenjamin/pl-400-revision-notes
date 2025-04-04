using Microsoft.AspNetCore.Http;
using Microsoft.AspNetCore.Mvc;
using Microsoft.Azure.Functions.Worker;
using Microsoft.Extensions.Logging;

namespace ArcadeScoresAPI.Functions;

public class TestConnectionsFunction
{
    [Function(nameof(TestConnection))]
    public static async Task<IActionResult> TestConnection(
        [HttpTrigger(AuthorizationLevel.Function, "get", Route = "test-connection")]
            HttpRequest req,
        FunctionContext context
    )
    {
        var logger = context.GetLogger(nameof(TestConnection));
        logger.LogInformation("Successful connection test");
        return new OkObjectResult("success");
    }
}
