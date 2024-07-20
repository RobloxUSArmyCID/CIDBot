namespace CIDBot.Models
{
    internal class GetPastUsernamesResponse
    {
        public string? PreviousPageCursor { get; set; }
        public string? NextPageCursor { get; set; }
        public IList<GetPastUsernamesResponseData>? Data { get; set; }
    }
}
