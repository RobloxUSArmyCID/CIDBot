namespace CIDBot.Models
{
    internal class ResponseData<T> where T : class
    {
        public IList<T>? Data { get; set; }
    }
}
