# Chadburn Configuration Guide

This guide explains how to configure Chadburn for your specific needs, covering configuration file formats, options, and best practices.

## Configuration Methods

Chadburn supports two primary configuration methods:

1. **Configuration File**: Traditional INI-style configuration file
2. **Docker Labels**: Dynamic configuration through container labels

These methods can be used together, allowing for flexible deployment scenarios.

## Configuration File

### File Location

By default, Chadburn looks for a configuration file at:

- `/etc/chadburn.conf`

You can specify a different location using the `-c` or `--config` flag:

```bash
chadburn daemon -c /path/to/config.ini
```

### File Format

Chadburn uses an INI-style configuration format with sections and key-value pairs:

```ini
[section-type "name"]
key1 = value1
key2 = value2
```

### Global Configuration

The `[global]` section defines settings that apply to all jobs:

```ini
[global]
timezone = UTC
concurrency = 10
log_level = info
metrics = true
metrics_path = /metrics
metrics_address = :8080
```

| Option | Description | Default |
|--------|-------------|---------|
| `timezone` | Default timezone for all jobs | System timezone |
| `concurrency` | Maximum number of concurrent jobs | 0 (unlimited) |
| `log_level` | Logging verbosity (debug, info, warn, error) | info |
| `metrics` | Enable Prometheus metrics | false |
| `metrics_path` | Path for metrics endpoint | /metrics |
| `metrics_address` | Address:port for metrics server | :8080 |

### Job Configuration

Each job is defined in its own section, with the section type indicating the job type:

```ini
[job-local "backup-database"]
schedule = @daily
command = /usr/local/bin/backup.sh
```

The general format is:

```ini
[job-type "job-name"]
option1 = value1
option2 = value2
```

## Job Types and Options

### job-local

Executes commands on the host machine:

```ini
[job-local "system-cleanup"]
schedule = @daily
command = /usr/local/bin/cleanup.sh
user = admin
```

| Option | Description | Required |
|--------|-------------|----------|
| `schedule` | When to run the job (cron format) | Yes |
| `command` | Command to execute | Yes |
| `user` | User to run as | No |
| `working_directory` | Directory to run in | No |
| `environment` | Environment variables | No |
| `timeout` | Maximum execution time | No |

### job-exec

Executes commands in an existing container:

```ini
[job-exec "cache-flush"]
schedule = @hourly
container = redis
command = redis-cli FLUSHALL
```

| Option | Description | Required |
|--------|-------------|----------|
| `schedule` | When to run the job | Yes |
| `container` | Container name or ID | Yes |
| `command` | Command to execute | Yes |
| `user` | User to run as | No |
| `environment` | Environment variables | No |
| `timeout` | Maximum execution time | No |

### job-run

Creates a new container to execute commands:

```ini
[job-run "database-backup"]
schedule = @daily
image = postgres:13
command = pg_dump -U postgres -d mydb > /backup/mydb.sql
volumes = ["/backup:/backup"]
```

| Option | Description | Required |
|--------|-------------|----------|
| `schedule` | When to run the job | Yes |
| `image` | Docker image to use | Yes |
| `command` | Command to execute | Yes |
| `volumes` | Volumes to mount | No |
| `networks` | Networks to connect to | No |
| `environment` | Environment variables | No |
| `user` | User to run as | No |
| `timeout` | Maximum execution time | No |

### job-service-run

Creates a service to execute commands (for Docker Swarm):

```ini
[job-service-run "batch-process"]
schedule = @daily
image = data-processor:latest
command = python /app/process.py
service_name = batch-processor
replicas = 3
```

| Option | Description | Required |
|--------|-------------|----------|
| `schedule` | When to run the job | Yes |
| `image` | Docker image to use | Yes |
| `command` | Command to execute | Yes |
| `service_name` | Name for the service | Yes |
| `replicas` | Number of replicas | No |
| `networks` | Networks to connect to | No |
| `environment` | Environment variables | No |
| `constraints` | Placement constraints | No |

### job-lifecycle

