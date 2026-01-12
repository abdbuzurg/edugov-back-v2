-- name: CreateEmployeePublication :one
INSERT INTO employee_publications(
  employee_id,
  language_code,
  publication_title,
  link_to_publication
) VALUES (
  $1, $2, $3, $4
) RETURNING *;

-- name: UpdateEmployeePublication :one
UPDATE employee_publications 
SET 
  publication_title = COALESCE(sqlc.narg('publication_title'), publication_title),
  link_to_publication = COALESCE(sqlc.narg('link_to_publication'), link_to_publication),
  updated_at = now()
WHERE id = sqlc.arg('id')
RETURNING *;

-- name: DeleteEmployeePublication :exec
DELETE FROM employee_publications
WHERE id = $1;

-- name: GetEmployeePublicationByID :one
SELECT *
FROM employee_publications
WHERE id = $1;

-- name: GetEmployeePublicationsByEmployeeIDAndLanguageCode :many
SELECT *
FROM employee_publications
WHERE employee_id = $1 AND language_code = $2;
