-- name: CreateEmployeeRefresherCourse :one
INSERT INTO employee_refresher_courses(
  employee_id,
  language_code,
  course_title,
  date_start,
  date_end
) VALUES (
  $1, $2, $3, $4, $5
) RETURNING *;

-- name: UpdateEmployeeRefresherCourse :one
UPDATE employee_refresher_courses
SET 
  course_title = COALESCE(sqlc.narg('course_title'), course_title),
  date_start = COALESCE(sqlc.narg('date_start'), date_start),
  date_end = COALESCE(sqlc.narg('date_end'), date_end),
  updated_at = now()
WHERE id = sqlc.arg('id')
RETURNING *;

-- name: DeleteEmployeeRefresherCourse :exec
DELETE FROM employee_refresher_courses
WHERE id = $1;

-- name: GetEmployeeRefresherCourseByID :one
SELECT *
FROM employee_refresher_courses
WHERE id = $1;

-- name: GetEmployeeRefresherCoursesByEmployeeIDAndLanguageCode :many
SELECT *
FROM employee_refresher_courses
WHERE employee_id = $1 AND language_code = $2;
