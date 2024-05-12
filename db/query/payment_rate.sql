-- name: GetPaymentAmount :one
SELECT amount FROM payment_rate WHERE (event_id = $1 AND user_type = $2 AND presentation_type = $3) LIMIT 1;

-- name: ListDetailsPaymentRates :many
SELECT payment_rate.created_at, events.name, user_type, presentation_type, amount
FROM payment_rate
INNER JOIN events ON payment_rate.event_id = events.id
LIMIT $1
OFFSET $2;