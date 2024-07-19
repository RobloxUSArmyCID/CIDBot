using Discord;
using Serilog;
using Discord.WebSocket;
using Serilog.Events;

namespace CIDBot;

class LoggingService
{
    private readonly DiscordSocketClient Client;

    public LoggingService(DiscordSocketClient client)
    {
        Client = client;
        Client.Log += LogAsync;
    }

    private async Task LogAsync(LogMessage msg)
    {

        Log.Logger = new LoggerConfiguration()
            .MinimumLevel.Information()
            .Enrich.FromLogContext()
            .WriteTo.Console()
            .WriteTo.File(new Serilog.Formatting.Json.JsonFormatter(), "log.json")
            .CreateLogger();

        LogEventLevel severity = msg.Severity switch
        {
            LogSeverity.Critical => LogEventLevel.Fatal,
            LogSeverity.Error => LogEventLevel.Error,
            LogSeverity.Warning => LogEventLevel.Warning,
            LogSeverity.Info => LogEventLevel.Information,
            LogSeverity.Verbose => LogEventLevel.Verbose,
            LogSeverity.Debug => LogEventLevel.Debug,
            _ => LogEventLevel.Information
        };

        Log.Write(severity, msg.Exception, "[{Source}] {Message}", msg.Source, msg.Message);

        await Task.CompletedTask;
    }
}