using CIDBot.Models;
using Discord;
using Discord.WebSocket;
using Microsoft.Extensions.DependencyInjection;
using System.Text.Json;

namespace CIDBot
{
    internal class OnSlashCommand(ServiceProvider serviceProvider)
    {

       // readonly DiscordSocketClient Client = serviceProvider.GetRequiredService<DiscordSocketClient>();
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
                GetUserInfoByUsernameRequest userInfoByUsernameRequest = new()
                {
                    Usernames = [
                        //Guranteed to be a string as that is what the command requires
                        (string) cmd.Data.Options.First().Value
                    ],
                    ExcludeBannedUsers = true
                };

                string userInfoByUsernameRequestStr = JsonSerializer.Serialize(userInfoByUsernameRequest, options: JsonOptions);

                var userInfoByUsernameResponseMessage = await UsersClient.PostAsync("/v1/usernames/users", new StringContent(userInfoByUsernameRequestStr));

                userInfoByUsernameResponseMessage.EnsureSuccessStatusCode();

                var userInfoByUsernameResponseStr = await userInfoByUsernameResponseMessage.Content.ReadAsStringAsync();

                var userInfo = JsonSerializer.Deserialize<GetUserInfoByUsernameResponse>
                    (userInfoByUsernameResponseStr, JsonOptions);

                // Both cannot be null but are set to nullable for compiler purposes.
                ulong userId = userInfo!.Data!.First().Id;

                var groupsResponseMessage = await GroupsClient.GetAsync($"/v2/users/{userId}/groups/roles?includeLocked=true&includeNotificationPreferences=false");
                groupsResponseMessage.EnsureSuccessStatusCode();
                var groupsResponseStr = await groupsResponseMessage.Content.ReadAsStringAsync();

                var groupsResponse = JsonSerializer.Deserialize<GetUserGroupsResponse>(groupsResponseStr, JsonOptions);
                int groupAmount = groupsResponse!.Data!.Count;

                const ulong USAR_GROUP_ID = 3108077;
                const int THIRTY_REQUIRED_MEMBERS = 30;
                const int USAR_E1_RANK = 5;

                List<Under30MembersGroup> groupsUnder30Members = [];

                bool isInUsar = false;
                bool isE1 = false;
                string usarRank = String.Empty;

                foreach (var data in groupsResponse.Data)
                {
                    if (data.Group!.MemberCount <= THIRTY_REQUIRED_MEMBERS)
                    {
                        var getGroupInfoMsg = await GroupsClient.GetAsync($"v2/groups?groupIds={data.Group!.Id}");
                        getGroupInfoMsg.EnsureSuccessStatusCode();
                        string getGroupInfoStr = await getGroupInfoMsg.Content.ReadAsStringAsync();

                        var groupInfo = JsonSerializer.Deserialize<GetGroupInfoByIdResponse>(getGroupInfoStr, JsonOptions);

                        ulong ownerId = groupInfo!.Data!.First()!.Owner!.Id;

                        var getOwnerInfoMsg = await UsersClient.GetAsync($"/v1/users/{ownerId}");
                        getOwnerInfoMsg.EnsureSuccessStatusCode();
                        string getOwnerInfoStr = await getOwnerInfoMsg.Content.ReadAsStringAsync();

                        var ownerInfo = JsonSerializer.Deserialize<GetUserInfoByIdResponse>
                            (getOwnerInfoStr, JsonOptions);

                        string ownerUsername = ownerInfo!.Name!;

                        groupsUnder30Members.Add(
                            new(
                                id: data.Group!.Id,
                                name: data.Group!.Name!,
                                memberCount: data.Group!.MemberCount,
                                hasVerifiedBadge: data.Group!.HasVerifiedBadge,
                                ownerId: ownerId,
                                ownerUsername: ownerUsername
                            )
                        );
                        //Skips checking if the under 30 members group is USAR. (0.01ms improvement :skull:)
                        continue;
                    }

                    if (data.Group!.Id == USAR_GROUP_ID)
                    {
                        isInUsar = true;
                        isE1 = data.Role!.Rank == USAR_E1_RANK;
                        usarRank = data.Role!.Name!;
                    }
                }



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
