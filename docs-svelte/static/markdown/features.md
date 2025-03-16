---
layout: default
title: Overview
nav_order: 3
has_children: true
---

# Chadburn Overview

Chadburn is a modern job scheduler designed specifically for Docker environments, offering a comprehensive set of features that make it a powerful alternative to traditional cron.

## Core Functionality

### [Jobs](jobs.html)
Chadburn supports multiple job types to fit your needs:
- **job-local**: Execute commands on the host machine
- **job-exec**: Run commands inside existing containers
- **job-run**: Create new containers to execute commands
- **job-service-run**: Create services to execute commands
- **job-lifecycle**: Execute commands on container lifecycle events

### [Docker Integration](docker-integration.html)
Seamlessly integrate with Docker:
- Configure jobs using Docker labels
- Manage container lifecycle events
- Dynamic configuration through container metadata

### [Security](docker-socket-permissions.html)
Ensure secure operation:
- Docker socket permission management
- Least privilege principles
- Secure container execution

## Additional Features

- **Flexible Scheduling**: Cron-like syntax with extended capabilities
- **Notifications**: Integrations with Slack, Email, and Gotify
- **Metrics**: Prometheus-compatible metrics endpoint for monitoring
- **Variable Substitution**: Use environment variables in your job definitions
- **Middlewares**: Extend functionality with custom middleware components 