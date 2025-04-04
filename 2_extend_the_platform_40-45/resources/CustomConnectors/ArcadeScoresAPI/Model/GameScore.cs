using Newtonsoft.Json;

namespace ArcadeScoresAPI.Model;

public class GameScore
{
    public Guid Id { get; set; }

    [JsonProperty("game")]
    public Guid? Game { get; set; }

    [JsonProperty("gamerTag")]
    public string? GamerTag { get; set; }

    [JsonProperty("teamName")]
    public string? TeamName { get; set; }

    [JsonProperty("score")]
    public int? Score { get; set; }

    [JsonProperty("duration")]
    public string? Duration { get; set; }

    public bool IsValid()
    {
        if (Game is null || GamerTag is null || TeamName is null || Score is null)
        {
            return false;
        }
        return true;
    }
}
