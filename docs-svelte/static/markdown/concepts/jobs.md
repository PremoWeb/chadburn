# Jobs in Chadburn

Jobs are the core building blocks of Chadburn. Each job represents a task that needs to be executed according to a schedule or in response to an event.

## Job Types

Chadburn supports several job types to accommodate different execution environments:

### job-local

Executes commands directly on the host machine where Chadburn is running.

```ini
[job-local "backup-database"]
schedule = @daily
command = /usr/local/bin/backup-script.sh
```

### job-exec

Executes commands inside an already running container.

```ini
[job-exec "update-cache"]
schedule = @every 1h
container = redis
command = redis-cli FLUSHALL
```

### job-run

Creates a new container to execute commands and removes it after completion.

```ini
[job-run "send-newsletter"]
schedule = @weekly
image = newsletter-sender:latest
command = node /app/send.js
```

### job-service-run

Creates a service to execute commands (useful in Docker Swarm environments).

```ini
[job-service-run "batch-process"]
schedule = @daily
image = data-processor:latest
command = python /app/process.py
service_name = batch-processor
```

### job-lifecycle

Executes commands in response to container lifecycle events (start/stop).

```ini
[job-lifecycle "log-container-events"]
container = web-app
event = start
command = echo "Container started at $(date)" >> /var/log/container-events.log
```

## Job Configuration

Each job type has its own set of configuration options, but they all share some common parameters:

### Common Parameters

| Parameter | Description | Example |
|-----------|-------------|---------|
| `schedule` | When the job should run (cron format) | `@daily` or `0 0 * * *` |
| `command` | The command to execute | `echo "Hello"` |
| `environment` | Environment variables for the job | `["KEY=value"]` |
| `user` | User to run the command as | `www-data` |
| `disabled` | Whether the job is disabled | `true` or `false` |

### Type-Specific Parameters

Different job types have additional parameters:

#### job-run and job-service-run

| Parameter | Description | Example |
|-----------|-------------|---------|
| `image` | Docker image to use | `alpine:latest` |
| `volumes` | Volumes to mount | `["/host:/container"]` |
| `networks` | Networks to connect to | `["backend"]` |

#### job-exec and job-lifecycle

| Parameter | Description | Example |
|-----------|-------------|---------|
| `container` | Target container name or ID | `web-server` |

## Job Execution

When a job is triggered (by schedule or event), Chadburn:

1. Prepares the execution environment based on the job type
2. Runs the specified command
3. Captures the output and exit code
4. Logs the results
5. Cleans up resources (for job-run)

## Best Practices

- Use descriptive job names that indicate their purpose
- Keep commands simple and idempotent when possible
- Use environment variables for configuration
- Set appropriate timeouts for long-running jobs
- Use labels for organizing and filtering jobs

## Examples

### Database Backup

```ini
[job-run "db-backup"]
schedule = @daily
image = postgres:13
volumes = ["/var/backups:/backups"]
environment = ["PGPASSWORD=secret"]
command = pg_dump -h db.example.com -U postgres -d mydb > /backups/mydb-$(date +%Y%m%d).sql
```

### Cleanup Job

```ini
[job-local "cleanup-temp"]
schedule = @hourly
command = find /tmp -type f -mtime +1 -delete
```

### Health Check

```ini
[job-exec "health-check"]
schedule = @every 5m
container = web-app
command = curl -s http://localhost/health
```

## Next Steps

- Learn about [Schedules](/docs/concepts/schedules) to control when jobs run
- Understand [Execution](/docs/concepts/execution) details and options
- Explore [Metrics](/docs/concepts/metrics) to monitor your jobs 