# Troubleshooting Chadburn

This guide provides solutions for common issues you might encounter when using Chadburn.

## Common Issues

### Docker Socket Permission Denied

**Problem**: When running Chadburn in a container, you see errors like:

```
ERROR Docker events error: permission denied while trying to connect to the Docker daemon socket
```

**Solution**: This issue is related to Docker socket permissions, which may be more prominent in recent versions due to the migration to the official Docker client library.

See the [Docker Socket Permissions Guide](docker-socket-permissions.md) for detailed solutions.

Quick fixes:
1. Add the `:z` suffix to the volume mount: `-v /var/run/docker.sock:/var/run/docker.sock:ro,z`
2. Run the container with the correct Docker group ID: 
   ```bash
   DOCKER_GID=$(getent group docker | cut -d: -f3)
   docker run -e DOCKER_GID=$DOCKER_GID --user "1000:$DOCKER_GID" ...
   ```
3. Run the container as root: `--user root`

**Note for upgraders**: If you're upgrading from an older version of Chadburn and suddenly experiencing this issue, it's likely due to the migration from `fsouza/go-dockerclient` to the official Docker client library. The official client may have different permission requirements.

### Jobs Not Running

**Problem**: Jobs are configured but not running.

**Solutions**:

1. **Check the schedule format**:
   - Ensure the schedule format is correct
   - For `@every` format, make sure there's a space between `@every` and the duration

2. **Check container labels**:
   - For `job-exec`, ensure the target container has `chadburn.enabled=true`
   - Verify label syntax: `chadburn.<JOB_TYPE>.<JOB_NAME>.<PARAMETER>=<VALUE>`

3. **Check container status**:
   - For `job-exec`, ensure the target container is running
   - For `job-run`, ensure the image exists and is accessible

4. **Check logs**:
   - Run `docker logs chadburn` to see if there are any error messages

### Container Not Found

**Problem**: You see errors like `container not found` when using `job-exec`.

**Solutions**:

1. Ensure the container name in the configuration matches exactly
2. Check if the container is running: `docker ps | grep <container-name>`
3. If using Docker Compose, ensure the container name isn't prefixed with the project name

### Command Execution Failures

**Problem**: Jobs run but the commands fail.

**Solutions**:

1. **Check command syntax**:
   - Ensure the command is valid for the container's environment
   - For complex commands, consider using a shell script inside the container

2. **Check permissions**:
   - Ensure the user running the command has the necessary permissions
   - For `job-exec`, consider specifying a user: `user = root`

3. **Check environment**:
   - Ensure any required environment variables are set
   - For `job-run`, specify environment variables: `environment = ["VAR1=value1", "VAR2=value2"]`

### Scheduling Issues

**Problem**: Jobs don't run at the expected times.

**Solutions**:

1. **Check timezone**:
   - Chadburn uses the container's timezone
   - Mount the host's timezone: `-v /etc/timezone:/etc/timezone:ro -v /etc/localtime:/etc/localtime:ro`

2. **Check schedule format**:
   - Ensure the cron expression is correct
   - For complex schedules, test with tools like [crontab.guru](https://crontab.guru/)

3. **Check for overlapping jobs**:
   - If a job takes longer than its schedule interval, it might appear to skip executions
   - Consider adding a timeout to prevent long-running jobs

## Logging and Debugging

### Enabling Debug Logging

To enable debug logging, set the `CHADBURN_LOG_LEVEL` environment variable:

```bash
docker run -d --name chadburn \
  -e CHADBURN_LOG_LEVEL=debug \
  -v /var/run/docker.sock:/var/run/docker.sock:ro,z \
  -v /path/to/config.ini:/etc/chadburn.conf \
  premoweb/chadburn:latest daemon
```

Available log levels:
- `debug`: Most verbose, shows all details
- `info`: Shows informational messages
- `notice`: Shows important notices
- `warning`: Shows warnings
- `error`: Shows errors only
- `critical`: Shows critical errors only

### Viewing Logs

To view Chadburn logs:

```bash
docker logs chadburn
```

To follow logs in real-time:

```bash
docker logs -f chadburn
```

To see only error logs:

```bash
docker logs chadburn 2>&1 | grep ERROR
```

### Validating Configuration

To validate your configuration file without starting the daemon:

```bash
docker run --rm \
  -v /path/to/config.ini:/etc/chadburn.conf \
  premoweb/chadburn:latest validate
```

### Metrics Endpoint

If you've enabled the metrics endpoint, you can check job statistics:

```bash
curl http://localhost:8080/metrics
```

This provides Prometheus-compatible metrics about job executions, errors, and durations.

## Common Error Messages

### "unable to start a empty scheduler"

**Cause**: No valid jobs were found in the configuration.

**Solution**: Check your configuration file or Docker labels to ensure at least one valid job is defined.

### "unable to add a job with a empty schedule"

**Cause**: A job was defined without a schedule.

**Solution**: Add a valid schedule to the job configuration.

### "error non-zero exit code"

**Cause**: The command executed by the job returned a non-zero exit code.

**Solution**: Check the command for errors, or modify it to always return 0 if the exit code isn't important.

### "the job has exceed the maximum allowed time running"

**Cause**: The job ran longer than the configured timeout.

**Solution**: Increase the timeout value or optimize the job to run faster.

## Getting Help

If you're still experiencing issues:

1. Check the [GitHub Issues](https://github.com/PremoWeb/Chadburn/issues) to see if others have encountered the same problem
2. Open a new issue with:
   - Your Chadburn version
   - Your configuration
   - Complete error logs
   - Steps to reproduce the issue 