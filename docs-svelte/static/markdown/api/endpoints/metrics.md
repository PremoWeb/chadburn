# Metrics API Endpoints

The Metrics API allows you to retrieve performance and operational metrics from Chadburn. These metrics can be used for monitoring, alerting, and performance analysis.

## Get Prometheus Metrics

```
GET /api/v1/metrics
```

Retrieves metrics in Prometheus format. This endpoint is compatible with Prometheus scraping.

### Response

```
# HELP chadburn_jobs_total Total number of jobs
# TYPE chadburn_jobs_total gauge
chadburn_jobs_total 42

# HELP chadburn_job_executions_total Total number of job executions
# TYPE chadburn_job_executions_total counter
chadburn_job_executions_total{status="success"} 156
chadburn_job_executions_total{status="failure"} 12

# HELP chadburn_job_execution_duration_seconds Duration of job executions in seconds
# TYPE chadburn_job_execution_duration_seconds histogram
chadburn_job_execution_duration_seconds_bucket{le="0.1"} 0
chadburn_job_execution_duration_seconds_bucket{le="0.5"} 12
chadburn_job_execution_duration_seconds_bucket{le="1"} 45
chadburn_job_execution_duration_seconds_bucket{le="5"} 134
chadburn_job_execution_duration_seconds_bucket{le="10"} 145
chadburn_job_execution_duration_seconds_bucket{le="30"} 156
chadburn_job_execution_duration_seconds_bucket{le="60"} 168
chadburn_job_execution_duration_seconds_bucket{le="+Inf"} 168
chadburn_job_execution_duration_seconds_sum 842.5
chadburn_job_execution_duration_seconds_count 168

# HELP chadburn_job_last_execution_timestamp_seconds Timestamp of the last execution of a job
# TYPE chadburn_job_last_execution_timestamp_seconds gauge
chadburn_job_last_execution_timestamp_seconds{job="backup-database",status="success"} 1672531200
chadburn_job_last_execution_timestamp_seconds{job="cleanup-logs",status="success"} 1672617600
```

## Get Metrics Summary

```
GET /api/v1/metrics/summary
```

Retrieves a summary of metrics in JSON format.

### Response

```json
{
  "data": {
    "jobs": {
      "total": 42,
      "active": 38,
      "paused": 4
    },
    "executions": {
      "total": 168,
      "success": 156,
      "failure": 12,
      "success_rate": 92.86
    },
    "performance": {
      "average_duration": 5.01,
      "p50_duration": 2.3,
      "p95_duration": 12.7,
      "p99_duration": 28.4
    },
    "system": {
      "uptime": 604800,
      "memory_usage_mb": 128.5,
      "cpu_usage_percent": 2.3
    }
  }
}
```

## Get Job-Specific Metrics

```
GET /api/v1/metrics/jobs/{id}
```

Retrieves metrics for a specific job.

### Path Parameters

| Parameter | Type   | Description |
|-----------|--------|-------------|
| `id`      | string | Job ID      |

### Response

```json
{
  "data": {
    "job_id": "job-123",
    "job_name": "backup-database",
    "executions": {
      "total": 30,
      "success": 28,
      "failure": 2,
      "success_rate": 93.33
    },
    "performance": {
      "average_duration": 8.5,
      "p50_duration": 7.2,
      "p95_duration": 15.3,
      "p99_duration": 22.1,
      "min_duration": 3.2,
      "max_duration": 25.7
    },
    "history": {
      "last_execution": "2023-01-02T00:00:00Z",
      "last_success": "2023-01-02T00:00:00Z",
      "last_failure": "2022-12-25T00:00:00Z",
      "next_scheduled": "2023-01-03T00:00:00Z"
    }
  }
}
```

## Get System Health Metrics

```
GET /api/v1/metrics/health
```

Retrieves system health metrics.

### Response

```json
{
  "data": {
    "status": "healthy",
    "uptime": 604800,
    "version": "1.9.0",
    "build_date": "2025-03-16T02:31:18Z",
    "resources": {
      "memory_usage_mb": 128.5,
      "memory_total_mb": 1024,
      "cpu_usage_percent": 2.3,
      "disk_usage_percent": 45.7
    },
    "components": {
      "scheduler": {
        "status": "healthy",
        "last_tick": "2023-01-02T12:34:56Z"
      },
      "database": {
        "status": "healthy",
        "connections": 5,
        "max_connections": 100
      },
      "api": {
        "status": "healthy",
        "requests_per_minute": 12.5
      }
    }
  }
}
```

## Get Metrics for a Time Period

```
GET /api/v1/metrics/history
```

Retrieves historical metrics for a specified time period.

### Query Parameters

| Parameter   | Type   | Description                                                |
|-------------|--------|------------------------------------------------------------|
| `start`     | string | Start time in ISO 8601 format (default: 24 hours ago)      |
| `end`       | string | End time in ISO 8601 format (default: now)                 |
| `interval`  | string | Aggregation interval (hourly, daily, weekly, default: hourly) |
| `job_id`    | string | Filter by job ID (optional)                                |

### Response

```json
{
  "data": {
    "start": "2023-01-01T00:00:00Z",
    "end": "2023-01-02T00:00:00Z",
    "interval": "hourly",
    "metrics": [
      {
        "timestamp": "2023-01-01T00:00:00Z",
        "executions": {
          "total": 5,
          "success": 5,
          "failure": 0
        },
        "average_duration": 4.2
      },
      {
        "timestamp": "2023-01-01T01:00:00Z",
        "executions": {
          "total": 3,
          "success": 2,
          "failure": 1
        },
        "average_duration": 6.7
      },
      // ... more intervals ...
      {
        "timestamp": "2023-01-01T23:00:00Z",
        "executions": {
          "total": 7,
          "success": 7,
          "failure": 0
        },
        "average_duration": 3.9
      }
    ]
  }
}
``` 