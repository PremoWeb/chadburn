# Advanced Topics

This document covers advanced features and concepts in Chadburn for users who want to get the most out of the system.

## Middlewares

Chadburn uses middlewares to extend job functionality. Middlewares are executed before and after job execution, allowing for notifications, logging, and other cross-cutting concerns.

### Built-in Middlewares

#### Slack Notifications

Sends job execution notifications to Slack:

```ini
[global]
slack-webhook = https://hooks.slack.com/services/XXX/YYY/ZZZ
slack-only-on-error = true
```

#### Email Notifications

Sends job execution notifications via email:

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

#### Execution Logging

Saves job execution logs to a directory:

```ini
[global]
save-folder = /var/log/chadburn
save-only-on-error = false
```

#### Gotify Notifications

Sends job execution notifications to a Gotify server:

```ini
[global]
gotify-url = https://gotify.example.com
gotify-token = YourGotifyToken
gotify-only-on-error = true
```

### Middleware Behavior

Middlewares are executed in the following order:

1. Global middlewares (defined in the `[global]` section)
2. Job-specific middlewares (if any)

Each middleware can:
- Perform actions before job execution
- Modify the execution context
- Perform actions after job execution
- Access execution results and logs

## Variable Substitution

Chadburn supports variable substitution in job commands, allowing for dynamic command generation.

### Available Variables

#### Container Variables

For `job-exec` and `job-run`:

- `{{.Container.Name}}`: Container name
- `{{.Container.ID}}`: Container ID
- `{{.Container.ImageName}}`: Container image name
- `{{.Container.State}}`: Container state

#### Time Variables

For all job types:

- `{{.Time.Now}}`: Current time (format: `2006-01-02T15:04:05Z07:00`)
- `{{.Time.NowUTC}}`: Current UTC time
- `{{.Time.NowFormat "2006-01-02"}}`: Current time with custom format

#### Environment Variables

For all job types:

- `{{.Env.VARIABLE_NAME}}`: Value of the environment variable

### Examples

```ini
[job-exec "backup-with-timestamp"]
schedule = @daily
container = app
command = tar -czf /backups/data-{{.Time.NowFormat "2006-01-02"}}.tar.gz /data

[job-exec "container-info"]
schedule = @hourly
container = app
command = echo "Running in {{.Container.Name}} ({{.Container.ID}})"

[job-local "env-example"]
schedule = @daily
command = echo "Database: {{.Env.DB_HOST}}"
```

## Metrics and Monitoring

Chadburn provides Prometheus-compatible metrics for monitoring job executions.

### Enabling Metrics

```bash
docker run -d --name chadburn \
  -v /var/run/docker.sock:/var/run/docker.sock:ro,z \
  -v /path/to/config.ini:/etc/chadburn.conf \
  -p 8080:8080 \
  premoweb/chadburn:latest daemon --metrics
```

### Available Metrics

- `chadburn_scheduler_jobs`: Active job count registered on the scheduler
- `chadburn_scheduler_register_errors_total`: Total number of failed scheduler registrations
- `chadburn_run_total`: Total number of completed job runs (labeled by job name)
- `chadburn_run_errors_total`: Total number of completed job runs that resulted in an error (labeled by job name)
- `chadburn_run_latest_timestamp`: Last time a job run completed (labeled by job name)
- `chadburn_run_duration_seconds`: Duration of all runs (labeled by job name)

### Prometheus Configuration

Example Prometheus configuration:

```yaml
scrape_configs:
  - job_name: 'chadburn'
    static_configs:
      - targets: ['chadburn:8080']
```

### Grafana Dashboard

You can create a Grafana dashboard to visualize Chadburn metrics:

- Job success/failure rates
- Job execution durations
- Job execution timestamps
- Error counts

## Job Execution Context

Each job execution has a context that contains information about the job, execution environment, and results.

### Context Properties

- `Job`: The job being executed
- `Execution`: Execution details (start time, end time, error, etc.)
- `Logger`: Logger instance for the job

### Execution Lifecycle

1. Job is scheduled based on its cron expression
2. Execution context is created
3. Pre-execution middlewares are run
4. Job command is executed
5. Post-execution middlewares are run
6. Execution results are stored

### Execution Timeouts

By default, job executions don't have a timeout. You can set a timeout for a job:

```ini
[job-exec "long-running-task"]
schedule = @daily
container = app
command = /path/to/long-script.sh
timeout = 1h
```

## Advanced Docker Integration

### Container Networks

For `job-run` and `job-service-run`, you can specify the network:

```ini
[job-run "network-example"]
schedule = @daily
image = ubuntu:latest
network = my-network
command = ping -c 3 other-container
```

### Volume Mounts

For `job-run` and `job-service-run`, you can mount volumes:

```ini
[job-run "volume-example"]
schedule = @daily
image = ubuntu:latest
volumes = ["/host/path:/container/path", "/tmp:/tmp"]
command = ls -la /container/path
```

### User Specification

For `job-exec`, `job-run`, and `job-service-run`, you can specify the user:

```ini
[job-exec "user-example"]
schedule = @daily
container = app
user = www-data
command = touch /var/www/file.txt
```

### Working Directory

For `job-exec`, `job-run`, and `job-service-run`, you can specify the working directory:

```ini
[job-exec "workdir-example"]
schedule = @daily
container = app
working-dir = /app
command = ./run-script.sh
```

## Hybrid Configuration Strategies

Chadburn supports both static (INI file) and dynamic (Docker labels) configuration. Here are some strategies for combining them:

### Core Jobs in INI, Dynamic Jobs in Labels

Use the INI file for critical jobs that should always run, and Docker labels for application-specific jobs that may change.

### Global Settings in INI, Jobs in Labels

Use the INI file for global settings (notifications, logging) and Docker labels for all job definitions.

### Development vs. Production

Use Docker labels in development for easy changes, and INI files in production for version-controlled configurations.

## Security Considerations

### Principle of Least Privilege

- Run Chadburn with read-only access to the Docker socket when possible
- Use specific users for job execution instead of root
- Limit the commands that can be executed

### Sensitive Information

- Avoid putting sensitive information in Docker labels (they're visible in `docker inspect`)
- Use environment variables or configuration files for sensitive data
- Consider using Docker secrets for credentials

### Network Isolation

- Run Chadburn in its own network namespace
- Only expose the metrics port if needed
- Use network policies to restrict access to the metrics endpoint 