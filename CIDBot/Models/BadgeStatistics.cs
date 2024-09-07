namespace CIDBot.Models
{
    internal class BadgeStatistics
    {
        public required ulong PastDayAwardedCount { get; set; }
        public required ulong AwardedCount { get; set; }
        public required double WinRatePercentage { get; set; }
    }
}
