-- name: GetUser :one
SELECT * FROM users
WHERE id = $1 LIMIT 1;

-- name: GetUserEmail :one
SELECT * FROM users
WHERE email = $1 LIMIT 1;

-- name: ListUsers :many
SELECT * FROM users
ORDER BY id
LIMIT $1
OFFSET $2;

-- name: CreateUser :one
INSERT INTO users(name, last_name, email, phone, hashed_password, user_type) 
VALUES ($1, $2, $3, $4, $5, $6) RETURNING *;

-- name: UpdateUser :one
UPDATE users
  SET 
    name = $2, 
    last_name = $3, 
    email = $4, 
    phone = $5, 
    hashed_password = $6, 
    user_type = $7, 
    is_email_verified = $8,
    role = $9,
    avatar_path = $10
WHERE id = $1 RETURNING *;

-- name: DeleteUser :exec
DELETE FROM users
WHERE id = $1;
