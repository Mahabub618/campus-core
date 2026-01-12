-- Fee Structures
CREATE TABLE IF NOT EXISTS fee_structures (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    institution_id UUID NOT NULL REFERENCES institutions(id),
    class_id UUID REFERENCES classes(id),
    name VARCHAR(100) NOT NULL,
    academic_year VARCHAR(20),
    total_amount DECIMAL(10,2),
    due_date DATE,
    is_active BOOLEAN DEFAULT true
);

-- Fee Payments
CREATE TABLE IF NOT EXISTS fee_payments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    institution_id UUID NOT NULL REFERENCES institutions(id),
    student_id UUID NOT NULL REFERENCES students(id),
    fee_structure_id UUID REFERENCES fee_structures(id),
    amount_paid DECIMAL(10,2),
    payment_date DATE,
    payment_mode VARCHAR(50), -- 'CASH', 'ONLINE', etc.
    transaction_id VARCHAR(100),
    collected_by UUID REFERENCES users(id),
    receipt_number VARCHAR(50)
);

CREATE INDEX IF NOT EXISTS idx_fee_payments_student ON fee_payments(student_id);

-- Expenses
CREATE TABLE IF NOT EXISTS expenses (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    institution_id UUID NOT NULL REFERENCES institutions(id),
    category VARCHAR(100),
    amount DECIMAL(10,2),
    description TEXT,
    expense_date DATE,
    approved_by UUID REFERENCES users(id),
    payment_mode VARCHAR(50),
    receipt_url VARCHAR(500)
);

-- Salaries
CREATE TABLE IF NOT EXISTS salaries (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    institution_id UUID NOT NULL REFERENCES institutions(id),
    employee_id UUID NOT NULL REFERENCES users(id),
    month VARCHAR(7), -- YYYY-MM
    basic_salary DECIMAL(10,2),
    allowances DECIMAL(10,2),
    deductions DECIMAL(10,2),
    net_salary DECIMAL(10,2),
    payment_status VARCHAR(20), -- 'PENDING', 'PAID'
    paid_date DATE,
    transaction_id VARCHAR(100)
);

-- Library: Book Categories
CREATE TABLE IF NOT EXISTS book_categories (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    institution_id UUID NOT NULL REFERENCES institutions(id),
    name VARCHAR(100) NOT NULL,
    description TEXT,
    parent_category_id UUID REFERENCES book_categories(id),
    is_active BOOLEAN DEFAULT true
);

-- Library: Books
CREATE TABLE IF NOT EXISTS books (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    institution_id UUID NOT NULL REFERENCES institutions(id),
    category_id UUID REFERENCES book_categories(id),
    title VARCHAR(255) NOT NULL,
    author VARCHAR(255),
    isbn VARCHAR(20),
    publisher VARCHAR(255),
    publication_year INTEGER,
    edition VARCHAR(50),
    description TEXT,
    cover_image_url VARCHAR(500),
    total_copies INTEGER DEFAULT 1,
    available_copies INTEGER DEFAULT 1,
    location VARCHAR(100),
    language VARCHAR(50) DEFAULT 'English',
    is_available BOOLEAN DEFAULT true
);

CREATE INDEX IF NOT EXISTS idx_books_institution_category ON books(institution_id, category_id);

-- Library: Borrowings
CREATE TABLE IF NOT EXISTS book_borrowings (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    institution_id UUID NOT NULL REFERENCES institutions(id),
    book_id UUID NOT NULL REFERENCES books(id),
    user_id UUID NOT NULL REFERENCES users(id),
    borrowed_date DATE NOT NULL,
    due_date DATE NOT NULL,
    returned_date DATE,
    status VARCHAR(20) DEFAULT 'BORROWED',
    issued_by UUID REFERENCES users(id),
    received_by UUID REFERENCES users(id),
    notes TEXT
);

-- Library: Fines
CREATE TABLE IF NOT EXISTS library_fines (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    institution_id UUID NOT NULL REFERENCES institutions(id),
    borrowing_id UUID REFERENCES book_borrowings(id),
    user_id UUID NOT NULL REFERENCES users(id),
    amount DECIMAL(10,2) NOT NULL,
    reason VARCHAR(50),
    days_overdue INTEGER,
    status VARCHAR(20) DEFAULT 'UNPAID',
    paid_at TIMESTAMP WITH TIME ZONE,
    collected_by UUID REFERENCES users(id),
    waived_by UUID REFERENCES users(id),
    waiver_reason TEXT
);

