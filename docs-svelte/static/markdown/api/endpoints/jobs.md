# Jobs API Endpoints

The Jobs API allows you to manage scheduled jobs in Chadburn. You can create, read, update, and delete jobs using these endpoints.

## List All Jobs

```
GET /api/v1/jobs
```

Retrieves a list of all jobs.

### Query Parameters

| Parameter | Type   | Description                                      |
|-----------|--------|--------------------------------------------------|
| `limit`   | int    | Maximum number of jobs to return (default: 100)  |
| `offset`  | int    | Number of jobs to skip (default: 0)              |
| `status`  | string | Filter by job status (active, paused, completed) |
| `tag`     | string | Filter by job tag                                |

### Response

```json
{
  "data": [
    {
      "id": "job-123",
      "name": "backup-database",
      "schedule": "0 0 * * *",
      "command": "pg_dump -U postgres mydb > /backups/mydb.sql",
      "status": "active",
      "tags": ["backup", "database"],
      "created_at": "2023-01-01T00:00:00Z",
      "updated_at": "2023-01-01T00:00:00Z"
    },
    {
      "id": "job-456",
      "name": "cleanup-logs",
      "schedule": "0 0 * * 0",
      "command": "find /var/log -type f -name \"*.log\" -mtime +7 -delete",
      "status": "active",
      "tags": ["cleanup", "logs"],
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

## Get a Specific Job

```
GET /api/v1/jobs/{id}
```

Retrieves a specific job by ID.

### Path Parameters

| Parameter | Type   | Description |
|-----------|--------|-------------|
| `id`      | string | Job ID      |

### Response

```json
{
  "data": {
    "id": "job-123",
    "name": "backup-database",
    "schedule": "0 0 * * *",
    "command": "pg_dump -U postgres mydb > /backups/mydb.sql",
    "status": "active",
    "tags": ["backup", "database"],
    "created_at": "2023-01-01T00:00:00Z",
    "updated_at": "2023-01-01T00:00:00Z",
    "last_run": {
      "status": "success",
      "started_at": "2023-01-02T00:00:00Z",
      "finished_at": "2023-01-02T00:00:05Z",
      "exit_code": 0,
      "output": "Database backup completed successfully"
    },
    "next_run": "2023-01-03T00:00:00Z"
  }
}
```

## Create a New Job

```
POST /api/v1/jobs
```

Creates a new job.

### Request Body

```json
{
  "name": "backup-database",
  "schedule": "0 0 * * *",
  "command": "pg_dump -U postgres mydb > /backups/mydb.sql",
  "tags": ["backup", "database"],
  "timeout": 3600,
  "retries": 3,
  "retry_delay": 300,
  "environment": {
    "PGPASSWORD": "secret"
  },
  "working_directory": "/opt/backups"
}
```

### Request Parameters

| Parameter           | Type             | Required | Description                                                |
|---------------------|------------------|----------|------------------------------------------------------------|
| `name`              | string           | Yes      | Job name                                                   |
| `schedule`          | string           | Yes      | Cron schedule expression                                   |
| `command`           | string           | Yes      | Command to execute                                         |
| `tags`              | array of strings | No       | Tags for categorizing jobs                                 |
| `timeout`           | int              | No       | Maximum execution time in seconds (default: 3600)          |
| `retries`           | int              | No       | Number of retry attempts on failure (default: 0)           |
| `retry_delay`       | int              | No       | Delay between retry attempts in seconds (default: 60)      |
| `environment`       | object           | No       | Environment variables for the job                          |
| `working_directory` | string           | No       | Working directory for the job                              |

### Response

```json
{
  "data": {
    "id": "job-123",
    "name": "backup-database",
    "schedule": "0 0 * * *",
    "command": "pg_dump -U postgres mydb > /backups/mydb.sql",
    "status": "active",
    "tags": ["backup", "database"],
    "timeout": 3600,
    "retries": 3,
    "retry_delay": 300,
    "environment": {
      "PGPASSWORD": "******"
    },
    "working_directory": "/opt/backups",
    "created_at": "2023-01-01T00:00:00Z",
    "updated_at": "2023-01-01T00:00:00Z",
    "next_run": "2023-01-02T00:00:00Z"
  }
}
```

## Update a Job

```
PUT /api/v1/jobs/{id}
```

Updates an existing job.

### Path Parameters

| Parameter | Type   | Description |
|-----------|--------|-------------|
| `id`      | string | Job ID      |

### Request Body

```json
{
  "name": "backup-database-daily",
  "schedule": "0 0 * * *",
  "command": "pg_dump -U postgres mydb > /backups/mydb-$(date +%Y%m%d).sql",
  "status": "paused",
  "tags": ["backup", "database", "daily"],
  "timeout": 7200
}
```

### Response

```json
{
  "data": {
    "id": "job-123",
    "name": "backup-database-daily",
    "schedule": "0 0 * * *",
    "command": "pg_dump -U postgres mydb > /backups/mydb-$(date +%Y%m%d).sql",
    "status": "paused",
    "tags": ["backup", "database", "daily"],
    "timeout": 7200,
    "retries": 3,
    "retry_delay": 300,
    "environment": {
      "PGPASSWORD": "******"
    },
    "working_directory": "/opt/backups",
    "created_at": "2023-01-01T00:00:00Z",
    "updated_at": "2023-01-03T12:34:56Z",
    "next_run": null
  }
}
```

## Delete a Job

```
DELETE /api/v1/jobs/{id}
```

Deletes a job.

### Path Parameters

| Parameter | Type   | Description |
|-----------|--------|-------------|
| `id`      | string | Job ID      |

### Response

```json
{
  "data": {
    "message": "Job deleted successfully"
  }
}
```

## Run a Job Immediately

```
POST /api/v1/jobs/{id}/run
```

Triggers a job to run immediately, regardless of its schedule.

### Path Parameters

| Parameter | Type   | Description |
|-----------|--------|-------------|
| `id`      | string | Job ID      |

### Response

```json
{
  "data": {
    "message": "Job triggered successfully",
    "execution_id": "exec-789"
  }
}
```

## Pause a Job

```
POST /api/v1/jobs/{id}/pause
```

Pauses a job, preventing it from running according to its schedule.

### Path Parameters

| Parameter | Type   | Description |
|-----------|--------|-------------|
| `id`      | string | Job ID      |

### Response

```json
{
  "data": {
    "message": "Job paused successfully",
    "id": "job-123",
    "status": "paused"
  }
}
```

## Resume a Job

```
POST /api/v1/jobs/{id}/resume
```

Resumes a paused job, allowing it to run according to its schedule.

### Path Parameters

| Parameter | Type   | Description |
|-----------|--------|-------------|
| `id`      | string | Job ID      |

### Response

```json
{
  "data": {
    "message": "Job resumed successfully",
    "id": "job-123",
    "status": "active",
    "next_run": "2023-01-04T00:00:00Z"
  }
}
```

## Get Job Execution History

```
GET /api/v1/jobs/{id}/executions
```

Retrieves the execution history of a job.

### Path Parameters

| Parameter | Type   | Description |
|-----------|--------|-------------|
| `id`      | string | Job ID      |

### Query Parameters

| Parameter | Type   | Description                                                |
|-----------|--------|------------------------------------------------------------|
| `limit`   | int    | Maximum number of executions to return (default: 10)       |
| `offset`  | int    | Number of executions to skip (default: 0)                  |
| `status`  | string | Filter by execution status (success, failure, in_progress) |

### Response

```json
{
  "data": [
    {
      "id": "exec-123",
      "job_id": "job-123",
      "status": "success",
      "started_at": "2023-01-02T00:00:00Z",
      "finished_at": "2023-01-02T00:00:05Z",
      "duration": 5,
      "exit_code": 0,
      "output": "Database backup completed successfully",
      "error": null,
      "retry_count": 0
    },
    {
      "id": "exec-456",
      "job_id": "job-123",
      "status": "failure",
      "started_at": "2023-01-01T00:00:00Z",
      "finished_at": "2023-01-01T00:00:02Z",
      "duration": 2,
      "exit_code": 1,
      "output": "",
      "error": "pg_dump: could not connect to database",
      "retry_count": 3
    }
  ],
  "pagination": {
    "total": 2,
    "limit": 10,
    "offset": 0
  }
}
``` 