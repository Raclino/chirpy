# Chirpy

Chirpy is a small REST API built in Go with PostgreSQL.

The project focuses on backend fundamentals: building HTTP endpoints with the Go standard library, working with a relational database, handling authentication, and managing a simple API lifecycle from user creation to protected resources.

## Features

- user registration
- login with email and password
- JWT access tokens
- refresh token flow with revocation
- chirp creation
- list chirps
- filter chirps by author
- sort chirps by creation date
- get a chirp by ID
- delete your own chirps
- webhook-based user upgrade to Chirpy Red

## Tech stack

- Go
- PostgreSQL
- Docker Compose
- Goose
- sqlc
- Argon2
- JWT
- slog

## Running the project

### Requirements

- Go installed
- Docker and Docker Compose
- Goose installed
- sqlc installed

### Start PostgreSQL

```bash
docker compose up -d
```

### Create a .env file

Example:

```bash
DB_URL=postgres://postgres:postgres@localhost:5433/chirpy?sslmode=disable
PLATFORM=dev
JWT_SECRET=replace_this_with_a_long_random_secret
POLKA_KEY=replace_this_with_your_polka_api_key
```

### Run migrations

```bash
make mig-up
```

### Generate sqlc code

```bash
make sqlc
```

### Start the server

```bash
make dev
```

The API runs on <http://localhost:8080>.

### Main endpoints

**Users**
`POST /api/users`
`PUT /api/users`

**Auth**
`POST /api/login`
`POST /api/refresh`
`POST /api/revoke`

**Chirps**
`POST /api/chirps`
`GET /api/chirps`
`GET /api/chirps/{chirpID}`
`DELETE /api/chirps/{chirpID}`

**Webhooks**
`POST /api/polka/webhooks`

#### Chirps query parameters

`GET /api/chirps supports:`

`author_id`
`sort=asc|desc`

Example:

`GET /api/chirps?author_id=<user-id>&sort=desc`

### Notes

This project was built as a backend practice project to improve my understanding of:

- HTTP servers in Go
- REST API design
- database migrations and typed SQL queries
- authentication and token-based security
- structured logging
- webhook handling
