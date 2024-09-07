namespace CIDBot.Models
{
    internal class OwnedGroup : IGroup
    {
        public required ulong Id { get; set; }
        public required string Name { get; set; }
        public required string Description { get; set; }
        public required GroupOwner Owner { get; set; }
        public required DateTime Created { get; set; }
        public required bool HasVerifiedBadge { get; set; }
    }
}
