using System.Net;
using Microsoft.Azure.Functions.Worker.Http;

namespace ArcadeScoresAPI.Utils;

public class Utilities
{
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
