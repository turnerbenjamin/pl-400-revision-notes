using Newtonsoft.Json;

namespace ArcadeScoresAPI.Model;

public class WebHookSubscription
{
    public Guid Id { get; set; }

    [JsonProperty(nameof(Url))]
    public string Url { get; set; } = "";
    public DateTime CreatedAt { get; set; }
    public bool IsActive { get; set; }
}
