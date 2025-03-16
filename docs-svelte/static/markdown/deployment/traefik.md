# Traefik Integration

Traefik is a modern HTTP reverse proxy and load balancer that makes deploying microservices easy. This document explains how to integrate Chadburn with Traefik for scheduled maintenance tasks.

## Overview

Chadburn can be used to automate various maintenance tasks for Traefik, such as:

- Health checks and automatic recovery
- Configuration backups
- Certificate management
- Log rotation and cleanup

## Standalone Traefik Container

Here's how to set up a standalone Traefik container with Chadburn labels:

```bash
docker run -d --name traefik \
  --label chadburn.enabled=true \
  --label chadburn.job-exec.check-health.schedule="@every 15m" \
  --label chadburn.job-exec.check-health.command="wget -q --spider http://localhost:8080/ping || kill -1 1" \
  --label chadburn.job-exec.backup-config.schedule="@daily" \
  --label chadburn.job-exec.backup-config.command="cp /etc/traefik/traefik.yml /backups/traefik.yml.$(date +%Y%m%d)" \
  --label chadburn.job-exec.clean-logs.schedule="@weekly" \
  --label chadburn.job-exec.clean-logs.command="find /var/log/traefik -type f -name '*.log' -mtime +30 -delete" \
  -p 80:80 -p 443:443 -p 8080:8080 \
  -v /var/run/docker.sock:/var/run/docker.sock:ro,z \
  -v ./traefik.yml:/etc/traefik/traefik.yml \
  -v ./acme.json:/acme.json \
  -v ./logs:/var/log/traefik \
  -v ./backups:/backups \
  traefik:latest
```

## Traefik with Docker Compose

Here's a Docker Compose example for Traefik with Chadburn integration:

```yaml
version: "3"
services:
  chadburn:
    image: premoweb/chadburn:latest
    command: daemon
    volumes:
      - /var/run/docker.sock:/var/run/docker.sock:ro,z
    depends_on:
      - traefik

  traefik:
    image: traefik:latest
    ports:
      - "80:80"
      - "443:443"
      - "8080:8080"
    volumes:
      - /var/run/docker.sock:/var/run/docker.sock:ro,z
      - ./traefik.yml:/etc/traefik/traefik.yml
      - ./dynamic_conf:/etc/traefik/dynamic
      - ./acme.json:/acme.json
      - ./logs:/var/log/traefik
      - ./backups:/backups
    labels:
      chadburn.enabled: "true"
      # Health check every 15 minutes
      chadburn.job-exec.check-health.schedule: "@every 15m"
      chadburn.job-exec.check-health.command: "wget -q --spider http://localhost:8080/ping || kill -1 1"
      # Backup configuration files daily
      chadburn.job-exec.backup-config.schedule: "@daily"
      chadburn.job-exec.backup-config.command: "tar -czf /backups/traefik-config-$(date +%Y%m%d).tar.gz /etc/traefik"
      # Backup ACME certificates weekly
      chadburn.job-exec.backup-certs.schedule: "@weekly"
      chadburn.job-exec.backup-certs.command: "cp /acme.json /backups/acme.json.$(date +%Y%m%d)"
      # Clean access logs older than 30 days
      chadburn.job-exec.clean-logs.schedule: "@monthly"
      chadburn.job-exec.clean-logs.command: "find /var/log/traefik -type f -mtime +30 -delete"
```

## Traefik Configuration Formats

Traefik can be configured in multiple formats. Here are examples of different configuration approaches:

### TOML Configuration

```toml
# traefik.toml

[global]
  checkNewVersion = true
  sendAnonymousUsage = false

[entryPoints]
  [entryPoints.web]
    address = ":80"
    [entryPoints.web.http.redirections.entryPoint]
      to = "websecure"
      scheme = "https"

  [entryPoints.websecure]
    address = ":443"

  [entryPoints.traefik]
    address = ":8080"

[api]
  dashboard = true
  insecure = true

[providers]
  [providers.docker]
    endpoint = "unix:///var/run/docker.sock"
    exposedByDefault = false
    watch = true
    
  [providers.file]
    directory = "/etc/traefik/dynamic"
    watch = true

[certificatesResolvers.letsencrypt.acme]
  email = "admin@example.com"
  storage = "/acme.json"
  [certificatesResolvers.letsencrypt.acme.httpChallenge]
    entryPoint = "web"

# Chadburn configuration in docker-compose.yml:
#
# chadburn:
#   image: premoweb/chadburn:latest
#   command: daemon
#   volumes:
#     - /var/run/docker.sock:/var/run/docker.sock:ro,z
#   depends_on:
#     - traefik
#
# traefik:
#   # ... other traefik configuration ...
#   labels:
#     chadburn.enabled: "true"
#     chadburn.job-exec.reload-config.schedule: "@every 1h"
#     chadburn.job-exec.reload-config.command: "kill -USR1 1"
#     chadburn.job-exec.backup-config.schedule: "@daily"
#     chadburn.job-exec.backup-config.command: "cp /etc/traefik/traefik.toml /backups/traefik.toml.$(date +%Y%m%d)"
```

### YAML Configuration

