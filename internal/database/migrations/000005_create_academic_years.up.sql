-- Academic Years
CREATE TABLE IF NOT EXISTS academic_years (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE,
    institution_id UUID NOT NULL REFERENCES institutions(id),
    name VARCHAR(50) NOT NULL,
    start_date DATE NOT NULL,
    end_date DATE NOT NULL,
    is_current BOOLEAN DEFAULT false,
    description TEXT
);

CREATE INDEX IF NOT EXISTS idx_academic_years_institution_id ON academic_years(institution_id);
CREATE INDEX IF NOT EXISTS idx_academic_years_is_current ON academic_years(institution_id, is_current);

-- Terms (within academic years)
CREATE TABLE IF NOT EXISTS terms (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE,
    academic_year_id UUID NOT NULL REFERENCES academic_years(id),
    name VARCHAR(50) NOT NULL,
    start_date DATE NOT NULL,
    end_date DATE NOT NULL,
    is_current BOOLEAN DEFAULT false
);

CREATE INDEX IF NOT EXISTS idx_terms_academic_year_id ON terms(academic_year_id);

-- Periods (time slots in a school day)
CREATE TABLE IF NOT EXISTS periods (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE,
    institution_id UUID NOT NULL REFERENCES institutions(id),
    name VARCHAR(50) NOT NULL,
    start_time TIME NOT NULL,
    end_time TIME NOT NULL,
    "order" INTEGER NOT NULL,
    is_break BOOLEAN DEFAULT false
);

CREATE INDEX IF NOT EXISTS idx_periods_institution_id ON periods(institution_id);

-- Add academic_year_id to timetables if not exists
ALTER TABLE timetables
    ADD COLUMN IF NOT EXISTS academic_year_id UUID REFERENCES academic_years(id),
    ADD COLUMN IF NOT EXISTS is_active BOOLEAN DEFAULT true,
    ALTER COLUMN day_of_week TYPE VARCHAR(20) USING
        CASE
            WHEN day_of_week = 1 THEN 'MONDAY'
            WHEN day_of_week = 2 THEN 'TUESDAY'
            WHEN day_of_week = 3 THEN 'WEDNESDAY'
            WHEN day_of_week = 4 THEN 'THURSDAY'
            WHEN day_of_week = 5 THEN 'FRIDAY'
            WHEN day_of_week = 6 THEN 'SATURDAY'
            WHEN day_of_week = 7 THEN 'SUNDAY'
            ELSE 'MONDAY'
        END,
    ALTER COLUMN start_time TYPE VARCHAR(10) USING start_time::text,
    ALTER COLUMN end_time TYPE VARCHAR(10) USING end_time::text;

-- Create index for academic year on timetables
CREATE INDEX IF NOT EXISTS idx_timetables_academic_year_id ON timetables(academic_year_id);

-- Add unique constraint for academic year name within institution
ALTER TABLE academic_years ADD CONSTRAINT unique_academic_year_name_per_institution
    UNIQUE (institution_id, name);

