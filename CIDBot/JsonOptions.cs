using System.Text.Json;

namespace CIDBot
{
    internal static class JsonOptions
    {
        public static JsonSerializerOptions Github = new()
        {
            WriteIndented = true,
            PropertyNameCaseInsensitive = true,
            PropertyNamingPolicy = JsonNamingPolicy.SnakeCaseLower,
        };

        public static JsonSerializerOptions OtherThanGithub = new()
        {
            WriteIndented = true,
            PropertyNameCaseInsensitive = true
        };
    }
}
