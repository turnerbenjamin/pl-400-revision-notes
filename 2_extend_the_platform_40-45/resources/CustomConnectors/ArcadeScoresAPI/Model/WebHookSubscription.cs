namespace ArcadeScoresAPI.Model;

public class WebHookSubscription
{
    public Guid Id { get; set; }

    public string Url { get; set; } = "";
    public DateTime CreatedAt { get; set; }
    public bool IsActive { get; set; }
}
