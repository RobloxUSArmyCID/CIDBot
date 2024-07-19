namespace CIDBot.Models
{
    internal class GetUserGroupsResponseDataModel
    {
        public GroupModel? Group { get; set; }
        public RoleModel? Role { get; set; }
        public bool IsNotificationsEnabled { get; set; }
    }
}
