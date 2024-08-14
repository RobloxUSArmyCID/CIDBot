namespace CIDBot.Models
{
    internal class GithubReleaseAsset
    {
        public string? Url { get; set; }
        public string? BrowserDownloadUrl { get; set; }
        public ulong Id { get; set; }
        public string? NodeId { get; set; }
        public string? Name { get; set; }
        public string? Label { get; set; }
        public string? State { get; set; }
        public string? ContentType { get; set; }
        public ulong Size { get; set; }
        public ulong DownloadCount { get; set; }
        public DateTime? CreatedAt { get; set; }
        public DateTime? UploadedAt { get; set; }
        public GithubReleaseAuthor? Uploader { get; set; }
    }
}
