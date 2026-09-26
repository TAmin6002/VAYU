CREATE TABLE events (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title           VARCHAR(200) NOT NULL,
    description     TEXT NOT NULL DEFAULT '',
    event_time      TIMESTAMPTZ NOT NULL,
    location        VARCHAR(200) NOT NULL,
    capacity        INTEGER NOT NULL CHECK (capacity > 0),
    created_by      UUID NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_events_event_time ON events(event_time);
CREATE INDEX idx_events_created_by ON events(created_by);
