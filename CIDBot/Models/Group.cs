namespace CIDBot.Models
{
    internal class Group
    {
        public ulong Id { get; set; }
        public string? Name { get; set; }
        public ulong MemberCount { get; set; }
        public bool HasVerifiedBadge { get; set; }

    }
}
