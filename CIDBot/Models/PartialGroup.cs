namespace CIDBot.Models
{
    internal partial class PartialGroup : IGroup
    {
        public required ulong Id { get; set; }
        public required string Name { get; set; }
        public required ulong MemberCount { get; set; }
        public required bool HasVerifiedBadge { get; set; }

    }
}
