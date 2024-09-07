namespace CIDBot.Models
{
    internal class Under30MembersGroup(
        ulong id,
        string name,

        bool hasVerifiedBadge,
        ulong ownerId,
        string ownerUsername) : IGroup
    {
        public ulong Id { get; private set; } = id;
        public string? Name { get; private set; } = name;
        public bool HasVerifiedBadge { get; private set; } = hasVerifiedBadge;
        public ulong OwnerId { get; set; } = ownerId;
        public string OwnerUsername { get; set; } = ownerUsername;
    }
}
