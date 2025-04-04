using System.Text;

namespace SwapiConnector;

// https://learn.microsoft.com/en-us/connectors/custom-connectors/write-code
public abstract class ScriptBase
{
    public IScriptContext Context { get; }
    public CancellationToken CancellationToken { get; }

    public static StringContent CreateJsonContent(string serializedJson)
    {
        return new StringContent(serializedJson, Encoding.UTF8, "application/json");
    }

    public abstract Task<HttpResponseMessage> ExecuteAsync();
}
