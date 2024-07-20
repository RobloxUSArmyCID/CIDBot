namespace CIDBot.Models
{
    internal class GetUserGroupsResponseData
    {
        public Group? Group { get; set; }
        public GroupRole? Role { get; set; }
        public bool? IsNotificationsEnabled { get; set; }
    }
}
