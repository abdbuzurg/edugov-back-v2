BEGIN;

-- 1) Make rf_institution_id nullable before dropping constraints/column
DO $$
BEGIN
  IF EXISTS (
    SELECT 1
    FROM information_schema.columns
    WHERE table_name='employee_degrees'
      AND column_name='rf_institution_id'
  ) THEN
    ALTER TABLE employee_degrees
      ALTER COLUMN rf_institution_id DROP NOT NULL;
  END IF;
END $$;

-- 2) Drop FK + column
ALTER TABLE employee_degrees
  DROP CONSTRAINT IF EXISTS fk_employee_degrees_rf_institution;

ALTER TABLE employee_degrees
  DROP COLUMN IF EXISTS rf_institution_id;

-- 3) Rename institution_name -> university_name
ALTER TABLE employee_degrees
  RENAME COLUMN institution_name TO university_name;

-- 4) Drop RF tables
DROP TABLE IF EXISTS rf_institution_names;
DROP TABLE IF EXISTS rf_institutions;

COMMIT;
