namespace CIDBot.Models
{
    internal class AvatarHeadshot
    {
        public required ulong TargetId { get; set; }
        public required string State { get; set; }
        public required string ImageUrl { get; set; }
        public required string Version { get; set; }
    }
}
