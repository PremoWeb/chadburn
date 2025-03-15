# Frequently Asked Questions

## General Questions

### What is Chadburn?

Chadburn is a lightweight job scheduler designed for Docker environments. It serves as a modern replacement for traditional cron, with enhanced features for container orchestration.

### Why use Chadburn instead of cron?

Chadburn offers several advantages over traditional cron:
- Native Docker integration
- Dynamic configuration through Docker labels
- Ability to execute commands inside containers
- Real-time updates without restarts
- Container lifecycle event handling
- Better logging and notification options
- Metrics and monitoring capabilities

### Is Chadburn production-ready?

Yes, Chadburn is used in production environments. It's built on the robust Go implementation of cron and uses the official Docker client library.

### How does Chadburn compare to other schedulers?

- **vs. Cron**: Chadburn adds Docker integration and dynamic configuration
- **vs. Ofelia**: Chadburn is a fork of Ofelia with additional features and improvements
- **vs. Kubernetes CronJobs**: Chadburn works in any Docker environment, not just Kubernetes

## Configuration Questions

### Can I use both INI files and Docker labels?

Yes, Chadburn supports hybrid configuration. You can define global settings and some jobs in the INI file, while using Docker labels for dynamic job configuration.

### How do I specify the schedule?

Chadburn uses the Go implementation of cron for scheduling. Examples include:
- `@every 1h30m`: Run every hour and thirty minutes
- `0 30 * * * *`: Run at 30 minutes past every hour
- `@hourly`: Run once an hour
- `@daily`: Run once a day

### Can I use environment variables in commands?

Yes, you can use environment variables in your commands:
```ini
command = /bin/bash -c "export DB_HOST=${DB_HOST:-localhost} && /app/backup-db.sh"
```

You can also use Chadburn's variable substitution:
```ini
command = echo "Running in container {{.Container.Name}}"
```

### How do I validate my configuration?

You can validate your configuration without starting the daemon:
```bash
chadburn validate --config=/path/to/config.ini
```

## Docker Integration Questions

### How does Chadburn access Docker?

Chadburn accesses Docker through the Docker socket (`/var/run/docker.sock`). When running Chadburn in a container, you need to mount this socket:
```bash
-v /var/run/docker.sock:/var/run/docker.sock:ro,z
```

### I'm getting "permission denied" errors with the Docker socket. How do I fix this?

See the [Docker Socket Permissions Guide](docker-socket-permissions.md) for detailed solutions. Common fixes include:
1. Adding the `:z` suffix to the volume mount
2. Running the container with the correct Docker group ID
3. Running the container as root

### Can Chadburn work with Docker Swarm?

Yes, Chadburn works well with Docker Swarm. It includes a special job type, `job-service-run`, designed specifically for Swarm environments.

### Can Chadburn work with Kubernetes?

Chadburn can work in Kubernetes, but with some limitations:
1. `job-exec` can only access containers on the same node
2. `job-local` runs on the node, not in the Kubernetes context
3. `job-service-run` is not applicable (use Kubernetes Jobs instead)

## Job Types Questions

### What job types does Chadburn support?

Chadburn supports five job types:
- `job-local`: Executes commands on the host machine
- `job-exec`: Executes commands inside a running container
- `job-run`: Creates a new container to execute commands
- `job-service-run`: Creates a service to execute commands (for Swarm)
- `job-lifecycle`: Executes commands on container lifecycle events

### How do I run a command inside a container?

Use the `job-exec` type:
```ini
[job-exec "my-job"]
schedule = @hourly
container = my-container
command = echo "Hello from inside the container"
```

### How do I start a container on a schedule?

Use the `job-run` type:
```ini
[job-run "my-job"]
schedule = @daily
image = ubuntu:latest
command = echo "Hello from a new container"
```

### How do I execute a command when a container starts or stops?

Use the `job-lifecycle` type:
```ini
[job-lifecycle "notify-on-start"]
container = my-container
event-type = start
command = echo "Container started"
```

## Operational Questions

### How do I check if Chadburn is running correctly?

Check the logs:
```bash
docker logs chadburn
```

If you've enabled metrics, you can also check the metrics endpoint:
```bash
curl http://localhost:8080/metrics
```

### How do I debug issues with Chadburn?

Enable debug logging:
```bash
docker run -d --name chadburn \
  -e CHADBURN_LOG_LEVEL=debug \
  -v /var/run/docker.sock:/var/run/docker.sock:ro,z \
  -v /path/to/config.ini:/etc/chadburn.conf \
  premoweb/chadburn:latest daemon
```

### How do I upgrade Chadburn?

Pull the latest image and restart:
```bash
docker pull premoweb/chadburn:latest
docker-compose up -d
```

### Can I run multiple instances of Chadburn?

Yes, but you should ensure they don't have overlapping job configurations to avoid duplicate job executions.

## Advanced Questions

### How do I set up notifications for job failures?

Chadburn supports several notification methods:

**Slack**:
```ini
[global]
slack-webhook = https://hooks.slack.com/services/XXX/YYY/ZZZ
slack-only-on-error = true
```

**Email**:
```ini
[global]
smtp-host = smtp.example.com
smtp-port = 587
smtp-user = user
smtp-password = password
email-to = alerts@example.com
email-from = chadburn@example.com
mail-only-on-error = true
```

**Gotify**:
```ini
[global]
gotify-url = https://gotify.example.com
gotify-token = YourGotifyToken
gotify-only-on-error = true
```

### How do I set a timeout for long-running jobs?

Add a timeout parameter to your job:
```ini
[job-exec "long-running-job"]
schedule = @daily
container = app
command = /app/long-process.sh
timeout = 1h
```

### Can I use Chadburn for high-availability setups?

Yes, but you need to consider:
1. Running Chadburn on a manager node in Swarm mode
2. Using a shared volume for configuration and logs
3. Setting up monitoring for the Chadburn container

### How do I contribute to Chadburn?

Contributions are welcome! Check out the [GitHub repository](https://github.com/PremoWeb/Chadburn) and the [Contributing Guide](contributing.md) for more information. 