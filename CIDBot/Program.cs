using Discord;
using Discord.WebSocket;
using Microsoft.Extensions.DependencyInjection;

namespace CIDBot;

class Program
{
    static readonly ServiceProvider _serviceProvider = CreateServices();

    static ServiceProvider CreateServices()
    {
        var clientConfig = new DiscordSocketConfig()
        {
            AlwaysDownloadUsers = true,
            GatewayIntents = GatewayIntents.AllUnprivileged | GatewayIntents.MessageContent
        };

        var collection = new ServiceCollection()
            .AddSingleton(clientConfig)
            .AddSingleton<DiscordSocketClient>();

        return collection.BuildServiceProvider();
    }

    public static async Task Main()
    {
        var client = _serviceProvider.GetRequiredService<DiscordSocketClient>();

        string token = Environment.GetEnvironmentVariable("CIDBot_TOKEN")
            ?? throw new NotImplementedException("Please define the CIDBot_TOKEN environment variable");

        await client.LoginAsync(TokenType.Bot, token);
        await client.StartAsync();

        await Task.Delay(-1);
    }
}