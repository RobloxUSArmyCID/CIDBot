namespace CIDBot.Models
{
    internal interface IGroup
    {
        public ulong Id { get; }
        public string Name { get; }
        public bool HasVerifiedBadge { get; }
    }
}
