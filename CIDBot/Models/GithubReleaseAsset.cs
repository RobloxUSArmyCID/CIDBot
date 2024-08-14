namespace CIDBot.Models
{
    internal class GithubReleaseAsset
    {
        public string? Url { get; set; }
        public string? Browser_Download_Url { get; set; }
        public ulong Id { get; set; }
        public string? Node_Id { get; set; }
        public string? Name { get; set; }
        public string? Label { get; set; }
        public string? State { get; set; }
        public string? Content_Type { get; set; }
        public ulong Size { get; set; }
        public ulong Download_Count { get; set; }
        public DateTime? Created_At { get; set; }
        public DateTime? Uploaded_At { get; set; }
        public GithubReleaseAuthor? Uploader { get; set; }
    }
}
