-- name: CreateVerifyEmail :one
INSERT INTO verify_emails (
  name,
  last_name,
  email,
  user_id,
  secret_code
) VALUES (
  $1, $2, $3, $4, $5
) RETURNING *;

-- name: UpdateVerifyEmail :one
UPDATE verify_emails 
SET
  is_used = TRUE
WHERE
  id = $1
  AND secret_code = $2
  AND is_used = FALSE
  AND expired_at > now()
RETURNING *;