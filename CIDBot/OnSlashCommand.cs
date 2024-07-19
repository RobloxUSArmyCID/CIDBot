using Discord.WebSocket;
using Microsoft.Extensions.DependencyInjection;

namespace CIDBot
{
    internal class OnSlashCommand(ServiceProvider serviceProvider)
    {
        readonly DiscordSocketClient Client = serviceProvider.GetRequiredService<DiscordSocketClient>();

        public async Task HandleSlashCommand(SocketSlashCommand cmd)
        {
            if (cmd.CommandName == "bgcheck") await OnBgcheckCommand(cmd);
        }

        async Task OnBgcheckCommand(SocketSlashCommand cmd)
        {
            await cmd.DeferAsync();
        }
    }
}
