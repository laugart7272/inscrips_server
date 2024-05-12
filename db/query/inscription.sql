-- name: GetInscription :one
SELECT * FROM inscriptions
WHERE id = $1 LIMIT 1;

-- name: ListInscriptions :many
SELECT * FROM inscriptions
ORDER BY id
LIMIT $1
OFFSET $2;

-- name: CreateInscription :one
INSERT INTO inscriptions(user_id, event_id) 
VALUES ($1, $2) RETURNING *;

-- name: UpdateInscription :one
UPDATE inscriptions
  set user_id = $2, event_id = $3
WHERE id = $1 RETURNING *;

-- name: DeleteInscription :exec
DELETE FROM inscriptions
WHERE id = $1;

-- name: FindInscription :one
SELECT * FROM inscriptions
WHERE (user_id = $1 AND event_id = $2) LIMIT 1; 

-- name: ListDetailsInscriptions :many
SELECT inscriptions.id, inscriptions.created_at, users.name, users.last_name, cientifics_works.title as work_title, events.name as event
FROM inscriptions
  INNER JOIN users ON inscriptions.user_id = users.id
  INNER JOIN cientifics_works ON inscriptions.id = cientifics_works.inscription_id
  INNER JOIN events ON inscriptions.event_id = events.id
ORDER BY 
  inscriptions.created_at
LIMIT $1
OFFSET $2;
