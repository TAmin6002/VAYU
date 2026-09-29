-- name: CreateStudent :one
INSERT INTO students (full_name, email, student_number, password_hash)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetStudentByEmail :one
SELECT * FROM students WHERE email = $1;

-- name: GetStudentByID :one
SELECT * FROM students WHERE id = $1;
