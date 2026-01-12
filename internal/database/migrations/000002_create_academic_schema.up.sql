-- Departments
CREATE TABLE IF NOT EXISTS departments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE,
    institution_id UUID NOT NULL REFERENCES institutions(id),
    name VARCHAR(100) NOT NULL,
    head_of_department_id UUID REFERENCES teachers(id),
    description TEXT
);

CREATE INDEX IF NOT EXISTS idx_departments_institution_id ON departments(institution_id);

-- Add constraint to teachers table if it exists
DO $$ 
BEGIN 
    IF EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'teachers') THEN
        ALTER TABLE teachers 
        ADD CONSTRAINT fk_teachers_department 
        FOREIGN KEY (department_id) 
        REFERENCES departments(id);
    END IF;
END $$;

-- Classes
CREATE TABLE IF NOT EXISTS classes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE,
    institution_id UUID NOT NULL REFERENCES institutions(id),
    name VARCHAR(50) NOT NULL,
    section_count INTEGER DEFAULT 1,
    class_teacher_id UUID REFERENCES teachers(id),
    capacity INTEGER
);

CREATE INDEX IF NOT EXISTS idx_classes_institution_id ON classes(institution_id);

-- Sections
CREATE TABLE IF NOT EXISTS sections (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE,
    class_id UUID NOT NULL REFERENCES classes(id),
    name VARCHAR(50) NOT NULL, -- 'A', 'B', 'Rose'
    room_number VARCHAR(20),
    capacity INTEGER
);

CREATE INDEX IF NOT EXISTS idx_sections_class_id ON sections(class_id);

-- Subjects
CREATE TABLE IF NOT EXISTS subjects (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE,
    institution_id UUID NOT NULL REFERENCES institutions(id),
    class_id UUID REFERENCES classes(id), -- Optional: Subject might be class-specific or global
    teacher_id UUID REFERENCES teachers(id),
    name VARCHAR(100) NOT NULL,
    code VARCHAR(20),
    is_elective BOOLEAN DEFAULT false,
    credit_hours DECIMAL(4,2)
);

CREATE INDEX IF NOT EXISTS idx_subjects_institution_id ON subjects(institution_id);
CREATE INDEX IF NOT EXISTS idx_subjects_class_id ON subjects(class_id);

-- Timetable
CREATE TABLE IF NOT EXISTS timetables (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    institution_id UUID NOT NULL REFERENCES institutions(id),
    class_id UUID NOT NULL REFERENCES classes(id),
    section_id UUID REFERENCES sections(id),
    day_of_week INTEGER, -- 1=Monday, 7=Sunday
    period_number INTEGER,
    subject_id UUID REFERENCES subjects(id),
    teacher_id UUID REFERENCES teachers(id),
    start_time TIME,
    end_time TIME,
    room_number VARCHAR(20)
);

CREATE INDEX IF NOT EXISTS idx_timetables_institution_id ON timetables(institution_id);
CREATE INDEX IF NOT EXISTS idx_timetables_class_section ON timetables(class_id, section_id);
CREATE INDEX IF NOT EXISTS idx_timetables_teacher_id ON timetables(teacher_id);
