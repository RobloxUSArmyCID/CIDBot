using Discord.WebSocket;
using Microsoft.Extensions.DependencyInjection;

namespace CIDBot
{
    internal class OnSlashCommand(ServiceProvider serviceProvider)
    {

        readonly DiscordSocketClient Client = serviceProvider.GetRequiredService<DiscordSocketClient>();
        readonly static HttpClient GroupsClient = new()
        {
            BaseAddress = new Uri("https://groups.roblox.com/")
        };


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
