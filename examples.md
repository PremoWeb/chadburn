# Chadburn Examples

This document provides practical examples of using Chadburn for various common scenarios.

## Database Backup Examples

### PostgreSQL Backup

```ini
[job-exec "postgres-backup"]
schedule = @daily
container = postgres
command = pg_dump -U postgres -d mydb > /backups/mydb-$(date +%Y%m%d).sql
```

Using Docker labels:

```bash
docker run -d --name postgres \
  --label chadburn.enabled=true \
  --label chadburn.job-exec.postgres-backup.schedule="@daily" \
  --label chadburn.job-exec.postgres-backup.command="pg_dump -U postgres -d mydb > /backups/mydb-$(date +%Y%m%d).sql" \
  -v backups:/backups \
  postgres:latest
```

### MySQL/MariaDB Backup

```ini
[job-exec "mysql-backup"]
schedule = @daily
container = mysql
command = mysqldump -u root -p"${MYSQL_ROOT_PASSWORD}" --all-databases > /backups/mysql-$(date +%Y%m%d).sql
```

Using Docker labels:

```bash
docker run -d --name mysql \
  --label chadburn.enabled=true \
  --label chadburn.job-exec.mysql-backup.schedule="@daily" \
  --label chadburn.job-exec.mysql-backup.command="mysqldump -u root -p\"${MYSQL_ROOT_PASSWORD}\" --all-databases > /backups/mysql-$(date +%Y%m%d).sql" \
  -v backups:/backups \
  -e MYSQL_ROOT_PASSWORD=secret \
  mysql:latest
```

### MongoDB Backup

```ini
[job-exec "mongodb-backup"]
schedule = @daily
container = mongodb
command = mongodump --out=/backups/mongodb-$(date +%Y%m%d)
```

## File Backup Examples

### Simple Directory Backup

```ini
[job-exec "app-data-backup"]
schedule = @daily
container = app
command = tar -czf /backups/app-data-$(date +%Y%m%d).tar.gz /app/data
```

### Backup with Rotation

```ini
[job-exec "backup-with-rotation"]
schedule = @daily
container = app
command = tar -czf /backups/app-data-$(date +%Y%m%d).tar.gz /app/data && find /backups -name "app-data-*.tar.gz" -mtime +7 -delete
```

This creates a daily backup and deletes backups older than 7 days.

## Cleanup Examples

### Docker Cleanup

```ini
[job-local "docker-cleanup"]
schedule = @weekly
command = docker system prune -f
```

### Log Rotation

```ini
[job-exec "log-rotation"]
schedule = @daily
container = app
command = find /app/logs -name "*.log" -mtime +30 -delete
```

## Monitoring Examples

### Health Check

```ini
[job-exec "health-check"]
schedule = @every 5m
container = app
command = curl -f http://localhost:8080/health || echo "Health check failed" | mail -s "Health Check Alert" admin@example.com
```

### Disk Space Check

```ini
[job-local "disk-space-check"]
schedule = @hourly
command = df -h | grep -E "([8-9][0-9]|100)%" && echo "Disk space critical" | mail -s "Disk Space Alert" admin@example.com
```

## Notification Examples

### Slack Notification

```ini
[global]
slack-webhook = https://hooks.slack.com/services/XXX/YYY/ZZZ

[job-local "notify-deployment"]
schedule = @daily
command = echo "Daily deployment starting"
```

### Email Report

```ini
[global]
smtp-host = smtp.example.com
smtp-port = 587
smtp-user = user
smtp-password = password
email-to = reports@example.com
email-from = chadburn@example.com

[job-exec "daily-report"]
schedule = @daily
container = app
command = /app/generate-report.sh
```

## Web Application Examples

### Cache Clearing

```ini
[job-exec "clear-cache"]
schedule = @hourly
container = webapp
command = rm -rf /app/cache/*
```

### Session Cleanup

```ini
[job-exec "session-cleanup"]
schedule = @daily
container = webapp
command = find /app/sessions -type f -mtime +7 -delete
```

## Data Processing Examples

### ETL Job

```ini
[job-exec "etl-process"]
schedule = 0 2 * * *
container = etl
command = /app/run-etl.sh
```

### Data Export

```ini
[job-exec "data-export"]
schedule = @weekly
container = app
command = /app/export-data.sh
```

## Container Lifecycle Examples

### Notification on Container Start

```ini
[job-lifecycle "notify-on-start"]
container = critical-app
event-type = start
command = echo "Container critical-app started at $(date)" | mail -s "Container Started" admin@example.com
```

### Cleanup on Container Stop

```ini
[job-lifecycle "cleanup-on-stop"]
container = app
event-type = stop
command = rm -rf /tmp/app-cache
```

## Advanced Examples

### Running Multiple Commands

```ini
[job-exec "multiple-commands"]
schedule = @daily
container = app
command = /bin/bash -c "echo 'Starting backup' && tar -czf /backups/data.tar.gz /data && echo 'Backup completed'"
```

### Using Environment Variables

```ini
[job-exec "env-example"]
schedule = @daily
container = app
command = /bin/bash -c "export DB_HOST=${DB_HOST:-localhost} && /app/backup-db.sh"
```

### Conditional Execution

```ini
[job-exec "conditional-job"]
schedule = @daily
container = app
command = /bin/bash -c "[ -f /app/needs-processing ] && /app/process.sh || echo 'Nothing to process'"
```

### Job with Timeout

```ini
[job-exec "long-running-job"]
schedule = @daily
container = app
command = /app/long-process.sh
timeout = 1h
```

### Variable Substitution

```ini
[job-exec "variable-example"]
schedule = @daily
container = app
command = echo "Running in container {{.Container.Name}} at {{.Time.Now}}"
```

## Docker Compose Examples

### Complete Application Stack

```yaml
version: "3"
services:
  chadburn:
    image: premoweb/chadburn:latest
    container_name: chadburn
    restart: unless-stopped
    command: daemon --config=/etc/chadburn.conf
    volumes:
      - /var/run/docker.sock:/var/run/docker.sock:ro,z
      - ./chadburn.conf:/etc/chadburn.conf
      - backups:/backups
    depends_on:
      - app
      - db

  app:
    image: your-app:latest
    container_name: app
    labels:
      chadburn.enabled: "true"
      chadburn.job-exec.app-backup.schedule: "@daily"
      chadburn.job-exec.app-backup.command: "tar -czf /backups/app-data.tar.gz /app/data"
    volumes:
      - app-data:/app/data
      - backups:/backups

  db:
    image: postgres:latest
    container_name: db
    labels:
      chadburn.enabled: "true"
      chadburn.job-exec.db-backup.schedule: "@daily"
      chadburn.job-exec.db-backup.command: "pg_dump -U postgres -d mydb > /backups/db-backup.sql"
    volumes:
      - db-data:/var/lib/postgresql/data
      - backups:/backups
    environment:
      - POSTGRES_PASSWORD=secret

volumes:
  app-data:
  db-data:
  backups:
```

With `chadburn.conf`:

```ini
[global]
slack-webhook = https://hooks.slack.com/services/XXX/YYY/ZZZ
slack-only-on-error = true

[job-local "system-check"]
schedule = @hourly
command = df -h | grep -E "([8-9][0-9]|100)%" && echo "Disk space critical" || echo "Disk space OK"
``` 