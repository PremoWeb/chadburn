# Schedules in Chadburn

Schedules determine when and how often jobs run in Chadburn. Understanding scheduling options is essential for effective job automation.

## Schedule Formats

Chadburn supports multiple schedule formats:

### Cron Expressions

The traditional cron format with six fields:

```
┌─────────── second (0 - 59)
│ ┌───────── minute (0 - 59)
│ │ ┌─────── hour (0 - 23)
│ │ │ ┌───── day of month (1 - 31)
│ │ │ │ ┌─── month (1 - 12)
│ │ │ │ │ ┌─ day of week (0 - 6) (Sunday to Saturday)
│ │ │ │ │ │
* * * * * *
```

Examples:
- `0 0 * * * *`: Run at the start of every hour
- `0 30 9 * * *`: Run at 9:30 AM every day
- `0 0 0 1 * *`: Run at midnight on the first day of each month

### Predefined Schedules

Chadburn provides convenient predefined schedules:

| Schedule | Description | Equivalent Cron |
|----------|-------------|-----------------|
| `@yearly` or `@annually` | Once a year at midnight on January 1st | `0 0 0 1 1 *` |
| `@monthly` | Once a month at midnight on the 1st | `0 0 0 1 * *` |
| `@weekly` | Once a week at midnight on Sunday | `0 0 0 * * 0` |
| `@daily` or `@midnight` | Once a day at midnight | `0 0 0 * * *` |
| `@hourly` | Once an hour at the beginning of the hour | `0 0 * * * *` |

### Interval Notation

For simple recurring schedules, you can use the `@every` notation:

```
@every <duration>
```

Where `<duration>` is a string like:
- `10s`: 10 seconds
- `1m30s`: 1 minute and 30 seconds
- `2h`: 2 hours
- `24h`: 24 hours (1 day)
- `168h`: 168 hours (1 week)

Examples:
- `@every 1h`: Run every hour
- `@every 30m`: Run every 30 minutes
- `@every 5m30s`: Run every 5 minutes and 30 seconds

## Special Schedules

### One-time Schedules

To run a job just once at a specific time:

```ini
[job-local "one-time-backup"]
schedule = 2023-12-31T23:59:00Z
command = /usr/local/bin/year-end-backup.sh
```

### Container Event Triggers

For `job-lifecycle` jobs, instead of a time-based schedule, you specify container events:

```ini
[job-lifecycle "container-started"]
container = web-app
event = start
command = echo "Container started" >> /var/log/events.log
```

Valid events are:
- `start`: Triggered when a container starts
- `stop`: Triggered when a container stops
- `die`: Triggered when a container dies unexpectedly
- `health_status`: Triggered when a container's health status changes

## Schedule Timezones

By default, Chadburn uses the system's local timezone for scheduling. You can specify a different timezone:

```ini
[job-local "daily-report"]
schedule = 0 0 8 * * *
timezone = America/New_York
command = /usr/local/bin/generate-report.sh
```

## Schedule Limitations

- The smallest scheduling interval is 1 second
- For very frequent jobs (sub-second), consider using a different tool
- Be aware of timezone changes, especially during daylight saving transitions

## Best Practices

### Staggering Jobs

To avoid resource contention, stagger similar jobs:

```ini
[job-local "backup-db1"]
schedule = 0 0 1 * * *  # 1:00 AM
command = backup-script.sh db1

[job-local "backup-db2"]
schedule = 0 0 2 * * *  # 2:00 AM
command = backup-script.sh db2
```

### Avoiding Resource Contention

For resource-intensive jobs, schedule them during off-peak hours:

```ini
[job-local "heavy-processing"]
schedule = 0 0 2 * * *  # 2:00 AM, typically low-traffic
command = /usr/local/bin/process-data.sh
```

### Handling Job Dependencies

For jobs that depend on others, use appropriate scheduling:

```ini
[job-local "generate-data"]
schedule = 0 0 1 * * *  # 1:00 AM
command = /usr/local/bin/generate.sh

[job-local "process-data"]
schedule = 0 0 2 * * *  # 2:00 AM (1 hour after generation)
command = /usr/local/bin/process.sh
```

## Examples

### Daily Backup at 2 AM

```ini
[job-local "daily-backup"]
schedule = 0 0 2 * * *
command = /usr/local/bin/backup.sh
```

### Weekday Business Hours Check

```ini
[job-exec "business-hours-check"]
schedule = 0 */15 9-17 * * 1-5  # Every 15 minutes, 9 AM to 5 PM, Monday to Friday
container = web-app
command = curl -s http://localhost/health
```

### Monthly Maintenance

```ini
[job-run "monthly-maintenance"]
schedule = @monthly
image = maintenance:latest
command = /app/run-maintenance.sh
```

## Next Steps

- Learn about [Execution](/docs/concepts/execution) details and options
- Explore [Metrics](/docs/concepts/metrics) to monitor your scheduled jobs
- Understand [Jobs](/docs/concepts/jobs) configuration in depth 