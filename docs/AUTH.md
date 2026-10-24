# Authentication

Speedtest Tracker can protect manual speedtest runs with an OpenID Connect (OIDC) provider. Authentication is disabled by default. When it is enabled, the UI shows a **Sign in** button and `POST /v1/speedtest/run` requires an authenticated user or a valid API token. Scheduled tests and read-only result endpoints remain available.

The application uses the authorization-code flow with PKCE.

## Configure an OIDC provider

Register this redirect URI with the provider exactly:

```text
${APP_URL}/v1/auth/callback
```

Then configure the application:

```dotenv
APP_URL=https://speedtest.example.com
AUTH_ENABLED=true
AUTH_ISSUER_URL=https://issuer.example.com
AUTH_CLIENT_ID=speedtest-tracker
AUTH_CLIENT_SECRET=replace-with-the-client-secret
```

| Variable | Default | Description |
| --- | --- | --- |
| `AUTH_ENABLED` | `false` | Set to `true` to enable OIDC authentication. |
| `AUTH_ISSUER_URL` | _required when enabled_ | OIDC issuer URL. Its discovery document provides the authorization and token endpoints. |
| `AUTH_CLIENT_ID` | _required when enabled_ | OIDC client ID. |
| `AUTH_CLIENT_SECRET` | _required when enabled_ | OIDC client secret. Keep it private. |
| `AUTH_SCOPES` | `openid profile email groups offline_access` | Space-separated requested scopes. `offline_access` permits session refresh when the provider supports refresh tokens. |
| `AUTH_GROUPS_CLAIM` | `groups` | ID-token claim containing the user's groups. |
| `AUTH_ADMIN_GROUP` | `admin` | Group whose members can perform administrator-only actions. |

All required OIDC settings must be present when `AUTH_ENABLED=true`.

## UI on a different origin

If the Nuxt UI is served separately, set `APP_UI_URL` to its browser origin, for example `http://localhost:3000`, and leave `APP_URL` pointing to the API, such as `http://localhost:8080`. After sign-in, the callback redirects the browser to `APP_UI_URL`.

Cross-origin UI/API requests must include credentials. The bundled UI is configured to do so.

## Sessions and refresh tokens

The browser receives an opaque, HTTP-only session cookie. OIDC tokens, the PKCE verifier, nonce, and authorization state remain in server-side actor state. Production deployments must use HTTPS because the cookie is marked `Secure`.

The login transaction expires after 10 minutes. An authenticated session is kept for seven days; while it remains valid, the app refreshes expiring OIDC tokens when the issuer supplied a refresh token. If refresh fails after a token expires, the stored session is invalidated and the user must sign in again.

## API tokens

Authenticated users can create API tokens through the authentication API. Send one with requests as a bearer token:

```text
Authorization: Bearer <token>
```

API tokens can be used for endpoints that require authentication, including manual speedtest runs. They do not represent an administrator session, so administrator-only endpoints still require an authenticated OIDC user who belongs to `AUTH_ADMIN_GROUP`.
