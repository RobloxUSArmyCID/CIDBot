namespace CIDBot.Models
{
    internal class BadgeAwarder : IUser
    {
        public required ulong Id { get; set; }
        public required string Type { get; set; }
        public required string Name { get; set; }
    }
}
