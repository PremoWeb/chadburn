# Quick Start Guide

This guide will help you get up and running with Chadburn in just a few minutes.

## Prerequisites

- Docker installed on your system
- Basic familiarity with command line operations

## 1. Create a Configuration File

Create a file named `chadburn.conf` with the following content:

```ini
[job-local "hello-world"]
schedule = @every 1m
command = echo "Hello, Chadburn!"
```

This configuration defines a simple job that prints "Hello, Chadburn!" every minute.

## 2. Run Chadburn

Start Chadburn using Docker:

```bash
docker run -d --name chadburn \
  -v /var/run/docker.sock:/var/run/docker.sock:ro,z \
  -v $(pwd)/chadburn.conf:/etc/chadburn.conf \
  premoweb/chadburn:latest daemon
```

This command:
- Runs Chadburn in a Docker container
- Mounts the Docker socket to allow Chadburn to interact with Docker
- Mounts your configuration file into the container

## 3. Check the Logs

Verify that Chadburn is running correctly:

```bash
docker logs chadburn
```

You should see output indicating that Chadburn has started and is running your job every minute.

## 4. Using Docker Labels (Alternative Approach)

Chadburn can also be configured using Docker labels. Here's an example:

```bash
docker run -d --name nginx \
  --label chadburn.enabled=true \
  --label chadburn.job-exec.print-date.schedule="@every 1m" \
  --label chadburn.job-exec.print-date.command="date" \
  nginx
```

Then run Chadburn with access to the Docker socket:

```bash
docker run -d --name chadburn \
  -v /var/run/docker.sock:/var/run/docker.sock:ro,z \
  premoweb/chadburn:latest daemon
```

## 5. Verify Jobs are Running

To verify that your jobs are running as expected:

```bash
# For jobs running in the Chadburn container
docker logs chadburn

# For jobs running in other containers (like the nginx example)
docker logs nginx
```

## Next Steps

Now that you have Chadburn up and running, you might want to:

- Learn more about [Configuration](/docs/guides/configuration)
- Explore different [Job Types](/docs/concepts/jobs)
- Understand [Docker Integration](/docs/docker-integration)
- Set up [Monitoring](/docs/guides/monitoring) for your jobs 