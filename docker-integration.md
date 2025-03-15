# Docker Integration

Chadburn is designed to work seamlessly with Docker, providing advanced scheduling capabilities for containerized environments. This document explains how Chadburn integrates with Docker and how to use its Docker-specific features.

## Docker Socket Access

Chadburn needs access to the Docker socket (`/var/run/docker.sock`) to interact with Docker containers. When running Chadburn in a container, you need to mount the Docker socket:

```bash
docker run -d --name chadburn \
  -v /var/run/docker.sock:/var/run/docker.sock:ro,z \
  -v /path/to/config.ini:/etc/chadburn.conf \
  premoweb/chadburn:latest daemon
```

> **Note**: If you encounter permission issues with the Docker socket, refer to our [Docker Socket Permissions Guide](docker-socket-permissions.md) for solutions.

## Docker Labels

Chadburn can be configured using Docker labels, allowing for dynamic configuration without modifying the Chadburn configuration file.

### Label Format

The label format is:

```
chadburn.<JOB_TYPE>.<JOB_NAME>.<JOB_PARAMETER>=<PARAMETER_VALUE>
```

### Required Labels

For `job-exec` jobs, the target container must have the label `chadburn.enabled=true`.

### Example: Running Commands in an Existing Container

```bash
docker run -d --name nginx \
  --label chadburn.enabled=true \
  --label chadburn.job-exec.print-date.schedule="@every 5m" \
  --label chadburn.job-exec.print-date.command="date" \
  nginx
```

### Example: Starting a Container on Schedule

```bash
docker run -d --name backup-container \
  --label chadburn.job-run.schedule="@daily" \
  -v /data:/data \
  backup-image
```

This will create a job that starts the `backup-container` container once a day.

### Example: Docker Compose

```yaml
version: "3"
services:
  chadburn:
    image: premoweb/chadburn:latest
    command: daemon
    volumes:
      - /var/run/docker.sock:/var/run/docker.sock:ro,z
    depends_on:
      - app

  app:
    image: myapp:latest
    labels:
      chadburn.enabled: "true"
      chadburn.job-exec.backup.schedule: "@daily"
      chadburn.job-exec.backup.command: "tar -czf /backup/data.tar.gz /app/data"
```

## Container Lifecycle Events

Chadburn can execute commands when containers start or stop using the `job-lifecycle` type.

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

### Using Docker Labels

```bash
docker run -d --name app \
  --label chadburn.job-lifecycle.notify.event-type="start" \
  --label chadburn.job-lifecycle.notify.command="echo 'Container started'" \
  --label chadburn.job-lifecycle.cleanup.event-type="stop" \
  --label chadburn.job-lifecycle.cleanup.command="echo 'Container stopped'" \
  myapp:latest
```

## Dynamic Configuration

One of Chadburn's key features is its ability to detect changes in Docker containers and update its configuration dynamically.

### Container Changes

When containers are started, stopped, or modified, Chadburn automatically:

1. Detects new containers with Chadburn labels
2. Updates existing job configurations if labels change
3. Removes jobs for containers that are stopped or removed

### No Restart Required

Unlike traditional cron, Chadburn doesn't require a restart when configuration changes. It continuously monitors Docker events and updates its job schedule accordingly.

### Example Workflow

1. Start Chadburn:
   ```bash
   docker run -d --name chadburn -v /var/run/docker.sock:/var/run/docker.sock:ro,z premoweb/chadburn:latest daemon
   ```

2. Start a container with Chadburn labels:
   ```bash
   docker run -d --name app1 \
     --label chadburn.enabled=true \
     --label chadburn.job-exec.task1.schedule="@every 5m" \
     --label chadburn.job-exec.task1.command="echo 'Task 1'" \
     myapp:latest
   ```

3. Chadburn automatically detects the new container and schedules the job.

4. Start another container:
   ```bash
   docker run -d --name app2 \
     --label chadburn.enabled=true \
     --label chadburn.job-exec.task2.schedule="@hourly" \
     --label chadburn.job-exec.task2.command="echo 'Task 2'" \
     myapp:latest
   ```

5. Chadburn automatically detects the new container and schedules the job.

6. Stop a container:
   ```bash
   docker stop app1
   ```

7. Chadburn automatically removes the jobs for the stopped container.

## Variable Substitution

Chadburn supports variable substitution in job commands, allowing you to reference container information dynamically:

```bash
docker run -d --name app \
  --label chadburn.enabled=true \
  --label chadburn.job-exec.log-info.schedule="@hourly" \
  --label chadburn.job-exec.log-info.command="echo 'Container {{.Container.Name}} ({{.Container.ID}}) is running'" \
  myapp:latest
```

Available variables:

- `{{.Container.Name}}`: Container name
- `{{.Container.ID}}`: Container ID
- `{{.Container.ImageName}}`: Container image name
- `{{.Container.State}}`: Container state
- `{{.Time.Now}}`: Current time

## Official Docker Client

Chadburn uses the official Docker client library for Docker integration, providing improved compatibility with the latest Docker features and better long-term maintainability. 