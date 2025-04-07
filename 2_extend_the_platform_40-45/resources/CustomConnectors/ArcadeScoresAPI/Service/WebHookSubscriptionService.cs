using ArcadeScoresAPI.Model;
using Microsoft.Data.SqlClient;

namespace ArcadeScoresAPI.Service;

public class WebHookSubscriptionService
{
    /// <summary>
    /// Inserts a new webhook subscription into the database.
    /// </summary>
    /// <param name="webhookSubscription">
    /// The WebHookSubscription object containing the details of the
    /// subscription to be inserted.
    /// </param>
    /// <returns>
    /// A task representing the asynchronous operation. Inserts the webhook
    /// subscription into the database.
    /// </returns>
    /// <exception cref="SqlException">
    /// Thrown if there is an issue executing the SQL command.
    /// </exception>
    /// <exception cref="InvalidOperationException">
    /// Thrown if the database connection cannot be established.
    /// </exception>
    public static async Task Post(WebHookSubscription webhookSubscription)
    {
        var connectionString = Environment.GetEnvironmentVariable(
            "SQLAZURECONNSTR_SqlConnectionString"
        );
        using var connection = new SqlConnection(connectionString);
        connection.Open();
        var query =
            "INSERT INTO dbo.WebhookSubscriptions (Id, Url, CreatedAt, IsActive) "
            + "VALUES (@Id, @Url, @CreatedAt, @IsActive)";

        using var command = new SqlCommand(query, connection);
        command.Parameters.AddWithValue("@Id", webhookSubscription.Id);
        command.Parameters.AddWithValue("@Url", webhookSubscription.Url);
        command.Parameters.AddWithValue("@CreatedAt", webhookSubscription.CreatedAt);
        command.Parameters.AddWithValue("@IsActive", webhookSubscription.IsActive);
        await command.ExecuteNonQueryAsync();
    }
}
