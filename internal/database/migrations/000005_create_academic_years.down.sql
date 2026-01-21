-- Remove academic_year_id from timetables
ALTER TABLE timetables
    DROP COLUMN IF EXISTS academic_year_id,
    DROP COLUMN IF EXISTS is_active;

-- Revert day_of_week back to INTEGER (this is a data-loss operation, be careful)
ALTER TABLE timetables
    ALTER COLUMN day_of_week TYPE INTEGER USING
        CASE
            WHEN day_of_week = 'MONDAY' THEN 1
            WHEN day_of_week = 'TUESDAY' THEN 2
            WHEN day_of_week = 'WEDNESDAY' THEN 3
            WHEN day_of_week = 'THURSDAY' THEN 4
            WHEN day_of_week = 'FRIDAY' THEN 5
            WHEN day_of_week = 'SATURDAY' THEN 6
            WHEN day_of_week = 'SUNDAY' THEN 7
            ELSE 1
        END,
    ALTER COLUMN start_time TYPE TIME USING start_time::TIME,
    ALTER COLUMN end_time TYPE TIME USING end_time::TIME;

-- Drop indexes
DROP INDEX IF EXISTS idx_timetables_academic_year_id;
DROP INDEX IF EXISTS idx_periods_institution_id;
DROP INDEX IF EXISTS idx_terms_academic_year_id;
DROP INDEX IF EXISTS idx_academic_years_is_current;
DROP INDEX IF EXISTS idx_academic_years_institution_id;

-- Drop unique constraint
ALTER TABLE academic_years DROP CONSTRAINT IF EXISTS unique_academic_year_name_per_institution;

-- Drop tables
DROP TABLE IF EXISTS periods;
DROP TABLE IF EXISTS terms;
DROP TABLE IF EXISTS academic_years;

