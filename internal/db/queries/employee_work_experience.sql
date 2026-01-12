-- name: CreateEmployeeWorkExperience :one
INSERT INTO employee_work_experiences(
  employee_id,
  language_code,
  workplace,
  job_title,
  description,
  date_start,
  date_end,
  on_going
) VALUES (
  $1, $2, $3, $4, $5, $6, $7, $8
) RETURNING *;

-- name: UpdateEmployeeWorkExperience :one
UPDATE employee_work_experiences
SET 
  workplace = COALESCE(sqlc.narg('workplace'), workplace),
  job_title = COALESCE(sqlc.narg('job_title'), job_title),
  description = COALESCE(sqlc.narg('description'), description),
  date_start = COALESCE(sqlc.narg('date_start'), date_start),
  date_end = COALESCE(sqlc.narg('date_end'), date_end),
  on_going = COALESCE(sqlc.narg('on_going'), on_going),
  updated_at = now()
WHERE id = sqlc.arg('id')
RETURNING *;

-- name: DeleteEmployeeWorkExperience :exec
delete from employee_work_experiences
where id = $1
;

-- name: GetEmployeeWorkExperienceByID :one
select *
from employee_work_experiences
where id = $1
;

-- name: GetEmployeeWorkExperiencesByEmployeeIDAndLanguageCode :many
select *
from employee_work_experiences
where employee_id = $1 and language_code = $2
order by
    employee_work_experiences.on_going desc, employee_work_experiences.date_end desc
;
