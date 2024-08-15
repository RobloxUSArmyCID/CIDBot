using System.Text.Json;

namespace CIDBot
{
    public class JsonOptions
    {
        public JsonSerializerOptions Github = new()
        {
            WriteIndented = true,
            PropertyNameCaseInsensitive = true,
            PropertyNamingPolicy = JsonNamingPolicy.SnakeCaseLower,
        };

        public JsonSerializerOptions OtherThanGithub = new()
        {
            WriteIndented = true,
            PropertyNameCaseInsensitive = true
        };
    }
}
