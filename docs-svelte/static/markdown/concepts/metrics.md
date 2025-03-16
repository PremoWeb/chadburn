# Metrics in Chadburn

Chadburn provides comprehensive metrics to help you monitor and understand the performance and health of your scheduled jobs. This document explains the available metrics, how to access them, and best practices for monitoring.

## Metrics Overview

Chadburn exposes metrics in Prometheus format, providing insights into:

- Job execution counts and durations
- Success and failure rates
- Scheduler performance
- System resource usage

These metrics help you:
- Identify problematic jobs
- Track execution trends
- Set up alerts for job failures
- Optimize scheduling

## Available Metrics

### Job Execution Metrics

| Metric | Type | Description |
|--------|------|-------------|
| `chadburn_job_executions_total` | Counter | Total number of job executions, labeled by job name and status |
| `chadburn_job_execution_duration_seconds` | Histogram | Duration of job executions in seconds |
| `chadburn_job_last_execution_timestamp` | Gauge | Timestamp of the last execution of each job |
| `chadburn_job_last_success_timestamp` | Gauge | Timestamp of the last successful execution |
| `chadburn_job_last_error_timestamp` | Gauge | Timestamp of the last failed execution |

### Scheduler Metrics

| Metric | Type | Description |
|--------|------|-------------|
| `chadburn_scheduler_jobs` | Gauge | Number of jobs currently registered |
| `chadburn_scheduler_next_run_seconds` | Gauge | Seconds until the next scheduled job |
| `chadburn_scheduler_overdue_jobs` | Gauge | Number of jobs that are overdue for execution |

### System Metrics

| Metric | Type | Description |
|--------|------|-------------|
| `chadburn_process_cpu_seconds_total` | Counter | Total CPU time spent by Chadburn process |
| `chadburn_process_memory_bytes` | Gauge | Memory usage of Chadburn process |
| `chadburn_process_open_fds` | Gauge | Number of open file descriptors |
| `chadburn_build_info` | Gauge | Build information including version |

## Enabling Metrics

Metrics are disabled by default. To enable them, add the following to your configuration:

```ini
[global]
metrics = true
metrics_path = /metrics
metrics_address = :8080
```

This configuration:
- Enables the metrics endpoint
- Sets the metrics path to `/metrics`
- Makes the metrics available on port 8080

## Accessing Metrics

Once enabled, metrics are available at the configured endpoint:

```
http://your-chadburn-host:8080/metrics
```

Example output:

```
# HELP chadburn_job_executions_total Total number of job executions
# TYPE chadburn_job_executions_total counter
chadburn_job_executions_total{job="backup-database",status="success"} 42
chadburn_job_executions_total{job="backup-database",status="error"} 2
chadburn_job_executions_total{job="cleanup-temp",status="success"} 168
...
```

## Integration with Prometheus

### Prometheus Configuration

Add Chadburn as a scrape target in your `prometheus.yml`:

```yaml
scrape_configs:
  - job_name: 'chadburn'
    static_configs:
      - targets: ['chadburn-host:8080']
```

### Example Queries

Monitor job success rate:
```
sum(rate(chadburn_job_executions_total{status="success"}[5m])) by (job) /
sum(rate(chadburn_job_executions_total[5m])) by (job)
```

Find slow jobs:
```
histogram_quantile(0.95, sum(rate(chadburn_job_execution_duration_seconds_bucket[5m])) by (job, le))
```

Check for jobs that haven't run recently:
```
time() - chadburn_job_last_execution_timestamp > 3600
```

## Grafana Dashboards

Chadburn provides a pre-built Grafana dashboard that you can import. The dashboard includes:

- Job execution success/failure rates
- Execution duration trends
- Scheduler performance metrics
- System resource usage

To import the dashboard:
1. Go to Grafana
2. Click "+" > "Import"
3. Enter dashboard ID `12345` or upload the JSON file from the Chadburn repository
4. Select your Prometheus data source
5. Click "Import"

## Alerting

### Prometheus Alerting Rules

Example alerting rules for Prometheus:

```yaml
groups:
- name: chadburn
  rules:
  - alert: ChadburnJobFailure
    expr: increase(chadburn_job_executions_total{status="error"}[1h]) > 0
    for: 5m
    labels:
      severity: warning
    annotations:
      summary: "Job {{ $labels.job }} has failed"
      description: "The job has failed at least once in the last hour."

  - alert: ChadburnJobMissing
    expr: time() - chadburn_job_last_execution_timestamp > 2 * 3600
    for: 5m
    labels:
      severity: warning
    annotations:
      summary: "Job {{ $labels.job }} hasn't run recently"
      description: "The job hasn't run in more than twice its expected interval."
```

### Grafana Alerts

You can also set up alerts directly in Grafana:

1. Edit a panel in your dashboard
2. Go to "Alert" tab
3. Configure alert conditions based on metrics
4. Set notification channels
5. Save the alert

## Best Practices

### Monitoring Strategy

1. **Track Success Rates**: Monitor the ratio of successful to total executions
2. **Watch Execution Times**: Alert on jobs that take longer than expected
3. **Check Execution Frequency**: Ensure jobs run at their expected intervals
4. **Monitor Resource Usage**: Watch for resource constraints affecting job execution

### Dashboard Organization

1. **Overview Dashboard**: High-level view of all jobs
2. **Job-Specific Dashboards**: Detailed metrics for critical jobs
3. **Alert Dashboard**: Shows current and recent alerts

### Alert Tuning

1. **Set Appropriate Thresholds**: Avoid alert fatigue from false positives
2. **Use Trend-Based Alerts**: Alert on changes in patterns, not just absolute values
3. **Implement Alert Severity Levels**: Distinguish between warnings and critical alerts

## Troubleshooting with Metrics

### Common Issues and Metrics to Check

| Issue | Metrics to Check |
|-------|------------------|
| Job failures | `chadburn_job_executions_total{status="error"}` |
| Slow jobs | `chadburn_job_execution_duration_seconds` |
| Missed schedules | `chadburn_job_last_execution_timestamp`, `chadburn_scheduler_overdue_jobs` |
| Resource constraints | `chadburn_process_cpu_seconds_total`, `chadburn_process_memory_bytes` |

## Next Steps

- Set up [Prometheus and Grafana](/docs/guides/monitoring) for monitoring
- Learn about [Jobs](/docs/concepts/jobs) configuration
- Understand [Execution](/docs/concepts/execution) details
- Explore [Schedules](/docs/concepts/schedules) to control when jobs run 