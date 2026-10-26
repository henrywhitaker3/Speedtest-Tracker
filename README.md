# Speedtest Tracker v2

[![Docker pulls](https://img.shields.io/docker/pulls/henrywhitaker3/speedtest-tracker?style=flat-square)](https://hub.docker.com/r/henrywhitaker3/speedtest-tracker)
[![GitHub Workflow Status](https://img.shields.io/github/actions/workflow/status/henrywhitaker3/Speedtest-Tracker/build.yaml?style=flat-square&label=Build)](https://github.com/henrywhitaker3/Speedtest-Tracker/actions/workflows/build.yaml)
[![last_commit](https://img.shields.io/github/last-commit/henrywhitaker3/Speedtest-Tracker?style=flat-square?branch=v2)](https://github.com/henrywhitaker3/Speedtest-Tracker/commits)
[![issues](https://img.shields.io/github/issues/henrywhitaker3/Speedtest-Tracker?style=flat-square)](https://github.com/henrywhitaker3/Speedtest-Tracker/issues)
[![commit_freq](https://img.shields.io/github/commit-activity/m/henrywhitaker3/Speedtest-Tracker?style=flat-square?branch=v2)](https://github.com/henrywhitaker3/Speedtest-Tracker/commits)
[![license](https://img.shields.io/github/license/henrywhitaker3/Speedtest-Tracker?style=flat-square?branch=v2)](https://github.com/henrywhitaker3/Speedtest-Tracker/blob/v2/LICENSE)

Speedtest Tracker runs scheduled [Ookla Speedtest CLI](https://www.speedtest.net/apps/cli) checks, stores the results, and presents them in a web UI. Version 2 is a Go service with a Nuxt frontend; it is a rewrite of the previous Laravel/React application.

> Ookla's EULA and privacy agreement must be accepted before the application will start.

<img width="45.5%" alt="homepage screenshot" src="https://github.com/user-attachments/assets/9d50880b-152e-4317-9cb4-c9adf4b78e9d" />
<img width="48%" alt="test result screenshot" src="https://github.com/user-attachments/assets/f66391eb-4fd6-4ea2-af08-b701f4a575b2" />



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

For a simple deployment, use SQLite with a persistent Docker volume. Create an empty directory and add this `compose.yaml` file:

```yaml
services:
  app:
    image: henrywhitaker3/speedtest-tracker:latest
    restart: unless-stopped
    ports:
      - "8080:8080"
    env_file: .env
    environment:
      ACCEPT_EULA: "true"
      DB_DRIVER: sqlite
      DB_URL: /data/speedtest-tracker.db
    volumes:
      - speedtest-tracker:/data

volumes:
  speedtest-tracker:
```

Create a `.env` file beside it with a private runtime key:

```dotenv
RUNTIME_PSK=replace-with-the-output-of-openssl-rand-hex-32
```

For example:

```sh
openssl rand -hex 32
```

Then start the application:

```sh
docker compose up -d
```

Open <http://localhost:8080>. The SQLite database is stored in the `speedtest-tracker` Docker volume and survives container recreation.

To stop the stack without removing its data:

```sh
docker compose down
```

The repository's `docker-compose.yaml` starts PostgreSQL and Jaeger too; it is intended for local development.

## Configuration

Configuration is read from environment variables (and a local `.env` file when running the binary directly). See the full [configuration reference](docs/CONFIGURATION.md), including result retention, scheduled tests, Healthchecks.io, and observability.

For optional OpenID Connect authentication setup, see [Authentication](docs/AUTH.md).

To import results from a Speedtest Tracker v1 JSON backup, see [Restore a v1 backup](docs/RESTORE.md).

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

This branch is a new implementation. The v1 Laravel/React configuration, authentication, notification, and external-service integration instructions do not apply to v2. v1 JSON result backups can be imported using the [restore process](docs/RESTORE.md).
