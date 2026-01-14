-- name: CreateEmployeeDegree :one
INSERT INTO employee_degrees(
  employee_id,
  language_code,
  rf_institution_id,
  degree_level,
  institution_name,
  speciality,
  date_start,
  date_end,
  given_by,
  date_degree_recieved
) VALUES (
  $1, $2, $3, $4, $5, $6, $7, $8, $9, $10
) RETURNING *;

-- name: UpdateEmployeeDegree :one
UPDATE employee_degrees
SET
  rf_institution_id    = COALESCE(sqlc.narg('rf_institution_id'), rf_institution_id),
  degree_level          = COALESCE(sqlc.narg('degree_level'), degree_level),
  institution_name      = COALESCE(sqlc.narg('institution_name'), institution_name),
  speciality            = COALESCE(sqlc.narg('speciality'), speciality),
  date_start            = COALESCE(sqlc.narg('date_start'), date_start),
  date_end              = COALESCE(sqlc.narg('date_end'), date_end),
  given_by              = COALESCE(sqlc.narg('given_by'), given_by),
  date_degree_recieved  = COALESCE(sqlc.narg('date_degree_recieved'), date_degree_recieved),
  updated_at = now()
WHERE id = sqlc.arg('id')
RETURNING *;

-- name: DeleteEmployeeDegree :exec
DELETE FROM employee_degrees
WHERE id = $1;

-- name: GetEmployeeDegreeByID :one
SELECT *
FROM employee_degrees
WHERE id = $1;

-- name: GetEmployeeDegreesByEmployeeIDAndLanguageCode :many
SELECT *
FROM employee_degrees
WHERE employee_id = $1 AND language_code = $2
ORDER BY employee_degrees.date_degree_recieved DESC;
