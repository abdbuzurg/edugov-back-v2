CREATE TABLE IF NOT EXISTS employee_publications_old (
  id BIGSERIAL,
  employee_id BIGINT NOT NULL,
  language_code CHAR(2) NOT NULL,
  publication_title VARCHAR(255) NOT NULL,
  link_to_publication VARCHAR(511) NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),

  CONSTRAINT employee_publications_old_pkey
    PRIMARY KEY (id),
  CONSTRAINT fk_employees_employee_publications_old
    FOREIGN KEY (employee_id)
    REFERENCES employees (id)
    ON DELETE CASCADE
);

INSERT INTO employee_publications_old (
  id,
  employee_id,
  language_code,
  publication_title,
  link_to_publication,
  created_at,
  updated_at
)
SELECT
  id,
  employee_id,
  language_code,
  name,
  link,
  created_at,
  updated_at
FROM employee_publications;

DROP TABLE employee_publications;

ALTER TABLE employee_publications_old
  RENAME TO employee_publications;

CREATE INDEX IF NOT EXISTS idx_employee_publications_language_code
  ON employee_publications (language_code);

DROP TABLE rf_publication_types;
