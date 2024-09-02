namespace CIDBot.Models
{
    internal class Under30MembersGroup : PartialGroup
    {
        public Under30MembersGroup(
            ulong id,
            string name,
            ulong memberCount,
            bool hasVerifiedBadge,
            ulong ownerId,
            string ownerUsername)
        {
            Id = id;
            Name = name;
            MemberCount = memberCount;
            HasVerifiedBadge = hasVerifiedBadge;
            OwnerId = ownerId;
            OwnerUsername = ownerUsername;
        }

        public ulong OwnerId { get; set; }
        public string OwnerUsername { get; set; }
    }
}
