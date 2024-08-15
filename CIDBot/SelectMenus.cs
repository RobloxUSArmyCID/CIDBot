using CIDBot.Models;
using Discord;
using Discord.Interactions;

namespace CIDBot
{
    public class SelectMenus : InteractionModuleBase
    {
        [ComponentInteraction("file-*")]
        public async Task FileSelection(string id, string[] filesSelected)
        {
            await RespondWithModalAsync<PunishmentRequestModal>($"punishment-request-{id}");
        }
    }
}
