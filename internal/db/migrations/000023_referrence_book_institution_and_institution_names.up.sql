-- 1) Create rf_institutions
CREATE TABLE IF NOT EXISTS rf_institutions (
  id BIGSERIAL PRIMARY KEY,

  postal_code TEXT,
  email TEXT,
  website TEXT,
  phone_number TEXT,

  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- 2) Create rf_institution_names
CREATE TABLE IF NOT EXISTS rf_institution_names (
  id BIGSERIAL PRIMARY KEY,
  rf_institution_id BIGINT NOT NULL
    REFERENCES rf_institutions(id) ON DELETE CASCADE,

  language_code CHAR(2) NOT NULL,
  name TEXT NOT NULL,
  name_norm TEXT NOT NULL,

  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),

  UNIQUE (rf_institution_id, language_code)
);

-- Indexes
CREATE INDEX IF NOT EXISTS idx_rf_institution_names_institution_id
  ON rf_institution_names (rf_institution_id);

CREATE INDEX IF NOT EXISTS idx_rf_institution_names_lang_norm
  ON rf_institution_names (language_code, name_norm);

-- 3) Add rf_institution_id column
ALTER TABLE employee_degrees
  ADD COLUMN IF NOT EXISTS rf_institution_id BIGINT;

-- Add FK constraint safely (Postgres has no ADD CONSTRAINT IF NOT EXISTS)
DO $$
BEGIN
  IF NOT EXISTS (
    SELECT 1 FROM pg_constraint WHERE conname = 'fk_employee_degrees_rf_institution'
  ) THEN
    ALTER TABLE employee_degrees
      ADD CONSTRAINT fk_employee_degrees_rf_institution
      FOREIGN KEY (rf_institution_id)
      REFERENCES rf_institutions(id);
  END IF;
END $$;

-- 4) Rename university_name -> institution_name
DO $$
BEGIN
  IF EXISTS (
    SELECT 1
    FROM information_schema.columns
    WHERE table_name='employee_degrees'
      AND column_name='university_name'
  ) AND NOT EXISTS (
    SELECT 1
    FROM information_schema.columns
    WHERE table_name='employee_degrees'
      AND column_name='institution_name'
  ) THEN
    ALTER TABLE employee_degrees RENAME COLUMN university_name TO institution_name;
  END IF;
END $$;

-- 5) Fill rf tables
DO $$
DECLARE
  rec RECORD;
  new_inst_id BIGINT;
BEGIN
  FOR rec IN
    SELECT DISTINCT
      ed.language_code,
      ed.institution_name AS name,
      lower(regexp_replace(trim(ed.institution_name), '\s+', ' ', 'g')) AS name_norm
    FROM employee_degrees ed
    WHERE ed.institution_name IS NOT NULL
      AND trim(ed.institution_name) <> ''
  LOOP
    INSERT INTO rf_institutions DEFAULT VALUES
    RETURNING id INTO new_inst_id;

    INSERT INTO rf_institution_names (
      rf_institution_id, language_code, name, name_norm
    ) VALUES (
      new_inst_id, rec.language_code, rec.name, rec.name_norm
    );
  END LOOP;
END $$;

-- 6) Backfill employee_degrees.rf_institution_id
UPDATE employee_degrees ed
SET rf_institution_id = rin.rf_institution_id,
    updated_at = now()
FROM rf_institution_names rin
WHERE rin.language_code = ed.language_code
  AND rin.name_norm = lower(regexp_replace(trim(ed.institution_name), '\s+', ' ', 'g'))
  AND ed.institution_name IS NOT NULL
  AND trim(ed.institution_name) <> '';

-- 7) Enforce NOT NULL (only if you’re sure every row matched)
ALTER TABLE employee_degrees
  ALTER COLUMN rf_institution_id SET NOT NULL;