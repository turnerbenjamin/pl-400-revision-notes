using Newtonsoft.Json;

namespace ArcadeScoresAPI.Model;

public class GameType
{
    public Guid Id { get; set; }

    [JsonProperty("label")]
    public string Label { get; set; } = "";

    [JsonProperty("doStoreDuration")]
    public bool DoStoreDuration { get; set; }

    public bool IsValid()
    {
        if (string.IsNullOrWhiteSpace(Label))
        {
            return false;
        }
        return true;
    }
}
