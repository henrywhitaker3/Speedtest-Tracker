# Restore a v1 backup

Speedtest Tracker v2 can import result backups produced by the v1 application. Restore is additive: existing results are never overwritten.

## Create the v1 backup

In the v1 UI:

1. Open **Settings**.
2. Select the **Backup/Restore** tab.
3. Open the **Backup** dropdown.
4. Select **JSON**. Your browser downloads the backup file.

## Before restoring

Keep a copy of the database before importing a backup. A restore cannot be undone from the UI.

The backup must be the JSON file downloaded from v1, whose top-level value is an array of result objects. Each result should include these fields:

```json
[
  {
    "id": 1,
    "ping": 12.4,
    "download": 450.5,
    "upload": 35.2,
    "created_at": "2024-01-15T10:30:00.123456789Z",
    "updated_at": "2024-01-15T10:30:00.123456789Z",
    "server_id": 1234,
    "server_name": "Example Server",
    "server_host": "speedtest.example.com",
    "url": "https://www.speedtest.net/result/123456",
    "failed": 0
  }
]
```

`created_at` and `updated_at` must use Go's `time.RFC3339Nano` format. The importer interprets `download` and `upload` as megabits per second, as used by the v1 backup format. A non-zero `failed` value imports the result as an errored result.

## Restore through the UI

1. Open **Settings**.
2. Under **Restore JSON Backup**, select **Restore**.
3. Choose the JSON backup file downloaded from v1.
4. Select **Restore** and wait for the completion message.

When authentication is enabled, restoring is an administrator-only action. Sign in with a user in `AUTH_ADMIN_GROUP` before opening Settings.

The completion message reports how many results were restored and skipped. A result is skipped when a result with the same `created_at` timestamp already exists, including a result that has been soft-deleted. Duplicate timestamps within the uploaded file are also skipped. v1 IDs are not retained; v2 assigns each imported result its own ID.

## API restore

The UI sends the backup array to `POST /v1/restore/json` with an `application/json` content type. On success, the endpoint returns:

```json
{
  "restored": 42,
  "skipped": 3
}
```

Use the generated OpenAPI documentation at `/docs` for the endpoint's current request and authentication details.
