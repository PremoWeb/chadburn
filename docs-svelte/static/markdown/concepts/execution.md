# Job Execution in Chadburn

This document explains how Chadburn executes jobs, including execution environments, output handling, error management, and advanced execution options.

## Execution Flow

When a job is triggered (by schedule or event), Chadburn follows this execution flow:

1. **Preparation**: Set up the execution environment based on job type
2. **Execution**: Run the specified command
3. **Monitoring**: Track execution status and collect output
4. **Completion**: Record results and handle success/failure
5. **Cleanup**: Release resources (especially for container-based jobs)

## Execution Environments

Each job type has a different execution environment:

### job-local

Executes directly on the host system where Chadburn is running.

- **Working Directory**: Chadburn's current working directory
- **Environment Variables**: Inherits Chadburn's environment plus any specified in the job
- **User**: By default, runs as the same user as Chadburn

### job-exec

Executes inside an existing container.

- **Working Directory**: Container's working directory
- **Environment Variables**: Container's environment plus job-specific variables
- **User**: Container's default user (can be overridden)

### job-run

Creates a new container for execution.

- **Working Directory**: Container image's default working directory
- **Environment Variables**: Container image's defaults plus job-specific variables
- **User**: Container image's default user (can be overridden)
- **Networking**: Inherits Docker's default network settings

### job-service-run

Creates a service in Docker Swarm.

- **Working Directory**: Container image's default working directory
- **Environment Variables**: Container image's defaults plus job-specific variables
- **User**: Container image's default user (can be overridden)
- **Networking**: Uses Docker Swarm's overlay network

## Command Execution

### Shell Execution

By default, commands are executed through a shell:

- For Linux: `/bin/sh -c "your command"`
- For Windows: `cmd /C "your command"`

This allows for shell features like:
- Environment variable expansion
- I/O redirection
- Command chaining
- Wildcards

### Direct Execution

For some job types, you can bypass the shell by providing an array of command arguments:

```ini
[job-run "direct-exec"]
schedule = @hourly
image = alpine:latest
command = ["echo", "Hello", "World"]
```

This is useful for:
- Avoiding shell parsing issues
- Passing arguments with spaces or special characters
- Slightly improved performance

## Output Handling

### Standard Output and Error

Chadburn captures both stdout and stderr from job executions:

- For `job-local`, output is captured directly
- For container-based jobs, output is retrieved from Docker logs

Output is:
1. Logged to Chadburn's logs
2. Available for inspection through the API
3. Stored temporarily for debugging

### Exit Codes

Exit codes determine job success or failure:

- **0**: Success
- **Non-zero**: Failure

Failed jobs are logged with their exit code and output for troubleshooting.

## Timeout Management

To prevent jobs from running indefinitely, you can set timeouts:

```ini
[job-local "backup"]
schedule = @daily
command = /usr/local/bin/backup.sh
timeout = 1h
```

If a job exceeds its timeout:
1. The process is terminated
2. The job is marked as failed
3. The failure is logged with a timeout message

## Concurrency Control

### Overlapping Executions

By default, Chadburn prevents overlapping executions of the same job. If a job is still running when it's scheduled to run again:

1. The new execution is skipped
2. A warning is logged
3. The job runs again at its next scheduled time

You can change this behavior with the `allow_overlapping` option:

```ini
[job-local "long-process"]
schedule = @hourly
command = /usr/local/bin/process.sh
allow_overlapping = true
```

### Concurrency Limits

For resource-intensive jobs, you can limit how many can run simultaneously:

```ini
[global]
concurrency = 5
```

This limits Chadburn to running at most 5 jobs at the same time across all job types.

## Error Handling

### Retry Logic

For transient failures, you can configure automatic retries:

```ini
[job-local "api-call"]
schedule = @hourly
command = /usr/local/bin/api-request.sh
retries = 3
retry_interval = 30s
```

This will:
1. Attempt the job up to 3 additional times after failure
2. Wait 30 seconds between retry attempts
3. Only consider the job failed after all retries are exhausted

### Failure Actions

You can specify actions to take when a job fails:

```ini
[job-local "critical-backup"]
schedule = @daily
command = /usr/local/bin/backup.sh
on_failure = /usr/local/bin/send-alert.sh "Backup failed"
```

The `on_failure` command runs when the job fails, allowing for:
- Notifications
- Cleanup actions
- Fallback procedures

## Advanced Execution Options

### Working Directory

Specify a working directory for job execution:

```ini
[job-local "compile"]
schedule = @daily
command = make all
working_directory = /path/to/source
```

### User and Group

Run jobs as a specific user:

```ini
[job-local "database-maintenance"]
schedule = @weekly
command = /usr/local/bin/db-maintenance.sh
user = postgres
```

### Environment Variables

Provide environment variables to jobs:

```ini
[job-local "deploy"]
schedule = @daily
command = /usr/local/bin/deploy.sh
environment = ["ENV=production", "DEBUG=false"]
```

### Resource Limits

For container-based jobs, limit resource usage:

```ini
[job-run "resource-intensive"]
schedule = @daily
image = data-processor:latest
command = python /app/process.py
memory = 2g
cpu_shares = 512
```

## Best Practices

1. **Keep Commands Idempotent**: Jobs should be safe to run multiple times
2. **Handle Failures Gracefully**: Include error handling in your scripts
3. **Set Appropriate Timeouts**: Prevent jobs from running indefinitely
4. **Use Retries Wisely**: Only retry for transient failures
5. **Monitor Job Execution**: Use Chadburn's metrics to track job performance
6. **Log Verbosely**: Include sufficient logging in your job commands
7. **Test Thoroughly**: Verify job execution in a test environment first

## Next Steps

- Learn about [Metrics](/docs/concepts/metrics) to monitor job execution
- Explore [Jobs](/docs/concepts/jobs) configuration in depth
- Understand [Schedules](/docs/concepts/schedules) to control when jobs run 