namespace CIDBot.Models
{
    internal class GetUserInfoByUsernameRequest
    {
        public IList<string>? Usernames { get; set; }
        public bool ExcludeBannedUsers { get; set; }
    }
}
