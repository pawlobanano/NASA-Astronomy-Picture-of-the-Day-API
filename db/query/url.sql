-- name: CreateURL :one
INSERT INTO url (
  id, url
) VALUES (
  $1, $2
) RETURNING *;