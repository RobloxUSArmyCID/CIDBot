using System.Text.Json.Serialization;

namespace CIDBot.Models
{
    internal class Badge
    {
        [JsonRequired]
        public required ulong Id { get; set; } 
        public required string Name { get; set; }
        public required string Description { get; set; }
        public required string DisplayName { get; set; }
        public required string DisplayDescription { get; set; }
        public required ulong IconImageId { get; set; }
        public required ulong DisplayIconImageId { get; set; }

        /*
         * 
         * BadgeAwarder and BadgeStatistics may be actually nullable - but I have no idea
         * I can't find any devforum post, however I decided to make them required in order to not have to handle any null errors at this moment
         * and in the future, if it becomes an issue (neither have any references at the moment), I can make them null-safe
         *
        */

        public required BadgeAwarder Awarder { get; set; }
        public required BadgeStatistics Statistics { get; set; }
        public required DateTime Created { get; set; }
        public required DateTime Updated { get; set; }

    }
}
