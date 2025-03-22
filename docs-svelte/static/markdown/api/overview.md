# Chadburn API Overview

> **Development Notice**: The Chadburn API is designed using an API-first approach, where the API specification is defined and documented before implementation. This ensures a consistent and well-thought-out API design that serves as the foundation for all client implementations.

Chadburn provides a comprehensive HTTP API that allows you to interact with the scheduler programmatically. This API enables you to create, read, update, and delete jobs, as well as monitor the status of your scheduled tasks.

## API Basics

The Chadburn API is a RESTful API that uses standard HTTP methods and returns JSON responses. All API endpoints are prefixed with `/api/v1`.

### Base URL

```
http://your-chadburn-instance:8080/api/v1
```

### Response Format

All API responses are returned in JSON format. Successful responses typically include a `data` field containing the requested information, while error responses include an `error` field with details about what went wrong.

Example success response:

```json
{
  "data": {
    "id": "job-123",
    "name": "backup-database",
    "schedule": "0 0 * * *",
    "command": "pg_dump -U postgres mydb > /backups/mydb.sql",
    "status": "active"
  }
}
```

Example error response:

```json
{
  "error": {
    "code": "invalid_schedule",
    "message": "Invalid cron schedule format"
  }
}
```

### Authentication

The Chadburn API supports several authentication methods:

- API Key Authentication
- Basic Authentication
- JWT Token Authentication

See the [Authentication](/api/authentication) section for more details.

## Available Endpoints

Chadburn's API is organized around the following resources:

### Jobs

- `GET /api/v1/jobs` - List all jobs
- `GET /api/v1/jobs/{id}` - Get a specific job
- `POST /api/v1/jobs` - Create a new job
- `PUT /api/v1/jobs/{id}` - Update a job
- `DELETE /api/v1/jobs/{id}` - Delete a job
- `POST /api/v1/jobs/{id}/run` - Run a job immediately

### Schedules

- `GET /api/v1/schedules` - List all schedules
- `GET /api/v1/schedules/{id}` - Get a specific schedule
- `POST /api/v1/schedules` - Create a new schedule
- `PUT /api/v1/schedules/{id}` - Update a schedule
- `DELETE /api/v1/schedules/{id}` - Delete a schedule
- `GET /api/v1/schedules/{id}/jobs` - Get jobs using a schedule
- `POST /api/v1/schedules/validate` - Validate a cron expression

### Users

- `GET /api/v1/users` - List all users
- `GET /api/v1/users/{id}` - Get a specific user
- `POST /api/v1/users` - Create a new user
- `PUT /api/v1/users/{id}` - Update a user
- `DELETE /api/v1/users/{id}` - Delete a user
- `POST /api/v1/users/{id}/password` - Change user password
- `GET /api/v1/users/{id}/activity` - Get user activity

### Settings

- `GET /api/v1/settings` - Get all settings
- `PUT /api/v1/settings` - Update settings
- `POST /api/v1/settings/notifications/email/test` - Test email settings
- `POST /api/v1/settings/notifications/slack/test` - Test Slack settings
- `GET /api/v1/settings/system-info` - Get system information

### Metrics

- `GET /api/v1/metrics` - Get metrics in Prometheus format
- `GET /api/v1/metrics/summary` - Get metrics summary

## Rate Limiting

To ensure fair usage and system stability, the Chadburn API implements rate limiting. By default, clients are limited to 100 requests per minute. See the [Rate Limiting](/api/rate-limiting) section for more details.

## Error Codes

The API uses standard HTTP status codes to indicate the success or failure of a request. In addition, specific error codes are provided in the response body to give more detailed information about what went wrong.

Common status codes:

- `200 OK` - The request was successful
- `201 Created` - A resource was successfully created
- `400 Bad Request` - The request was invalid
- `401 Unauthorized` - Authentication failed
- `403 Forbidden` - The client does not have permission to access the requested resource
- `404 Not Found` - The requested resource was not found
- `429 Too Many Requests` - Rate limit exceeded
- `500 Internal Server Error` - An error occurred on the server

## Getting Started

To get started with the Chadburn API, check out the [Basic Usage](/api/examples/basic) examples. 