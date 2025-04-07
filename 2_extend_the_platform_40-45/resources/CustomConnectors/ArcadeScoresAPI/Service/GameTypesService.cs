using ArcadeScoresAPI.Model;
using Microsoft.Data.SqlClient;

namespace ArcadeScoresAPI.Service;

public class GameTypesService
{
    /// <summary>
    /// Inserts a new game type into the database.
    /// </summary>
    /// <param name="gameType">
    /// The GameType object containing the details of the game type to be
    /// inserted.
    /// </param>
    /// <returns>
    /// A task representing the asynchronous operation. Inserts the game type
    /// into the database.
    /// </returns>
    /// <exception cref="SqlException">
    /// Thrown if there is an issue executing the SQL command.
    /// </exception>
    /// <exception cref="InvalidOperationException">
    /// Thrown if the database connection cannot be established.
    /// </exception>
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
