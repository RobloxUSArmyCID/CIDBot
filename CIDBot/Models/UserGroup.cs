namespace CIDBot.Models
{
    internal class UserGroup
    {
        public PartialGroup? Group { get; set; }
        public GroupRole? Role { get; set; }
        public bool? IsNotificationsEnabled { get; set; }
    }
}
