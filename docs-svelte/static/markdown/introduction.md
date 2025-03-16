---
layout: default
title: Introduction
nav_order: 1
---

# Introduction to Chadburn

Chadburn is a powerful job scheduler for Docker environments, designed to make it easy to run periodic tasks in your containerized applications. Named after the telegraph system used on ships to communicate between the bridge and the engine room, Chadburn helps you coordinate and schedule tasks across your Docker containers.

## Key Features

### Simple Configuration

Chadburn uses a straightforward YAML configuration format that makes it easy to define and manage your scheduled jobs:

```yaml
jobs:
  - name: "backup-db"
    command: "pg_dump -U postgres my_db > /backups/backup_$(date +%Y%m%d_%H%M%S).sql"
    schedule: "0 0 * * *"  # Run at midnight every day
    container: "postgres"
```

### Flexible Scheduling

Chadburn supports standard cron syntax for job scheduling, allowing you to define complex schedules with ease. You can schedule jobs to run:

- At specific times of day
- On certain days of the week or month
- At regular intervals (every X minutes/hours)
- With custom expressions for complex scheduling needs

### Docker Integration

Chadburn is designed specifically for Docker environments and offers seamless integration:

- Run commands in any container in your Docker environment
- No need to install cron inside your containers
- Works with Docker Compose and standalone Docker deployments
- Supports Docker Swarm for distributed scheduling

### Metrics and Monitoring

With Chadburn's metrics support, you can monitor your scheduled jobs and track their performance:

- Prometheus metrics for job execution
- Track successful and failed jobs
- Monitor execution time and resource usage
- Integrate with Grafana for visualization

## Getting Started

To get started with Chadburn, check out the [Installation Guide](/installation) and [Configuration Guide](/configuration). For examples of how to use Chadburn in different scenarios, see the [Examples](/examples) section. 