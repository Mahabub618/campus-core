-- Enable pgcrypto extension for UUID generation
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- Institutions Table
CREATE TABLE IF NOT EXISTS institutions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE,
    name VARCHAR(255) NOT NULL,
    code VARCHAR(50) NOT NULL UNIQUE,
    address TEXT,
    phone VARCHAR(20),
    email VARCHAR(255),
    principal_name VARCHAR(255),
    established_year INTEGER,
    logo_url VARCHAR(500),
    academic_year VARCHAR(20),
    is_active BOOLEAN DEFAULT true
);

CREATE INDEX idx_institutions_deleted_at ON institutions(deleted_at);

-- Users Table
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE,
    email VARCHAR(255),
    phone VARCHAR(20),
    password_hash VARCHAR(255),
    role VARCHAR(50) NOT NULL,
    is_active BOOLEAN DEFAULT true,
    last_login_at TIMESTAMP WITH TIME ZONE,
    refresh_token VARCHAR(500),
    reset_token VARCHAR(255),
    reset_token_expiry TIMESTAMP WITH TIME ZONE
);

CREATE UNIQUE INDEX idx_users_email ON users(email) WHERE deleted_at IS NULL;
CREATE INDEX idx_users_deleted_at ON users(deleted_at);

-- User Profiles Table
CREATE TABLE IF NOT EXISTS user_profiles (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE,
    user_id UUID NOT NULL REFERENCES users(id),
    institution_id UUID REFERENCES institutions(id),
    first_name VARCHAR(100),
    last_name VARCHAR(100),
    date_of_birth TIMESTAMP WITH TIME ZONE,
    gender VARCHAR(10),
    address TEXT,
    profile_image_url VARCHAR(500),
    employee_id VARCHAR(50),
    admission_number VARCHAR(50),
    occupation VARCHAR(100)
);

CREATE UNIQUE INDEX idx_user_profiles_user_id ON user_profiles(user_id);
CREATE INDEX idx_user_profiles_institution_id ON user_profiles(institution_id);
CREATE INDEX idx_user_profiles_deleted_at ON user_profiles(deleted_at);

-- Teachers Table
CREATE TABLE IF NOT EXISTS teachers (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE,
    institution_id UUID NOT NULL REFERENCES institutions(id),
    user_id UUID NOT NULL REFERENCES users(id),
    qualifications TEXT[],
    joining_date TIMESTAMP WITH TIME ZONE,
    department_id UUID -- References departments table (to be created)
);

CREATE UNIQUE INDEX idx_teachers_user_id ON teachers(user_id);
CREATE INDEX idx_teachers_institution_id ON teachers(institution_id);
CREATE INDEX idx_teachers_deleted_at ON teachers(deleted_at);

-- Students Table
CREATE TABLE IF NOT EXISTS students (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE,
    institution_id UUID NOT NULL REFERENCES institutions(id),
    user_id UUID NOT NULL REFERENCES users(id),
    class_id UUID, -- References classes table (to be created)
    section_id UUID, -- References sections table (to be created)
    roll_number INTEGER,
    admission_date TIMESTAMP WITH TIME ZONE,
    blood_group VARCHAR(5),
    medical_info TEXT
);

CREATE UNIQUE INDEX idx_students_user_id ON students(user_id);
CREATE INDEX idx_students_institution_id ON students(institution_id);
CREATE INDEX idx_students_deleted_at ON students(deleted_at);

-- Parents Table
CREATE TABLE IF NOT EXISTS parents (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE,
    institution_id UUID NOT NULL REFERENCES institutions(id),
    user_id UUID NOT NULL REFERENCES users(id),
    occupation VARCHAR(100),
    office_address TEXT,
    emergency_contact VARCHAR(20)
);

CREATE UNIQUE INDEX idx_parents_user_id ON parents(user_id);
CREATE INDEX idx_parents_institution_id ON parents(institution_id);
CREATE INDEX idx_parents_deleted_at ON parents(deleted_at);

-- Parent Student Relations Table
CREATE TABLE IF NOT EXISTS parent_student_relations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE,
    parent_id UUID NOT NULL REFERENCES parents(id),
    student_id UUID NOT NULL REFERENCES students(id),
    relationship VARCHAR(50),
    is_primary BOOLEAN DEFAULT false
);

CREATE INDEX idx_parent_student_relations_parent_id ON parent_student_relations(parent_id);
CREATE INDEX idx_parent_student_relations_student_id ON parent_student_relations(student_id);
CREATE INDEX idx_parent_student_relations_deleted_at ON parent_student_relations(deleted_at);

-- Accountants Table
CREATE TABLE IF NOT EXISTS accountants (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE,
    institution_id UUID NOT NULL REFERENCES institutions(id),
    user_id UUID NOT NULL REFERENCES users(id),
    joining_date TIMESTAMP WITH TIME ZONE,
    qualification VARCHAR(255)
);

CREATE UNIQUE INDEX idx_accountants_user_id ON accountants(user_id);
CREATE INDEX idx_accountants_institution_id ON accountants(institution_id);
CREATE INDEX idx_accountants_deleted_at ON accountants(deleted_at);
