-- name: GetAdminByEmail :one
SELECT * FROM admins WHERE email = $1;

-- name: GetAdminByID :one
SELECT * FROM admins WHERE id = $1;
