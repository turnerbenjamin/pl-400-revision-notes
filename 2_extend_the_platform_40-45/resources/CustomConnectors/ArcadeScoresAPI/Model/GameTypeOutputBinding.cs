using Microsoft.Azure.Functions.Worker.Extensions.Sql;
using Microsoft.Azure.Functions.Worker.Http;

namespace ArcadeScoresAPI.Model;

public class GameTypeOutputBinding
{
    [SqlOutput("dbo.GameTypes", connectionStringSetting: "SqlConnectionString")]
    public GameScore? GameType { get; set; }
    public HttpResponseData? HttpResponse { get; set; }
}
