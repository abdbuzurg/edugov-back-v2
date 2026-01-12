-- name: CreateEmployeePatent :one
INSERT INTO employee_patents(
  employee_id,
  language_code,
  patent_title,
  description
) VALUES (
  $1, $2, $3, $4
) RETURNING *;

-- name: UpdateEmployeePatent :one
UPDATE employee_patents 
SET 
  patent_title = COALESCE(sqlc.narg('patent_title'), patent_title),
  description = COALESCE(sqlc.narg('description'), description),
  updated_at = now()
WHERE id = sqlc.arg('id')
RETURNING *;

-- name: DeleteEmployeePatent :exec
DELETE FROM employee_patents
WHERE id = $1;

-- name: GetEmployeePatentByID :one
SELECT *
FROM employee_patents
WHERE id = $1;

-- name: GetEmployeePatentsByEmployeeIDAndLanguageCode :many
SELECT *
FROM employee_patents
WHERE employee_id = $1 AND language_code = $2;
