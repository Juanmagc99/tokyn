# Tokyn – API Key Management Service

**Tokyn** is a lightweight, internal-use API key manager built in Go. It provides a simple and secure way to create, revoke, delete, and verify API keys using a RESTful interface. It is designed to run behind a reverse proxy and be consumed by internal services only.

---

## ⚙️ Features

- RESTful API for managing API keys
- SQLite as the database backend
- Redis-based rate limiting and cache
- Internal security layer using a custom header
- Lightweight and easy to deploy (Docker-ready)
- Token verification via JSON body

---

## 🚀 Getting Started

You can run the app via Docker Compose:

```bash
docker-compose up --build
```

Environment variables:

| Variable       | Description                         | Default              |
|----------------|-------------------------------------|----------------------|
| `SQLITE_DB`    | Path to SQLite DB file              | `data.db`            |
| `REDIS_ADDR`   | Redis connection address            | `localhost:6379`     |
| `REDIS_PASS`   | Redis password (optional)           | `""`                 |
| `APP_ADDR`     | App bind address (host:port)        | `0.0.0.0:8080`       |
| `INTERNAL_API_TOKEN` | Internal access token (optional) | _not set_            |

---

## 🔐 Internal Security

This service is intended to run **behind a reverse proxy** (e.g. NGINX, Traefik) and be used by internal services.

Additionally, the service supports a lightweight internal authentication mechanism using a custom header:

- Set `INTERNAL_API_TOKEN` in the environment.
- All requests to `/apikeys/*` must then include the header:

```http
X-Internal-Token: your_token_here
```

If `INTERNAL_API_TOKEN` is not set, the check is skipped.

---

## 📖 API Reference

All endpoints are grouped under `/apikeys`, protected by internal middleware.

### POST `/apikeys`

**Create a new API key.**

#### Request body:
```json
{
  "name": "Service A"
}
```

#### Response:
```json
{
  "id": "123",
  "token": "generated-token",
  "name": "Service A",
  "created_at": "..."
}
```

---

### GET `/apikeys/:id`

**Get details for a specific API key by ID.**

#### Response:
```json
{
  "id": "123",
  "name": "Service A",
  "revoked": false,
  "revoked_at": null
}
```

---

### PATCH `/apikeys/:id/revocation`

**Soft-revoke an API key (it remains in the DB but becomes invalid).**

#### Response:
HTTP 204 No Content

---

### DELETE `/apikeys/:id`

**Permanently delete an API key.**

#### Response:
HTTP 204 No Content

---

### GET `/apikeys/verification`

**Verify the validity of an API key.**

This endpoint uses a JSON body (not query or path param for security reasons).

#### Request body:
```json
{
  "token": "abcdef123456"
}
```

#### Success response:
```json
{
  "name": "Service A"
}
```

#### Error responses:
- `401 Unauthorized` – token is missing, revoked, or invalid
- `429 Too Many Requests` – rate limit exceeded

---

## 📦 Technologies Used

- [Go](https://golang.org/)
- [Echo](https://echo.labstack.com/)
- [GORM](https://gorm.io/) + SQLite
- [Redis](https://redis.io/)
- Docker & Docker Compose

---

## 🧩 Notes

- Tokens are stored hashed in the database.
- SQLite database is persisted via Docker volume.
- The API is intended for internal use only – it should not be exposed directly to the internet.

---

## 📝 License

MIT
