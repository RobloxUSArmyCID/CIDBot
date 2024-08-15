using CIDBot.Models;
using Discord;
using Discord.Interactions;
using Discord.WebSocket;
using System.Text.Json;

namespace CIDBot
{
    public class ModalResponses : InteractionModuleBase
    {

        [ModalInteraction("punishment-request-*")]
        public async Task RespondToPunishmentRequestModal(string id, PunishmentRequestModal modal)
        {

            StreamReader str = File.OpenText("temp.json");
            string json = str.ReadToEnd();
            TemporaryJson? temp = JsonSerializer.Deserialize<TemporaryJson>(json);
            str.Close();

            ITextChannel? ssPunishmentChannel = await Context.Client.GetChannelAsync(temp!.SectionStaffPunishmentChannel) as ITextChannel;
            ITextChannel? bsPunishmentChannel = await Context.Client.GetChannelAsync(temp!.BattalionStaffPunishmentChannel) as ITextChannel;
            IRole ssRole = ssPunishmentChannel!.Guild.GetRole(temp.SectionStaffRole);
            IRole bsRole = ssPunishmentChannel.Guild.GetRole(temp.BattalionStaffRole);

            var user = Context.User as SocketGuildUser;

            Embed sectionStaffEmbed = new EmbedBuilder()
                .WithAuthor(Context.User)
                .WithTitle($"New punishment request: {user!.Nickname}")
                .WithDescription($"""
                - Investigator: {user.Nickname}
                - Ticket ID: {modal.TicketId}
                """)
                .WithCurrentTimestamp()
                .WithColor(Color.Gold)
                .Build();

            await ssPunishmentChannel.SendMessageAsync(embed: sectionStaffEmbed, text: ssRole.Mention);

            Embed battalionStaffEmbed = new EmbedBuilder()
                .WithAuthor(Context.User)
                .WithTitle($"New punishment request: {modal.OffenderUsername}")
                .WithDescription($"""
                - Offender: {modal.OffenderUsername}
                - Investigator: {user.Nickname}
                - Ticket ID: {modal.TicketId}
                - Punishment Request Link: {modal.PunishmentRequestUrl}
                - Case File Link: {modal.CaseFileUrl}
                """)
                .WithColor(Color.DarkPurple)
                .WithCurrentTimestamp()
                .Build();

            MessageComponent components = new ComponentBuilder()
                .WithButton("Approve", $"approve-punishment-{id}", ButtonStyle.Success, Emoji.Parse(":white_check_mark:"))
                .WithButton("Deny", $"deny-punishment-{id}", ButtonStyle.Danger, Emoji.Parse(":x:"))
                .Build();

            await bsPunishmentChannel!.SendMessageAsync(text: bsRole.Mention, components: components, embed: battalionStaffEmbed);

            await RespondAsync("x", ephemeral: true);
        }
    }
}
