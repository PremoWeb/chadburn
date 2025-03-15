# Deployment Guide

This guide covers various deployment scenarios for Chadburn in different environments.

## Docker Compose

Docker Compose is the simplest way to deploy Chadburn alongside your application containers.

### Basic Deployment

Create a `docker-compose.yml` file:

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
    environment:
      - CHADBURN_LOG_LEVEL=info
```

Create a `chadburn.conf` file:

```ini
[job-local "example"]
schedule = @every 1m
command = echo "Hello from Chadburn!"
```

Start the services:

```bash
docker-compose up -d
```

### With Application Containers

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
    depends_on:
      - app
      - db

  app:
    image: your-app:latest
    container_name: app
    labels:
      chadburn.enabled: "true"
      chadburn.job-exec.backup.schedule: "@daily"
      chadburn.job-exec.backup.command: "tar -czf /backups/app-data.tar.gz /app/data"
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

### With Metrics

```yaml
version: "3"
services:
  chadburn:
    image: premoweb/chadburn:latest
    container_name: chadburn
    restart: unless-stopped
    command: daemon --config=/etc/chadburn.conf --metrics
    ports:
      - "8080:8080"
    volumes:
      - /var/run/docker.sock:/var/run/docker.sock:ro,z
      - ./chadburn.conf:/etc/chadburn.conf

  prometheus:
    image: prom/prometheus:latest
    container_name: prometheus
    volumes:
      - ./prometheus.yml:/etc/prometheus/prometheus.yml
    ports:
      - "9090:9090"

  grafana:
    image: grafana/grafana:latest
    container_name: grafana
    ports:
      - "3000:3000"
    depends_on:
      - prometheus
```

Create a `prometheus.yml` file:

```yaml
global:
  scrape_interval: 15s

scrape_configs:
  - job_name: 'chadburn'
    static_configs:
      - targets: ['chadburn:8080']
```

## Kubernetes

Deploying Chadburn in Kubernetes requires special consideration for accessing the Docker socket.

### Using the Host Docker Socket

Create a `chadburn.yaml` file:

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: chadburn-config
data:
  chadburn.conf: |
    [job-local "example"]
    schedule = @every 1m
    command = echo "Hello from Chadburn!"
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: chadburn
spec:
  replicas: 1
  selector:
    matchLabels:
      app: chadburn
  template:
    metadata:
      labels:
        app: chadburn
    spec:
      containers:
      - name: chadburn
        image: premoweb/chadburn:latest
        args: ["daemon", "--config=/etc/chadburn/chadburn.conf"]
        volumeMounts:
        - name: docker-socket
          mountPath: /var/run/docker.sock
        - name: config
          mountPath: /etc/chadburn
      volumes:
      - name: docker-socket
        hostPath:
          path: /var/run/docker.sock
          type: Socket
      - name: config
        configMap:
          name: chadburn-config
```

Apply the configuration:

```bash
kubectl apply -f chadburn.yaml
```

### Using DaemonSet

For a more Kubernetes-native approach, deploy Chadburn as a DaemonSet:

```yaml
apiVersion: apps/v1
kind: DaemonSet
metadata:
  name: chadburn
spec:
  selector:
    matchLabels:
      app: chadburn
  template:
    metadata:
      labels:
        app: chadburn
    spec:
      containers:
      - name: chadburn
        image: premoweb/chadburn:latest
        args: ["daemon", "--config=/etc/chadburn/chadburn.conf"]
        volumeMounts:
        - name: docker-socket
          mountPath: /var/run/docker.sock
        - name: config
          mountPath: /etc/chadburn
      volumes:
      - name: docker-socket
        hostPath:
          path: /var/run/docker.sock
          type: Socket
      - name: config
        configMap:
          name: chadburn-config
```

### Limitations in Kubernetes

When running in Kubernetes:

1. `job-exec` can only access containers on the same node
2. `job-local` runs on the node, not in the Kubernetes context
3. `job-service-run` is not applicable (use Kubernetes Jobs instead)

## Docker Swarm

Chadburn works well in Docker Swarm mode, especially with the `job-service-run` job type.

### Basic Swarm Deployment

```yaml
version: "3.8"
services:
  chadburn:
    image: premoweb/chadburn:latest
    command: daemon --config=/etc/chadburn.conf
    volumes:
      - /var/run/docker.sock:/var/run/docker.sock:ro,z
      - ./chadburn.conf:/etc/chadburn.conf
    deploy:
      mode: replicated
      replicas: 1
      placement:
        constraints:
          - node.role == manager
    networks:
      - backend

networks:
  backend:
```

### Using job-service-run

The `job-service-run` job type is specifically designed for Swarm mode:

```ini
[job-service-run "backup-service"]
schedule = @daily
image = backup:latest
network = backend
command = /backup.sh
```

This creates a new service for each job execution, which is then removed after completion.

## Production Considerations

### High Availability

For high availability:

1. Use Docker Swarm with placement constraints to ensure Chadburn runs on a manager node
2. Use a shared volume for configuration and logs
3. Set up monitoring for the Chadburn container

### Resource Limits

Set appropriate resource limits:

```yaml
services:
  chadburn:
    image: premoweb/chadburn:latest
    deploy:
      resources:
        limits:
          cpus: '0.5'
          memory: 256M
        reservations:
          cpus: '0.1'
          memory: 128M
```

### Logging

Configure logging for production:

```yaml
services:
  chadburn:
    image: premoweb/chadburn:latest
    logging:
      driver: "json-file"
      options:
        max-size: "10m"
        max-file: "3"
```

### Security

Secure your Chadburn deployment:

1. Run with read-only access to the Docker socket
2. Use a non-root user when possible
3. Store sensitive configuration in secrets
4. Limit network access to the metrics endpoint

## Upgrading

To upgrade Chadburn:

1. Pull the latest image:
   ```bash
   docker pull premoweb/chadburn:latest
   ```

2. Update your deployment:
   ```bash
   docker-compose up -d
   ```

3. Check the logs for any issues:
   ```bash
   docker-compose logs chadburn
   ```

## Backup and Restore

### Backing Up Configuration

Regularly back up your configuration files:

```bash
cp chadburn.conf chadburn.conf.backup
```

### Backing Up Logs

If you're using the save middleware, back up the log directory:

```bash
tar -czf chadburn-logs-$(date +%Y%m%d).tar.gz /var/log/chadburn
```

### Disaster Recovery

In case of a failure:

1. Restore your configuration files
2. Redeploy Chadburn
3. Verify that jobs are running correctly 