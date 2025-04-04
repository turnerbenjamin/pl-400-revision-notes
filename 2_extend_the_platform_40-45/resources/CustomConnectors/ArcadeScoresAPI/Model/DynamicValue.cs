namespace ArcadeScoresAPI.Model;

public class DynamicValue(string label, string value)
{
    public string Value { get; } = value;
    public string Label { get; } = label;
}
