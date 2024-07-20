namespace CIDBot.Models
{
    internal class GetUserInfoByIdResponse
    {
        public string? Description { get; set; }
        public DateTime? Created { get; set; }
        public bool IsBanned { get; set; }
        public string? ExternalAppDisplayName { get; set; }
        public bool HasVerifiedBadge { get; set; }
        public ulong Id { get; set; }
        public string? Name { get; set; }
        public string? DisplayName { get; set; }
    }
}
