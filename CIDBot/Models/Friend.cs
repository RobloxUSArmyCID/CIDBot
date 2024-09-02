namespace CIDBot.Models
{
    internal class Friend
    {
        public bool IsOnline { get; set; }
        public long PresenceType { get; set; }
        public bool IsDeleted { get; set; }
        public double FriendFrequentScore { get; set; }
        public double FriendFrequentRank { get; set; }
        public bool HasVerifiedBadge { get; set; }
        public string? Description { get; set; }
        public DateTime? Created { get; set; }
        public bool IsBanned { get; set; }
        public string? ExternalAppDisplayName { get; set; }
        public ulong Id { get; set; }
        public string? Name { get; set; }
        public string? DisplayName { get; set; }
    }
}
