using System.Text.Json.Serialization;

namespace CIDBot.Models
{
    internal class GithubRelease
    {
        public required string Url { get; set; }
        public required string HtmlUrl { get; set; }
        public required string AssetsUrl { get; set; }
        public required string UploadUrl { get; set; }
        public required string TarballUrl { get; set; }
        public required string ZipballUrl { get; set; }
        public string? DiscussionUrl { get; set; }
        public required ulong Id { get; set; }
        public required string NodeId { get; set; }
        public required string TagName { get; set; }
        public required string TargetCommitish { get; set; }
        public required string Name { get; set; }
        public string? Body { get; set; }
        public required bool Draft { get; set; }
        public required bool Prerelease { get; set; }
        public required DateTime CreatedAt { get; set; }
        public required DateTime PublishedAt { get; set; }
        public required GithubReleaseAuthor Author { get; set; }
        public required IList<GithubReleaseAsset> Assets { get; set; }
    }
}
