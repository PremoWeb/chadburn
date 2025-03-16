# Caddy Integration

Caddy is a powerful, enterprise-ready, open source web server with automatic HTTPS. This document explains how to integrate Chadburn with Caddy for scheduled maintenance tasks.

## Overview

Chadburn can be used to automate various maintenance tasks for Caddy, such as:

- Configuration reloading
- Log rotation and cleanup
- Configuration backups
- Certificate management
- Health checks

## Standalone Caddy Container

Here's how to set up a standalone Caddy container with Chadburn labels for common maintenance tasks:

```bash
docker run -d --name caddy \
  --label chadburn.enabled=true \
  --label chadburn.job-exec.reload-config.schedule="@every 1h" \
  --label chadburn.job-exec.reload-config.command="caddy reload --config /etc/caddy/Caddyfile" \
  --label chadburn.job-exec.clean-logs.schedule="@daily" \
  --label chadburn.job-exec.clean-logs.command="find /var/log/caddy -type f -name '*.log' -mtime +7 -delete" \
  --label chadburn.job-exec.backup-caddyfile.schedule="@weekly" \
  --label chadburn.job-exec.backup-caddyfile.command="cp /etc/caddy/Caddyfile /backup/Caddyfile.$(date +%Y%m%d)" \
  -p 80:80 -p 443:443 \
  -v caddy_data:/data \
  -v caddy_config:/config \
  -v ./Caddyfile:/etc/caddy/Caddyfile \
  caddy:latest
```

## Caddy with Docker Compose

Here's a Docker Compose example for Caddy with Chadburn integration:

```yaml
version: "3"
services:
  chadburn:
    image: premoweb/chadburn:latest
    command: daemon
    volumes:
      - /var/run/docker.sock:/var/run/docker.sock:ro,z
    depends_on:
      - caddy

  caddy:
    image: caddy:latest
    ports:
      - "80:80"
      - "443:443"
    volumes:
      - ./Caddyfile:/etc/caddy/Caddyfile
      - caddy_data:/data
      - caddy_config:/config
      - ./site:/srv
      - ./backups:/backups
    labels:
      chadburn.enabled: "true"
      # Reload Caddy configuration hourly
      chadburn.job-exec.reload-config.schedule: "@every 1h"
      chadburn.job-exec.reload-config.command: "caddy reload --config /etc/caddy/Caddyfile"
      # Backup Caddyfile weekly
      chadburn.job-exec.backup-caddyfile.schedule: "@weekly"
      chadburn.job-exec.backup-caddyfile.command: "cp /etc/caddy/Caddyfile /backups/Caddyfile.$(date +%Y%m%d)"
      # Clean access logs older than 7 days
      chadburn.job-exec.clean-logs.schedule: "@daily"
      chadburn.job-exec.clean-logs.command: "find /data/caddy/logs -type f -name 'access.log.*' -mtime +7 -delete"
      # Check TLS certificates and reload if needed
      chadburn.job-exec.check-certs.schedule: "@daily"
      chadburn.job-exec.check-certs.command: "caddy validate --config /etc/caddy/Caddyfile && caddy reload --config /etc/caddy/Caddyfile"

volumes:
  caddy_data:
  caddy_config:
```

## Common Caddy Maintenance Tasks

Here are some useful Chadburn jobs for Caddy:

1. **Reload Configuration**: Periodically reload Caddy to pick up configuration changes
   ```
   chadburn.job-exec.reload-config.schedule: "@hourly"
   chadburn.job-exec.reload-config.command: "caddy reload --config /etc/caddy/Caddyfile"
   ```

2. **Backup Configuration**: Create regular backups of your Caddyfile
   ```
   chadburn.job-exec.backup-config.schedule: "@daily"
   chadburn.job-exec.backup-config.command: "cp /etc/caddy/Caddyfile /backups/Caddyfile.$(date +%Y%m%d)"
   ```

3. **Log Rotation**: Clean up old logs to prevent disk space issues
   ```
   chadburn.job-exec.rotate-logs.schedule: "@daily"
   chadburn.job-exec.rotate-logs.command: "find /data/caddy/logs -type f -mtime +7 -delete"
   ```

4. **Certificate Validation**: Regularly check TLS certificates
   ```
   chadburn.job-exec.check-certs.schedule: "@daily"
   chadburn.job-exec.check-certs.command: "caddy validate"
   ```

5. **Health Check**: Ensure Caddy is responding properly
   ```
   chadburn.job-exec.health-check.schedule: "@every 5m"
   chadburn.job-exec.health-check.command: "wget -q --spider http://localhost:80 || caddy reload --config /etc/caddy/Caddyfile"
   ```

## Advanced Configuration

### Automatic Site Updates

You can use Chadburn to automatically update your website content:

```yaml
services:
  caddy:
    # ... other configuration ...
    labels:
      chadburn.enabled: "true"
      chadburn.job-exec.update-site.schedule: "@hourly"
      chadburn.job-exec.update-site.command: |
        cd /srv && git pull origin main && caddy reload --config /etc/caddy/Caddyfile
```

### Database Backups for Web Applications

If your Caddy server hosts web applications with databases:

```yaml
services:
  caddy:
    # ... other configuration ...
    labels:
      chadburn.enabled: "true"
      chadburn.job-exec.backup-db.schedule: "@daily"
      chadburn.job-exec.backup-db.command: |
        mysqldump -h db -u user -ppassword myapp > /backups/myapp-$(date +%Y%m%d).sql
```

### Monitoring and Alerting

Set up monitoring and alerting for your Caddy server:

```yaml
services:
  caddy:
    # ... other configuration ...
    labels:
      chadburn.enabled: "true"
      chadburn.job-exec.monitor.schedule: "@every 10m"
      chadburn.job-exec.monitor.command: |
        if ! curl -s --head http://localhost:80 | grep "200 OK"; then
          echo "Caddy server is not responding properly" | mail -s "Caddy Alert" admin@example.com
        fi
```

## Troubleshooting

### Common Issues

1. **Permission Denied for Caddy Files**
   
   Ensure proper permissions for Caddy files:
   
   ```bash
   chown -R 1000:1000 /path/to/Caddyfile
   chmod 644 /path/to/Caddyfile
   ```

2. **Caddy Not Reloading Configuration**
   
   Verify that the Caddyfile is valid:
   
   ```bash
   caddy validate --config /etc/caddy/Caddyfile
   ```

3. **TLS Certificate Issues**
   
   Check Caddy's certificate storage:
   
   ```bash
   caddy list-certificates
   ```

## Security Considerations

When integrating Chadburn with Caddy, consider these security best practices:

1. Use read-only mounts for configuration files
2. Separate volumes for different types of data (config, site, logs)
3. Regularly validate configuration before reloading
4. Use environment variables for sensitive information
5. Implement proper backup rotation to prevent disk space issues

## Example Caddyfile

Here's an example Caddyfile that works well with Chadburn:

```
example.com {
    root * /srv
    file_server
    encode gzip
    log {
        output file /data/caddy/logs/access.log {
            roll_size 10MB
            roll_keep 10
            roll_keep_for 168h
        }
    }
    tls admin@example.com
}
```

This configuration includes:
- Basic file serving
- Gzip compression
- Structured logging with rotation
- Automatic HTTPS with Let's Encrypt

With Chadburn, you can automate the maintenance of this setup, ensuring your Caddy server runs smoothly with minimal manual intervention. 