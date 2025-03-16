---
layout: default
title: Installation
nav_order: 2
---

# Installing Chadburn

Chadburn is distributed as a Docker image, making it easy to deploy in any environment that supports Docker. This guide will walk you through the installation process.

## Prerequisites

Before installing Chadburn, ensure you have:

- Docker installed on your system
- Basic knowledge of Docker and container orchestration
- Access to the Docker socket (for Docker integration)

## Installation Methods

### Using Docker Run

The simplest way to run Chadburn is using the `docker run` command:

```bash
docker run -d \
  --name chadburn \
  -v /var/run/docker.sock:/var/run/docker.sock \
  -v $(pwd)/config.yml:/etc/chadburn/config.yml \
  premoweb/chadburn:latest
```

### Using Docker Compose

For a more maintainable setup, you can use Docker Compose:

```yaml
version: '3'

services:
  chadburn:
    image: premoweb/chadburn:latest
    container_name: chadburn
    volumes:
      - /var/run/docker.sock:/var/run/docker.sock
      - ./config.yml:/etc/chadburn/config.yml
    restart: unless-stopped
```

Save this to a file named `docker-compose.yml` and run:

```bash
docker-compose up -d
```

## Configuration

After installation, you'll need to create a configuration file. Create a file named `config.yml` with the following basic structure:

```yaml
jobs:
  - name: "example-job"
    command: "echo 'Hello, world!'"
    schedule: "*/5 * * * *"  # Run every 5 minutes
    container: "my-container"
```

For more detailed configuration options, see the [Configuration Guide](/configuration).

## Verifying Installation

To verify that Chadburn is running correctly:

1. Check the container status:
   ```bash
   docker ps | grep chadburn
   ```

2. View the logs to ensure it's scheduling jobs:
   ```bash
   docker logs chadburn
   ```

You should see output indicating that Chadburn has started and is scheduling jobs according to your configuration.

## Next Steps

Now that you have Chadburn installed, you can:

- Learn more about [configuration options](/configuration)
- Explore [example use cases](/examples)
- Set up [metrics and monitoring](/metrics)
- Troubleshoot common issues in the [troubleshooting guide](/troubleshooting) 