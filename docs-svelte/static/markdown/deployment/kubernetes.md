# Kubernetes Deployment

This guide covers how to deploy Chadburn in a Kubernetes environment.

## Prerequisites

- A Kubernetes cluster (v1.16+)
- `kubectl` configured to communicate with your cluster
- Basic understanding of Kubernetes concepts

## Kubernetes and Container Runtimes

> **Note:** As of Kubernetes 1.24 (released in May 2022), Docker Engine support has been removed from Kubernetes. Kubernetes now uses the Container Runtime Interface (CRI) with containerd as the default runtime.

Kubernetes originally used Docker as its container runtime, but has since transitioned to using containerd directly. This change is transparent to users, and Docker-built images continue to work perfectly in Kubernetes. The `docker build` command creates OCI (Open Container Initiative) compliant images that are compatible with all container runtimes that Kubernetes supports.

If you're using Docker to build your Chadburn images, they will work without any modifications in your Kubernetes cluster.

## Deployment Steps

### 1. Create a Namespace

First, create a dedicated namespace for Chadburn:

```bash
kubectl create namespace chadburn
```

### 2. Create ConfigMap

Create a ConfigMap for Chadburn configuration:

```bash
cat <<EOF | kubectl apply -f -
apiVersion: v1
kind: ConfigMap
metadata:
  name: chadburn-config
  namespace: chadburn
data:
  config.yml: |
    timezone: UTC
    schedule:
      - name: "Example job"
        command: "echo 'Hello from Kubernetes!'"
        schedule: "*/5 * * * *"
EOF
```

### 3. Create Deployment

Create a Deployment for Chadburn:

```bash
cat <<EOF | kubectl apply -f -
apiVersion: apps/v1
kind: Deployment
metadata:
  name: chadburn
  namespace: chadburn
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
        image: ghcr.io/artefactual-labs/chadburn:latest
        resources:
          limits:
            cpu: "0.5"
            memory: "512Mi"
          requests:
            cpu: "0.2"
            memory: "256Mi"
        volumeMounts:
        - name: config-volume
          mountPath: /etc/chadburn
      volumes:
      - name: config-volume
        configMap:
          name: chadburn-config
EOF
```

### 4. Verify Deployment

Check if the deployment was successful:

```bash
kubectl get pods -n chadburn
```

You should see a running pod for Chadburn.

## Advanced Configuration

### Using Secrets

For sensitive information, you can use Kubernetes Secrets:

```bash
kubectl create secret generic chadburn-secrets \
  --namespace chadburn \
  --from-literal=API_KEY=your-api-key
```

Then mount the secret in your deployment:

```yaml
volumeMounts:
- name: secrets-volume
  mountPath: /etc/chadburn/secrets
volumes:
- name: secrets-volume
  secret:
    secretName: chadburn-secrets
```

### Persistent Storage

If you need persistent storage for logs or data:

```yaml
volumeMounts:
- name: data-volume
  mountPath: /var/lib/chadburn
volumes:
- name: data-volume
  persistentVolumeClaim:
    claimName: chadburn-pvc
```

## Container Runtime Considerations

If your Chadburn workload needs to interact with a container runtime (for example, to spawn containers as part of its jobs), you'll need to consider the following:

1. **Using containerd**: For newer Kubernetes clusters, you should use the containerd API directly or a tool like `nerdctl` that works with containerd.

2. **Using Docker-in-Docker (DinD)**: If you specifically need Docker, you can use a Docker-in-Docker approach by mounting the Docker socket or running a Docker daemon in a sidecar container.

Example of mounting the Docker socket (if your node still has Docker installed):

```yaml
volumeMounts:
- name: docker-socket
  mountPath: /var/run/docker.sock
volumes:
- name: docker-socket
  hostPath:
    path: /var/run/docker.sock
    type: Socket
```

## Troubleshooting

If you encounter issues with your Chadburn deployment:

1. Check pod logs:
   ```bash
   kubectl logs -n chadburn <pod-name>
   ```

2. Verify the ConfigMap is correctly mounted:
   ```bash
   kubectl describe pod -n chadburn <pod-name>
   ```

3. Check for events in the namespace:
   ```bash
   kubectl get events -n chadburn
   ```

For more detailed troubleshooting, refer to the [Troubleshooting](/troubleshooting) guide.
