# API Authentication

Securing your Chadburn API is essential to prevent unauthorized access to your job scheduling system. Chadburn provides multiple authentication methods to suit different security requirements and integration scenarios.

## Authentication Methods

### API Key Authentication

API key authentication is the simplest method to implement. Each client is assigned a unique API key that must be included in every request.

To use API key authentication:

1. Generate an API key in the Chadburn configuration file:

```yaml
api:
  enabled: true
  keys:
    - name: "service-a"
      key: "sk_chadburn_a1b2c3d4e5f6g7h8i9j0"
    - name: "service-b"
      key: "sk_chadburn_z9y8x7w6v5u4t3s2r1q0"
```

2. Include the API key in your requests using the `X-API-Key` header:

```bash
curl -X GET "http://your-chadburn-instance:8080/api/v1/jobs" \
  -H "X-API-Key: sk_chadburn_a1b2c3d4e5f6g7h8i9j0"
```

### Basic Authentication

Basic authentication uses a username and password combination. This method is widely supported but should only be used over HTTPS to ensure credentials are encrypted during transmission.

To use basic authentication:

1. Configure basic authentication in the Chadburn configuration file:

```yaml
api:
  enabled: true
  auth:
    basic:
      enabled: true
      users:
        - username: "admin"
          password: "$2a$10$JqSzUQOZrV2U5w6kfVQXCe6jd6Qj/Hl1wqwABOsZnzqyoC5xOW8jK" # hashed password
        - username: "readonly"
          password: "$2a$10$7tJxF4Rl8Yy3U/Jv1JVXDuZJfBVQlM7oeWMYQMO7zJ2Bj0qOHLOjS" # hashed password
```

2. Include the credentials in your requests using the `Authorization` header:

```bash
curl -X GET "http://your-chadburn-instance:8080/api/v1/jobs" \
  -H "Authorization: Basic $(echo -n 'admin:password' | base64)"
```

### JWT Token Authentication

JWT (JSON Web Token) authentication provides a more secure and flexible authentication mechanism. It allows for token expiration, role-based access control, and more.

To use JWT authentication:

1. Configure JWT authentication in the Chadburn configuration file:

```yaml
api:
  enabled: true
  auth:
    jwt:
      enabled: true
      secret: "your-secret-key"
      expiration: 3600 # token expiration in seconds
```

2. Obtain a JWT token by authenticating with your credentials:

```bash
curl -X POST "http://your-chadburn-instance:8080/api/v1/auth/token" \
  -H "Content-Type: application/json" \
  -d '{"username": "admin", "password": "password"}'
```

3. Include the JWT token in your requests using the `Authorization` header:

```bash
curl -X GET "http://your-chadburn-instance:8080/api/v1/jobs" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
```

## Role-Based Access Control

Chadburn supports role-based access control (RBAC) to restrict access to certain API endpoints based on user roles.

To configure RBAC:

```yaml
api:
  enabled: true
  rbac:
    roles:
      - name: "admin"
        permissions: ["*"] # all permissions
      - name: "readonly"
        permissions: ["jobs:read", "schedules:read", "metrics:read"]
    users:
      - username: "admin"
        roles: ["admin"]
      - username: "readonly"
        roles: ["readonly"]
```

## Best Practices

1. **Use HTTPS**: Always use HTTPS to encrypt API traffic and protect authentication credentials.
2. **Rotate API Keys**: Regularly rotate API keys to minimize the impact of key exposure.
3. **Implement Least Privilege**: Assign the minimum necessary permissions to each user or service.
4. **Use Strong Passwords**: Enforce strong password policies for basic authentication.
5. **Set Short Token Expiration**: For JWT authentication, set short expiration times and implement token refresh mechanisms.
6. **Monitor API Usage**: Keep track of API usage to detect suspicious activity.

## Troubleshooting

### Common Authentication Errors

- **401 Unauthorized**: The provided credentials are invalid or missing.
- **403 Forbidden**: The authenticated user does not have permission to access the requested resource.

### Debugging Authentication Issues

To debug authentication issues, enable debug logging in the Chadburn configuration:

```yaml
logging:
  level: "debug"
```

This will provide more detailed information about authentication failures in the Chadburn logs. 