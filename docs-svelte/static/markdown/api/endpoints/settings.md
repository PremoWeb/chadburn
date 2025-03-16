# Settings API Endpoints

The Settings API allows you to manage system-wide configuration settings for Chadburn. This includes general settings, notification settings, security settings, and more.

## Get All Settings

```
GET /api/v1/settings
```

Retrieves all system settings.

### Response

```json
{
  "data": {
    "general": {
      "instance_name": "Production Scheduler",
      "default_timezone": "UTC",
      "date_format": "YYYY-MM-DD",
      "time_format": "HH:mm:ss",
      "log_level": "info",
      "max_log_retention_days": 30
    },
    "security": {
      "session_timeout_minutes": 60,
      "max_login_attempts": 5,
      "lockout_duration_minutes": 15,
      "password_policy": {
        "min_length": 10,
        "require_uppercase": true,
        "require_lowercase": true,
        "require_numbers": true,
        "require_special_chars": true,
        "max_age_days": 90
      },
      "two_factor_auth": {
        "enabled": true,
        "required_for_admins": true
      }
    },
    "notifications": {
      "email": {
        "enabled": true,
        "from_address": "chadburn@example.com",
        "smtp_server": "smtp.example.com",
        "smtp_port": 587,
        "smtp_username": "••••••••••••••••",
        "smtp_use_tls": true
      },
      "slack": {
        "enabled": true,
        "webhook_url": "https://hooks.slack.com/services/••••••••••••••••",
        "default_channel": "#chadburn-alerts"
      },
      "webhooks": {
        "enabled": true,
        "max_retries": 3,
        "retry_interval_seconds": 60
      }
    },
    "execution": {
      "max_concurrent_jobs": 10,
      "job_timeout_seconds": 3600,
      "default_retry_attempts": 3,
      "default_retry_delay_seconds": 60,
      "allow_overlapping_jobs": false
    },
    "ui": {
      "theme": "light",
      "items_per_page": 20,
      "dashboard_refresh_interval_seconds": 30
    },
    "storage": {
      "log_storage_path": "/var/log/chadburn",
      "max_log_size_mb": 100,
      "backup_enabled": true,
      "backup_interval_hours": 24,
      "backup_retention_count": 7
    },
    "metrics": {
      "enabled": true,
      "prometheus_endpoint_enabled": true,
      "collection_interval_seconds": 60
    }
  }
}
```

## Get a Specific Setting Category

```
GET /api/v1/settings/{category}
```

Retrieves settings for a specific category.

### Path Parameters

| Parameter | Type   | Description                                                                                |
|-----------|--------|--------------------------------------------------------------------------------------------|
| `category` | string | Setting category: 'general', 'security', 'notifications', 'execution', 'ui', 'storage', 'metrics' |

### Response

```json
{
  "data": {
    "email": {
      "enabled": true,
      "from_address": "chadburn@example.com",
      "smtp_server": "smtp.example.com",
      "smtp_port": 587,
      "smtp_username": "••••••••••••••••",
      "smtp_use_tls": true
    },
    "slack": {
      "enabled": true,
      "webhook_url": "https://hooks.slack.com/services/••••••••••••••••",
      "default_channel": "#chadburn-alerts"
    },
    "webhooks": {
      "enabled": true,
      "max_retries": 3,
      "retry_interval_seconds": 60
    }
  }
}
```

## Update Settings

```
PUT /api/v1/settings/{category}
```

Updates settings for a specific category.

### Path Parameters

| Parameter | Type   | Description                                                                                |
|-----------|--------|--------------------------------------------------------------------------------------------|
| `category` | string | Setting category: 'general', 'security', 'notifications', 'execution', 'ui', 'storage', 'metrics' |

### Request Body

```json
{
  "email": {
    "enabled": true,
    "from_address": "new-email@example.com",
    "smtp_server": "smtp.newserver.com",
    "smtp_port": 465,
    "smtp_username": "new-username",
    "smtp_password": "new-password",
    "smtp_use_tls": true
  },
  "slack": {
    "enabled": true,
    "webhook_url": "https://hooks.slack.com/services/NEW_WEBHOOK_URL",
    "default_channel": "#new-channel"
  }
}
```

### Response

```json
{
  "data": {
    "email": {
      "enabled": true,
      "from_address": "new-email@example.com",
      "smtp_server": "smtp.newserver.com",
      "smtp_port": 465,
      "smtp_username": "••••••••••••••••",
      "smtp_use_tls": true
    },
    "slack": {
      "enabled": true,
      "webhook_url": "https://hooks.slack.com/services/••••••••••••••••",
      "default_channel": "#new-channel"
    },
    "webhooks": {
      "enabled": true,
      "max_retries": 3,
      "retry_interval_seconds": 60
    }
  }
}
```

## Test Email Settings

```
POST /api/v1/settings/notifications/email/test
```

Tests the email notification settings by sending a test email.

### Request Body

```json
{
  "recipient": "test@example.com"
}
```

### Response

```json
{
  "data": {
    "success": true,
    "message": "Test email sent successfully to test@example.com"
  }
}
```

## Test Slack Settings

```
POST /api/v1/settings/notifications/slack/test
```

