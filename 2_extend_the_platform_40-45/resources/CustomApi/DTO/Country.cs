// Root myDeserializedClass = JsonSerializer.Deserialize<List<Root>>(myJsonResponse);
using System.Collections.Generic;
using Newtonsoft.Json;

namespace CustomApi.DTO
{
    public class Country
    {
        [JsonProperty("capital")]
        public List<string> Capital { get; set; }

        [JsonProperty("population")]
        public int Population { get; set; }
    }
}
