---
layout: default
title: Docker Socket Permissions
parent: Features
nav_order: 4
---

# Docker Socket Permission Issues

This guide addresses common permission issues when accessing the Docker socket (`/var/run/docker.sock`) from within containers, particularly in environments with SELinux enabled.

## Problem Description

When running Chadburn in a container that needs to communicate with the Docker daemon on the host, you might encounter permission errors like:

```
ERROR Docker events error: permission denied while trying to connect to the Docker daemon socket at unix:///var/run/docker.sock: Get "http://%2Fvar%2Frun%2Fdocker.sock/v1.48/events": dial unix /var/run/docker.sock: connect: permission denied
```

This happens because:

1. The Docker socket is owned by the `root` user and the `docker` group on the host
2. The container user doesn't have the same permissions
3. SELinux may be enforcing additional security policies

## Solutions

Here are several approaches to solve this issue, from most to least secure:

### 1. Use SELinux Volume Labels (Recommended)

Add the `:z` or `:Z` suffix to your volume mount in `docker-compose.yml`:

```yaml
volumes:
  - /var/run/docker.sock:/var/run/docker.sock:ro,z
```

- `:z` - Tells SELinux to relabel the content with a shared label
- `:Z` - Tells SELinux to relabel with a private unshared label

### 2. Match the Docker Group ID

Ensure the container runs with the same group ID as the Docker group on the host:

```yaml
services:
  scheduler:
    # ... other configuration
    user: "root:${DOCKER_GID:-999}"
    environment:
      - DOCKER_GID=${DOCKER_GID:-999}
```

When starting the container:

```bash
# Get the Docker GID
DOCKER_GID=$(getent group docker | cut -d: -f3)

# Start with the correct GID
DOCKER_GID=$DOCKER_GID docker-compose up -d
```

### 3. Modify Dockerfile to Add Docker Group

```dockerfile
# Add docker group with same GID as host
RUN addgroup -S -g 969 docker && \
    adduser -S -D -H -h /app -s /sbin/nologin -G docker -u 1000 appuser

# Use the appuser
USER appuser
```

Replace `969` with your host's Docker group ID.

### 4. Run as Root User (Less Secure)

```yaml
services:
  scheduler:
    # ... other configuration
    user: "root"
```

### 5. Change Docker Socket Permissions (Not Recommended for Production)

```bash
sudo chmod 666 /var/run/docker.sock
```

This makes the socket readable and writable by all users, which is a security risk.

### 6. Temporarily Disable SELinux (For Testing Only)

```bash
sudo setenforce 0
```

To re-enable:

```bash
sudo setenforce 1
```

**Note**: This is temporary and will reset after a system reboot.

### 7. Create a Custom SELinux Policy (Advanced)

For production environments, create a custom SELinux policy:

```bash
# Create a policy module
ausearch -c 'docker' --raw | audit2allow -M my-docker

# Install the policy
semodule -i my-docker.pp
```

## Best Practices

1. **Principle of Least Privilege**: Only grant the minimum permissions needed
2. **Use Read-Only Access**: Use `:ro` for the socket mount when possible
3. **Avoid Disabling SELinux**: Work with SELinux rather than disabling it
4. **Container User**: Avoid running as root when possible
5. **Regular Audits**: Regularly review your permission settings

## Troubleshooting

### Check Docker Socket Ownership

```bash
ls -la /var/run/docker.sock
```

### Check Docker Group ID

```bash
getent group docker
```

### Check SELinux Status

```bash
getenforce
```

### Check SELinux Contexts

```bash
ls -Z /var/run/docker.sock
```

### Check Container Logs

```bash
docker-compose logs scheduler
```

## Conclusion

When dealing with Docker socket permissions, always prioritize security. The recommended approach is to use SELinux volume labels (`:z` or `:Z`) combined with proper user/group mapping. For production environments, consider creating a custom SELinux policy for more granular control. 