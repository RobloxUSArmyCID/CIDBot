namespace CIDBot.Models
{
    internal class User
    {
        public required string RequestedUsername { get; set; }
        public required bool HasVerifiedBadge { get; set; }
        public required ulong Id { get; set; }
        public required string Name { get; set; }
        public required string DisplayName { get; set; }
    }
}
