-- name: LockEventCapacity :one
SELECT capacity FROM events WHERE id = $1 FOR UPDATE;

-- name: CountRegistered :one
SELECT COUNT(*) FROM participations WHERE event_id = $1 AND status = 'registered';

-- name: CreateParticipation :one
INSERT INTO participations (full_name, email, event_id, status)
VALUES ($1, $2, $3, 'registered')
RETURNING *;

-- name: ListParticipants :many
SELECT full_name, email, status, registered_at
FROM participations
WHERE event_id = $1
ORDER BY registered_at ASC;
