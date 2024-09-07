using System.Text.Json;
using System.Text.Json.Serialization;

namespace CIDBot.Models
{
    internal interface IUser
    {
        public ulong Id { get; set; }

        [JsonRequired]
        public string Name { get; set; }
    }
}
