---
layout: default
title: Home
nav_order: 1
---

# Chadburn Documentation

Welcome to the Chadburn documentation. This guide provides comprehensive information about installing, configuring, and using Chadburn, a modern job scheduler for Docker environments.

![Chadburn Logo](assets/images/chadburn-logo.png)

## What is Chadburn?

Chadburn is a lightweight job scheduler designed for Docker environments. It serves as a modern replacement for traditional cron, with enhanced features for container orchestration.

### Key Features

- **Docker Integration**: Native support for Docker containers
- **Dynamic Configuration**: Configure jobs using Docker labels
- **Multiple Job Types**: Run commands locally, in existing containers, or create new containers
- **Container Lifecycle Events**: Execute commands when containers start or stop
- **Notifications**: Slack, Email, and Gotify integration
- **Metrics**: Prometheus-compatible metrics endpoint

## Getting Started

To get started with Chadburn, check out the [Getting Started Guide](getting-started.html).

## Table of Contents

1. [Getting Started](getting-started.html)
   - Installation
   - Quick Start
   - Basic Concepts

2. [Configuration](configuration.html)
   - Configuration File Format
   - Environment Variables
   - Command Line Options

3. [Job Types](jobs.html)
   - job-exec
   - job-run
   - job-local
   - job-service-run
   - job-lifecycle

4. [Docker Integration](docker-integration.html)
   - Docker Labels
   - Container Lifecycle Events
   - Dynamic Configuration

5. [Troubleshooting](troubleshooting.html)
   - Common Issues
   - Logging
   - Debugging

6. [Advanced Topics](advanced-topics.html)
   - Middlewares
   - Variable Substitution
   - Metrics and Monitoring

7. [Deployment](deployment.html)
   - Docker Compose
   - Kubernetes
   - Swarm Mode

8. [Security](security.html)
   - Docker Socket Permissions
   - Best Practices

9. [Examples](examples.html)
   - Common Use Cases
   - Docker Compose Examples

10. [FAQ](faq.html)
    - Frequently Asked Questions

11. [Contributing](contributing.html)
    - Development Setup
    - Code Style
    - Testing

## License

Chadburn is licensed under the [MIT License](https://github.com/PremoWeb/Chadburn/blob/main/LICENSE). 