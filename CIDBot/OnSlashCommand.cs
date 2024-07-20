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

        readonly static HttpClient BadgesClient = new()
        {
            BaseAddress = new Uri("https://badges.roblox.com")
        };

        readonly static HttpClient FriendsClient = new()
        {
            BaseAddress = new Uri("https://friends.roblox.com")
        };

        readonly static HttpClient ThumbnailsClient = new()
        {
            BaseAddress = new Uri("https://thumbnails.roblox.com")
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

                        // group owners can be null
                        // wtf roblox
                        // skips groups without owners to not fuck shit up
                        if (groupInfo!.Data!.First()!.Owner is null) continue;
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

                bool has200OrMoreBadges = false;
                int badges = 0;

                var first100BadgesMsg = await BadgesClient.GetAsync($"/v1/users/{userId}/badges?limit=100");
                first100BadgesMsg.EnsureSuccessStatusCode();
                string first100BadgesStr = await first100BadgesMsg.Content.ReadAsStringAsync();

                var first100Badges = JsonSerializer.Deserialize<GetOwnedBadgesByIdResponse>(first100BadgesStr, JsonOptions);

                if (first100Badges!.NextPageCursor is null)
                {
                    badges = first100Badges!.Data!.Count;
                }
                else
                {
                    var next100BadgesMsg = await BadgesClient.GetAsync(
                        $"/v1/users/{userId}/badges?limit=100&cursor={first100Badges.NextPageCursor}");
                    next100BadgesMsg.EnsureSuccessStatusCode();
                    string next100BadgesStr = await next100BadgesMsg.Content.ReadAsStringAsync();

                    var next100Badges = JsonSerializer.Deserialize<GetOwnedBadgesByIdResponse>(next100BadgesStr, JsonOptions);
                    if (next100Badges!.Data!.Count != 100)
                    {
                        badges = next100Badges.Data.Count + 100;
                    }
                    else
                    {
                        badges = 200;
                        has200OrMoreBadges = true;
                    }
                }

                List<string> pastUsernames = [];
                string? pastUsernamesNextPageCursor = null;
                while (true)
                {
                    HttpResponseMessage? getPastUsernamesMsg = null;

                    if (pastUsernamesNextPageCursor is not null)
                    {
                        getPastUsernamesMsg = await UsersClient.GetAsync($"/v1/users/{userId}/username-history?limit=100&cursor={pastUsernamesNextPageCursor}");
                    }
                    else
                    {
                        getPastUsernamesMsg = await UsersClient.GetAsync($"/v1/users/{userId}/username-history?limit=100");
                    }

                    getPastUsernamesMsg.EnsureSuccessStatusCode();
                    string getPastUsernamesStr = await getPastUsernamesMsg.Content.ReadAsStringAsync();

                    var pastUsernamesJson = JsonSerializer.Deserialize<GetPastUsernamesResponse>(getPastUsernamesStr, JsonOptions);

                    foreach (var pastUsername in pastUsernamesJson!.Data!)
                    {
                        pastUsernames.Add(pastUsername.Name!);
                    }

                    if (pastUsernamesJson.NextPageCursor is null) break;
                    else pastUsernamesNextPageCursor = pastUsernamesJson.NextPageCursor;
                }

                var userInfoByIdMsg = await UsersClient.GetAsync($"/v1/users/{userId}");
                userInfoByIdMsg.EnsureSuccessStatusCode();
                string userInfoByIdStr = await userInfoByIdMsg.Content.ReadAsStringAsync();

                var userInfoById = JsonSerializer.Deserialize<GetUserInfoByIdResponse>(userInfoByIdStr, JsonOptions);

                var createdDateTime = userInfoById!.Created;
                var todayToCreatedSpan = DateTime.Now - createdDateTime;
                
                int daysFromCreated = todayToCreatedSpan!.Value.Days;

                var friendsCountMsg = await FriendsClient.GetAsync($"/v1/users/{userId}/friends/count");
                friendsCountMsg.EnsureSuccessStatusCode();
                string friendsCountStr = await friendsCountMsg.Content.ReadAsStringAsync();

                var friendsCount = JsonSerializer.Deserialize<FriendsCount>(friendsCountStr, JsonOptions);
                int amountOfFriends = friendsCount!.Count;

                var avatarHeadshotMsg = await ThumbnailsClient.GetAsync($"/v1/users/avatar-headshot?userIds={userId}&size=150x150&format=Webp&isCircular=false");
                avatarHeadshotMsg.EnsureSuccessStatusCode();
                string avatarHeadshotStr = await avatarHeadshotMsg.Content.ReadAsStringAsync();

                var avatarHeadshot = JsonSerializer.Deserialize<GetAvatarHeadshotResponse>(avatarHeadshotStr, JsonOptions);
                string thumbnailUrl = avatarHeadshot!.Data!.First()!.ImageUrl!;

                
            }
            catch (Exception ex)
            {
                Embed failureEmbed = new EmbedBuilder()
                    .WithAuthor(cmd.User)
                    .WithColor(Color.DarkRed)
                    .WithTitle(":x: | An error occured!")
                    .WithDescription($"Unhandled exception:\n```\n{ex}```")
                    .WithFooter("If you believe this is an error, contact the Investigatory Director.")
                    .WithCurrentTimestamp()
                    .Build();
                await cmd.FollowupAsync(embed: failureEmbed);
                throw;
            }
        }
    }
}
