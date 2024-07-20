namespace CIDBot.Models
{
    internal class GetOwnedBadgesByIdResponse
    {
        public string? PreviousPageCursor { get; set; }
        public string? NextPageCursor { get; set; }
        public IList<GetOwnedBadgesByIdResponseData>? Data { get; set; }

    }
}
