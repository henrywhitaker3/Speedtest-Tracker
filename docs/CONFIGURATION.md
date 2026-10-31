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
| `TEST_SERVER` | unset | Optional comma-separated Ookla server IDs. One is selected at random for each speedtest. When unset, a nearby server is selected at random from Ookla's server list. |
| `RETENTION_DAYS` | `0` (disabled) | Retention duration for results. Expired results are permanently deleted hourly. Use a Go-style duration such as `90d`; `0` disables cleanup. |

## Notifications

Configure each notification with a zero-based, contiguous index. Every notification requires its kind and phases:

| Variable | Description |
| --- | --- |
| `NOTIFICATIONS_<index>_KIND` | Notification type: `healthcheck`, `webhook`, `telegram`, `slack`, or `discord`. |
| `NOTIFICATIONS_<index>_PHASES` | Comma-separated phases: `pre-run`, `post-run`, or both. |
| `NOTIFICATIONS_<index>_URL` | Destination URL, required for `healthcheck`, `webhook`, `slack`, and `discord`. |
| `NOTIFICATIONS_<index>_HEADERS` | Optional JSON object of custom HTTP headers for a `webhook` notification. |
| `NOTIFICATIONS_TELEGRAM_BOT_TOKEN` | Telegram bot token, required when a `telegram` notification is configured. |
| `NOTIFICATIONS_TELEGRAM_CHAT_ID` | Telegram chat ID, required when a `telegram` notification is configured. |

For example, this configures a Healthchecks.io notification before and after the speedtest, a webhook after it finishes, and Telegram, Slack, and Discord messages before and after each test:

```dotenv
NOTIFICATIONS_0_KIND=healthcheck
NOTIFICATIONS_0_PHASES=pre-run,post-run
NOTIFICATIONS_0_URL=https://hc-ping.com/your-check-id

NOTIFICATIONS_1_KIND=webhook
NOTIFICATIONS_1_PHASES=post-run
NOTIFICATIONS_1_URL=https://example.com/speedtest-complete
NOTIFICATIONS_1_HEADERS='{"Authorization":"Bearer replace-with-token","X-Source":"speedtest-tracker"}'

NOTIFICATIONS_2_KIND=telegram
NOTIFICATIONS_2_PHASES=pre-run,post-run
NOTIFICATIONS_TELEGRAM_BOT_TOKEN=123456:replace-with-your-bot-token
NOTIFICATIONS_TELEGRAM_CHAT_ID=-1001234567890

NOTIFICATIONS_3_KIND=slack
NOTIFICATIONS_3_PHASES=pre-run,post-run
NOTIFICATIONS_3_URL=https://hooks.slack.com/services/your-webhook-id

NOTIFICATIONS_4_KIND=discord
NOTIFICATIONS_4_PHASES=pre-run,post-run
NOTIFICATIONS_4_URL=https://discord.com/api/webhooks/your-webhook-id/your-webhook-token
```

Telegram credentials are shared by every Telegram notification. The post-run message contains the result measurements or failure reason and links to the result in the UI.

Webhook headers are sent only for `webhook` notifications. Set `NOTIFICATIONS_<index>_HEADERS` to a JSON object whose keys and values are header names and values, respectively. The configured headers are added to the JSON `POST` request.

Slack notifications use an [Incoming Webhook](https://api.slack.com/messaging/webhooks) and render native Block Kit messages. Discord notifications use a standard [Discord webhook](https://discord.com/developers/docs/resources/webhook#execute-webhook), without the `/slack` suffix, and render native embeds. Both include the speedtest result link; successful post-run messages show ping, download, and upload, while failed tests show the error.

`NOTIFICATIONS_HEALTHCHECKS` remains supported for compatibility. It creates one `healthcheck` notification for both `pre-run` and `post-run` phases.

## Authentication

The authentication variables are documented in [Authentication](AUTH.md).

## Observability

| Variable | Default | Description |
| --- | --- | --- |
| `OTEL_EXPORTER_OTLP_ENDPOINT` | unset | OTLP endpoint used to export traces. |

