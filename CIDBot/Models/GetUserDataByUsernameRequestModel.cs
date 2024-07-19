namespace CIDBot.Models
{
    internal class GetUserDataByUsernameRequestModel
    {
        public IList<string>? Usernames { get; set; }
        public bool ExcludeBannedUsers { get; set; }
    }
}
