# University Events

A web platform for announcing and registering for faculty events (workshops, contests, talks). Anyone can browse events; students sign up and register; an admin manages events.

## Tech Stack

| Layer    | Technology                              |
|----------|-----------------------------------------|
| Backend  | Go 1.22, [chi](https://github.com/go-chi/chi) |
| Database | PostgreSQL 16                           |
| Auth     | Stateless JWT (role stored in token)    |
| Frontend | Plain HTML / CSS / JavaScript           |

## Quick Start (Docker only)

**Prerequisite:** Docker with Docker Compose. Nothing else to install.

```bash
docker compose up --build
```

Open **http://localhost:8080**. This starts PostgreSQL and the app together. On the first start the database is created automatically from `database/init/init.sql` (all tables + the default admin).

```bash
docker compose down        # stop (data is kept)
docker compose down -v     # stop and wipe the database (init.sql runs again next time)
```

> `init.sql` is the concatenation of `database/migrations/*.up.sql`. If you change a migration, regenerate it:
> `cat database/migrations/*.up.sql > database/init/init.sql`

## Project Structure

```
cmd/server/          entry point and routes
cmd/hashpassword/    CLI to generate a bcrypt hash
Dockerfile, docker-compose.yml   run app + PostgreSQL locally
internal/config/     env / .env configuration
database/            connection, init/ (Docker DB setup), migrations/, queries/ (sqlc, not yet used)
repository/          database access (incl. capacity transaction)
handler/             HTTP layer
middleware/          JWT auth, RequireRole, CORS
models/              data structs
pkg/jwt, pkg/response  JWT helper, JSON envelope
web/                 static frontend
```

Request flow: `middleware (CORS, Auth) → handler → repository → PostgreSQL`

## Roles

| Role        | Table      | Can do |
|-------------|------------|--------|
| **Admin**   | `admins`   | Create / edit / delete events, view participants. Cannot register for events. |
| **Student** | `students` | Sign up, log in/out, register for events (name, email, student number), cancel own registrations. |
| Guest       | —          | Browse event list and details only. |

**Default admin** (created by `init.sql`): username `admin`, password `12345678`.
Change it before any real deployment. To create another admin, run
`go run ./cmd/hashpassword "<password>"` and insert the hash into `admins`.

## Capacity Control

Registration runs in one transaction (`repository/participation_repo.go`):

1. Lock the event row (`SELECT ... FOR UPDATE`).
2. Count active registrations.
3. If full, return `event_full`; otherwise insert and commit.

This prevents overbooking under concurrent requests. A partial unique index on `(event_id, student_id)` prevents duplicate active registrations.

## API

Base path: `/api/v1`. Protected routes need `Authorization: Bearer <token>`.
Responses: `{ "data": ... }` or `{ "error": "message" }`.

**Public**

| Method | Path            | Description |
|--------|-----------------|-------------|
| POST   | `/auth/login`   | `{username, password}` → `{token, role, user}` (students use email as username) |
| POST   | `/auth/signup`  | `{full_name, email, student_number, password}` |
| GET    | `/events`       | List events (soonest first, with registered count) |
| GET    | `/events/{id}`  | Event details |

**Any logged-in user**

| Method | Path            | Description |
|--------|-----------------|-------------|
| POST   | `/auth/logout`  | Client discards the token (JWTs are stateless) |

**Student only**

| Method | Path                      | Description |
|--------|---------------------------|-------------|
| GET    | `/student/me`             | Current profile |
| GET    | `/student/registrations`  | Events I'm registered for |
| POST   | `/events/{id}/register`   | `{full_name, email, student_number}` (5–20 digits) |
| POST   | `/events/{id}/cancel`     | Cancel my registration |

**Admin only**

| Method | Path                        | Description |
|--------|-----------------------------|-------------|
| GET    | `/admin/me`                 | Current profile |
| POST   | `/events`                   | `{title, description, location, capacity, event_time}` |
| PUT    | `/events/{id}`              | Update event |
| DELETE | `/events/{id}`              | Delete event |
| GET    | `/events/{id}/participants` | List participants |

## Next Steps

- Escape user-supplied values in the frontend (XSS) and tighten CORS.
- Generate repository code with `sqlc` instead of hand-written queries.
- Add tests, especially for the concurrent-registration scenario.
- Add rate limiting and stronger input validation.
- Add a `department` field to prepare for multiple departments.