using Microsoft.Extensions.Logging;

namespace SwapiConnector;

// https://learn.microsoft.com/en-us/connectors/custom-connectors/write-code
public interface IScriptContext
{
    string CorrelationId { get; }
    string OperationId { get; }
    HttpRequestMessage Request { get; }
    ILogger Logger { get; }

    Task<HttpResponseMessage> SendAsync(
        HttpRequestMessage request,
        CancellationToken cancellationToken
    );
}
