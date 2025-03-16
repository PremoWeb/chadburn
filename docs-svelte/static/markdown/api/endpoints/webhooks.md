# Webhooks API Endpoints

The Webhooks API allows you to create, manage, and monitor webhooks that are triggered by events in Chadburn. Webhooks provide a way to integrate Chadburn with other systems by sending HTTP notifications when specific events occur.

## List All Webhooks

```
GET /api/v1/webhooks
```

Retrieves a list of all webhooks.

### Query Parameters

| Parameter  | Type    | Description                                                |
|------------|---------|-----------------------------------------------------------|
| `page`     | integer | Page number for pagination (default: 1)                    |
| `per_page` | integer | Number of webhooks per page (default: 20, max: 100)        |
| `status`   | string  | Filter by status: 'active', 'inactive' (optional)          |
| `event`    | string  | Filter by event type (optional)                            |

### Response

```json
{
  "data": [
    {
      "id": "webhook-123",
      "name": "Slack Notification",
      "url": "https://hooks.slack.com/services/T00000000/B00000000/XXXXXXXXXXXXXXXXXXXXXXXX",
      "events": ["job.success", "job.failure"],
      "status": "active",
      "created_at": "2023-01-01T00:00:00Z",
      "updated_at": "2023-01-02T00:00:00Z",
      "secret": "••••••••••••••••",
      "headers": {
        "Content-Type": "application/json"
      },
      "retry_config": {
        "max_retries": 3,
        "retry_interval": 60
      }
    },
    {
      "id": "webhook-456",
      "name": "Monitoring System",
      "url": "https://api.monitoring.example.com/webhooks/chadburn",
      "events": ["job.failure", "system.error"],
      "status": "active",
      "created_at": "2023-01-01T12:00:00Z",
      "updated_at": "2023-01-01T12:00:00Z",
      "secret": "••••••••••••••••",
      "headers": {
        "Content-Type": "application/json",
        "X-API-Key": "••••••••••••••••"
      },
      "retry_config": {
        "max_retries": 5,
        "retry_interval": 30
      }
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

## Get a Specific Webhook

```
GET /api/v1/webhooks/{id}
```

Retrieves details for a specific webhook.

### Path Parameters

| Parameter | Type   | Description |
|-----------|--------|-------------|
| `id`      | string | Webhook ID  |

### Response

```json
{
  "data": {
    "id": "webhook-123",
    "name": "Slack Notification",
    "url": "https://hooks.slack.com/services/T00000000/B00000000/XXXXXXXXXXXXXXXXXXXXXXXX",
    "events": ["job.success", "job.failure"],
    "status": "active",
    "created_at": "2023-01-01T00:00:00Z",
    "updated_at": "2023-01-02T00:00:00Z",
    "secret": "••••••••••••••••",
    "headers": {
      "Content-Type": "application/json"
    },
    "retry_config": {
      "max_retries": 3,
      "retry_interval": 60
    },
    "filters": {
      "job_ids": ["job-123", "job-456"],
      "job_tags": ["production", "critical"]
    }
  }
}
```

## Create a New Webhook

```
POST /api/v1/webhooks
```

Creates a new webhook.

### Request Body

```json
{
  "name": "New Webhook",
  "url": "https://api.example.com/webhooks/chadburn",
  "events": ["job.success", "job.failure", "job.started"],
  "status": "active",
  "secret": "your-webhook-secret",
  "headers": {
    "Content-Type": "application/json",
    "X-API-Key": "your-api-key"
  },
  "retry_config": {
    "max_retries": 3,
    "retry_interval": 60
  },
  "filters": {
    "job_ids": ["job-123", "job-456"],
    "job_tags": ["production"]
  }
}
```

### Response

```json
{
  "data": {
    "id": "webhook-789",
    "name": "New Webhook",
    "url": "https://api.example.com/webhooks/chadburn",
    "events": ["job.success", "job.failure", "job.started"],
    "status": "active",
    "created_at": "2023-01-03T00:00:00Z",
    "updated_at": "2023-01-03T00:00:00Z",
    "secret": "••••••••••••••••",
    "headers": {
      "Content-Type": "application/json",
      "X-API-Key": "••••••••••••••••"
    },
    "retry_config": {
      "max_retries": 3,
      "retry_interval": 60
    },
    "filters": {
      "job_ids": ["job-123", "job-456"],
      "job_tags": ["production"]
    }
  }
}
```

## Update a Webhook

```
PUT /api/v1/webhooks/{id}
```

Updates an existing webhook.

### Path Parameters

| Parameter | Type   | Description |
|-----------|--------|-------------|
| `id`      | string | Webhook ID  |

### Request Body

```json
{
  "name": "Updated Webhook Name",
  "url": "https://api.example.com/webhooks/chadburn-updated",
  "events": ["job.success", "job.failure", "job.started", "system.error"],
  "status": "active",
  "headers": {
    "Content-Type": "application/json",
    "X-API-Key": "new-api-key"
  },
  "retry_config": {
    "max_retries": 5,
    "retry_interval": 30
  },
  "filters": {
    "job_ids": ["job-123", "job-456", "job-789"],
    "job_tags": ["production", "critical"]
  }
}
```

### Response

```json
{
  "data": {
    "id": "webhook-123",
    "name": "Updated Webhook Name",
    "url": "https://api.example.com/webhooks/chadburn-updated",
    "events": ["job.success", "job.failure", "job.started", "system.error"],
    "status": "active",
    "created_at": "2023-01-01T00:00:00Z",
    "updated_at": "2023-01-03T12:00:00Z",
    "secret": "••••••••••••••••",
    "headers": {
      "Content-Type": "application/json",
      "X-API-Key": "••••••••••••••••"
    },
    "retry_config": {
      "max_retries": 5,
      "retry_interval": 30
    },
    "filters": {
      "job_ids": ["job-123", "job-456", "job-789"],
      "job_tags": ["production", "critical"]
    }
  }
}
```

## Delete a Webhook

```
DELETE /api/v1/webhooks/{id}
```

Deletes a webhook.

### Path Parameters

| Parameter | Type   | Description |
|-----------|--------|-------------|
| `id`      | string | Webhook ID  |

### Response

```json
{
  "data": {
    "message": "Webhook deleted successfully"
  }
}
```

## Get Webhook Delivery History

```
GET /api/v1/webhooks/{id}/deliveries
```

Retrieves the delivery history for a specific webhook.

### Path Parameters

| Parameter | Type   | Description |
|-----------|--------|-------------|
| `id`      | string | Webhook ID  |

### Query Parameters

| Parameter  | Type    | Description                                                |
|------------|---------|-----------------------------------------------------------|
| `page`     | integer | Page number for pagination (default: 1)                    |
| `per_page` | integer | Number of deliveries per page (default: 20, max: 100)      |
| `status`   | string  | Filter by status: 'success', 'failure' (optional)          |
| `event`    | string  | Filter by event type (optional)                            |
| `start`    | string  | Start time in ISO 8601 format (optional)                   |
| `end`      | string  | End time in ISO 8601 format (optional)                     |

### Response

```json
{
  "data": [
    {
      "id": "delivery-123",
      "webhook_id": "webhook-123",
      "event": "job.success",
      "status": "success",
      "request": {
        "url": "https://hooks.slack.com/services/T00000000/B00000000/XXXXXXXXXXXXXXXXXXXXXXXX",
        "method": "POST",
        "headers": {
          "Content-Type": "application/json",
          "X-Chadburn-Signature": "sha256=..."
        },
        "body": "{ ... }"
      },
      "response": {
        "status_code": 200,
        "headers": {
          "Content-Type": "application/json"
        },
        "body": "{ \"ok\": true }"
      },
      "created_at": "2023-01-02T12:34:56Z",
      "duration_ms": 245
    },
    {
      "id": "delivery-456",
      "webhook_id": "webhook-123",
      "event": "job.failure",
      "status": "failure",
      "request": {
        "url": "https://hooks.slack.com/services/T00000000/B00000000/XXXXXXXXXXXXXXXXXXXXXXXX",
        "method": "POST",
        "headers": {
          "Content-Type": "application/json",
          "X-Chadburn-Signature": "sha256=..."
        },
        "body": "{ ... }"
      },
      "response": {
        "status_code": 500,
        "headers": {
          "Content-Type": "application/json"
        },
        "body": "{ \"ok\": false, \"error\": \"server_error\" }"
      },
      "created_at": "2023-01-02T13:45:67Z",
      "duration_ms": 356,
      "retry_count": 2,
      "next_retry_at": "2023-01-02T13:55:67Z"
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

## Retry a Failed Webhook Delivery

```
POST /api/v1/webhooks/deliveries/{id}/retry
```

Retries a failed webhook delivery.

### Path Parameters

| Parameter | Type   | Description        |
|-----------|--------|--------------------|
| `id`      | string | Delivery ID        |

### Response

```json
{
  "data": {
    "id": "delivery-789",
    "webhook_id": "webhook-123",
    "event": "job.failure",
    "status": "pending",
    "original_delivery_id": "delivery-456",
    "created_at": "2023-01-02T14:00:00Z"
  }
}
```

## Available Webhook Events

The following events can be subscribed to via webhooks:

| Event Name           | Description                                           |
|----------------------|-------------------------------------------------------|
| `job.created`        | Triggered when a new job is created                   |
| `job.updated`        | Triggered when a job is updated                       |
| `job.deleted`        | Triggered when a job is deleted                       |
| `job.started`        | Triggered when a job execution starts                 |
| `job.success`        | Triggered when a job execution completes successfully |
| `job.failure`        | Triggered when a job execution fails                  |
| `job.paused`         | Triggered when a job is paused                        |
| `job.resumed`        | Triggered when a job is resumed                       |
| `schedule.created`   | Triggered when a new schedule is created              |
| `schedule.updated`   | Triggered when a schedule is updated                  |
| `schedule.deleted`   | Triggered when a schedule is deleted                  |
| `system.error`       | Triggered when a system error occurs                  |
| `system.started`     | Triggered when the Chadburn system starts             |
| `system.shutdown`    | Triggered when the Chadburn system shuts down         |

## Webhook Payload Format

Each webhook delivery includes a JSON payload with information about the event. The payload format varies depending on the event type, but all payloads include the following common fields:

```json
{
  "id": "event-123",
  "type": "job.success",
  "created_at": "2023-01-02T12:34:56Z",
  "data": {
    // Event-specific data
  }
}
```

### Example: Job Success Event

```json
{
  "id": "event-123",
  "type": "job.success",
  "created_at": "2023-01-02T12:34:56Z",
  "data": {
    "job": {
      "id": "job-123",
      "name": "Database Backup",
      "command": "pg_dump -U postgres -d mydb > /backups/mydb.sql",
      "schedule": "0 0 * * *"
    },
    "execution": {
      "id": "exec-456",
      "started_at": "2023-01-02T12:30:00Z",
      "finished_at": "2023-01-02T12:34:56Z",
      "duration": 296,
      "exit_code": 0,
      "output": "Database backup completed successfully"
    }
  }
}
```

## Webhook Security

Chadburn signs all webhook requests with a signature in the `X-Chadburn-Signature` header. The signature is a HMAC SHA-256 hash of the request body, using the webhook's secret as the key.

To verify the signature:

1. Get the signature from the `X-Chadburn-Signature` header (format: `sha256=<signature>`)
2. Compute an HMAC with the SHA256 hash function
3. Use the webhook secret as the HMAC key
4. Use the request body as the HMAC message
5. Compare the computed signature with the signature in the header

Example verification in Node.js:

```javascript
const crypto = require('crypto');

function verifySignature(body, signature, secret) {
  const hmac = crypto.createHmac('sha256', secret);
  const digest = hmac.update(body).digest('hex');
  return crypto.timingSafeEqual(
    Buffer.from(signature),
    Buffer.from(`sha256=${digest}`)
  );
}

// Usage
const isValid = verifySignature(
  requestBody,
  request.headers['x-chadburn-signature'],
  webhookSecret
);
``` 