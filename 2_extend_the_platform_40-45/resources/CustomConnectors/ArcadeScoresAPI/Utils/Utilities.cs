using System.Net;
using Microsoft.Azure.Functions.Worker.Http;

namespace ArcadeScoresAPI.Utils;

public class Utilities
{
    /// <summary>
    /// Builds an HTTP response with the specified status code and message.
    /// </summary>
    /// <param name="req">
    /// The HTTP request data used to create the response.
    /// </param>
    /// <param name="status">
    /// The HTTP status code to set for the response.
    /// </param>
    /// <param name="message">
    /// The message to include in the response body.
    /// </param>
    /// <returns>
    /// A task representing the asynchronous operation. Returns an
    /// HttpResponseData object with the specified status and message.
    /// </returns>
    public static async Task<HttpResponseData> BuildResponse(
        HttpRequestData req,
        HttpStatusCode status,
        string message
    )
    {
        var res = req.CreateResponse(status);
        await res.WriteStringAsync(message);
        return res;
    }
}
