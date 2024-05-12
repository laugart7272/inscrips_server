-- name: GetPayment :one
SELECT * FROM payment
WHERE id = $1 LIMIT 1;

-- name: ListPayments :many
SELECT * FROM payment
ORDER BY id
LIMIT $1
OFFSET $2;

-- name: CreatePayment :one
INSERT INTO payment(user_id, event_id, user_type, presentation_type, amount, inscription_id) 
VALUES ($1, $2, $3, $4, $5, $6) RETURNING *;

-- name: UpdatePayment :one
UPDATE payment
  set user_id = $2, event_id = $3, user_type = $4, presentation_type = $5, amount = $6, inscription_id = $7
WHERE id = $1 RETURNING *;

-- name: DeletePayment :exec
DELETE FROM payment
WHERE id = $1;

-- name: ListDetailsPayments :many
SELECT payment.created_at, users.name, users.last_name, events.name as event, payment.user_type as user_type, presentation_type, amount, inscription_id
FROM payment
INNER JOIN users ON payment.user_id = users.id
INNER JOIN events ON payment.event_id = events.id
ORDER BY 
  payment.created_at
LIMIT $1
OFFSET $2;
