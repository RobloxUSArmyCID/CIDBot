using System.Text.Json.Serialization;

namespace CIDBot.Models
{
    internal class GetLatestGithubReleaseResponse
    {
        public string? Url { get; set; }
        public string? HtmlUrl { get; set; }
        public string? AssetsUrl { get; set; }
        public string? UploadUrl { get; set; }
        public string? TarballUrl { get; set; }
        public string? ZipballUrl { get; set; }
        public string? DiscussionUrl { get; set; }
        public ulong Id { get; set; }
        public string? NodeId { get; set; }
        //[JsonPropertyName("tag_name")]
        public string? TagName { get; set; }
        public string? TargetCommitish { get; set; }
        public string? Name { get; set; }
        public string? Body { get; set; }
        public bool Draft { get; set; }
        public bool Prerelease { get; set; }
        public DateTime? CreatedAt { get; set; }
        public DateTime? PublishedAt { get; set; }
        public GithubReleaseAuthor? Author { get; set; }
        public IList<GithubReleaseAsset>? Assets { get; set; }
    }
}
