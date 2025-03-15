# Chadburn Configuration

This document describes the various ways to configure Chadburn, including the configuration file format, environment variables, and command line options.

## Configuration File Format

Chadburn uses the INI file format for its configuration. The default location is `/etc/chadburn.conf`, but you can specify a different location using the `--config` flag.

### Global Section

The `[global]` section contains settings that apply to all jobs:

```ini
[global]
# Slack integration
slack-webhook = https://hooks.slack.com/services/XXX/YYY/ZZZ
slack-only-on-error = true

# Email notifications
smtp-host = smtp.example.com
smtp-port = 587
smtp-user = user
smtp-password = password
email-to = alerts@example.com
email-from = chadburn@example.com
mail-only-on-error = true

# Log saving
save-folder = /var/log/chadburn
save-only-on-error = false

# Gotify notifications
gotify-url = https://gotify.example.com
gotify-token = YourGotifyToken
gotify-only-on-error = true
```

### Job Sections

Each job type has its own section format:

#### job-exec

Executes a command inside a running container:

```ini
[job-exec "job-name"]
schedule = @hourly
container = my-container
command = touch /tmp/example
```

#### job-run

Runs a command in a new container:

```ini
[job-run "job-name"]
schedule = @hourly
image = ubuntu:latest
command = touch /tmp/example
```

#### job-local

Executes a command on the host running Chadburn:

```ini
[job-local "job-name"]
schedule = @hourly
command = touch /tmp/example
```

#### job-service-run

Runs a command inside a new "run-once" service (for Swarm environments):

```ini
[job-service-run "job-name"]
schedule = 0,20,40 * * * *
image = ubuntu
network = swarm_network
command = touch /tmp/example
```

#### job-lifecycle

Executes commands when a container starts or stops:

```ini
[job-lifecycle "notify-on-start"]
container = my-container
event-type = start
command = echo "Container started" | mail -s "Container Event" admin@example.com
```

## Environment Variables

Chadburn supports the following environment variables:

| Variable | Description | Default |
|----------|-------------|---------|
| `CHADBURN_LOG_LEVEL` | Log level (debug, info, notice, warning, error, critical) | `info` |
| `DOCKER_GID` | Docker group ID for socket access | `999` |

## Command Line Options

Chadburn supports the following command line options:

### Global Options

```
--help, -h             Show help
--version, -v          Show version
```

### Daemon Command

```
chadburn daemon [options]
```

Options:

```
--config=FILE          Configuration file (default: /etc/chadburn.conf)
--metrics              Enable Prometheus compatible metrics endpoint
--listen-address=ADDR  Metrics endpoint listen address (default: :8080)
--disable-docker       Disable docker integration (only job-local will work)
```

### Validate Command

```
chadburn validate [options]
```

Options:

```
--config=FILE          Configuration file to validate (default: /etc/chadburn.conf)
```

## Docker Labels Configuration

When using Docker labels, the format is:

```
chadburn.<JOB_TYPE>.<JOB_NAME>.<JOB_PARAMETER>=<PARAMETER_VALUE>
```

For example:

```
chadburn.enabled=true
chadburn.job-exec.backup.schedule=@daily
chadburn.job-exec.backup.command=tar -czf /backup/data.tar.gz /data
```

The `chadburn.enabled=true` label is required for containers that will have `job-exec` tasks executed on them.

## Configuration Precedence

When both INI file and Docker labels are used, they are merged with the following rules:

1. Jobs from the INI file are loaded first
2. Jobs from Docker labels are added or update existing jobs
3. If a job from Docker labels has the same name as a job from the INI file, the Docker label version takes precedence

## Hybrid Configuration

You can use both INI files and Docker labels together:

```ini
# config.ini
[global]
slack-webhook = https://hooks.slack.com/services/XXX/YYY/ZZZ

[job-local "system-backup"]
schedule = @daily
command = tar -czf /backup/system.tar.gz /etc
```

Combined with Docker labels:

```
chadburn.enabled=true
chadburn.job-exec.app-backup.schedule=@daily
chadburn.job-exec.app-backup.command=tar -czf /backup/app.tar.gz /app/data
```

This approach allows for a mix of static and dynamic configuration. 