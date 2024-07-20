namespace CIDBot.Models
{
    internal class GetOwnedBadgesByIdResponseData
    {
        public ulong Id { get; set; } 
        public string? Name { get; set; }
        public string? Description { get; set; }
        public string? DisplayName { get; set; }
        public string? DisplayDescription { get; set; }
        public ulong IconImageId { get; set; }
        public ulong DisplayIconImageId { get; set; }
        public BadgeAwarder? Awarder { get; set; }
        public BadgeStatistics? Statistics { get; set; }
        public DateTime? Created { get; set; }
        public DateTime? Updated { get; set; }

    }
}
