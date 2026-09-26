CREATE TYPE participation_status AS ENUM ('registered', 'cancelled');


CREATE TABLE participations (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    event_id        UUID NOT NULL REFERENCES events(id) ON DELETE CASCADE,
    full_name       VARCHAR(150) NOT NULL,
    email           VARCHAR(150) NOT NULL,
    status          participation_status NOT NULL DEFAULT 'registered',
    registered_at   TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_participations_event_id ON participations(event_id);
CREATE INDEX idx_participations_email ON participations(email);


CREATE UNIQUE INDEX participations_event_email_active_key
    ON participations (event_id, email)
    WHERE status = 'registered';
