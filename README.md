# edugov-back-v2

Backend service for the EduGov employee domain. It exposes HTTP endpoints to
manage employee profile data and supports JWT-based authentication.

## Structure

- `cmd/api/main.go`: application entrypoint and HTTP server setup.
- `internal/config`: configuration loading via Viper and `config.yaml`.
- `internal/http`: routing, middleware, and request handlers.
- `internal/service`: business logic and validation.
- `internal/dto`: request/response payload shapes.
- `internal/apperr`: centralized error types and mapping.
- `internal/db/db.go`: database pool and transaction helper.
- `storage/employee/profile_picture`: profile picture storage.

## Configuration

Configuration is loaded from `internal/config/config.yaml` and can be overridden
via environment variables. Viper maps `server.addr` to `SERVER_ADDR`, etc.

Required keys:

- `server.addr`
- `db.url`
- `jwt.secret`
- `jwt.access_ttl`
- `jwt.refresh_ttl`

## Auth

Auth uses JWT access and refresh tokens:

- Access tokens are short-lived (default 2 hours).
- Refresh tokens are stored in `user_session` and rotated on refresh.
- Protected routes use the JWT middleware in `internal/http/middleware/jwt.go`.

## Routes

Routes are configured in `internal/http/router.go`. Most groups expose a public
`GET` endpoint, while mutating endpoints are protected with JWT middleware.

## Profile Pictures

Profile pictures are served from `storage/employee/profile_picture` via:

`GET /employee/profile-picture/{uid}`

Uploads are handled by:

`PUT /employee/profile-picture/{uid}` with multipart `file`.

