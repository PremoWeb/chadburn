# Job Types

Chadburn supports several job types, each designed for different use cases. This document explains the different job types and how to configure them.

## job-exec

The `job-exec` job type executes a command inside a running container.

### Configuration

```ini
[job-exec "job-name"]
schedule = @hourly
container = my-container
command = touch /tmp/example
```

### Parameters

- `schedule`: When to run the job (cron expression or interval)
- `container`: The name of the container to execute the command in
- `command`: The command to execute
- `user`: (Optional) The user to run the command as
- `working-dir`: (Optional) The working directory for the command
- `timeout`: (Optional) Maximum time the command can run

### Docker Labels

```bash
docker run -d --name my-container \
  --label chadburn.enabled=true \
  --label chadburn.job-exec.backup.schedule="@daily" \
  --label chadburn.job-exec.backup.command="tar -czf /backup/data.tar.gz /data" \
  my-image
```

## job-run

The `job-run` job type creates a new container to execute a command.

### Configuration

```ini
[job-run "job-name"]
schedule = @hourly
image = ubuntu:latest
command = touch /tmp/example
```

### Parameters

- `schedule`: When to run the job (cron expression or interval)
- `image`: The Docker image to use
- `command`: The command to execute
- `user`: (Optional) The user to run the command as
- `working-dir`: (Optional) The working directory for the command
- `environment`: (Optional) Environment variables as a list of "KEY=VALUE" strings
- `volumes`: (Optional) Volumes to mount as a list of "HOST:CONTAINER" strings
- `network`: (Optional) The network to connect the container to
- `timeout`: (Optional) Maximum time the command can run

### Docker Labels

```bash
docker run -d --name scheduler \
  --label chadburn.job-run.backup.schedule="@daily" \
  --label chadburn.job-run.backup.image="ubuntu:latest" \
  --label chadburn.job-run.backup.command="tar -czf /backup/data.tar.gz /data" \
  --label chadburn.job-run.backup.volumes="['/data:/data', '/backup:/backup']" \
  my-scheduler-image
```

## job-local

The `job-local` job type executes a command on the host running Chadburn.

### Configuration

```ini
[job-local "job-name"]
schedule = @hourly
command = touch /tmp/example
```

### Parameters

- `schedule`: When to run the job (cron expression or interval)
- `command`: The command to execute
- `timeout`: (Optional) Maximum time the command can run

### Docker Labels

```bash
docker run -d --name scheduler \
  --label chadburn.job-local.cleanup.schedule="@daily" \
  --label chadburn.job-local.cleanup.command="find /tmp -type f -mtime +7 -delete" \
  my-scheduler-image
```

## job-service-run

The `job-service-run` job type creates a service to execute commands (for Swarm environments).

### Configuration

```ini
[job-service-run "job-name"]
schedule = 0,20,40 * * * *
image = ubuntu
network = swarm_network
command = touch /tmp/example
```

### Parameters

- `schedule`: When to run the job (cron expression or interval)
- `image`: The Docker image to use
- `command`: The command to execute
- `user`: (Optional) The user to run the command as
- `working-dir`: (Optional) The working directory for the command
- `environment`: (Optional) Environment variables as a list of "KEY=VALUE" strings
- `volumes`: (Optional) Volumes to mount as a list of "HOST:CONTAINER" strings
- `network`: (Optional) The network to connect the service to
- `timeout`: (Optional) Maximum time the command can run

### Docker Labels

```bash
docker run -d --name scheduler \
  --label chadburn.job-service-run.backup.schedule="@daily" \
  --label chadburn.job-service-run.backup.image="ubuntu:latest" \
  --label chadburn.job-service-run.backup.command="tar -czf /backup/data.tar.gz /data" \
  --label chadburn.job-service-run.backup.network="swarm_network" \
  my-scheduler-image
```

## job-lifecycle

The `job-lifecycle` job type executes commands when a container starts or stops.

### Configuration

```ini
[job-lifecycle "notify-on-start"]
container = my-container
event-type = start
command = echo "Container started" | mail -s "Container Event" admin@example.com

[job-lifecycle "cleanup-on-stop"]
container = my-container
event-type = stop
command = rm -rf /tmp/cache
```

### Parameters

- `container`: The name of the container to monitor
- `event-type`: The event type to trigger on (`start` or `stop`)
- `command`: The command to execute
- `timeout`: (Optional) Maximum time the command can run

### Docker Labels

```bash
docker run -d --name my-container \
  --label chadburn.job-lifecycle.notify.event-type="start" \
  --label chadburn.job-lifecycle.notify.command="echo 'Container started'" \
  --label chadburn.job-lifecycle.cleanup.event-type="stop" \
  --label chadburn.job-lifecycle.cleanup.command="echo 'Container stopped'" \
  my-image
```

## Scheduling Format

Chadburn uses the Go implementation of cron for scheduling. Examples include:

- `@every 1h30m`: Run every hour and thirty minutes
- `0 30 * * * *`: Run at 30 minutes past every hour
- `@hourly`: Run once an hour
- `@daily`: Run once a day
- `@midnight`: Run once a day at midnight
- `@weekly`: Run once a week 