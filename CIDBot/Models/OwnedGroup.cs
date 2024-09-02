namespace CIDBot.Models
{
    internal class OwnedGroup : IGroup
    {
        public ulong Id { get; set; }
        public string? Name { get; set; }
        public string? Description { get; set; }
        public GroupOwner? Owner { get; set; }
        public DateTime? Created { get; set; }
        public bool HasVerifiedBadge { get; set; }

    }
}
