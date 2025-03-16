# Tags API Endpoints

The Tags API allows you to manage tags for jobs and schedules in Chadburn. Tags provide a way to categorize and filter jobs and schedules, making it easier to manage large numbers of scheduled tasks.

## List All Tags

```
GET /api/v1/tags
```

Retrieves a list of all tags used in the system.

### Query Parameters

| Parameter  | Type    | Description                                                |
|------------|---------|-----------------------------------------------------------|
| `page`     | integer | Page number for pagination (default: 1)                    |
| `per_page` | integer | Number of tags per page (default: 50, max: 200)            |
| `sort`     | string  | Sort by: 'name', 'count' (default: 'name')                 |
| `order`    | string  | Sort order: 'asc', 'desc' (default: 'asc')                 |
| `search`   | string  | Search term to filter tags by name (optional)              |

### Response

```json
{
  "data": [
    {
      "name": "backup",
      "count": 5,
      "created_at": "2023-01-01T00:00:00Z",
      "updated_at": "2023-01-02T00:00:00Z"
    },
    {
      "name": "critical",
      "count": 8,
      "created_at": "2023-01-01T00:00:00Z",
      "updated_at": "2023-01-03T00:00:00Z"
    },
    {
      "name": "maintenance",
      "count": 3,
      "created_at": "2023-01-01T12:00:00Z",
      "updated_at": "2023-01-01T12:00:00Z"
    },
    {
      "name": "production",
      "count": 15,
      "created_at": "2023-01-01T00:00:00Z",
      "updated_at": "2023-01-03T12:00:00Z"
    }
  ],
  "meta": {
    "page": 1,
    "per_page": 50,
    "total": 4,
    "total_pages": 1
  }
}
```

## Get Jobs by Tag

```
GET /api/v1/tags/{name}/jobs
```

Retrieves all jobs that have the specified tag.

### Path Parameters

| Parameter | Type   | Description |
|-----------|--------|-------------|
| `name`    | string | Tag name    |

### Query Parameters

| Parameter  | Type    | Description                                                |
|------------|---------|-----------------------------------------------------------|
| `page`     | integer | Page number for pagination (default: 1)                    |
| `per_page` | integer | Number of jobs per page (default: 20, max: 100)            |
| `status`   | string  | Filter by status: 'active', 'paused' (optional)            |

### Response

```json
{
  "data": [
    {
      "id": "job-123",
      "name": "Database Backup",
      "command": "pg_dump -U postgres -d mydb > /backups/mydb.sql",
      "schedule": "0 0 * * *",
      "tags": ["backup", "production", "critical"],
      "status": "active",
      "created_at": "2023-01-01T00:00:00Z",
      "updated_at": "2023-01-02T00:00:00Z"
    },
    {
      "id": "job-456",
      "name": "Log Rotation",
      "command": "/usr/local/bin/rotate-logs.sh",
      "schedule": "0 0 * * 0",
      "tags": ["maintenance", "production"],
      "status": "active",
      "created_at": "2023-01-01T12:00:00Z",
      "updated_at": "2023-01-01T12:00:00Z"
    }
  ],
  "meta": {
    "page": 1,
    "per_page": 20,
    "total": 2,
    "total_pages": 1
  }
}
```

## Get Schedules by Tag

```
GET /api/v1/tags/{name}/schedules
```

Retrieves all schedules that have the specified tag.

### Path Parameters

| Parameter | Type   | Description |
|-----------|--------|-------------|
| `name`    | string | Tag name    |

### Query Parameters

| Parameter  | Type    | Description                                                |
|------------|---------|-----------------------------------------------------------|
| `page`     | integer | Page number for pagination (default: 1)                    |
| `per_page` | integer | Number of schedules per page (default: 20, max: 100)       |

### Response

```json
{
  "data": [
    {
      "id": "schedule-123",
      "name": "Daily Backup Schedule",
      "cron": "0 0 * * *",
      "timezone": "UTC",
      "tags": ["backup", "production"],
      "created_at": "2023-01-01T00:00:00Z",
      "updated_at": "2023-01-02T00:00:00Z"
    },
    {
      "id": "schedule-456",
      "name": "Weekly Maintenance",
      "cron": "0 0 * * 0",
      "timezone": "America/New_York",
      "tags": ["maintenance", "production"],
      "created_at": "2023-01-01T12:00:00Z",
      "updated_at": "2023-01-01T12:00:00Z"
    }
  ],
  "meta": {
    "page": 1,
    "per_page": 20,
    "total": 2,
    "total_pages": 1
  }
}
```

