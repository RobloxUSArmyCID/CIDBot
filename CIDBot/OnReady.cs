using CIDBot.Models;
using Discord;
using Discord.WebSocket;
using System.Text.Json;
using NuGet.Versioning;
using Discord.Interactions;

namespace CIDBot
{
    public class OnReady
    {
        public OnReady(DiscordSocketClient client, string githubToken, JsonOptions jsonOptions, SemanticVersion botVersion, LoggingService logging, InteractionService interactions)
        {
            Client = client;
            GithubToken = githubToken;
            GithubJsonOptions = jsonOptions.Github;
            BotVersion = botVersion;
            Logging = logging;
            Interactions = interactions;
        }

        readonly DiscordSocketClient Client;
        readonly string GithubToken;
        readonly JsonSerializerOptions GithubJsonOptions;
        readonly SemanticVersion BotVersion;
        readonly LoggingService Logging;
        readonly InteractionService Interactions;

        readonly HttpClient GithubClient = new() 
        {
            BaseAddress = new Uri("https://api.github.com/")
        };

        public bool IsOlderVersion { get; private set; }
        public TaskCompletionSource<bool> ReadyTaskCompletionSource { get; } = new();

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

                GithubRelease latestRelease = JsonSerializer.Deserialize<GithubRelease>(getLatestReleaseStr, GithubJsonOptions)!;

                // I am NOT stripping the "v" because some library doesn't understand NAMING
                // IF A DEVELOPER AFTER ME REMOVES THE "v" YOU WILL NOT COMPILE
                // AND YOU HAVE A SMALL ONE
                latestRelease.TagName = latestRelease!.TagName![1..];
                IsOlderVersion = SemanticVersion.Parse(latestRelease!.TagName) > BotVersion;

                await Client.SetActivityAsync(new CustomStatusGame("Background checking..."));

                await Interactions.RegisterCommandsGloballyAsync();

                ReadyTaskCompletionSource.SetResult(true);
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
