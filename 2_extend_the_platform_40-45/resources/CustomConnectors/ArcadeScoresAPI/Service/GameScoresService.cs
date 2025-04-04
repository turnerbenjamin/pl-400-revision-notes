using ArcadeScoresAPI.Model;
using Microsoft.Data.SqlClient;

namespace ArcadeScoresAPI.Service;

public class GameScoresService
{
    public static async Task Post(GameScore gameScore)
    {
        var connectionString = Environment.GetEnvironmentVariable(
            "SQLAZURECONNSTR_SqlConnectionString"
        );
        using var connection = new SqlConnection(connectionString);
        connection.Open();
        var query =
            "INSERT INTO dbo.GameScores (Id, game, gamerTag, teamName, score, duration) "
            + "VALUES (@Id, @game, @gamerTag, @teamName, @score, @duration)";

        using var command = new SqlCommand(query, connection);
        command.Parameters.AddWithValue("@Id", gameScore.Id);
        command.Parameters.AddWithValue("@game", gameScore.Game);
        command.Parameters.AddWithValue("@gamerTag", gameScore.GamerTag);
        command.Parameters.AddWithValue("@teamName", gameScore.TeamName);
        command.Parameters.AddWithValue("@score", gameScore.Score);
        command.Parameters.AddWithValue(
            "@duration",
            !string.IsNullOrEmpty(gameScore.Duration) ? gameScore.Duration : DBNull.Value
        );

        await command.ExecuteNonQueryAsync();
    }
}
