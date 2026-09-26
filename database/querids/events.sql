-- name: CreateEvent :one
INSERT INTO events (title, description, event_time, location, capacity, created_by)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: UpdateEvent :one
UPDATE events
SET title = $1, description = $2, event_time = $3, location = $4, capacity = $5, updated_at = now()
WHERE id = $6
RETURNING *;

-- name: DeleteEvent :exec
DELETE FROM events WHERE id = $1;

-- name: GetEventByID :one
SELECT * FROM events WHERE id = $1;

-- name: ListEvents :many
SELECT
    e.*,
    COALESCE(p.registered_count, 0) AS registered_count
FROM events e
LEFT JOIN (
    SELECT event_id, COUNT(*) AS registered_count
    FROM participations
    WHERE status = 'registered'
    GROUP BY event_id
) p ON p.event_id = e.id
ORDER BY e.event_time ASC;
