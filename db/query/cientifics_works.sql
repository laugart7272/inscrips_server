-- name: GetCientificWork :one
SELECT * FROM cientifics_works
WHERE id = $1 LIMIT 1;

-- name: GetCientificWorkByAuthorId :many
SELECT * FROM cientifics_works
WHERE author_id = $1
LIMIT $2
OFFSET $3;

-- name: ListCientificWorks :many
SELECT * FROM cientifics_works
ORDER BY id
LIMIT $1
OFFSET $2;

-- name: CreateCientificWork :one
INSERT INTO cientifics_works(title, author_id, inscription_id, resume, file, presentation_type, exposition_type) 
VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING *;

-- name: UpdateCientificWork :one
UPDATE cientifics_works
  set title = $2, author_id = $3, resume = $4, file = $5, presentation_type = $6, exposition_type = $7, inscription_id = $8
WHERE id = $1 RETURNING *;

-- name: SetCientificWorkAdmitted :one
UPDATE cientifics_works
  set admitted=$2
WHERE id = $1 RETURNING *;

-- name: SetCientificWorkRecomendations :one
UPDATE cientifics_works
  set recommendations=$2
WHERE id = $1 RETURNING *;

-- name: GetCientificWorkAdmitted :one
SELECT * FROM cientifics_works WHERE (id = $1 AND admitted = true) LIMIT 1;

-- name: ListCientificsWorksAdmitted :one
SELECT * FROM cientifics_works WHERE (admitted = true) LIMIT 1;

-- #Pendig Event filter

-- name: DeleteCientificWork :exec
DELETE FROM cientifics_works
WHERE id = $1;
