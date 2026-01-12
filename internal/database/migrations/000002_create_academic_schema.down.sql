DROP TABLE IF EXISTS timetables;
DROP TABLE IF EXISTS subjects;
DROP TABLE IF EXISTS sections;
DROP TABLE IF EXISTS classes;

-- Drop foreign key from teachers before dropping departments
DO $$ 
BEGIN 
    IF EXISTS (SELECT 1 FROM information_schema.table_constraints WHERE constraint_name = 'fk_teachers_department') THEN
        ALTER TABLE teachers DROP CONSTRAINT fk_teachers_department;
    END IF;
END $$;

DROP TABLE IF EXISTS departments;
