using CIDBot.Models;
using Discord;
using System.Text;
using System.Text.Json;
using F23.StringSimilarity;
using Discord.Interactions;
using Discord.WebSocket;

namespace CIDBot
{
    public class SlashCommands : InteractionModuleBase
    {
        public OnReady _onReady;
        public JsonOptions _jsonOptions;
        public SlashCommands(OnReady onReady, JsonOptions jsonOptions)
        {
            _onReady = onReady;
            _jsonOptions = jsonOptions;
            RobloxJsonOptions = _jsonOptions.OtherThanGithub;
        }

        public JsonSerializerOptions RobloxJsonOptions;
        const string NEW_VERSION_AVAILABLE_UPDATE = ":arrows_counterclockwise: | A newer version is available! Please update at https://github.com/RobloxUSArmyCID/CIDBot/releases/latest and run the newer version of the bot.";

        //readonly InteractionService _interactions = serviceProvider.GetRequiredService<InteractionService>();

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

        readonly static NormalizedLevenshtein _levenshtein = new();

        bool NewVersionAvailable = false;

        public override async Task BeforeExecuteAsync(ICommandInfo command)
        {
            await _onReady.ReadyTaskCompletionSource.Task;
            NewVersionAvailable = _onReady.IsOlderVersion;
        }

        public async Task SlashCommandErrored(Exception ex)
        {
            Embed failureEmbed = new EmbedBuilder()
                    .WithAuthor(Context.User)
                    .WithColor(Color.DarkRed)
                    .WithTitle(":x: | An error occured!")
                    .WithDescription($"Unhandled exception:\n```\n{ex}```")
                    .WithFooter("If you believe this is an error, contact the Investigatory Director.")
                    .WithCurrentTimestamp()
                    .Build();
            await FollowupAsync(embed: failureEmbed);
        }

        async Task NoUsernameFoundAsync(string username)
        {
            Embed embed = new EmbedBuilder()
                .WithAuthor(Context.User)
                .WithColor(Color.Red)
                .WithCurrentTimestamp()
                .WithTitle(":x: | No user found!")
                .WithDescription($"The user `{username}` doesn't exist or is banned on Roblox. Please verify the spelling.")
                .Build();
            await FollowupAsync(embed: embed);
        }

