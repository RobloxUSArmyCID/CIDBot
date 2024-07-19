using Discord;
using Discord.WebSocket;
using Microsoft.Extensions.DependencyInjection;

namespace CIDBot
{
    internal class OnReady(ServiceProvider serviceProvider)
    {
        readonly DiscordSocketClient Client = serviceProvider.GetRequiredService<DiscordSocketClient>();

        public async Task ClientReadyAsync()
        {
            await Client.SetActivityAsync(new CustomStatusGame("Background checking..."));

            var bgcheckCommand = new SlashCommandBuilder()
                .WithName("bgcheck")
                .WithDescription("Background check a Roblox user")
                .AddOption(
                    name: "username",
                    description: "The username of the user you wish to background check",
                    isRequired: true,
                    type: ApplicationCommandOptionType.String
                )
                .Build();

            await Client.CreateGlobalApplicationCommandAsync(bgcheckCommand);
        }
    }
}
