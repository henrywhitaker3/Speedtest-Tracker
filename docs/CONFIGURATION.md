# Configuration

Speedtest Tracker reads configuration from environment variables. When running the binary directly, it also loads a local `.env` file. Docker Compose reads its `.env` file when starting the stack.

## Application

| Variable | Default | Description |
| --- | --- | --- |
| `ACCEPT_EULA` | _required_ | Set to `true` to confirm acceptance of Ookla's EULA and privacy agreement. |
| `APP_URL` | `http://localhost:8080` | Public application URL, including scheme and port. It is used for OpenAPI documentation and the OIDC callback URL. |
| `APP_UI_URL` | `APP_URL` | Browser URL for the UI. Set this when the UI is served from a different origin, such as a Nuxt development server. |
| `APP_UI_DIR` | `ui` | Directory containing the built frontend. The container image provides it at `/ui`. |
| `LOG_LEVEL` | application default | Logging verbosity. |

The Speedtest binary receives Ookla's EULA and GDPR flags when it runs. `ACCEPT_EULA=true` is still required as this application's explicit startup acknowledgement.

## Database and runtime

| Variable | Default | Description |
| --- | --- | --- |
| `DB_DRIVER` | `sqlite` | Database driver: `sqlite` or `postgres`. |
| `DB_URL` | `db.sqlite` | SQLite database filename or PostgreSQL connection string. |
| `RUNTIME_HOSTNAME` | unset | Optional runtime hostname. |
| `RUNTIME_PSK` | unset | Private key used by the local runtime. Set a strong, stable secret for deployments. |

## Speedtests and retention

| Variable | Default | Description |
| --- | --- | --- |
| `TEST_SCHEDULE` | `5 */2 * * *` | Cron expression for scheduled tests: by default, five minutes past every second hour. |
| `TEST_SERVER` | unset | Optional Ookla server ID used for speedtests. When unset, Ookla selects a server. |
| `RETENTION_DAYS` | `0` (disabled) | Retention duration for results. Expired results are permanently deleted hourly. Use a Go-style duration such as `90d`; `0` disables cleanup. |

## Notifications

| Variable | Default | Description |
| --- | --- | --- |
| `NOTIFICATIONS_HEALTHCHECKS` | unset | Healthchecks.io check URL. The workflow sends a start ping and a success or failure ping; a failed ping does not stop a speedtest. |

## Authentication

The authentication variables are documented in [Authentication](AUTH.md).

## Observability

| Variable | Default | Description |
| --- | --- | --- |
| `OTEL_EXPORTER_OTLP_ENDPOINT` | unset | OTLP endpoint used to export traces. |

