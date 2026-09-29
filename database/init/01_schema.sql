-- Database schema. Postgres runs every *.sql file in this folder
-- automatically, in alphabetical order, the FIRST time the database
-- container starts with an empty data volume (see docker-compose.yml).
-- No migration tool is needed.

-- Account type 1: the admin. Can only manage events.
CREATE TABLE admins (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    full_name       VARCHAR(150) NOT NULL,
    email           VARCHAR(150) NOT NULL UNIQUE,   -- also used as the login username
    password_hash   VARCHAR(255) NOT NULL,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- Account type 2: students. Can sign up, log in/out, register for events
-- and cancel their own registrations.
CREATE TABLE students (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    full_name       VARCHAR(150) NOT NULL,
    email           VARCHAR(150) NOT NULL UNIQUE,
    student_number  VARCHAR(20)  NOT NULL UNIQUE,
    password_hash   VARCHAR(255) NOT NULL,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE events (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title           VARCHAR(200) NOT NULL,
    description     TEXT NOT NULL DEFAULT '',
    event_time      TIMESTAMPTZ NOT NULL,
    location        VARCHAR(200) NOT NULL,
    capacity        INTEGER NOT NULL CHECK (capacity > 0),
    created_by      UUID NOT NULL REFERENCES admins(id) ON DELETE RESTRICT,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_events_event_time ON events(event_time);
CREATE INDEX idx_events_created_by ON events(created_by);

CREATE TYPE participation_status AS ENUM ('registered', 'cancelled');

-- A student's registration for an event. The form's name / email /
-- student number are stored with the registration.
CREATE TABLE participations (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    event_id        UUID NOT NULL REFERENCES events(id) ON DELETE CASCADE,
    student_id      UUID REFERENCES students(id) ON DELETE CASCADE,
    full_name       VARCHAR(150) NOT NULL,
    email           VARCHAR(150) NOT NULL,
    student_number  VARCHAR(20)  NOT NULL,
    status          participation_status NOT NULL DEFAULT 'registered',
    registered_at   TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_participations_event_id ON participations(event_id);
CREATE INDEX idx_participations_student_id ON participations(student_id);
CREATE INDEX idx_participations_email ON participations(email);
CREATE INDEX idx_participations_student_number ON participations(student_number);

-- Only *active* registrations must be unique per (event, student).
-- Cancelled rows are kept for history and must not block the student
-- from registering again later.
CREATE UNIQUE INDEX participations_event_student_active_key
    ON participations (event_id, student_id)
    WHERE status = 'registered' AND student_id IS NOT NULL;
