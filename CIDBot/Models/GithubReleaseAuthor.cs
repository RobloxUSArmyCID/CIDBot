namespace CIDBot.Models
{
    internal class GithubReleaseAuthor
    {
        public required string Login { get; set; }
        public required ulong Id { get; set; }
        public required string NodeId { get; set; }
        public required string AvatarUrl { get; set; }
        public required string GravatarId { get; set; }
        public required string Url { get; set; }
        public required string HtmlUrl { get; set; }
        public required string FollowersUrl { get; set; }
        public required string FollowingUrl { get; set; }
        public required string GistsUrl { get; set; }
        public required string StarredUrl { get; set; }
        public required string SubscriptionsUrl { get; set; }
        public required string OrganizationsUrl { get; set; }
        public required string ReposUrl { get; set; }
        public required string EventsUrl { get; set; }
        public required string ReceivedEventsUrl { get; set; }
        public required string Type { get; set; }
        public required bool SiteAdmin { get; set; }

    }
}