        [SlashCommand("bgcheck", "Background check a Roblox user")]
        async Task BgcheckCommand
        (
            [Summary(description: "The username of the user you wish to background check")]
            string username
        )
        {
            try
            {
                
                await DeferAsync();

                if (NewVersionAvailable)
                {
                    await FollowupAsync(NEW_VERSION_AVAILABLE_UPDATE);
                    return;
                }

                GetUserInfoByUsernameRequest userInfoByUsernameRequest = new()
                {
                    Usernames = [
                        username
                    ],
                    ExcludeBannedUsers = true
                };

                string userInfoByUsernameRequestStr = JsonSerializer.Serialize(userInfoByUsernameRequest, options: RobloxJsonOptions);

                var userInfoByUsernameResponseMessage = await UsersClient.PostAsync("/v1/usernames/users", new StringContent(userInfoByUsernameRequestStr));
                userInfoByUsernameResponseMessage.EnsureSuccessStatusCode();
                var userInfoByUsernameResponseStr = await userInfoByUsernameResponseMessage.Content.ReadAsStringAsync();

                var userInfo = JsonSerializer.Deserialize<ResponseData<RequestedUser>>
                    (userInfoByUsernameResponseStr, RobloxJsonOptions) ?? throw new InvalidOperationException("The data for getting the user's ID is null.");
                

                if (userInfo.Data.Count == 0)
                {
                    await NoUsernameFoundAsync(username);
                    return;
                }


                // Both cannot be null but are set to nullable for compiler purposes.
                ulong userId = userInfo.Data.First().Id;
                username = userInfo.Data.First().Name!;

                var concurrentTasks = new[]
                {
                    GroupsClient.GetAsync($"/v2/users/{userId}/groups/roles?includeLocked=true&includeNotificationPreferences=false"),
                    BadgesClient.GetAsync($"/v1/users/{userId}/badges?limit=100"),
                    UsersClient.GetAsync($"/v1/users/{userId}"),
                    FriendsClient.GetAsync($"/v1/users/{userId}/friends/count"),
                    ThumbnailsClient.GetAsync($"/v1/users/avatar-headshot?userIds={userId}&size=150x150&format=Webp&isCircular=false"),
                    FriendsClient.GetAsync($"/v1/users/{userId}/friends"),
                };

                var completedTasks = await Task.WhenAll(concurrentTasks);

                foreach (var task in completedTasks)
                {
                    task.EnsureSuccessStatusCode();
                }

                var responses = new List<string>();

                completedTasks
                    .ToList()
                    .ForEach(async (response) => 
                    {
                        responses.Add(await response.Content.ReadAsStringAsync());
                    });
                
                var deserializationClasses = new[] 
                {
                    typeof(ResponseData<UserGroup>),
                    typeof(ResponseData<Badge>),
                    typeof(User),
                    typeof(FriendsCount),
                    typeof(ResponseData<AvatarHeadshot>),
                    typeof(ResponseData<Friend>)
                };

                List<object> deserializedResponses = [];

                for (int i = 0; i < responses.Count; i++) 
                {
                    object? response = null;
                    Console.WriteLine(responses[i]);
                    try 
                    {
                        response = JsonSerializer.Deserialize(responses[i], deserializationClasses[i], RobloxJsonOptions);
                    }
                    catch (JsonException ex) 
                    {
                        Console.WriteLine(ex);
                        return;
                    }
                    deserializedResponses.Add(response ?? throw new InvalidOperationException($"An API returned a null object."));
                }


                // these are already null-checked, we can safely apply the null forgiving operator
                var groups = (deserializedResponses[0] as ResponseData<UserGroup>)!;
                var first100Badges = (deserializedResponses[1] as ResponseData<Badge>)!;
                var userInfoById = (deserializedResponses[2] as User)!;
                var friendsCount = (deserializedResponses[3] as FriendsCount)!;
                var avatarHeadshot = (deserializedResponses[4] as ResponseData<AvatarHeadshot>)!;
                var friends = (deserializedResponses[5] as ResponseData<Friend>)!;
        
                int groupAmount = groups.Data.Count;

                const ulong USAR_GROUP_ID = 3108077;
                const int THIRTY_REQUIRED_MEMBERS = 30;
                const int USAR_E1_RANK = 5;

                List<Under30MembersGroup> groupsUnder30Members = [];

                bool isInUsar = false;
                bool isE1 = false;
                string usarRank = string.Empty;

                foreach (var group in groups.Data)
                {
                    if (group.Group.MemberCount <= THIRTY_REQUIRED_MEMBERS)
                    {
                        var getGroupInfoMsg = await GroupsClient.GetAsync($"v2/groups?groupIds={group.Group!.Id}");
                        getGroupInfoMsg.EnsureSuccessStatusCode();
                        string getGroupInfoStr = await getGroupInfoMsg.Content.ReadAsStringAsync();

                        var groupInfo = (JsonSerializer.Deserialize(getGroupInfoStr, typeof(ResponseData<Group>), RobloxJsonOptions) as ResponseData<Group>)!;

                        var groupOwner = groupInfo.Data.First().Owner;
                        if (groupOwner is null) continue;

                        ulong ownerId = groupOwner.Id;

                        groupsUnder30Members.Add(
                            new(
                                id: group.Group.Id,
                                name: group.Group.Name,
                                memberCount: group.Group.MemberCount,
                                hasVerifiedBadge: group.Group.HasVerifiedBadge,
                                ownerId: ownerId
                            )
                        );
                        
                        continue;
                    }

                    if (group.Group!.Id == USAR_GROUP_ID)
                    {
                        isInUsar = true;
                        isE1 = group.Role.Rank == USAR_E1_RANK;
                        usarRank = group.Role.Name;
                    }
                }

                bool has200OrMoreBadges = false;
                int badges = 0;

                if (first100Badges.NextPageCursor is null)
                {
                    badges = first100Badges!.Data!.Count;
                }
                else
                {
                    var next100BadgesMsg = await BadgesClient.GetAsync(
                        $"/v1/users/{userId}/badges?limit=100&cursor={first100Badges.NextPageCursor}");
                    next100BadgesMsg.EnsureSuccessStatusCode();
                    string next100BadgesStr = await next100BadgesMsg.Content.ReadAsStringAsync();

                    var next100Badges = JsonSerializer.Deserialize<ResponseData<Badge>>(next100BadgesStr, RobloxJsonOptions);
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
                    HttpResponseMessage? getPastUsernamesMsg;

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

                    var pastUsernamesJson = JsonSerializer.Deserialize<ResponseData<PastUsername>>(getPastUsernamesStr, RobloxJsonOptions);

                    foreach (var pastUsername in pastUsernamesJson!.Data!)
                    {
                        pastUsernames.Add(pastUsername.Name!);
                    }

                    if (pastUsernamesJson.NextPageCursor is null) break;
                    else pastUsernamesNextPageCursor = pastUsernamesJson.NextPageCursor;
                }


                var createdDateTime = userInfoById.Created;
                var todayToCreatedSpan = DateTime.Now - createdDateTime;
                
                int daysFromCreated = todayToCreatedSpan.Days;

                int amountOfFriends = friendsCount.Count;

                string thumbnailUrl = avatarHeadshot.Data.First().ImageUrl;


                // CLOSER TO 1 - THE LESS SIMILAR
                // CLOSER TO 0 - THE MORE SIMILAR
                // 0.72 WAS PICKED ALONGSIDE CID HICOM DUE TO THE ALGORITHM
                List<string?> usernamesOfSuspiciousFriends = friends.Data.Where(friend =>
                {
                    return _levenshtein.Distance(friend.Name, username) <= 0.72;
                })
                .Select(friend => friend.Name)
                .ToList();
                

                StringBuilder descriptionBuilder = new();

                bool failedBackgroundCheck = false;

                if (!isInUsar)
                {
                    descriptionBuilder.AppendLine("- ⚠ Not in USAR ⚠ ");
                    failedBackgroundCheck = true;
                }

                if (isE1)
                {
                    descriptionBuilder.AppendLine("- ⚠ E1 ⚠ ");
                    failedBackgroundCheck = true;
                }

                if (!has200OrMoreBadges)
                {
                    descriptionBuilder.AppendLine($"- ⚠ Less than 200 badges ({badges}) ⚠ ");
                    failedBackgroundCheck = true;
                }

                const int NINETY_REQUIRED_DAYS_FOR_ENTRANCE = 90;
                const int ONE_YEAR_IN_DAYS = 365;
                
                if (daysFromCreated < ONE_YEAR_IN_DAYS)
                {
                    if (daysFromCreated < NINETY_REQUIRED_DAYS_FOR_ENTRANCE)
                    {
                        descriptionBuilder.AppendLine($"- ⚠ Account age under 90 days old ({daysFromCreated}) (failing) ⚠");
                        failedBackgroundCheck = true;
                    }
                    else
                    {
                        descriptionBuilder.AppendLine($"- ⚠ Account age under 365 days old (suspicious, not failing) ⚠ ");
                    }
                }

                if (amountOfFriends <= 3)
                {
                    descriptionBuilder.AppendLine($"- ⚠ 3 or less friends. ⚠");
                    failedBackgroundCheck = true;
                }
                
                if (groupAmount <= 15)
                {
                    descriptionBuilder.AppendLine($"- ⚠ In 15 or less groups ({groupAmount}) ⚠");
                    failedBackgroundCheck = true;
                }

                if (groupsUnder30Members.Count > 0)
                {
                    groupsUnder30Members.ForEach(
                        group => descriptionBuilder.AppendLine($"- ⚠ Suspicious group: {group.Name} ({group.MemberCount} member{(group.MemberCount != 1 ? "s" : "")}) ⚠"));
                }

                if (pastUsernames.Count > 0)
                {
                    descriptionBuilder.AppendLine($"• Past username{(pastUsernames.Count != 1 ? "s" : "")}: {String.Join(", ", pastUsernames)}");
                }

                if (usernamesOfSuspiciousFriends.Count > 0)
                {
                    descriptionBuilder.AppendLine($"- ⚠ Suspected Alt{(usernamesOfSuspiciousFriends.Count != 1 ? "s" : "")} in Friends List: {String.Join(", ", usernamesOfSuspiciousFriends)} ⚠");
                }

                if (descriptionBuilder.Length == 0)
                {
                    descriptionBuilder.AppendLine("+ No concerns found! (Verify punishments and criteria not checked by the bot.)");
                }


                string description = descriptionBuilder.ToString();

                if (description.Length > 4096)
                {
                    description = "- ❌ Embed description too long! Please manually background check the user! ❌";
                }

                Embed embed = new EmbedBuilder()
                    .WithAuthor(Context.User)
                    .WithTitle(":white_check_mark: | Background check finished!")
                    .WithUrl($"https://www.roblox.com/users/{userId}/profile")
                    .WithDescription($"```diff\n{description}```")
                    .WithColor(Color.Blue)
                    .WithCurrentTimestamp()
                    .WithThumbnailUrl(thumbnailUrl)
                    .WithFields(
                        [
                            new EmbedFieldBuilder()
                                .WithName("Username:")
                                .WithValue($"```{username}```")
                                .WithIsInline(true),
                            new EmbedFieldBuilder()
                                .WithName("ID:")
                                .WithValue($"```{userId}```")
                                .WithIsInline(true),
                            new EmbedFieldBuilder()
                                .WithName("Failed:")
                                .WithValue($"```diff\n{(failedBackgroundCheck ? "- Yes" : "+ No")}```")
                                .WithIsInline(true),
                            new EmbedFieldBuilder()
                                .WithName("USAR Rank:")
                                .WithValue($"```{(isInUsar ? usarRank : "N/A")}```")
                                .WithIsInline(true),
                            new EmbedFieldBuilder()
                                .WithName("Account Age:")
                                .WithValue($"```{daysFromCreated} day{(daysFromCreated != 1 ? "s" : "")} old```")
                                .WithIsInline(true),
                        ]
                    )
                    .Build();


                await FollowupAsync(embed: embed);

            }
            catch (Exception ex)
            {
                await SlashCommandErrored(ex);
                throw;
            }
        }