-- Library: Settings
CREATE TABLE IF NOT EXISTS library_settings (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    institution_id UUID NOT NULL REFERENCES institutions(id) UNIQUE,
    max_books_per_student INTEGER DEFAULT 2,
    max_books_per_teacher INTEGER DEFAULT 5,
    loan_period_days INTEGER DEFAULT 14,
    fine_per_day DECIMAL(5,2) DEFAULT 1.00,
    max_renewals INTEGER DEFAULT 2
);

-- Events
CREATE TABLE IF NOT EXISTS events (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    institution_id UUID NOT NULL REFERENCES institutions(id),
    title VARCHAR(255) NOT NULL,
    description TEXT,
    event_type VARCHAR(50),
    start_datetime TIMESTAMP WITH TIME ZONE NOT NULL,
    end_datetime TIMESTAMP WITH TIME ZONE NOT NULL,
    location VARCHAR(255),
    is_all_day BOOLEAN DEFAULT false,
    target_audience VARCHAR(50)[],
    target_classes UUID[],
    organizer_id UUID REFERENCES users(id),
    attachment_urls VARCHAR(500)[],
    is_mandatory BOOLEAN DEFAULT false,
    is_active BOOLEAN DEFAULT true
);

CREATE INDEX IF NOT EXISTS idx_events_institution_date ON events(institution_id, start_datetime);

-- Event Participants
CREATE TABLE IF NOT EXISTS event_participants (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    event_id UUID NOT NULL REFERENCES events(id),
    user_id UUID NOT NULL REFERENCES users(id),
    status VARCHAR(20) DEFAULT 'INVITED',
    responded_at TIMESTAMP WITH TIME ZONE,
    attended_at TIMESTAMP WITH TIME ZONE,
    UNIQUE(event_id, user_id)
);

-- Holidays
CREATE TABLE IF NOT EXISTS holidays (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    institution_id UUID NOT NULL REFERENCES institutions(id),
    name VARCHAR(255) NOT NULL,
    date DATE NOT NULL,
    description TEXT,
    holiday_type VARCHAR(50),
    is_recurring BOOLEAN DEFAULT false,
    academic_year VARCHAR(20),
    UNIQUE(institution_id, date, name)
);

-- Academic Calendar
CREATE TABLE IF NOT EXISTS academic_calendar (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    institution_id UUID NOT NULL REFERENCES institutions(id),
    academic_year VARCHAR(20),
    event_type VARCHAR(50),
    title VARCHAR(255),
    start_date DATE NOT NULL,
    end_date DATE,
    description TEXT,
    is_active BOOLEAN DEFAULT true
);

-- Communication: Notices
CREATE TABLE IF NOT EXISTS notices (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    institution_id UUID NOT NULL REFERENCES institutions(id),
    title VARCHAR(255),
    content TEXT,
    priority VARCHAR(20),
    target_audience VARCHAR(50)[],
    published_by UUID REFERENCES users(id),
    published_at TIMESTAMP WITH TIME ZONE,
    expiry_date DATE,
    attachment_urls VARCHAR(500)[]
);

-- Communication: Messages
CREATE TABLE IF NOT EXISTS messages (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    sender_id UUID REFERENCES users(id),
    receiver_id UUID REFERENCES users(id),
    message TEXT,
    sent_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    read_at TIMESTAMP WITH TIME ZONE,
    message_type VARCHAR(20)
);

CREATE INDEX IF NOT EXISTS idx_messages_sender_receiver ON messages(sender_id, receiver_id);

-- Communication: Announcements
CREATE TABLE IF NOT EXISTS announcements (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    institution_id UUID NOT NULL REFERENCES institutions(id),
    title VARCHAR(255),
    content TEXT,
    announced_by UUID REFERENCES users(id),
    announcement_date DATE,
    is_active BOOLEAN DEFAULT true
);
