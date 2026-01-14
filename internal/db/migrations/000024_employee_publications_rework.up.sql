CREATE TABLE IF NOT EXISTS rf_publication_types (
  id BIGSERIAL,
  language_code CHAR(2) NOT NULL,
  name VARCHAR(255) NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),

  CONSTRAINT rf_publication_types_pkey
    PRIMARY KEY (id)
);

CREATE INDEX IF NOT EXISTS idx_rf_publication_types_language_code
  ON rf_publication_types (language_code);

INSERT INTO rf_publication_types (language_code, name) VALUES
  ('ru', 'Научная статья в журнале'),
  ('ru', 'Электронный препринт научной статьи'),
  ('ru', 'Дисертация'),
  ('ru', 'Монография'),
  ('ru', 'Глава в книге'),
  ('ru', 'Статья или тезис в сборнике трудов конференции'),
  ('ru', 'Учебное пособие и методические указания'),
  ('ru', 'Отчёт о НИР/НИОКР'),
  ('ru', 'Брошюра'),
  ('ru', 'Научно-популярная статья'),
  ('ru', 'Интервью'),
  ('tg', 'Мақолаи илмӣ дар маҷалла'),
  ('tg', 'Пешнашри электронии мақолаи илмӣ'),
  ('tg', 'Диссертатсия'),
  ('tg', 'Монография'),
  ('tg', 'Боб дар китоб'),
  ('tg', 'Мақола ё тезис дар маҷмуаи маводҳои конференсия'),
  ('tg', 'Васоити таълимӣ ва дастурҳои методӣ'),
  ('tg', 'Ҳисобот оид ба КИТ/КИТ ва ТК'),
  ('tg', 'Брошюра'),
  ('tg', 'Мақолаи илмӣ-оммавӣ'),
  ('tg', 'Мусоҳиба'),
  ('en', 'Scientific journal article'),
  ('en', 'Electronic preprint of a scientific article'),
  ('en', 'Dissertation'),
  ('en', 'Monograph'),
  ('en', 'Book chapter'),
  ('en', 'Article or thesis in conference proceedings'),
  ('en', 'Study guide and methodological guidelines'),
  ('en', 'R&D report'),
  ('en', 'Brochure'),
  ('en', 'Popular science article'),
  ('en', 'Interview');

CREATE TABLE IF NOT EXISTS employee_publications_new (
  id BIGSERIAL,
  employee_id BIGINT NOT NULL,
  language_code CHAR(2) NOT NULL,
  rf_publication_type_id BIGINT NOT NULL,
  name VARCHAR(255) NOT NULL,
  type VARCHAR(255) NOT NULL,
  authors TEXT,
  journal_name TEXT,
  volume TEXT,
  number TEXT,
  pages TEXT,
  year INTEGER,
  link VARCHAR(1023) NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),

  CONSTRAINT employee_publications_new_pkey
    PRIMARY KEY (id),
  CONSTRAINT fk_employees_employee_publications_new
    FOREIGN KEY (employee_id)
    REFERENCES employees (id)
    ON DELETE CASCADE,
  CONSTRAINT fk_rf_publication_types_employee_publications_new
    FOREIGN KEY (rf_publication_type_id)
    REFERENCES rf_publication_types (id)
    ON DELETE RESTRICT
);

INSERT INTO employee_publications_new (
  id,
  employee_id,
  language_code,
  rf_publication_type_id,
  name,
  type,
  link,
  created_at,
  updated_at
)
SELECT
  id,
  employee_id,
  'tg',
  (SELECT id FROM rf_publication_types WHERE language_code = 'tg' AND name = 'Мақолаи илмӣ дар маҷалла' LIMIT 1),
  publication_title,
  'Мақолаи илмӣ дар маҷалла',
  link_to_publication,
  created_at,
  updated_at
FROM employee_publications;

DROP TABLE employee_publications;

ALTER TABLE employee_publications_new
  RENAME TO employee_publications;

CREATE INDEX IF NOT EXISTS idx_employee_publications_language_code
  ON employee_publications (language_code);
