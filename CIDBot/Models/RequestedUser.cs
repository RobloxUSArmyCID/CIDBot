using System.Text.Json.Serialization;

namespace CIDBot.Models
{
    internal class RequestedUser : IUser
    {
        public required string RequestedUsername { get; set; }
        public required bool HasVerifiedBadge { get; set; }
        public required ulong Id { get; set; }
        public string? Name { get; set; }
        public required string DisplayName { get; set; }
    }
}
