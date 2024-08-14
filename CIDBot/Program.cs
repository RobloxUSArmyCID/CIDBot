using Discord;
using Discord.WebSocket;
using Microsoft.Extensions.DependencyInjection;
using NuGet.Versioning;
using System.Text.Json;

namespace CIDBot;

internal sealed class Program
{
    static readonly ServiceProvider _serviceProvider = CreateServices();

    static ServiceProvider CreateServices()
    {
        var clientConfig = new DiscordSocketConfig()
        {
            AlwaysDownloadUsers = false,
            GatewayIntents = GatewayIntents.AllUnprivileged
                // & represents a binary AND operation
                // which combined with the | binary OR operation
                // results in the intents being removed
                & GatewayIntents.GuildScheduledEvents
                & GatewayIntents.GuildInvites
        };

        var githubSerializerOptions = new JsonSerializerOptions()
        {
            WriteIndented = true,
            PropertyNameCaseInsensitive = true,
            PropertyNamingPolicy = JsonNamingPolicy.SnakeCaseLower,
        };

        var jsonSerializerOptions = new JsonSerializerOptions()
        {
            WriteIndented = true,
            PropertyNameCaseInsensitive = true,
            PropertyNamingPolicy = JsonNamingPolicy.SnakeCaseLower,
        };

        var githubToken = "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";

        SemanticVersion version = new(1, 1, 1);

        var collection = new ServiceCollection()
            .AddSingleton(clientConfig)
            .AddSingleton<DiscordSocketClient>()
            .AddSingleton<LoggingService>()
            .AddSingleton(githubToken)
            .AddSingleton(version);

        return collection.BuildServiceProvider();
    }

    public static async Task Main()
    {
        var client = _serviceProvider.GetRequiredService<DiscordSocketClient>();

        // Required for the logging to actually work, but isn't assigned a variable name
        // as it isn't referred to later in the code.
        _ = _serviceProvider.GetRequiredService<LoggingService>();

        string token = Environment.GetEnvironmentVariable("CIDBot_TOKEN")
            ?? throw new NotImplementedException("Please define the CIDBot_TOKEN environment variable.");

        await client.LoginAsync(TokenType.Bot, token);
        await client.StartAsync();

        var onReady = new OnReady(_serviceProvider);
        client.Ready += onReady.ClientReadyAsync;

        var onSlashCommand = new OnSlashCommand(onReady.IsOlderVersion);
        client.SlashCommandExecuted += onSlashCommand.HandleSlashCommand;

        await Task.Delay(-1);
    }
}