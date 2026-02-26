# 16: IP Geolocation - Mapping the Attackers 🗺️

In cybersec, knowing *who* attacked you is great, but knowing *where* they are is even better. 

## Why use Geolocation?
1. **Threat Intelligence:** If you see 90% of your SQLi attacks coming from a specific country, you might want to block that entire region at the firewall.
2. **Compliance:** Some companies are only allowed to serve traffic to specific regions.
3. **Cool Visuals:** It makes your dashboard feel like a "War Room."

## Implementation Choice: Public API
In a high-traffic production system, you would use a local database (like MaxMind GeoLite2). But for API Sentinel, we will use a **Public API lookup** to keep it simple and effective.

- **API:** `http://ip-api.com/json/{ip}`
- **Data:** We will extract the `country` and `city`.

## Performance Consideration: Caching
Looking up an IP on every request is slow. We will only look up an IP **after it has been blocked** (so we don't slow down legitimate users) and we will **cache** the result in memory so we don't spam the API.

Next, we will implement the Geolocation client!
