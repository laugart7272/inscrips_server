-- name: GetEvent :one
SELECT * FROM events
WHERE id = $1 LIMIT 1;

-- name: ListEvents :many
SELECT * FROM events
ORDER BY id
LIMIT $1
OFFSET $2;

-- name: CreateEvent :one
INSERT INTO events(name, init_date, end_date) 
VALUES ($1, $2, $3) RETURNING *;

-- name: UpdateEvent :one
UPDATE events
  set name = $2, init_date = $3, end_date = $4, created_at = now()
WHERE id = $1 RETURNING *;

-- name: DeleteEvent :exec
DELETE FROM events
WHERE id = $1;

-- name: ActiveEvent :one
SELECT * FROM events WHERE active = true;
