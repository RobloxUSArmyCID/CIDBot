namespace CIDBot.Models
{
    internal class ResponseData<T> 
    {
        public string? PreviousPageCursor { get; set; }
        public string? NextPageCursor { get; set; }

        public required IList<T> Data { get; set; }
    }
}