Executes commands on container lifecycle events:

```ini
[job-lifecycle "log-container-start"]
container = web-app
event = start
command = echo "Container started at $(date)" >> /var/log/container-events.log
```

| Option | Description | Required |
|--------|-------------|----------|
| `container` | Container name or ID | Yes |
| `event` | Event type (start, stop, die, health_status) | Yes |
| `command` | Command to execute | Yes |
| `user` | User to run as | No |

## Docker Labels Configuration

You can configure Chadburn using Docker labels on containers:

```bash
docker run -d --name web-app \
  --label chadburn.enabled=true \
  --label chadburn.job-exec.healthcheck.schedule="@every 5m" \
  --label chadburn.job-exec.healthcheck.command="curl -s http://localhost/health" \
  nginx
```

The label format is:

```
chadburn.job-type.job-name.option=value
```

For example:
- `chadburn.enabled=true`: Enable Chadburn for this container
- `chadburn.job-exec.healthcheck.schedule="@every 5m"`: Schedule for the job
- `chadburn.job-exec.healthcheck.command="curl -s http://localhost/health"`: Command to run

## Advanced Configuration

### Environment Variables

You can use environment variables in your configuration:

```ini
[job-local "backup"]
schedule = @daily
command = /usr/local/bin/backup.sh
environment = ["DB_HOST=${DB_HOST}", "DB_PASSWORD=${DB_PASSWORD}"]
```

### Templating

Chadburn supports basic templating in commands:

```ini
[job-local "dated-backup"]
schedule = @daily
command = /usr/local/bin/backup.sh --date={{ now.Format "2006-01-02" }}
```

Available template variables:
- `now`: Current time
- `job`: Job information

### Includes

You can split your configuration across multiple files:

```ini
[include]
files = ["/etc/chadburn.d/*.conf"]
```

This will include all `.conf` files in the `/etc/chadburn.d/` directory.

## Best Practices

### Organization

1. **Use Descriptive Names**: Choose clear, descriptive job names
2. **Group Related Jobs**: Use naming conventions to group related jobs
3. **Split Configuration**: Use includes for large configurations

### Security

1. **Least Privilege**: Run jobs with the minimum required permissions
2. **Avoid Sensitive Data**: Don't store secrets directly in the configuration
3. **Use Environment Variables**: Pass sensitive data via environment variables

### Maintenance

1. **Version Control**: Keep your configuration in version control
2. **Documentation**: Add comments to explain complex jobs
3. **Validation**: Use `chadburn validate` to check your configuration

## Example Configurations

### Basic Configuration

```ini
[global]
timezone = UTC
log_level = info

[job-local "cleanup"]
schedule = @daily
command = find /tmp -type f -mtime +7 -delete

[job-exec "redis-backup"]
schedule = 0 0 * * * *
container = redis
command = redis-cli SAVE
```

### Advanced Configuration

```ini
[global]
timezone = America/New_York
concurrency = 5
metrics = true
metrics_address = :9090

[job-local "system-stats"]
schedule = @every 5m
command = /usr/local/bin/collect-stats.sh
timeout = 30s
environment = ["OUTPUT_DIR=/var/stats"]

[job-run "database-backup"]
schedule = 0 0 1 * * *  # 1:00 AM
image = postgres:13
command = pg_dump -U postgres -d mydb > /backup/mydb-{{ now.Format "2006-01-02" }}.sql
volumes = ["/backup:/backup"]
environment = ["PGPASSWORD=${DB_PASSWORD}"]
timeout = 1h

[job-lifecycle "log-container-events"]
container = web-app
event = health_status
command = /usr/local/bin/log-health.sh {{ .container.name }} {{ .container.health }}
```

## Next Steps

- Learn about [Job Types](/docs/concepts/jobs) in detail
- Understand [Schedules](/docs/concepts/schedules) to control when jobs run
- Explore [Metrics](/docs/concepts/metrics) to monitor your jobs
- Set up [Monitoring](/docs/guides/monitoring) for your Chadburn instance 