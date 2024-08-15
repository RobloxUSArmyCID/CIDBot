using Discord.WebSocket;

namespace CIDBot.Models
{
    internal class TemporaryJson
    {
        public ulong SectionStaffPunishmentChannel { get; set; }
        public ulong BattalionStaffPunishmentChannel { get; set; }
        public ulong SectionStaffRole { get; set; }
        public ulong BattalionStaffRole { get; set; }
    }
}
