namespace CIDBot.Models
{
    internal class GetUserInfoByUsernameResponseData
    {
        public string? RequestedUsername { get; set; }
        public bool HasVerifiedBadge { get; set; }
        public ulong Id { get; set; }
        public string? Name { get; set; }
        public string? DisplayName { get; set; }
    }
}
