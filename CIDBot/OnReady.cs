using CIDBot.Models;
using Discord;
using Discord.WebSocket;
using Microsoft.Extensions.DependencyInjection;
using System.Text.Json;
using NuGet.Versioning;

namespace CIDBot
{
    internal class OnReady(ServiceProvider serviceProvider)
    {
        readonly DiscordSocketClient Client = serviceProvider.GetRequiredService<DiscordSocketClient>();
        readonly string GithubToken = serviceProvider.GetRequiredService<string>();
        readonly JsonSerializerOptions JsonOptions = serviceProvider.GetRequiredService<JsonSerializerOptions>();
        readonly SemanticVersion BotVersion = serviceProvider.GetRequiredService<SemanticVersion>();

        readonly HttpClient GithubClient = new() 
        {
            BaseAddress = new Uri("https://api.github.com/")
        };

        public bool IsOlderVersion { get; set; }

        public async Task ClientReadyAsync()
        {
            if (GithubClient.DefaultRequestHeaders.Authorization is null)
            {
                GithubClient.DefaultRequestHeaders.Authorization = new(GithubToken);
            }

            string repoOwner = "RobloxUSArmyCIDBot";
            string repoName = "CIDBot";

            var getLatestReleaseMsg = await GithubClient.GetAsync($"repos/{repoOwner}/{repoName}/releases/latest");
            getLatestReleaseMsg.EnsureSuccessStatusCode();
            string getLatestReleaseStr = await getLatestReleaseMsg.Content.ReadAsStringAsync();

            var latestRelease = JsonSerializer.Deserialize<GetLatestGithubReleaseResponse>(getLatestReleaseStr);

            IsOlderVersion = SemanticVersion.Parse(latestRelease!.Tag_Name!) > BotVersion;

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
