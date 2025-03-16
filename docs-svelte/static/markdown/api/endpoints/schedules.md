# Schedules API Endpoints

The Schedules API allows you to manage schedule definitions in Chadburn. Schedules can be reused across multiple jobs, making it easier to maintain consistent timing across your system.

## List All Schedules

```
GET /api/v1/schedules
```

Retrieves a list of all schedules.

### Query Parameters

| Parameter | Type   | Description                                      |
|-----------|--------|--------------------------------------------------|
| `limit`   | int    | Maximum number of schedules to return (default: 100) |
| `offset`  | int    | Number of schedules to skip (default: 0)         |

### Response

```json
{
  "data": [
    {
      "id": "schedule-123",
      "name": "daily-midnight",
      "cron": "0 0 * * *",
      "description": "Runs every day at midnight",
      "created_at": "2023-01-01T00:00:00Z",
      "updated_at": "2023-01-01T00:00:00Z"
    },
    {
      "id": "schedule-456",
      "name": "weekly-sunday",
      "cron": "0 0 * * 0",
      "description": "Runs every Sunday at midnight",
      "created_at": "2023-01-02T00:00:00Z",
      "updated_at": "2023-01-02T00:00:00Z"
    }
  ],
  "pagination": {
    "total": 2,
    "limit": 100,
    "offset": 0
  }
}
```

## Get a Specific Schedule

```
GET /api/v1/schedules/{id}
```

Retrieves a specific schedule by ID.

### Path Parameters

| Parameter | Type   | Description |
|-----------|--------|-------------|
| `id`      | string | Schedule ID |

### Response

```json
{
  "data": {
    "id": "schedule-123",
    "name": "daily-midnight",
    "cron": "0 0 * * *",
    "description": "Runs every day at midnight",
    "created_at": "2023-01-01T00:00:00Z",
    "updated_at": "2023-01-01T00:00:00Z",
    "jobs": [
      {
        "id": "job-123",
        "name": "backup-database"
      },
      {
        "id": "job-456",
        "name": "cleanup-logs"
      }
    ],
    "next_run": "2023-01-03T00:00:00Z"
  }
}
```

## Create a New Schedule

```
POST /api/v1/schedules
```

Creates a new schedule.

### Request Body

```json
{
  "name": "daily-midnight",
  "cron": "0 0 * * *",
  "description": "Runs every day at midnight"
}
```

### Request Parameters

| Parameter     | Type   | Required | Description                                |
|---------------|--------|----------|--------------------------------------------|
| `name`        | string | Yes      | Schedule name                              |
| `cron`        | string | Yes      | Cron schedule expression                   |
| `description` | string | No       | Human-readable description of the schedule |

### Response

```json
{
  "data": {
    "id": "schedule-123",
    "name": "daily-midnight",
    "cron": "0 0 * * *",
    "description": "Runs every day at midnight",
    "created_at": "2023-01-01T00:00:00Z",
    "updated_at": "2023-01-01T00:00:00Z",
    "next_run": "2023-01-02T00:00:00Z"
  }
}
```

## Update a Schedule

```
PUT /api/v1/schedules/{id}
```

Updates an existing schedule.

### Path Parameters

| Parameter | Type   | Description |
|-----------|--------|-------------|
| `id`      | string | Schedule ID |

### Request Body

```json
{
  "name": "daily-1am",
  "cron": "0 1 * * *",
  "description": "Runs every day at 1 AM"
}
```

### Response

```json
{
  "data": {
    "id": "schedule-123",
    "name": "daily-1am",
    "cron": "0 1 * * *",
    "description": "Runs every day at 1 AM",
    "created_at": "2023-01-01T00:00:00Z",
    "updated_at": "2023-01-03T12:34:56Z",
    "next_run": "2023-01-04T01:00:00Z"
  }
}
```

## Delete a Schedule

```
DELETE /api/v1/schedules/{id}
```

Deletes a schedule. This will fail if the schedule is currently in use by any jobs.

### Path Parameters

| Parameter | Type   | Description |
|-----------|--------|-------------|
| `id`      | string | Schedule ID |

### Response

```json
{
  "data": {
    "message": "Schedule deleted successfully"
  }
}
```

## Get Jobs Using a Schedule

```
GET /api/v1/schedules/{id}/jobs
```

Retrieves all jobs that use a specific schedule.

### Path Parameters

| Parameter | Type   | Description |
|-----------|--------|-------------|
| `id`      | string | Schedule ID |

### Response

```json
{
  "data": [
    {
      "id": "job-123",
      "name": "backup-database",
      "command": "pg_dump -U postgres mydb > /backups/mydb.sql",
      "status": "active",
      "created_at": "2023-01-01T00:00:00Z",
      "updated_at": "2023-01-01T00:00:00Z"
    },
    {
      "id": "job-456",
      "name": "cleanup-logs",
      "command": "find /var/log -type f -name \"*.log\" -mtime +7 -delete",
      "status": "active",
      "created_at": "2023-01-02T00:00:00Z",
      "updated_at": "2023-01-02T00:00:00Z"
    }
  ],
  "pagination": {
    "total": 2,
    "limit": 100,
    "offset": 0
  }
}
```

## Validate a Cron Expression

```
POST /api/v1/schedules/validate
```

Validates a cron expression and returns the next run times.

### Request Body

```json
{
  "cron": "0 0 * * *"
}
```

### Response

```json
{
  "data": {
    "valid": true,
    "next_runs": [
      "2023-01-02T00:00:00Z",
      "2023-01-03T00:00:00Z",
      "2023-01-04T00:00:00Z",
      "2023-01-05T00:00:00Z",
      "2023-01-06T00:00:00Z"
    ]
  }
}
```

If the cron expression is invalid:

```json
{
  "data": {
    "valid": false,
    "error": "Invalid cron expression: field 'minute' has invalid value '60'"
  }
}
``` 