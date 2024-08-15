using Discord;
using Serilog;
using Discord.WebSocket;
using Serilog.Events;
using Discord.Interactions;

namespace CIDBot;

public class LoggingService
{
    public LoggingService(DiscordSocketClient client, InteractionService interactions)
    {
        client.Log += LogAsync;
        interactions.Log += LogAsync;
    }

    public async Task LogAsync(LogMessage msg)
    {

        Log.Logger = new LoggerConfiguration()
#if DEBUG // Only use the Debug log level if the program is run in a development environment.
            .MinimumLevel.Debug()
#else
            .MinimumLevel.Information()
#endif
            .Enrich.FromLogContext()
            .WriteTo.Console()
            .WriteTo.File(new Serilog.Formatting.Json.JsonFormatter(), $"{DateTime.Now:yyyy-MM-dd}-CIDBot-Log.json")
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