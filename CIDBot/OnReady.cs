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
        readonly JsonSerializerOptions GithubJsonOptions = JsonOptions.Github;
        readonly SemanticVersion BotVersion = serviceProvider.GetRequiredService<SemanticVersion>();
        readonly LoggingService Logging = serviceProvider.GetRequiredService<LoggingService>();

        readonly HttpClient GithubClient = new() 
        {
            BaseAddress = new Uri("https://api.github.com/")
        };

        public bool IsOlderVersion { get; private set; }

        public async Task ClientReadyAsync()
        {
            try
            {
                GithubClient.DefaultRequestHeaders.Authorization = new("Bearer", GithubToken);
                GithubClient.DefaultRequestHeaders.UserAgent.Add(new("RobloxUSArmyCID", BotVersion.ToString()));

                string repoOwner = "RobloxUSArmyCID";
                string repoName = "CIDBot";

                var getLatestReleaseMsg = await GithubClient.GetAsync($"repos/{repoOwner}/{repoName}/releases/latest");
                getLatestReleaseMsg.EnsureSuccessStatusCode();
                var getLatestReleaseStr = await getLatestReleaseMsg.Content.ReadAsStringAsync();

                var latestRelease = JsonSerializer.Deserialize<GetLatestGithubReleaseResponse>(getLatestReleaseStr, GithubJsonOptions);

                // I am NOT stripping the "v" because some library doesn't understand NAMING
                // IF A DEVELOPER AFTER ME REMOVES THE "v" YOU WILL NOT COMPILE
                // AND YOU HAVE A SMALL ONE
                latestRelease!.TagName = latestRelease!.TagName!.Substring(1);
                IsOlderVersion = SemanticVersion.Parse(latestRelease!.TagName) > BotVersion;

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
            catch (Exception ex)
            {
                await Logging.LogAsync(new LogMessage(LogSeverity.Critical, ex.Source, ex.Message, ex));
                Console.WriteLine("Press ENTER to exit.");
                Console.ReadLine();
                Environment.Exit(-1);
            }
        }
    }
}
