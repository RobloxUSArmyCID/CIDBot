namespace CIDBot.Models
{
    internal class User : IUser
    {
        public required string Description { get; set; }
        public required DateTime Created { get; set; }
        public required bool IsBanned { get; set; }
        public required string ExternalAppDisplayName { get; set; }
        public required bool HasVerifiedBadge { get; set; }
        public required ulong Id { get; set; }
        public string? Name { get; set; }
        public required string DisplayName { get; set; }
    }
}