```yaml
# traefik.yml

global:
  checkNewVersion: true
  sendAnonymousUsage: false

entryPoints:
  web:
    address: ":80"
    http:
      redirections:
        entryPoint:
          to: websecure
          scheme: https
  
  websecure:
    address: ":443"
  
  traefik:
    address: ":8080"

api:
  dashboard: true
  insecure: true

providers:
  docker:
    endpoint: "unix:///var/run/docker.sock"
    exposedByDefault: false
    watch: true
  
  file:
    directory: "/etc/traefik/dynamic"
    watch: true

certificatesResolvers:
  letsencrypt:
    acme:
      email: "admin@example.com"
      storage: "/acme.json"
      httpChallenge:
        entryPoint: web

# Chadburn configuration in docker-compose.yml:
#
# chadburn:
#   image: premoweb/chadburn:latest
#   command: daemon
#   volumes:
#     - /var/run/docker.sock:/var/run/docker.sock:ro,z
#   depends_on:
#     - traefik
#
# traefik:
#   # ... other traefik configuration ...
#   labels:
#     chadburn.enabled: "true"
#     chadburn.job-exec.reload-config.schedule: "@every 1h"
#     chadburn.job-exec.reload-config.command: "kill -USR1 1"
#     chadburn.job-exec.backup-config.schedule: "@daily"
#     chadburn.job-exec.backup-config.command: "cp /etc/traefik/traefik.yml /backups/traefik.yml.$(date +%Y%m%d)"
```

## Common Traefik Maintenance Tasks

Here are some useful Chadburn jobs for Traefik:

1. **Health Check**: Periodically check if Traefik is responding and restart if needed
   ```
   chadburn.job-exec.health-check.schedule: "@every 5m"
   chadburn.job-exec.health-check.command: "wget -q --spider http://localhost:8080/ping || kill -1 1"
   ```

2. **Backup ACME Certificates**: Regularly backup Let's Encrypt certificates
   ```
   chadburn.job-exec.backup-certs.schedule: "@weekly"
   chadburn.job-exec.backup-certs.command: "cp /acme.json /backups/acme.json.$(date +%Y%m%d)"
   ```

3. **Configuration Backup**: Create backups of your Traefik configuration
   ```
   chadburn.job-exec.backup-config.schedule: "@daily"
   chadburn.job-exec.backup-config.command: "tar -czf /backups/traefik-config-$(date +%Y%m%d).tar.gz /etc/traefik"
   ```

4. **Log Management**: Clean up old logs to prevent disk space issues
   ```
   chadburn.job-exec.clean-logs.schedule: "@weekly"
   chadburn.job-exec.clean-logs.command: "find /var/log/traefik -type f -mtime +30 -delete"
   ```

5. **Reload Configuration**: Force Traefik to reload its configuration
   ```
   chadburn.job-exec.reload-config.schedule: "@hourly"
   chadburn.job-exec.reload-config.command: "kill -USR1 1"
   ```

## Advanced Configuration

### Dynamic Configuration Reloading

Traefik supports hot-reloading of configuration, which can be automated with Chadburn:

```yaml
services:
  traefik:
    # ... other configuration ...
    labels:
      chadburn.enabled: "true"
      chadburn.job-exec.watch-config.schedule: "@every 5m"
      chadburn.job-exec.watch-config.command: |
        if [ -f /etc/traefik/dynamic/updated-config.yml ]; then
          mv /etc/traefik/dynamic/updated-config.yml /etc/traefik/dynamic/config.yml
          echo "Configuration updated at $(date)"
        fi
```

### Certificate Renewal Monitoring

Monitor Let's Encrypt certificate renewal:

```yaml
services:
  traefik:
    # ... other configuration ...
    labels:
      chadburn.enabled: "true"
      chadburn.job-exec.check-certs.schedule: "@daily"
      chadburn.job-exec.check-certs.command: |
        if openssl x509 -checkend 604800 -noout -in /path/to/cert.pem; then
          echo "Certificate valid for more than 7 days"
        else
          echo "Certificate will expire within 7 days" | mail -s "Certificate Expiry Warning" admin@example.com
        fi
```

## Troubleshooting

### Common Issues

1. **Permission Denied for Docker Socket**
   
   If Chadburn cannot access the Docker socket, ensure the proper permissions:
   
   ```bash
   chmod 666 /var/run/docker.sock
   ```
   
   Or use a more secure approach with proper group permissions.

2. **Traefik Not Responding to Health Checks**
   
   Ensure the Traefik API is properly configured and accessible:
   
   ```toml
   [api]
     dashboard = true
     insecure = true  # Only for testing, use secure settings in production
   ```

3. **Configuration Not Being Applied**
   
   Check that your configuration files are properly mounted and have the correct permissions:
   
   ```yaml
   volumes:
     - ./traefik.yml:/etc/traefik/traefik.yml:ro
     - ./dynamic:/etc/traefik/dynamic:ro
   ```

## Security Considerations

When integrating Chadburn with Traefik, consider these security best practices:

1. Use read-only mounts for configuration files
2. Limit Docker socket access with proper permissions
3. Use secure API settings in production
4. Regularly rotate backup files to prevent disk space issues
5. Use environment variables for sensitive information instead of hardcoding in configuration files 