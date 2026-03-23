# Identity Service

This is a cloud-native authentication service written in Go.

## Overview
This service provides user registration, login, and JWT generation/validation. It is designed as a true microservice, operating independently with its own database and cache.

## Architectural Decisions
- **True microservice**: Own DB, own Redis, public interface via JWT.
- **Stateless JWT**: Tokens carry `user_id` and claims, validated locally by other services.
- **Redis for token cache**: Isolated Redis for refresh tokens and blacklists.
- **Isolated users database**: Postgres exclusive to this service.
- **Stateless by default**: Go process holds no in-memory state.

## Stack
- **Language**: Go
- **HTTP Framework**: `chi`
- **Database**: Postgres
- **Cache / Sessions**: Redis
- **Tokens**: JWT (`golang-jwt/jwt`)
- **Containerization**: Docker
- **Orchestration**: Kubernetes (future)

## Folder Structure
```
/cmd
  /server         # main entrypoint
/internal
  /auth           # domain logic: registration, login, validation
  /token          # JWT generation and validation
  /user           # user repository (Postgres access)
  /cache          # Redis access
/pkg
  /middleware     # HTTP middlewares (token validation)
/config           # env var configuration
/migrations       # Postgres migrations
```

## Main Flows

### Registration
`POST /auth/register`
Body: `{ email, password }`
- Validate input
- Hash password (bcrypt)
- Persist to Postgres
- Return 201

### Login
`POST /auth/login`
Body: `{ email, password }`
- Fetch user from Postgres
- Compare hash
- Generate access token (JWT, short-lived)
- Generate refresh token (JWT, long-lived → save to Redis)
- Return `{ access_token, refresh_token }`

### Validation (used by other services)
`GET /auth/validate`
Header: `Authorization: Bearer <token>`
- Validate JWT signature
- Return `{ user_id, claims }` or 401

## Environment Variables

| Variable          | Default                 | Description                                      |
| :---------------- | :---------------------- | :----------------------------------------------- |
| `PORT`            | `8080`                  | Port for the HTTP server to listen on.           |
| `DATABASE_URL`    | `postgres://user:pass@localhost:5432/authdb` | Connection string for the PostgreSQL database.   |
| `REDIS_URL`       | `redis://localhost:6379`| Connection string for the Redis instance.        |
| `JWT_SECRET`      | `change-this-in-production` | Secret key for signing and verifying JWTs.       |
| `JWT_ACCESS_TTL`  | `15m`                   | Time-to-live for access tokens (e.g., `15m`, `1h`). |
| `JWT_REFRESH_TTL` | `7d`                    | Time-to-live for refresh tokens (e.g., `7d`, `30d`). |

## Development
To get started, copy the `.env.example` to `.env` and fill in the necessary details.