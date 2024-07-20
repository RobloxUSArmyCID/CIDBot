namespace CIDBot.Models
{
    internal class GetUserGroupsResponseData
    {
        public Group? Group { get; set; }
        public Role? Role { get; set; }
        public bool? IsNotificationsEnabled { get; set; }
    }
}
