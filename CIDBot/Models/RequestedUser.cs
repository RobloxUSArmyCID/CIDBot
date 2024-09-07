using System.Text.Json.Serialization;

namespace CIDBot.Models
{
    internal class RequestedUser : IUser
    {
        [JsonRequired]
        public required string RequestedUsername { get; set; }

        [JsonRequired]
        public bool HasVerifiedBadge { get; set; }
        public ulong Id { get; set; }
        public string? Name { get; set; }
        public string? DisplayName { get; set; }
    }
}
