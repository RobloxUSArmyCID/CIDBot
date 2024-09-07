namespace CIDBot.Models
{
    internal class GetUserInfoByUsernameRequest
    {
        public required IList<string>? Usernames { get; set; }
        public required bool ExcludeBannedUsers { get; set; }
    }
}
