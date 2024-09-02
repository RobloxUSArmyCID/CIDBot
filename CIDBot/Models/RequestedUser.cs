namespace CIDBot.Models
{
    internal class RequestedUser : IUser
    {
        public string? RequestedUsername { get; set; }
        public bool HasVerifiedBadge { get; set; }
        public ulong Id { get; set; }
        public string? Name { get; set; }
        public string? DisplayName { get; set; }
    }
}
