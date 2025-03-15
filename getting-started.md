# Getting Started with Chadburn

Chadburn is a lightweight job scheduler designed for Docker environments. It serves as a modern replacement for traditional cron, with enhanced features for container orchestration.

## Installation

### Using Docker (Recommended)

The easiest way to use Chadburn is with Docker:

```bash
docker run -d --name chadburn \
  -v /var/run/docker.sock:/var/run/docker.sock:ro,z \
  -v /path/to/config.ini:/etc/chadburn.conf \
  premoweb/chadburn:latest daemon
```

### From Source

If you prefer to build from source:

```bash
git clone https://github.com/PremoWeb/Chadburn.git
cd Chadburn
go build -o chadburn .
```

## Quick Start

### 1. Create a Configuration File

Create a file named `chadburn.conf` with the following content:

```ini
[job-local "hello-world"]
schedule = @every 1m
command = echo "Hello, Chadburn!"
```

### 2. Run Chadburn

```bash
docker run -d --name chadburn \
  -v /var/run/docker.sock:/var/run/docker.sock:ro,z \
  -v $(pwd)/chadburn.conf:/etc/chadburn.conf \
  premoweb/chadburn:latest daemon
```

### 3. Check the Logs

```bash
docker logs chadburn
```

You should see output indicating that Chadburn has started and is running your job every minute.

## Using Docker Labels

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

## Basic Concepts

### Job Types

Chadburn supports several job types:

- **job-local**: Executes commands on the host machine
- **job-exec**: Executes commands inside a running container
- **job-run**: Creates a new container to execute commands
- **job-service-run**: Creates a service to execute commands
- **job-lifecycle**: Executes commands on container lifecycle events (start/stop)

### Scheduling Format

Chadburn uses the Go implementation of cron for scheduling. Examples include:

- `@every 1h30m`: Run every hour and thirty minutes
- `0 30 * * * *`: Run at 30 minutes past every hour
- `@hourly`: Run once an hour
- `@daily`: Run once a day
- `@midnight`: Run once a day at midnight
- `@weekly`: Run once a week

### Configuration Methods

Chadburn supports two configuration methods:

1. **INI File**: Traditional configuration file
2. **Docker Labels**: Dynamic configuration through container labels

These methods can be used together, allowing for flexible deployment scenarios.

## Next Steps

- Learn more about [Configuration](configuration.md)
- Explore different [Job Types](jobs.md)
- Understand [Docker Integration](docker-integration.md) 