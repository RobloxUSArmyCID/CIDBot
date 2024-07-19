namespace CIDBot.Models
{
    internal class GetUserGroupsDataModel
    {
        public Group? Group { get; set; }
        public Role? Role { get; set; }
        public bool IsNotificationsEnabled { get; set; }
    }
}
