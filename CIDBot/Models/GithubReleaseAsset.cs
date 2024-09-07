namespace CIDBot.Models
{
    internal class GithubReleaseAsset
    {
        public required string Url { get; set; }
        public required string BrowserDownloadUrl { get; set; }
        public required ulong Id { get; set; }
        public required string NodeId { get; set; }
        public required string Name { get; set; }
        public required string Label { get; set; }
        public required string State { get; set; }
        public required string ContentType { get; set; }
        public required ulong Size { get; set; }
        public required ulong DownloadCount { get; set; }
        public required DateTime CreatedAt { get; set; }
        public required DateTime UploadedAt { get; set; }
        public required GithubReleaseAuthor Uploader { get; set; }
    }
}
