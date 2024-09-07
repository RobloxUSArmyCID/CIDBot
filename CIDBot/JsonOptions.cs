using System.Text.Json;

namespace CIDBot
{
    public class JsonOptions
    {
        // Github SUCKS at JSON
        public JsonSerializerOptions Github = new()
        {
            WriteIndented = true,
            PropertyNameCaseInsensitive = true,
            PropertyNamingPolicy = JsonNamingPolicy.SnakeCaseLower,
        };

        // camel case lower FTW
        public JsonSerializerOptions OtherThanGithub = new()
        {
            WriteIndented = true,
            PropertyNameCaseInsensitive = true
        };
    }
}
