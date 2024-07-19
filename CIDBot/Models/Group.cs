namespace CIDBot.Models
{
    internal class Group
    {
        public int Id { get; set; }
        public string? Name { get; set; }
        public int MemberCount { get; set; }
        public bool HasVerifiedBadge { get; set; }

    }
}
