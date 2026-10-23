# Speedtest Tracker v2

[![Docker pulls](https://img.shields.io/docker/pulls/henrywhitaker3/speedtest-tracker?style=flat-square)](https://hub.docker.com/r/henrywhitaker3/speedtest-tracker)
[![GitHub Workflow Status](https://img.shields.io/github/actions/workflow/status/henrywhitaker3/Speedtest-Tracker/build.yaml?style=flat-square&label=Build)](https://github.com/henrywhitaker3/Speedtest-Tracker/actions/workflows/build.yaml)
[![last_commit](https://img.shields.io/github/last-commit/henrywhitaker3/Speedtest-Tracker?style=flat-square?branch=v2)](https://github.com/henrywhitaker3/Speedtest-Tracker/commits)
[![issues](https://img.shields.io/github/issues/henrywhitaker3/Speedtest-Tracker?style=flat-square)](https://github.com/henrywhitaker3/Speedtest-Tracker/issues)
[![commit_freq](https://img.shields.io/github/commit-activity/m/henrywhitaker3/Speedtest-Tracker?style=flat-square?branch=v2)](https://github.com/henrywhitaker3/Speedtest-Tracker/commits)
[![license](https://img.shields.io/github/license/henrywhitaker3/Speedtest-Tracker?style=flat-square?branch=v2)](https://github.com/henrywhitaker3/Speedtest-Tracker/blob/v2/LICENSE)

Speedtest Tracker runs scheduled [Ookla Speedtest CLI](https://www.speedtest.net/apps/cli) checks, stores the results, and presents them in a web UI. Version 2 is a Go service with a Nuxt frontend; it is a rewrite of the previous Laravel/React application.

> Ookla's EULA and privacy agreement must be accepted before the application will start.

<img width="2730" height="1974" alt="image" src="https://github.com/user-attachments/assets/fa39015d-f3b2-47e6-bdfd-34e54ce2d728" />

## Features

- Scheduled speed tests using a cron expression, plus on-demand runs from the UI or API.
- A static Nuxt web UI served by the Go application.
- SQLite and PostgresSQL database support
- Paginated results, latest-result, and time-series API endpoints.
- Generated OpenAPI schema and Swagger UI.
- Prometheus metrics, including a native histogram of Speedtest duration.
- OpenTelemetry tracing support.
- Optional OpenID Connect (OIDC) authentication.

## Quick start with Docker Compose

Docker Compose is the recommended way to run v2. It starts the application, PostgreSQL, and Jaeger for local trace inspection.

Create a `.env` file beside `docker-compose.yaml` with a private runtime key:

```dotenv
RUNTIME_PSK=replace-with-the-output-of-openssl-rand-hex-32
```

For example:

```sh
openssl rand -hex 32
```

Then build and start the stack:

```sh
docker compose up --build -d
```

Open <http://localhost:8080>. The included PostgreSQL volume retains results between restarts. Jaeger's UI is available at <http://localhost:16686>.

To stop the stack without removing its data:

```sh
docker compose down
```

## Configuration

Configuration is read from environment variables (and a local `.env` file when running the binary directly).

| Variable | Default | Description |
| --- | --- | --- |
| `ACCEPT_EULA` | _required_ | Set to `true` to confirm acceptance of Ookla's EULA and privacy agreement. |
| `APP_URL` | `http://localhost:8080` | Public application URL; used when generating OpenAPI documentation. |
| `APP_UI_URL` | `APP_URL` | Browser URL for the UI; used as the post-login redirect destination when the frontend runs on a different port or origin. |
| `APP_UI_DIR` | `ui` | Directory containing the built frontend. The container image provides it at `/ui`. |
| `AUTH_ENABLED` | `false` | Set to `true` to enable OIDC authentication. |
| `AUTH_ISSUER_URL` | _required when auth is enabled_ | OIDC issuer URL. Its discovery document is used to find the authorization and token endpoints. |
| `AUTH_CLIENT_ID` | _required when auth is enabled_ | OIDC client ID. |
| `AUTH_CLIENT_SECRET` | _required when auth is enabled_ | OIDC client secret. Keep this private. |
| `AUTH_SCOPES` | `openid profile email groups offline_access` | Space-separated scopes requested from the issuer. `offline_access` permits session refresh when the issuer supports refresh tokens. |
| `AUTH_GROUPS_CLAIM` | `groups` | ID-token claim containing the user's groups. |
| `AUTH_ADMIN_GROUP` | `admin` | Group treated as administrators by the UI. |
| `DB_DRIVER` | `sqlite` | Database driver: `sqlite` or `postgres`. |
| `DB_URL` | `db.sqlite` | SQLite filename or PostgreSQL connection string. |
| `RUNTIME_PSK` | empty | Private key used by the local runtime. Set a strong, stable value for deployments. |
| `TEST_SCHEDULE` | `5 */2 * * *` | Cron expression controlling scheduled tests (by default, five minutes past every second hour). |
| `RETENTION_DAYS` | `0` (disabled) | Number of days to retain results. Older results are permanently deleted hourly. E.g. `90d` |
| `LOG_LEVEL` | application default | Logging verbosity. |
| `OTEL_EXPORTER_OTLP_ENDPOINT` | unset | Optional OTLP endpoint for traces. |

The Speedtest binary accepts Ookla's EULA and GDPR flags when it runs. `ACCEPT_EULA=true` is still required as the application's explicit startup acknowledgement.

### Optional OIDC authentication

Authentication is disabled by default. When enabled, the UI shows a **Sign in** button and `POST /v1/speedtest/run` requires a valid OIDC session; scheduled tests and read-only result endpoints remain available. The application uses the authorization-code flow with PKCE.

Register this redirect URI with the OIDC provider exactly:

```text
${APP_URL}/v1/auth/callback
```

For example, a deployment whose API and UI share `https://speedtest.example.com` can use:

```dotenv
APP_URL=https://speedtest.example.com
AUTH_ENABLED=true
AUTH_ISSUER_URL=https://issuer.example.com
AUTH_CLIENT_ID=speedtest-tracker
AUTH_CLIENT_SECRET=replace-with-the-client-secret
```

If the Nuxt UI is served separately during development, set `APP_UI_URL` to its browser origin (for example, `http://localhost:3000`) and leave `APP_URL` pointing to the API (`http://localhost:8080`). The callback returns the browser to `APP_UI_URL` after a successful login.

The browser receives only an opaque, HTTP-only session cookie. OIDC tokens, PKCE verifier, nonce, and authorization state are kept in server-side actor state. The cookie is marked `Secure`, so production deployments must use HTTPS. Cross-origin UI/API requests must include credentials; the bundled UI is configured to do so.

## Local development

The project uses [mise](https://mise.jdx.dev/) to provision the backend tooling. The frontend has its own `frontend/mise.toml` for Node.js and pnpm.

```sh
mise install

cd frontend
mise install
mise exec -- pnpm install --frozen-lockfile
cd ..

mise run build-ui
```

Create a local `.env` file:

```dotenv
ACCEPT_EULA=true
RUNTIME_PSK=replace-with-a-random-value
APP_UI_DIR=frontend/.output/public
DB_DRIVER=sqlite
DB_URL=speedtest-tracker.db
```

Start the application and its development dependencies:

```sh
mise run serve
```

For iterative frontend development, run the Nuxt development server separately:

```sh
cd frontend
mise run ui
```

The frontend API base URL defaults to `http://localhost:8080`. Generate its TypeScript client from the running service with:

```sh
cd frontend
mise run client
```

Run backend tests with:

```sh
mise run test
```

## Building the container image

The Dockerfile builds the Nuxt UI first, copies its static output to `/ui`, then builds the Go binary into a small non-root runtime image.

```sh
docker build -t speedtest-tracker:local .
```

The image expects a database and the configuration described above. Use the supplied Compose stack when you need a ready-to-run PostgreSQL deployment.

## Notes on v2

This branch is a new implementation. The legacy Laravel/React configuration, authentication, notification, backup, and external-service integration instructions from v1 do not apply to v2.
