using CIDBot.Models;
using Discord;
using Discord.WebSocket;
using Microsoft.Extensions.DependencyInjection;
using Microsoft.VisualBasic;
using System.Runtime.InteropServices;
using System.Text.Json;

namespace CIDBot
{
    internal class OnSlashCommand(ServiceProvider serviceProvider)
    {

        readonly DiscordSocketClient Client = serviceProvider.GetRequiredService<DiscordSocketClient>();
        readonly JsonSerializerOptions JsonOptions = serviceProvider.GetRequiredService<JsonSerializerOptions>();

        readonly static HttpClient GroupsClient = new()
        {
            BaseAddress = new Uri("https://groups.roblox.com")
        };

        readonly static HttpClient UsersClient = new()
        {
            BaseAddress = new Uri("https://users.roblox.com")
        };


        public async Task HandleSlashCommand(SocketSlashCommand cmd)
        {
            if (cmd.CommandName == "bgcheck") await OnBgcheckCommand(cmd);
        }

        async Task OnBgcheckCommand(SocketSlashCommand cmd)
        {
            try
            {
                await cmd.DeferAsync();
                GetUserDataByUsernameRequestModel usernameRequest = new()
                {
                    Usernames = [
                        //Guranteed to be a string as that is what the command requires
                        (string) cmd.Data.Options.First().Value
                    ],
                    ExcludeBannedUsers = true
                };

                string usernameRequestJson = JsonSerializer.Serialize(usernameRequest, options: JsonOptions);

                var usernameResponseMessage = await UsersClient.PostAsync("/v1/usernames/users", new StringContent(usernameRequestJson));

                usernameResponseMessage.EnsureSuccessStatusCode();

                var usernameResponseStr = await usernameResponseMessage.Content.ReadAsStringAsync();

                var usernameResponse = JsonSerializer.Deserialize<GetUserDataByUsernameResponseModel>
                    (usernameResponseStr, JsonOptions);

                // Both cannot be null but are set to nullable for compiler purposes.
                ulong userId = usernameResponse!.Data!.First().Id;

                

            }
            catch (Exception ex)
            {
                Embed failureEmbed = new EmbedBuilder()
                    .WithAuthor(cmd.User)
                    .WithColor(Color.DarkRed)
                    .WithTitle(":x: | An error occured!")
                    .WithDescription($"Unhandled exception:\n\n```\n{ex.Message}```")
                    .WithFooter("If you believe this is an error, contact the Investigatory Director.")
                    .WithCurrentTimestamp()
                    .Build();
                await cmd.FollowupAsync(embed: failureEmbed);
                throw;
            }
        }
    }
}
