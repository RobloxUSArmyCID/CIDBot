using Discord;
using Discord.Interactions;
using Discord.WebSocket;
using Microsoft.Extensions.DependencyInjection;
using NuGet.Versioning;

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

        var githubToken = "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";

        var client = new DiscordSocketClient(clientConfig);

        var interactionsConfig = new InteractionServiceConfig()
        {
            UseCompiledLambda = true
        };

        var interactions = new InteractionService(client.Rest, interactionsConfig);

        SemanticVersion version = new(1, 0, 0);

        var collection = new ServiceCollection()
            .AddSingleton(client)
            .AddSingleton(interactions)
            .AddSingleton<LoggingService>()
            .AddSingleton(githubToken)
            .AddSingleton(version)
            .AddSingleton<JsonOptions>()
            .AddSingleton<OnReady>();

        return collection.BuildServiceProvider();
    }

    public static async Task Main()
    {
        var client = _serviceProvider.GetRequiredService<DiscordSocketClient>();
        var onReady = _serviceProvider.GetRequiredService<OnReady>();
        var interactions = _serviceProvider.GetRequiredService<InteractionService>();

        // Required for the logging to actually work, but isn't assigned a variable name
        // as it isn't referred to later in the code.
        _ = _serviceProvider.GetRequiredService<LoggingService>();

        string token = Environment.GetEnvironmentVariable("CIDBot_TOKEN")
            ?? throw new NotImplementedException("Please define the CIDBot_TOKEN environment variable.");

        await client.LoginAsync(TokenType.Bot, token);
        await client.StartAsync();

        client.Ready += onReady.ClientReadyAsync;

        await interactions.AddModulesAsync(typeof(SlashCommands).Assembly, _serviceProvider);

        client.InteractionCreated += async (interaction) =>
        {
            var ctx = new SocketInteractionContext(client, interaction);
            await interactions.ExecuteCommandAsync(ctx, _serviceProvider);
        };

        await Task.Delay(-1);
    }
}