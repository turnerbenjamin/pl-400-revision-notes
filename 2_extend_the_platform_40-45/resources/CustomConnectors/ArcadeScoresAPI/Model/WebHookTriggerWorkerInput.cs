namespace ArcadeScoresAPI.Model;

public class WebHookTriggerWorkerInput(string url, string message)
{
    public string Url { get; set; } = url;
    public string Message { get; set; } = message;
}