## Create a Tag

```
POST /api/v1/tags
```

Creates a new tag.

### Request Body

```json
{
  "name": "development"
}
```

### Response

```json
{
  "data": {
    "name": "development",
    "count": 0,
    "created_at": "2023-01-04T00:00:00Z",
    "updated_at": "2023-01-04T00:00:00Z"
  }
}
```

## Delete a Tag

```
DELETE /api/v1/tags/{name}
```

Deletes a tag. This will remove the tag from all jobs and schedules that have it.

### Path Parameters

| Parameter | Type   | Description |
|-----------|--------|-------------|
| `name`    | string | Tag name    |

### Response

```json
{
  "data": {
    "message": "Tag deleted successfully",
    "affected_jobs": 3,
    "affected_schedules": 2
  }
}
```

## Add Tag to Job

```
POST /api/v1/jobs/{id}/tags
```

Adds a tag to a job.

### Path Parameters

| Parameter | Type   | Description |
|-----------|--------|-------------|
| `id`      | string | Job ID      |

### Request Body

```json
{
  "name": "development"
}
```

### Response

```json
{
  "data": {
    "job_id": "job-123",
    "tags": ["backup", "production", "critical", "development"]
  }
}
```

## Remove Tag from Job

```
DELETE /api/v1/jobs/{id}/tags/{name}
```

Removes a tag from a job.

### Path Parameters

| Parameter | Type   | Description |
|-----------|--------|-------------|
| `id`      | string | Job ID      |
| `name`    | string | Tag name    |

### Response

```json
{
  "data": {
    "job_id": "job-123",
    "tags": ["backup", "production", "critical"]
  }
}
```

## Add Tag to Schedule

```
POST /api/v1/schedules/{id}/tags
```

Adds a tag to a schedule.

### Path Parameters

| Parameter | Type   | Description |
|-----------|--------|-------------|
| `id`      | string | Schedule ID |

### Request Body

```json
{
  "name": "development"
}
```

### Response

```json
{
  "data": {
    "schedule_id": "schedule-123",
    "tags": ["backup", "production", "development"]
  }
}
```

## Remove Tag from Schedule

```
DELETE /api/v1/schedules/{id}/tags/{name}
```

Removes a tag from a schedule.

### Path Parameters

| Parameter | Type   | Description |
|-----------|--------|-------------|
| `id`      | string | Schedule ID |
| `name`    | string | Tag name    |

### Response

```json
{
  "data": {
    "schedule_id": "schedule-123",
    "tags": ["backup", "production"]
  }
}
```

## Batch Update Tags

```
POST /api/v1/tags/batch
```

Updates tags for multiple jobs or schedules in a single operation.

### Request Body

```json
{
  "operation": "add",
  "tags": ["high-priority", "monitored"],
  "jobs": ["job-123", "job-456"],
  "schedules": ["schedule-123"]
}
```

Supported operations: `add`, `remove`, `replace`

### Response

```json
{
  "data": {
    "message": "Tags updated successfully",
    "affected_jobs": 2,
    "affected_schedules": 1
  }
}
```

## Get Tag Statistics

```
GET /api/v1/tags/stats
```

Retrieves statistics about tag usage.

### Response

```json
{
  "data": {
    "total_tags": 12,
    "most_used": [
      {
        "name": "production",
        "count": 15
      },
      {
        "name": "critical",
        "count": 8
      },
      {
        "name": "backup",
        "count": 5
      }
    ],
    "recently_added": [
      {
        "name": "development",
        "created_at": "2023-01-04T00:00:00Z"
      },
      {
        "name": "testing",
        "created_at": "2023-01-03T12:00:00Z"
      }
    ],
    "usage_by_resource_type": {
      "jobs": 42,
      "schedules": 18
    }
  }
}
``` 