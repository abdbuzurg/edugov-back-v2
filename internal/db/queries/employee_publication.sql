-- name: CreateEmployeePublication :one
INSERT INTO employee_publications(
  employee_id,
  language_code,
  rf_publication_type_id,
  name,
  type,
  authors,
  journal_name,
  volume,
  number,
  pages,
  year,
  link
) VALUES (
  $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12
) RETURNING *;

-- name: UpdateEmployeePublication :one
UPDATE employee_publications 
SET 
  rf_publication_type_id = COALESCE(sqlc.narg('rf_publication_type_id'), rf_publication_type_id),
  name = COALESCE(sqlc.narg('name'), name),
  type = COALESCE(sqlc.narg('type'), type),
  authors = COALESCE(sqlc.narg('authors'), authors),
  journal_name = COALESCE(sqlc.narg('journal_name'), journal_name),
  volume = COALESCE(sqlc.narg('volume'), volume),
  number = COALESCE(sqlc.narg('number'), number),
  pages = COALESCE(sqlc.narg('pages'), pages),
  year = COALESCE(sqlc.narg('year'), year),
  link = COALESCE(sqlc.narg('link'), link),
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
