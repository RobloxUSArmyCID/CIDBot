namespace CIDBot.Models
{
    internal class Group : IGroup
    {
        public required ulong Id { get; set; }
        public required string Name { get; set; }
        public required string Description { get; set; }

        // yes a group owner can be null lmfao
        public GroupOwner? Owner { get; set; }
        public required DateTime Created { get; set; }
        public required bool HasVerifiedBadge { get; set; }
    }
}
