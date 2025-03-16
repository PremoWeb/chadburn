# Users API Endpoints

The Users API allows you to manage user accounts in Chadburn. This includes creating, updating, and deleting users, as well as managing their roles and permissions.

## List All Users

```
GET /api/v1/users
```

Retrieves a list of all users.

### Query Parameters

| Parameter  | Type    | Description                                                |
|------------|---------|-----------------------------------------------------------|
| `page`     | integer | Page number for pagination (default: 1)                    |
| `per_page` | integer | Number of users per page (default: 20, max: 100)           |
| `sort`     | string  | Sort by: 'name', 'email', 'created_at' (default: 'name')   |
| `order`    | string  | Sort order: 'asc', 'desc' (default: 'asc')                 |
| `role`     | string  | Filter by role: 'admin', 'user', 'viewer' (optional)       |
| `search`   | string  | Search term to filter users by name or email (optional)    |

### Response

```json
{
  "data": [
    {
      "id": "user-123",
      "name": "John Doe",
      "email": "john.doe@example.com",
      "role": "admin",
      "active": true,
      "created_at": "2023-01-01T00:00:00Z",
      "updated_at": "2023-01-02T00:00:00Z",
      "last_login_at": "2023-01-03T12:34:56Z"
    },
    {
      "id": "user-456",
      "name": "Jane Smith",
      "email": "jane.smith@example.com",
      "role": "user",
      "active": true,
      "created_at": "2023-01-01T12:00:00Z",
      "updated_at": "2023-01-01T12:00:00Z",
      "last_login_at": "2023-01-02T09:12:34Z"
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

## Get a Specific User

```
GET /api/v1/users/{id}
```

Retrieves details for a specific user.

### Path Parameters

| Parameter | Type   | Description |
|-----------|--------|-------------|
| `id`      | string | User ID     |

### Response

```json
{
  "data": {
    "id": "user-123",
    "name": "John Doe",
    "email": "john.doe@example.com",
    "role": "admin",
    "active": true,
    "created_at": "2023-01-01T00:00:00Z",
    "updated_at": "2023-01-02T00:00:00Z",
    "last_login_at": "2023-01-03T12:34:56Z",
    "permissions": [
      "jobs:read",
      "jobs:write",
      "jobs:delete",
      "schedules:read",
      "schedules:write",
      "schedules:delete",
      "users:read",
      "users:write",
      "users:delete",
      "settings:read",
      "settings:write"
    ],
    "preferences": {
      "timezone": "America/New_York",
      "theme": "dark",
      "notifications": {
        "email": true,
        "slack": false
      }
    }
  }
}
```

## Create a New User

```
POST /api/v1/users
```

Creates a new user.

### Request Body

```json
{
  "name": "New User",
  "email": "new.user@example.com",
  "password": "securePassword123!",
  "role": "user",
  "active": true,
  "permissions": [
    "jobs:read",
    "jobs:write",
    "schedules:read"
  ],
  "preferences": {
    "timezone": "UTC",
    "theme": "light",
    "notifications": {
      "email": true,
      "slack": true
    }
  }
}
```

### Response

```json
{
  "data": {
    "id": "user-789",
    "name": "New User",
    "email": "new.user@example.com",
    "role": "user",
    "active": true,
    "created_at": "2023-01-04T00:00:00Z",
    "updated_at": "2023-01-04T00:00:00Z",
    "permissions": [
      "jobs:read",
      "jobs:write",
      "schedules:read"
    ],
    "preferences": {
      "timezone": "UTC",
      "theme": "light",
      "notifications": {
        "email": true,
        "slack": true
      }
    }
  }
}
```

## Update a User

```
PUT /api/v1/users/{id}
```

Updates an existing user.

### Path Parameters

| Parameter | Type   | Description |
|-----------|--------|-------------|
| `id`      | string | User ID     |

### Request Body

```json
{
  "name": "Updated User Name",
  "email": "updated.email@example.com",
  "role": "admin",
  "active": true,
  "permissions": [
    "jobs:read",
    "jobs:write",
    "jobs:delete",
    "schedules:read",
    "schedules:write",
    "users:read"
  ],
  "preferences": {
    "timezone": "Europe/London",
    "theme": "dark",
    "notifications": {
      "email": false,
      "slack": true
    }
  }
}
```

### Response

```json
{
  "data": {
    "id": "user-123",
    "name": "Updated User Name",
    "email": "updated.email@example.com",
    "role": "admin",
    "active": true,
    "created_at": "2023-01-01T00:00:00Z",
    "updated_at": "2023-01-04T12:00:00Z",
    "last_login_at": "2023-01-03T12:34:56Z",
    "permissions": [
      "jobs:read",
      "jobs:write",
      "jobs:delete",
      "schedules:read",
      "schedules:write",
      "users:read"
    ],
    "preferences": {
      "timezone": "Europe/London",
      "theme": "dark",
      "notifications": {
        "email": false,
        "slack": true
      }
    }
  }
}
```

## Delete a User

```
DELETE /api/v1/users/{id}
```

Deletes a user.

### Path Parameters

| Parameter | Type   | Description |
|-----------|--------|-------------|
| `id`      | string | User ID     |

### Response

```json
{
  "data": {
    "message": "User deleted successfully"
  }
}
```

## Change User Password

```
POST /api/v1/users/{id}/password
```

Changes the password for a user.

### Path Parameters

| Parameter | Type   | Description |
|-----------|--------|-------------|
| `id`      | string | User ID     |

### Request Body

```json
{
  "current_password": "oldPassword123!",
  "new_password": "newSecurePassword456!"
}
```

### Response

```json
{
  "data": {
    "message": "Password changed successfully"
  }
}
```

## Get User Activity

```
GET /api/v1/users/{id}/activity
```

Retrieves the activity history for a specific user.

### Path Parameters

| Parameter | Type   | Description |
|-----------|--------|-------------|
| `id`      | string | User ID     |

### Query Parameters

| Parameter  | Type    | Description                                                |
|------------|---------|-----------------------------------------------------------|
| `page`     | integer | Page number for pagination (default: 1)                    |
| `per_page` | integer | Number of activities per page (default: 20, max: 100)      |
| `start`    | string  | Start time in ISO 8601 format (optional)                   |
| `end`      | string  | End time in ISO 8601 format (optional)                     |
| `type`     | string  | Filter by activity type (optional)                         |

### Response

```json
{
  "data": [
    {
      "id": "activity-123",
      "user_id": "user-123",
      "type": "job_created",
      "details": {
        "job_id": "job-789",
        "job_name": "New Backup Job"
      },
      "ip_address": "192.168.1.1",
      "user_agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36",
      "created_at": "2023-01-03T10:15:30Z"
    },
    {
      "id": "activity-456",
      "user_id": "user-123",
      "type": "login",
      "details": {
        "method": "password"
      },
      "ip_address": "192.168.1.1",
      "user_agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36",
      "created_at": "2023-01-03T09:00:00Z"
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

## Get Current User

```
GET /api/v1/users/me
```

Retrieves details for the currently authenticated user.

### Response

```json
{
  "data": {
    "id": "user-123",
    "name": "John Doe",
    "email": "john.doe@example.com",
    "role": "admin",
    "active": true,
    "created_at": "2023-01-01T00:00:00Z",
    "updated_at": "2023-01-02T00:00:00Z",
    "last_login_at": "2023-01-03T12:34:56Z",
    "permissions": [
      "jobs:read",
      "jobs:write",
      "jobs:delete",
      "schedules:read",
      "schedules:write",
      "schedules:delete",
      "users:read",
      "users:write",
      "users:delete",
      "settings:read",
      "settings:write"
    ],
    "preferences": {
      "timezone": "America/New_York",
      "theme": "dark",
      "notifications": {
        "email": true,
        "slack": false
      }
    }
  }
}
```

## Update Current User Preferences

```
PATCH /api/v1/users/me/preferences
```

Updates preferences for the currently authenticated user.

### Request Body

```json
{
  "timezone": "Europe/Paris",
  "theme": "light",
  "notifications": {
    "email": true,
    "slack": true
  }
}
```

### Response

```json
{
  "data": {
    "preferences": {
      "timezone": "Europe/Paris",
      "theme": "light",
      "notifications": {
        "email": true,
        "slack": true
      }
    }
  }
}
```

## List User Roles

```
GET /api/v1/roles
```

Retrieves a list of all available user roles.

### Response

```json
{
  "data": [
    {
      "name": "admin",
      "description": "Full system access",
      "permissions": [
        "jobs:read",
        "jobs:write",
        "jobs:delete",
        "schedules:read",
        "schedules:write",
        "schedules:delete",
        "users:read",
        "users:write",
        "users:delete",
        "settings:read",
        "settings:write"
      ]
    },
    {
      "name": "user",
      "description": "Standard user access",
      "permissions": [
        "jobs:read",
        "jobs:write",
        "schedules:read",
        "schedules:write"
      ]
    },
    {
      "name": "viewer",
      "description": "Read-only access",
      "permissions": [
        "jobs:read",
        "schedules:read"
      ]
    }
  ]
}
```

## Get a Specific Role

```
GET /api/v1/roles/{name}
```

Retrieves details for a specific role.

### Path Parameters

| Parameter | Type   | Description |
|-----------|--------|-------------|
| `name`    | string | Role name   |

### Response

```json
{
  "data": {
    "name": "admin",
    "description": "Full system access",
    "permissions": [
      "jobs:read",
      "jobs:write",
      "jobs:delete",
      "schedules:read",
      "schedules:write",
      "schedules:delete",
      "users:read",
      "users:write",
      "users:delete",
      "settings:read",
      "settings:write"
    ],
    "users_count": 2
  }
}
``` 