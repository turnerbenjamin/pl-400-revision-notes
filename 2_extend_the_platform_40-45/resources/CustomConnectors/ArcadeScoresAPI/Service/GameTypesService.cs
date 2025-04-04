using ArcadeScoresAPI.Model;
using Microsoft.Data.SqlClient;

namespace ArcadeScoresAPI.Service;

public class GameTypesService
{
    public static async Task Post(GameType gameType)
    {
        var connectionString = Environment.GetEnvironmentVariable(
            "SQLAZURECONNSTR_SqlConnectionString"
        );
        using var connection = new SqlConnection(connectionString);
        connection.Open();
        var query =
            "INSERT INTO dbo.GameTypes (Id, label, doStoreDuration) "
            + "VALUES (@Id, @label, @doStoreDuration)";

        using var command = new SqlCommand(query, connection);
        command.Parameters.AddWithValue("@Id", gameType.Id);
        command.Parameters.AddWithValue("@label", gameType.Label);
        command.Parameters.AddWithValue("@doStoreDuration", gameType.DoStoreDuration);
        await command.ExecuteNonQueryAsync();
    }
}
