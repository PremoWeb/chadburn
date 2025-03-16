# API Rate Limiting

To ensure fair usage and protect the Chadburn API from abuse or excessive load, rate limiting is implemented. This page explains how rate limiting works and how to handle rate limit errors.

## Rate Limit Configuration

By default, the Chadburn API limits clients to 100 requests per minute. This limit can be configured in the Chadburn configuration file:

```yaml
api:
  enabled: true
  rate_limiting:
    enabled: true
    requests_per_minute: 100
    burst: 20
```

The configuration options are:

- `enabled`: Whether rate limiting is enabled (default: `true`)
- `requests_per_minute`: The maximum number of requests allowed per minute (default: `100`)
- `burst`: The maximum number of requests that can be made in a short burst (default: `20`)

## Rate Limit Headers

When you make a request to the Chadburn API, the response includes headers that provide information about your current rate limit status:

- `X-RateLimit-Limit`: The maximum number of requests you can make per minute
- `X-RateLimit-Remaining`: The number of requests remaining in the current rate limit window
- `X-RateLimit-Reset`: The time at which the current rate limit window resets, in Unix epoch seconds

Example response headers:

```
HTTP/1.1 200 OK
Content-Type: application/json
X-RateLimit-Limit: 100
X-RateLimit-Remaining: 95
X-RateLimit-Reset: 1625097600
```

## Handling Rate Limit Errors

If you exceed the rate limit, the API will respond with a `429 Too Many Requests` status code and a JSON error message:

```json
{
  "error": {
    "code": "rate_limit_exceeded",
    "message": "Rate limit exceeded. Please try again later."
  }
}
```

The response will also include the `Retry-After` header, which indicates the number of seconds to wait before making another request:

```
HTTP/1.1 429 Too Many Requests
Content-Type: application/json
Retry-After: 30
```

### Best Practices for Handling Rate Limits

1. **Monitor rate limit headers**: Keep track of the `X-RateLimit-Remaining` header to avoid hitting the limit.

2. **Implement exponential backoff**: When you receive a 429 response, wait for the time specified in the `Retry-After` header before retrying, and increase the wait time exponentially for consecutive failures.

3. **Optimize your requests**: Reduce the number of API calls by batching operations when possible.

4. **Cache responses**: Cache API responses when appropriate to reduce the number of requests.

## Rate Limit Exemptions

In some cases, you may need to exempt certain clients from rate limiting. This can be configured in the Chadburn configuration file:

```yaml
api:
  enabled: true
  rate_limiting:
    enabled: true
    requests_per_minute: 100
    exempt_ips:
      - "10.0.0.1"
      - "192.168.1.0/24"
    exempt_api_keys:
      - "sk_chadburn_a1b2c3d4e5f6g7h8i9j0"
```

The exemption options are:

- `exempt_ips`: A list of IP addresses or CIDR ranges that are exempt from rate limiting
- `exempt_api_keys`: A list of API keys that are exempt from rate limiting

## Custom Rate Limits

You can also configure different rate limits for different API endpoints or client types:

```yaml
api:
  enabled: true
  rate_limiting:
    enabled: true
    default:
      requests_per_minute: 100
    custom:
      - path: "/api/v1/jobs"
        method: "GET"
        requests_per_minute: 200
      - path: "/api/v1/metrics"
        method: "*"
        requests_per_minute: 50
```

This configuration sets:

- A default limit of 100 requests per minute for all endpoints
- A custom limit of 200 requests per minute for GET requests to `/api/v1/jobs`
- A custom limit of 50 requests per minute for all methods on `/api/v1/metrics`

## Monitoring Rate Limit Usage

Chadburn provides metrics for monitoring rate limit usage:

- `chadburn_api_rate_limit_exceeded_total`: The total number of requests that exceeded the rate limit
- `chadburn_api_requests_total`: The total number of API requests, labeled by endpoint and status

These metrics can be accessed through the Prometheus metrics endpoint at `/metrics`. 