Tests the Slack notification settings by sending a test message.

### Request Body

```json
{
  "channel": "#test-channel"
}
```

### Response

```json
{
  "data": {
    "success": true,
    "message": "Test message sent successfully to #test-channel"
  }
}
```

## Get System Information

```
GET /api/v1/settings/system-info
```

Retrieves system information.

### Response

```json
{
  "data": {
    "version": "1.9.0",
    "build_date": "2025-03-16T02:31:18Z",
    "uptime": 604800,
    "system": {
      "os": "Linux",
      "os_version": "5.15.0-1019-aws",
      "architecture": "x86_64",
      "cpu_cores": 4,
      "memory_total_mb": 8192,
      "memory_used_mb": 2048,
      "disk_total_gb": 100,
      "disk_used_gb": 45
    },
    "database": {
      "type": "PostgreSQL",
      "version": "14.5",
      "connection_pool_size": 10,
      "active_connections": 3
    },
    "jobs": {
      "total": 42,
      "active": 38,
      "paused": 4,
      "currently_running": 2
    }
  }
}
```

## Get License Information

```
GET /api/v1/settings/license
```

Retrieves license information.

### Response

```json
{
  "data": {
    "type": "enterprise",
    "status": "active",
    "issued_to": "Example Corp",
    "issued_at": "2023-01-01T00:00:00Z",
    "expires_at": "2024-01-01T00:00:00Z",
    "features": [
      "unlimited_jobs",
      "advanced_scheduling",
      "user_management",
      "audit_logs",
      "sla_monitoring"
    ],
    "restrictions": {
      "max_users": 50,
      "max_jobs": null
    }
  }
}
```

## Update License

```
PUT /api/v1/settings/license
```

Updates the license.

### Request Body

```json
{
  "license_key": "XXXX-XXXX-XXXX-XXXX-XXXX"
}
```

### Response

```json
{
  "data": {
    "type": "enterprise",
    "status": "active",
    "issued_to": "Example Corp",
    "issued_at": "2023-01-01T00:00:00Z",
    "expires_at": "2025-01-01T00:00:00Z",
    "features": [
      "unlimited_jobs",
      "advanced_scheduling",
      "user_management",
      "audit_logs",
      "sla_monitoring",
      "high_availability"
    ],
    "restrictions": {
      "max_users": 100,
      "max_jobs": null
    }
  }
}
```

## Get Audit Logs

```
GET /api/v1/settings/audit-logs
```

Retrieves audit logs for system settings changes.

### Query Parameters

| Parameter  | Type    | Description                                                |
|------------|---------|-----------------------------------------------------------|
| `page`     | integer | Page number for pagination (default: 1)                    |
| `per_page` | integer | Number of logs per page (default: 20, max: 100)            |
| `start`    | string  | Start time in ISO 8601 format (optional)                   |
| `end`      | string  | End time in ISO 8601 format (optional)                     |
| `user_id`  | string  | Filter by user ID (optional)                               |
| `category` | string  | Filter by setting category (optional)                      |

### Response

```json
{
  "data": [
    {
      "id": "log-123",
      "user_id": "user-123",
      "user_name": "John Doe",
      "action": "update",
      "category": "notifications",
      "setting": "email.smtp_server",
      "old_value": "smtp.example.com",
      "new_value": "smtp.newserver.com",
      "ip_address": "192.168.1.1",
      "timestamp": "2023-01-03T12:34:56Z"
    },
    {
      "id": "log-456",
      "user_id": "user-123",
      "user_name": "John Doe",
      "action": "update",
      "category": "security",
      "setting": "password_policy.min_length",
      "old_value": "8",
      "new_value": "10",
      "ip_address": "192.168.1.1",
      "timestamp": "2023-01-02T09:12:34Z"
    }
  ],
  "meta": {
    "page": 1,
    "per_page": 20,
    "total": 2,
    "total_pages": 1
  }
}
```

## Backup Settings

```
POST /api/v1/settings/backup
```

Creates a backup of all system settings.

### Response

```json
{
  "data": {
    "backup_id": "backup-123",
    "timestamp": "2023-01-04T00:00:00Z",
    "size_bytes": 15360,
    "download_url": "/api/v1/settings/backup/backup-123/download",
    "expires_at": "2023-01-11T00:00:00Z"
  }
}
```

## Restore Settings

```
POST /api/v1/settings/restore
```

Restores system settings from a backup.

### Request Body

```json
{
  "backup_id": "backup-123"
}
```

### Response

```json
{
  "data": {
    "success": true,
    "message": "Settings restored successfully from backup-123",
    "timestamp": "2023-01-04T12:34:56Z",
    "categories_restored": [
      "general",
      "security",
      "notifications",
      "execution",
      "ui",
      "storage",
      "metrics"
    ]
  }
}
```

## Reset Settings to Default

```
POST /api/v1/settings/reset
```

Resets all system settings to their default values.

### Request Body

```json
{
  "categories": ["ui", "notifications"],
  "confirm": true
}
```

### Response

```json
{
  "data": {
    "success": true,
    "message": "Settings reset successfully",
    "categories_reset": [
      "ui",
      "notifications"
    ]
  }
}
``` 