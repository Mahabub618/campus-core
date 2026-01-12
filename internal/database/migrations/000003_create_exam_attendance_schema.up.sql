-- Attendance
CREATE TABLE IF NOT EXISTS attendance (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    institution_id UUID NOT NULL REFERENCES institutions(id),
    student_id UUID REFERENCES students(id),
    date DATE NOT NULL,
    status VARCHAR(20) NOT NULL, -- 'PRESENT', 'ABSENT', 'LATE', 'HALF_DAY'
    marked_by UUID REFERENCES users(id),
    remarks TEXT
);

CREATE INDEX IF NOT EXISTS idx_attendance_student_date ON attendance(student_id, date);
CREATE INDEX IF NOT EXISTS idx_attendance_institution_date ON attendance(institution_id, date);

-- Leave Types
CREATE TABLE IF NOT EXISTS leave_types (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    institution_id UUID NOT NULL REFERENCES institutions(id),
    name VARCHAR(100) NOT NULL,
    description TEXT,
    max_days_per_year INTEGER DEFAULT 0,
    is_paid BOOLEAN DEFAULT true,
    applicable_to VARCHAR(50)[], -- ['TEACHER', 'STUDENT', 'STAFF']
    requires_document BOOLEAN DEFAULT false,
    is_active BOOLEAN DEFAULT true
);

-- Leave Applications
CREATE TABLE IF NOT EXISTS leaves (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    institution_id UUID NOT NULL REFERENCES institutions(id),
    user_id UUID NOT NULL REFERENCES users(id),
    leave_type_id UUID REFERENCES leave_types(id),
    start_date DATE NOT NULL,
    end_date DATE NOT NULL,
    total_days INTEGER NOT NULL,
    reason TEXT NOT NULL,
    document_urls VARCHAR(500)[],
    status VARCHAR(20) DEFAULT 'PENDING',
    applied_for_user_id UUID REFERENCES users(id),
    approved_by UUID REFERENCES users(id),
    approved_at TIMESTAMP WITH TIME ZONE,
    rejection_reason TEXT
);

CREATE INDEX IF NOT EXISTS idx_leaves_user_id ON leaves(user_id);
CREATE INDEX IF NOT EXISTS idx_leaves_institution_status ON leaves(institution_id, status);

-- Leave Balances
CREATE TABLE IF NOT EXISTS leave_balances (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    institution_id UUID NOT NULL REFERENCES institutions(id),
    user_id UUID NOT NULL REFERENCES users(id),
    leave_type_id UUID NOT NULL REFERENCES leave_types(id),
    academic_year VARCHAR(20) NOT NULL,
    total_allowed INTEGER,
    used INTEGER DEFAULT 0,
    -- remaining generated column not supported in all postgres versions or complex, easier to manage in app logic or trigger
    remaining INTEGER, 
    UNIQUE(user_id, leave_type_id, academic_year)
);

-- Exams
CREATE TABLE IF NOT EXISTS exams (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    institution_id UUID NOT NULL REFERENCES institutions(id),
    class_id UUID REFERENCES classes(id),
    name VARCHAR(100) NOT NULL,
    exam_type VARCHAR(50), -- 'TERM', 'UNIT', 'FINAL'
    start_date DATE,
    end_date DATE,
    total_marks DECIMAL(6,2)
);

CREATE INDEX IF NOT EXISTS idx_exams_institution_class ON exams(institution_id, class_id);

-- Exam Results
CREATE TABLE IF NOT EXISTS exam_results (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    institution_id UUID NOT NULL REFERENCES institutions(id),
    exam_id UUID NOT NULL REFERENCES exams(id),
    student_id UUID NOT NULL REFERENCES students(id),
    subject_id UUID NOT NULL REFERENCES subjects(id),
    marks_obtained DECIMAL(5,2),
    grade VARCHAR(5),
    percentage DECIMAL(5,2),
    rank_in_class INTEGER,
    remarks TEXT
);

CREATE INDEX IF NOT EXISTS idx_exam_results_exam_student ON exam_results(exam_id, student_id);

-- Assignments
CREATE TABLE IF NOT EXISTS assignments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    institution_id UUID NOT NULL REFERENCES institutions(id),
    class_id UUID REFERENCES classes(id),
    subject_id UUID REFERENCES subjects(id),
    teacher_id UUID REFERENCES teachers(id),
    title VARCHAR(255) NOT NULL,
    description TEXT,
    due_date TIMESTAMP WITH TIME ZONE,
    total_marks DECIMAL(5,2),
    attachment_urls VARCHAR(500)[]
);

-- Student Assignments (Submissions)
CREATE TABLE IF NOT EXISTS student_assignments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    assignment_id UUID NOT NULL REFERENCES assignments(id),
    student_id UUID NOT NULL REFERENCES students(id),
    submitted_at TIMESTAMP WITH TIME ZONE,
    marks_obtained DECIMAL(5,2),
    feedback TEXT,
    submission_url VARCHAR(500),
    status VARCHAR(20) DEFAULT 'SUBMITTED'
);
