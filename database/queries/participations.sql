-- name: LockEventCapacity :one
SELECT capacity FROM events WHERE id = $1 FOR UPDATE;

-- name: CountRegistered :one
SELECT COUNT(*) FROM participations WHERE event_id = $1 AND status = 'registered';

-- name: CreateParticipation :one
INSERT INTO participations (full_name, email, student_number, event_id, student_id, status)
VALUES ($1, $2, $3, $4, $5, 'registered')
RETURNING *;

-- name: ListParticipants :many
SELECT full_name, email, student_number, status, registered_at
FROM participations
WHERE event_id = $1
ORDER BY registered_at ASC;

-- name: CancelParticipation :execrows
UPDATE participations SET status = 'cancelled'
WHERE event_id = $1 AND student_id = $2 AND status = 'registered';

-- name: ListStudentRegistrations :many
SELECT p.event_id, p.registered_at
FROM participations p
JOIN events e ON e.id = p.event_id
WHERE p.student_id = $1 AND p.status = 'registered'
ORDER BY e.event_time ASC;