        //[RequireUserPermission(GuildPermission.Administrator)]
        [SlashCommand("setup", "Setup CID Bot (temporary)")]
        public async Task SetupCommand
        (
            [Summary("submissions_channel", "Channel with the message for submissions")]
            SocketTextChannel submissionsChannel,
            [Summary("section_staff_punishments_log", "Channel with the punishment logs visible to Section Staff")]
            SocketTextChannel sectionStaffPunishmentsLogChannel,
            [Summary("battalion_staff_punishments_log", "Channel with the punishment logs visible to Battalion Staff")]
            SocketTextChannel battalionStaffPunishmentsLogChannel,
            [Summary("section_staff_role", "Section Staff Role")]
            SocketRole sectionStaffRole,
            [Summary("battalion_staff_role", "Battalion Staff Role")]
            SocketRole battalionStaffRole
        )
        {
            try
            {
                await DeferAsync();
                Console.WriteLine("setup command started");

                if (NewVersionAvailable) 
                {
                    await FollowupAsync(NEW_VERSION_AVAILABLE_UPDATE);
                    return;
                }

                Embed submissionEmbed = new EmbedBuilder()
                    .WithAuthor("Criminal Investigation Division")
                    .WithColor(Color.Magenta)
                    .WithCurrentTimestamp()
                    .WithTitle("FILE SUBMISSION")
                    .WithDescription("Please select the action file you wish to submit.")
                    .Build();
                var fileMenuBuilder = new SelectMenuBuilder()
                    .WithPlaceholder("Select an action file")
                    .WithCustomId($"file-{Context.Interaction.User.Id}")
                    .WithMinValues(1)
                    .WithMaxValues(1)
                    .AddOption("Punishment Request", "punishment", "Temporairly the only option");
                var components = new ComponentBuilder()
                    .WithSelectMenu(fileMenuBuilder)
                    .Build();


                // TEMPORARY SOLUTION PRE-DATABASE
                TemporaryJson temp = new()
                {
                    BattalionStaffPunishmentChannel = battalionStaffPunishmentsLogChannel.Id,
                    SectionStaffPunishmentChannel = sectionStaffPunishmentsLogChannel.Id,
                    BattalionStaffRole = battalionStaffRole.Id,
                    SectionStaffRole = sectionStaffRole.Id,
                };
                string json = JsonSerializer.Serialize(temp);
                StreamWriter str = File.CreateText("temp.json");
                str.WriteLine(json);
                str.Close();

                await submissionsChannel.SendMessageAsync(embed: submissionEmbed, components: components);
                await FollowupAsync("Finished!");

            }
            catch (Exception ex)
            {
                Console.WriteLine(ex);
                await SlashCommandErrored(ex);
                throw;
            }

        }
    }
}
