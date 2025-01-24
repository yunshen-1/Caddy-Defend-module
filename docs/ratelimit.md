# Defender Module - Rate Limiting Integration Guide

**Feature:** Match requests by IP range and apply rate limiting using [caddy-ratelimit](https://github.com/mholt/caddy-ratelimit)

## Configuration Overview

### Caddyfile Syntax
```caddy
defender ratelimit {
    ranges <cidr_or_predefined...>
    # Optional rate limit marker header (default: X-Defender-RateLimit)
    rate_limit_header <header-name>
}

ratelimit {
    # Match requests marked by Defender
    match header <header-name> <value>
    
    # Rate limiting parameters
    rate  <requests-per-second>
    burst <burst-size>
    key   <rate-limit-key>
}
```

### JSON Configuration
```json
{
    "handler": "defender",
    "raw_responder": "ratelimit",
    "ranges": ["aws", "10.0.0.0/8"],
    "rate_limit_header": "X-Defender-RateLimit"
}
```

## Example Configurations

### Basic Configuration
```caddy
example.com {
    defender ratelimit {
        ranges cloudflare openai
    }
    
    ratelimit {
        match header X-Defender-RateLimit true
        rate  5r/s
        burst 10
        key   {http.request.remote.host}
    }
    
    respond "Hello World"
}
```

### Advanced Configuration
```caddy
api.example.com {
    defender ratelimit {
        ranges 192.168.1.0/24 azure
        rate_limit_header X-API-RateLimit
    }
    
    ratelimit {
        match header X-API-RateLimit true
        rate  10r/s
        burst 20
        key   {http.request.uri.path}
        
        # Optional: Custom response
        respond {
            status_code 429
            body "Too Many Requests - Try Again Later"
        }
    }
    
    reverse_proxy localhost:3000
}
```

## Documentation

### Directives

**Defender Rate Limit Responder:**
- `ranges` - IP ranges to apply rate limiting (CIDR or predefined)
- `rate_limit_header` (optional) - Header to mark requests for rate limiting (default: `X-Defender-RateLimit`)

**Rate Limit Module:**
- `match header` - Match the header set by Defender
- `rate` - Requests per second (e.g., `10r/s`)
- `burst` - Allow temporary bursts of requests
- `key` - Rate limit key (typically client IP or path)

### How It Works
1. **IP Matching:** Defender checks if client IP matches configured ranges
2. **Header Marking:** Matching requests get a header (`X-Defender-RateLimit: true`)
3. **Rate Limiting:** caddy-ratelimit applies limits only to marked requests
4. **Request Processing:** Non-matched requests bypass rate limiting

### Use Cases
- Protect API endpoints from scraping
- Mitigate brute force attacks
- Enforce different rate limits for:
  - Different geographic regions
  - Known bot networks
  - Internal vs external traffic

### Requirements

- [caddy-ratelimit](https://github.com/mholt/caddy-ratelimit) module installed
- [caddy-defender](https://github.com/JasonLovesDoggo/caddy-defender) v0.5.0+

### Notes
1. **Order Matters:** Defender must come before ratelimit in handler chain
2. **Header Customization:** Change header name if conflicts occur
3. **Combination with Other Protections:**
   ```caddy
   defender ratelimit {
       ranges aws
   }
   
   ratelimit {
       match header X-Defender-RateLimit true
       rate 2r/s
   }
   
   defender block {
       ranges known-bad-ips
   }
   ```

### Troubleshooting
1. **Check Headers:**
   ```bash
   curl -I http://example.com
   ```
2. **Verify Handler Order:** Defender → Ratelimit → Other handlers
3. **Test Rate Limits:**
   ```bash
   # Simulate requests from blocked range
   for i in {1..20}; do
       curl -H "X-Forwarded-For: 20.202.43.67" http://example.com
   done
   ```
