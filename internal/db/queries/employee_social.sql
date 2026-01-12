-- name: CreateEmployeeSocial :one
INSERT INTO employee_socials(
  employee_id,
  social_name,
  link_to_social
) VALUES (
  $1, $2, $3
) RETURNING *;

-- name: UpdateEmployeeSocial :one
UPDATE employee_socials 
SET 
  social_name = COALESCE(sqlc.narg('social_name'), social_name),
  link_to_social = COALESCE(sqlc.narg('link_to_social'), link_to_social),
  updated_at = now()
WHERE id = sqlc.arg('id')
RETURNING *;

-- name: DeleteEmployeeSocial :exec
DELETE FROM employee_socials
WHERE id = $1;

-- name: GetEmployeeSocialByID :one
SELECT *
FROM employee_socials
WHERE id = $1;

-- name: GetEmployeeSocialsByEmployeeID :many
SELECT *
FROM employee_socials
WHERE employee_id = $1; 
