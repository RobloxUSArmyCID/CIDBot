namespace CIDBot.Models
{
    internal class Under30MembersGroup(
        ulong id,
        string name,
        ulong memberCount,
        bool hasVerifiedBadge,
        ulong ownerId) : IGroup
    {
        public ulong Id { get; private set; } = id;
        public string Name { get; private set; } = name;
        public ulong MemberCount { get; private set; } = memberCount;
        public bool HasVerifiedBadge { get; private set; } = hasVerifiedBadge;
        public ulong OwnerId { get; set; } = ownerId;
    }
}
