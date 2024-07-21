using Discord;
using Discord.WebSocket;
using Microsoft.Extensions.DependencyInjection;
using System.Text.Json;

namespace CIDBot;

internal sealed class Program
{
    static readonly ServiceProvider _serviceProvider = CreateServices();

    static ServiceProvider CreateServices()
    {
        var clientConfig = new DiscordSocketConfig()
        {
            AlwaysDownloadUsers = true,
            GatewayIntents = GatewayIntents.AllUnprivileged
        };

        var jsonSerializerOptions = new JsonSerializerOptions()
        {
            WriteIndented = true,
            PropertyNameCaseInsensitive = true,
        };

        var collection = new ServiceCollection()
            .AddSingleton(clientConfig)
            .AddSingleton<DiscordSocketClient>()
            .AddSingleton<LoggingService>()
            .AddSingleton(jsonSerializerOptions);

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

        var onSlashCommand = new OnSlashCommand(_serviceProvider);
        client.SlashCommandExecuted += onSlashCommand.HandleSlashCommand;

        await Task.Delay(-1);
    }
}