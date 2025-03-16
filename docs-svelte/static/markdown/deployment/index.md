---
layout: default
title: Deployment
nav_order: 4
has_children: true
---

# Chadburn Deployment

Chadburn can be deployed in various environments and configurations to suit your needs.

## Deployment Options

Chadburn can be deployed in several ways:

- **Docker**: Run as a standalone container
- **Docker Compose**: Deploy alongside your application stack
- **Kubernetes**: Run as a DaemonSet or Deployment
- **Swarm Mode**: Deploy as a global service

## Docker Integrations

Chadburn works seamlessly with popular Docker-based services:

- [Docker Deployment](/deployment/docker): Basic Docker deployment guide
- [Caddy Integration](/deployment/docker#integration-with-caddy): Schedule maintenance tasks for Caddy web server
- [Traefik Integration](/deployment/docker#integration-with-traefik): Automate maintenance for Traefik proxy

## Kubernetes Deployment

For container orchestration with Kubernetes:

- [Kubernetes Deployment](/deployment/kubernetes): Deploy Chadburn in Kubernetes
- [Container Runtime Considerations](/deployment/kubernetes#container-runtime-considerations): Working with containerd and Docker in Kubernetes

## Deployment Documentation

Explore the detailed documentation for deployment:

- [Deployment Options](deployment.html): Learn about different ways to deploy Chadburn
- [Advanced Topics](advanced-topics.html): Explore advanced deployment configurations and scenarios 