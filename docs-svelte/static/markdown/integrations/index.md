# Integrations

Chadburn can be integrated with various third-party services and tools to enhance its functionality and provide automated maintenance for your infrastructure. This section covers the available integrations and how to configure them.

## Available Integrations

Chadburn currently offers the following integrations:

### Web Servers

- [Traefik](traefik.md) - Integrate Chadburn with Traefik for automated maintenance tasks
- [Caddy](caddy.md) - Integrate Chadburn with Caddy for automated maintenance tasks

## Integration Benefits

Integrating Chadburn with other services provides several benefits:

1. **Automated Maintenance**: Schedule routine maintenance tasks to run automatically
2. **Centralized Management**: Manage scheduled tasks for multiple services from a single configuration
3. **Consistent Approach**: Use the same scheduling mechanism across your entire infrastructure
4. **Docker-Native**: Leverage Docker labels for easy configuration in containerized environments
5. **Flexible Scheduling**: Use cron expressions or interval notation for precise scheduling

## General Integration Approach

Most Chadburn integrations follow a similar pattern:

1. **Label-Based Configuration**: Add Chadburn labels to your containers
2. **Service-Specific Commands**: Configure commands specific to the service you're integrating with
3. **Scheduling**: Set up appropriate schedules for different maintenance tasks
4. **Monitoring**: Optionally configure health checks and monitoring

For specific integration instructions, refer to the individual integration guides linked above. 