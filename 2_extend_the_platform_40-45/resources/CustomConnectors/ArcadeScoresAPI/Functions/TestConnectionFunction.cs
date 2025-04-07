using Microsoft.AspNetCore.Http;
using Microsoft.AspNetCore.Mvc;
using Microsoft.Azure.Functions.Worker;
using Microsoft.Extensions.Logging;

namespace ArcadeScoresAPI.Functions;

public class TestConnectionsFunction
{
    /// <summary>
    /// Tests the connection to the API by returning a success message.
    /// </summary>
    /// <param name="req">
    /// The HTTP request data, triggered by an HTTP GET request.
    /// </param>
    /// <param name="context">
    /// The function execution context, used for logging and other runtime
    /// operations.
    /// </param>
    /// <returns>
    /// An HTTP response with a status code of 200 (OK) and a success message.
    /// </returns>
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
