namespace CIDBot.Models
{
    internal class UserGroup
    {
        public required PartialGroup Group { get; set; }
        public required GroupRole Role { get; set; }
        public bool? IsNotificationsEnabled { get; set; }
    }
}
