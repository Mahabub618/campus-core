# Smart School Management System (SSMS) - Backend Implementation Plan

## Project Overview

A comprehensive multi-tenant SaaS school management system backend built with Golang, featuring RESTful APIs, WebSocket support for real-time communication, JWT-based authentication, and role-based access control for 6 user roles.

### Tech Stack

| Category | Technology |
|----------|------------|
| **Language** | Go 1.24 |
| **Web Framework** | Gin / Echo (recommended: Gin) |
| **Database** | PostgreSQL 16+ |
| **ORM** | GORM |
| **Authentication** | JWT (golang-jwt/jwt/v5) |
| **WebSocket** | Gorilla WebSocket |
| **Caching** | Redis |
| **Validation** | go-playground/validator |
| **API Documentation** | Swaggo/swag (Swagger) |
| **Migration** | golang-migrate/migrate |
| **Testing** | testify, mockery |
| **Configuration** | Viper |
| **Logging** | Zap / Zerolog |
| **Task Queue** | Asynq (Redis-based) |

### Key Features

- ✅ Multi-Institution Support (SaaS model with tenant isolation)
- ✅ Role-based Access Control (6 roles: Super Admin, Admin, Teacher, Student, Parent, Accountant)
- ✅ JWT Authentication with Refresh Tokens
- ✅ Real-time Notifications via WebSocket
- ✅ RESTful API (Base URL: `/api/v1`)
- ✅ Database Migrations
- ✅ Request Validation & Sanitization
- ✅ Comprehensive Error Handling
- ✅ API Rate Limiting
- ✅ Audit Logging

---

## 1. Project Structure

```
campus-core/
├── cmd/
│   └── server/
│       └── main.go                    # Application entry point
│
├── internal/
│   ├── config/
│   │   └── config.go                  # Configuration management (Viper)
│   │
│   ├── database/
│   │   ├── database.go                # Database connection
│   │   ├── redis.go                   # Redis connection
│   │   └── migrations/                # SQL migration files
│   │       ├── 000001_create_users_table.up.sql
│   │       ├── 000001_create_users_table.down.sql
│   │       └── ...
│   │
│   ├── models/                        # GORM models
│   │   ├── user.go
│   │   ├── institution.go
│   │   ├── teacher.go
│   │   ├── student.go
│   │   ├── parent.go
│   │   ├── class.go
│   │   ├── section.go
│   │   ├── subject.go
│   │   ├── attendance.go
│   │   ├── assignment.go
│   │   ├── exam.go
│   │   ├── fee.go
│   │   ├── notice.go
│   │   ├── message.go
│   │   ├── leave.go
│   │   ├── library.go
│   │   └── event.go
│   │
│   ├── dto/                           # Data Transfer Objects
│   │   ├── request/
│   │   │   ├── auth_request.go
│   │   │   ├── user_request.go
│   │   │   └── ...
│   │   └── response/
│   │       ├── auth_response.go
│   │       ├── user_response.go
│   │       └── ...
│   │
│   ├── repository/                    # Data access layer
│   │   ├── repository.go              # Base repository interface
│   │   ├── user_repository.go
│   │   ├── institution_repository.go
│   │   ├── teacher_repository.go
│   │   ├── student_repository.go
│   │   ├── parent_repository.go
│   │   ├── academic_repository.go
│   │   ├── attendance_repository.go
│   │   ├── assessment_repository.go
│   │   ├── finance_repository.go
│   │   ├── communication_repository.go
│   │   ├── leave_repository.go
│   │   ├── library_repository.go
│   │   └── event_repository.go
│   │
│   ├── service/                       # Business logic layer
│   │   ├── auth_service.go
│   │   ├── user_service.go
│   │   ├── institution_service.go
│   │   ├── teacher_service.go
│   │   ├── student_service.go
│   │   ├── parent_service.go
│   │   ├── academic_service.go
│   │   ├── attendance_service.go
│   │   ├── assessment_service.go
│   │   ├── finance_service.go
│   │   ├── communication_service.go
│   │   ├── leave_service.go
│   │   ├── library_service.go
│   │   ├── event_service.go
│   │   ├── report_service.go
│   │   └── websocket_service.go
│   │
│   ├── handler/                       # HTTP handlers (controllers)
│   │   ├── auth_handler.go
│   │   ├── user_handler.go
│   │   ├── institution_handler.go
│   │   ├── teacher_handler.go
│   │   ├── student_handler.go
│   │   ├── parent_handler.go
│   │   ├── academic_handler.go
│   │   ├── attendance_handler.go
│   │   ├── assessment_handler.go
│   │   ├── finance_handler.go
│   │   ├── communication_handler.go
│   │   ├── leave_handler.go
│   │   ├── library_handler.go
│   │   ├── event_handler.go
│   │   ├── report_handler.go
│   │   └── websocket_handler.go
│   │
│   ├── middleware/
│   │   ├── auth_middleware.go         # JWT authentication
│   │   ├── rbac_middleware.go         # Role-based access control
│   │   ├── tenant_middleware.go       # Multi-tenancy (X-Institution-ID)
│   │   ├── cors_middleware.go         # CORS handling
│   │   ├── logging_middleware.go      # Request/response logging
│   │   ├── ratelimit_middleware.go    # Rate limiting
│   │   └── recovery_middleware.go     # Panic recovery
│   │
│   ├── router/
│   │   ├── router.go                  # Main router setup
│   │   ├── auth_routes.go
│   │   ├── user_routes.go
│   │   ├── institution_routes.go
│   │   ├── academic_routes.go
│   │   ├── attendance_routes.go
│   │   ├── assessment_routes.go
│   │   ├── finance_routes.go
│   │   ├── communication_routes.go
│   │   ├── leave_routes.go
│   │   ├── library_routes.go
│   │   ├── event_routes.go
│   │   └── report_routes.go
│   │
│   ├── websocket/
│   │   ├── hub.go                     # WebSocket connection manager
│   │   ├── client.go                  # WebSocket client
│   │   └── message.go                 # WebSocket message types
│   │
│   └── utils/
│       ├── jwt.go                     # JWT utilities
│       ├── password.go                # Password hashing (bcrypt)
│       ├── validator.go               # Custom validators
│       ├── pagination.go              # Pagination helpers
│       ├── response.go                # Standard API responses
│       └── errors.go                  # Custom error types
│
├── pkg/                               # Shared packages (if needed)
│   └── logger/
│       └── logger.go
│
├── docs/                              # Swagger documentation
│   └── swagger.json
│
├── scripts/
│   ├── migrate.sh                     # Migration script
│   └── seed.sh                        # Database seeder
│
├── .env.example
├── .gitignore
├── Dockerfile
├── docker-compose.yml
├── Makefile
├── go.mod
├── go.sum
└── README.md
```

---

## 2. Multi-Tenancy Architecture

### Approach: Schema-per-Tenant (Recommended)

For a SaaS school management system, we'll use **shared database with tenant identifier** approach:

```go
// Every tenant-specific table includes institution_id
type Student struct {
    ID            uuid.UUID `gorm:"type:uuid;primary_key"`
    InstitutionID uuid.UUID `gorm:"type:uuid;not null;index"` // Tenant identifier
    UserID        uuid.UUID `gorm:"type:uuid;not null"`
    // ... other fields
}
```

### Tenant Resolution Middleware

```go
// middleware/tenant_middleware.go
func TenantMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        // Super Admin doesn't need institution context
        if userRole == "SUPER_ADMIN" {
            c.Next()
            return
        }
        
        institutionID := c.GetHeader("X-Institution-ID")
        if institutionID == "" {
            c.AbortWithStatusJSON(400, gin.H{"error": "X-Institution-ID header required"})
            return
        }
        
        c.Set("institution_id", institutionID)
        c.Next()
    }
}
```

---

## 3. Authentication & Authorization

### JWT Token Structure

```go
type Claims struct {
    UserID        string   `json:"user_id"`
    Email         string   `json:"email"`
    Role          string   `json:"role"`
    InstitutionID string   `json:"institution_id,omitempty"`
    Permissions   []string `json:"permissions"`
    jwt.RegisteredClaims
}

// Token expiry
const (
    AccessTokenExpiry  = 15 * time.Minute
    RefreshTokenExpiry = 7 * 24 * time.Hour
)
```

### Role-Based Access Control

```go
// Permission constants (aligned with permission-matrix.txt)
const (
    // Super Admin - Full access
    PermFullSystemAccess = "*"
    
    // User Management
    PermUserCreate = "USER_CREATE"
    PermUserUpdate = "USER_UPDATE"
    PermUserDelete = "USER_DELETE"
    PermUserView   = "USER_VIEW"
    
    // Academic
    PermClassManage    = "CLASS_MANAGE"
    PermSectionManage  = "SECTION_MANAGE"
    PermSubjectManage  = "SUBJECT_MANAGE"
    PermTimetableManage = "TIMETABLE_MANAGE"
    
    // Attendance
    PermAttendanceMark = "ATTENDANCE_MARK"
    PermAttendanceView = "ATTENDANCE_VIEW"
    
    // Assessment
    PermAssignmentCreate = "ASSIGNMENT_CREATE"
    PermAssignmentGrade  = "ASSIGNMENT_GRADE"
    PermExamCreate       = "EXAM_CREATE"
    PermResultEnter      = "RESULT_ENTER"
    
    // Finance
    PermFeeCollect   = "FEE_COLLECT"
    PermFeeViewAll   = "FEE_VIEW_ALL"
    PermExpenseManage = "EXPENSE_MANAGE"
    PermSalaryProcess = "SALARY_PROCESS"
    
    // ... (complete list from permission-matrix.txt)
)

// Role to permissions mapping
var RolePermissions = map[string][]string{
    "SUPER_ADMIN": {"*"},
    "ADMIN": {
        PermUserCreate, PermUserUpdate, PermUserDelete, PermUserView,
        PermClassManage, PermSectionManage, PermSubjectManage,
        PermTimetableManage, PermNoticePublish, PermReportGenerate,
        PermLeaveApprove, PermLibraryManage, PermEventManage,
        // ...
    },
    "TEACHER": {
        PermAttendanceMark, PermAttendanceView, PermAssignmentCreate,
        PermAssignmentGrade, PermExamCreate, PermResultEnter,
        PermStudentProgressView, PermParentCommunicate, PermLeaveApply,
        // ...
    },
    // ... (other roles)
}
```

### RBAC Middleware

```go
// middleware/rbac_middleware.go
func RequireRole(roles ...string) gin.HandlerFunc {
    return func(c *gin.Context) {
        userRole := c.GetString("user_role")
        for _, role := range roles {
            if userRole == role {
                c.Next()
                return
            }
        }
        c.AbortWithStatusJSON(403, gin.H{"error": "Insufficient permissions"})
    }
}

func RequirePermission(permissions ...string) gin.HandlerFunc {
    return func(c *gin.Context) {
        userPerms := c.GetStringSlice("user_permissions")
        
        // Super Admin has all permissions
        if contains(userPerms, "*") {
            c.Next()
            return
        }
        
        for _, required := range permissions {
            if !contains(userPerms, required) {
                c.AbortWithStatusJSON(403, gin.H{"error": "Permission denied"})
                return
            }
        }
        c.Next()
    }
}
```

---

## 4. API Response Standards

### Success Response

```go
type APIResponse struct {
    Success bool        `json:"success"`
    Message string      `json:"message,omitempty"`
    Data    interface{} `json:"data,omitempty"`
}

type PaginatedResponse struct {
    Success    bool        `json:"success"`
    Data       interface{} `json:"data"`
    Pagination Pagination  `json:"pagination"`
}

type Pagination struct {
    CurrentPage int   `json:"current_page"`
    PerPage     int   `json:"per_page"`
    TotalItems  int64 `json:"total_items"`
    TotalPages  int   `json:"total_pages"`
}
```

### Error Response

```go
type ErrorResponse struct {
    Success bool              `json:"success"`
    Error   string            `json:"error"`
    Code    string            `json:"code,omitempty"`
    Details map[string]string `json:"details,omitempty"` // Validation errors
}
```

### Standard HTTP Status Codes

| Code | Usage |
|------|-------|
| 200 | Successful GET, PUT, PATCH |
| 201 | Successful POST (resource created) |
| 204 | Successful DELETE (no content) |
| 400 | Bad Request (validation error) |
| 401 | Unauthorized (missing/invalid token) |
| 403 | Forbidden (insufficient permissions) |
| 404 | Resource not found |
| 409 | Conflict (duplicate resource) |
| 422 | Unprocessable Entity |
| 429 | Too Many Requests (rate limited) |
| 500 | Internal Server Error |

---

## 5. Database Models

### Base Model

```go
type BaseModel struct {
    ID        uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
    CreatedAt time.Time      `gorm:"autoCreateTime"`
    UpdatedAt time.Time      `gorm:"autoUpdateTime"`
    DeletedAt gorm.DeletedAt `gorm:"index"`
}

// For tenant-specific models
type TenantBaseModel struct {
    BaseModel
    InstitutionID uuid.UUID `gorm:"type:uuid;not null;index"`
}
```

### Core Models (Aligned with database-schema)

```go
// Institution (schools table)
type Institution struct {
    BaseModel
    Name            string `gorm:"size:255;not null"`
    Code            string `gorm:"size:50;unique;not null"`
    Address         string `gorm:"type:text"`
    Phone           string `gorm:"size:20"`
    Email           string `gorm:"size:255"`
    PrincipalName   string `gorm:"size:255"`
    EstablishedYear int
    LogoURL         string `gorm:"size:500"`
    AcademicYear    string `gorm:"size:20"`
    IsActive        bool   `gorm:"default:true"`
}

// User
type User struct {
    BaseModel
    Email        string `gorm:"size:255;uniqueIndex"`
    Phone        string `gorm:"size:20"`
    PasswordHash string `gorm:"size:255"`
    Role         string `gorm:"size:50;not null"` // SUPER_ADMIN, ADMIN, TEACHER, STUDENT, PARENT, ACCOUNTANT
    IsActive     bool   `gorm:"default:true"`
    Profile      UserProfile
}

// UserProfile
type UserProfile struct {
    BaseModel
    UserID          uuid.UUID `gorm:"type:uuid;not null"`
    InstitutionID   uuid.UUID `gorm:"type:uuid"`
    FirstName       string    `gorm:"size:100"`
    LastName        string    `gorm:"size:100"`
    DateOfBirth     *time.Time
    Gender          string    `gorm:"size:10"`
    Address         string    `gorm:"type:text"`
    ProfileImageURL string    `gorm:"size:500"`
    EmployeeID      string    `gorm:"size:50"`
    AdmissionNumber string    `gorm:"size:50"`
    Occupation      string    `gorm:"size:100"`
}
```

---

## 6. WebSocket Implementation

### Hub Architecture

```go
// websocket/hub.go
type Hub struct {
    clients    map[*Client]bool
    broadcast  chan *Message
    register   chan *Client
    unregister chan *Client
    rooms      map[string]map[*Client]bool // Channel -> Clients
    mu         sync.RWMutex
}

type Client struct {
    hub           *Hub
    conn          *websocket.Conn
    send          chan []byte
    userID        string
    institutionID string
    subscriptions []string
}

type Message struct {
    Type      string      `json:"type"`
    Channel   string      `json:"channel"`
    Data      interface{} `json:"data"`
    Timestamp time.Time   `json:"timestamp"`
}
```

### Supported Channels

| Channel | Purpose | Subscribe |
|---------|---------|-----------|
| `notifications` | User notifications | Auto on connect |
| `messages:{userId}` | Direct messages | Auto on connect |
| `announcements:{institutionId}` | Institution announcements | Auto on connect |
| `attendance:{classId}` | Live attendance | Manual subscribe |

---

## 7. Development Phases

### Phase 1: Foundation Setup (Week 1-2)

**Goal:** Set up project structure, database, and core infrastructure

#### Week 1 Tasks

- [ ] **Project Initialization**
  - Initialize Go module with dependencies
  - Set up project folder structure
  - Configure Viper for environment management
  - Create `.env.example` and configuration loading

- [ ] **Database Setup**
  - Set up PostgreSQL connection with GORM
  - Set up Redis connection for caching/sessions
  - Create migration files for core tables:
    - `users`, `user_profiles`
    - `institutions` (schools)
    - `teachers`, `students`, `parents`
    - `parent_student_relations`
  - Create database seeder for test data

- [ ] **Utilities**
  - Password hashing utilities (bcrypt)
  - JWT token generation/validation
  - Standard API response helpers
  - Custom error types
  - Pagination helpers
  - Validation setup (go-playground/validator)

#### Week 2 Tasks

- [ ] **Core Middleware**
  - CORS middleware
  - Request logging middleware (Zap)
  - Recovery/panic handler middleware
  - Rate limiting middleware

- [ ] **Authentication System**
  - Auth models & DTOs
  - Auth repository
  - Auth service (login, register, refresh token, logout)
  - Auth handler (API endpoints)
  - JWT middleware
  - Password reset functionality

- [ ] **Basic Router Setup**
  - Main router configuration
  - Auth routes (`/api/v1/auth/*`)
  - Health check endpoint

**API Endpoints Delivered:**
```
POST   /api/v1/auth/login
POST   /api/v1/auth/refresh-token
POST   /api/v1/auth/logout
POST   /api/v1/auth/forgot-password
POST   /api/v1/auth/reset-password
POST   /api/v1/auth/change-password
GET    /api/v1/health
```

### Phase 2: Multi-Tenancy & User Management (Week 3-4)

**Goal:** Implement multi-tenant architecture and complete user management

#### Week 3 Tasks

- [ ] **Multi-Tenancy Middleware**
  - Tenant resolution from `X-Institution-ID` header
  - Automatic tenant filtering in repositories
  - Tenant context in all queries

- [ ] **RBAC Implementation**
  - Role middleware
  - Permission middleware
  - Permission constants and role mappings

- [ ] **Institution Management (Super Admin)**
  - Institution model, repository, service, handler
  - CRUD operations for institutions
  - Institution statistics endpoint
  - Enable/disable institution

**API Endpoints Delivered:**
```
GET    /api/v1/institutions
POST   /api/v1/institutions
GET    /api/v1/institutions/:id
PUT    /api/v1/institutions/:id
DELETE /api/v1/institutions/:id
PATCH  /api/v1/institutions/:id/status
GET    /api/v1/institutions/:id/stats
GET    /api/v1/institutions/:id/admins
POST   /api/v1/institutions/:id/admins
```

#### Week 4 Tasks

- [ ] **User Management**
  - User models, DTOs
  - User repository with tenant filtering
  - User service (CRUD, status management)
  - User handler

- [ ] **Profile Management**
  - Profile model, repository, service
  - Get/update own profile
  - Avatar upload (file handling)

- [ ] **Role-Specific User Management**
  - Teacher CRUD
  - Student CRUD (with parent linking)
  - Parent CRUD (with child linking)
  - Accountant CRUD

**API Endpoints Delivered:**
```
# Users
GET    /api/v1/users
GET    /api/v1/users/:id
POST   /api/v1/users
PUT    /api/v1/users/:id
DELETE /api/v1/users/:id
PATCH  /api/v1/users/:id/status

# Profile
GET    /api/v1/profile
PUT    /api/v1/profile
PUT    /api/v1/profile/avatar
PUT    /api/v1/profile/password

# Teachers
GET    /api/v1/teachers
GET    /api/v1/teachers/:id
POST   /api/v1/teachers
PUT    /api/v1/teachers/:id
GET    /api/v1/teachers/:id/classes
GET    /api/v1/teachers/:id/subjects

# Students
GET    /api/v1/students
GET    /api/v1/students/:id
POST   /api/v1/students
PUT    /api/v1/students/:id
GET    /api/v1/students/:id/parents
POST   /api/v1/students/:id/parents
DELETE /api/v1/students/:id/parents/:parentId

# Parents
GET    /api/v1/parents
GET    /api/v1/parents/:id
POST   /api/v1/parents
PUT    /api/v1/parents/:id
GET    /api/v1/parents/:id/children

# Accountants
GET    /api/v1/accountants
GET    /api/v1/accountants/:id
POST   /api/v1/accountants
PUT    /api/v1/accountants/:id
```

---

### Phase 3: Academic Management (Week 5-6)

**Goal:** Implement academic structure management

#### Week 5 Tasks

- [ ] **Database Migrations**
  - `academic_years`, `departments`
  - `classes`, `sections`
  - `subjects`, `timetable`

- [ ] **Academic Year Management**
  - Model, repository, service, handler
  - Set current academic year

- [ ] **Department Management**
  - CRUD operations
  - Staff assignment to departments

- [ ] **Class Management**
  - CRUD operations
  - Students and teachers in class

#### Week 6 Tasks

- [ ] **Section Management**
  - CRUD operations
  - Students in section

- [ ] **Subject Management**
  - CRUD operations
  - Subject-teacher assignment

- [ ] **Timetable Management**
  - Timetable CRUD
  - Class timetable view
  - Teacher timetable view
  - Section timetable view
  - Conflict detection

**API Endpoints Delivered:**
```
# Academic Years
GET    /api/v1/academic-years
POST   /api/v1/academic-years
GET    /api/v1/academic-years/:id
PUT    /api/v1/academic-years/:id
PATCH  /api/v1/academic-years/:id/activate
GET    /api/v1/academic-years/current

# Classes
GET    /api/v1/classes
POST   /api/v1/classes
GET    /api/v1/classes/:id
PUT    /api/v1/classes/:id
DELETE /api/v1/classes/:id
GET    /api/v1/classes/:id/students
GET    /api/v1/classes/:id/teachers

# Sections
GET    /api/v1/classes/:classId/sections
POST   /api/v1/classes/:classId/sections
PUT    /api/v1/sections/:id
DELETE /api/v1/sections/:id
GET    /api/v1/sections/:id/students

# Departments
GET    /api/v1/departments
POST   /api/v1/departments
GET    /api/v1/departments/:id
PUT    /api/v1/departments/:id
DELETE /api/v1/departments/:id
GET    /api/v1/departments/:id/staff

# Subjects
GET    /api/v1/subjects
POST   /api/v1/subjects
GET    /api/v1/subjects/:id
PUT    /api/v1/subjects/:id
DELETE /api/v1/subjects/:id
GET    /api/v1/subjects/class/:classId

# Timetable
GET    /api/v1/timetable
POST   /api/v1/timetable
PUT    /api/v1/timetable/:id
DELETE /api/v1/timetable/:id
GET    /api/v1/timetable/class/:classId
GET    /api/v1/timetable/teacher/:teacherId
GET    /api/v1/timetable/section/:sectionId
```

---

### Phase 4: Attendance Management (Week 7)

**Goal:** Implement attendance marking and tracking

- [ ] **Database Migrations**
  - `attendance` table

- [ ] **Attendance Models & DTOs**
  - Attendance model
  - Request/response DTOs
  - Bulk attendance request

- [ ] **Attendance Repository**
  - Mark attendance (single/batch)
  - Get attendance by student
  - Get attendance by class/date
  - Attendance statistics

- [ ] **Attendance Service**
  - Business logic for attendance
  - Validation (duplicate marking prevention)
  - Attendance report generation

- [ ] **Attendance Handler**
  - API endpoints

**API Endpoints Delivered:**
```
POST   /api/v1/attendance/mark
GET    /api/v1/attendance
GET    /api/v1/attendance/student/:studentId
GET    /api/v1/attendance/class/:classId/date/:date
PUT    /api/v1/attendance/:id
GET    /api/v1/attendance/report
POST   /api/v1/attendance/bulk-mark
GET    /api/v1/attendance/export
```

---

### Phase 5: Assessment Management (Week 8)

**Goal:** Implement assignments, exams, and results

- [ ] **Database Migrations**
  - `assignments`, `student_assignments`
  - `exams`, `exam_results`

- [ ] **Assignment Management**
  - Assignment CRUD
  - Assignment submission
  - Grading system

- [ ] **Exam Management**
  - Exam CRUD
  - Exam scheduling

- [ ] **Results Management**
  - Result entry (batch)
  - Grade calculation
  - Student results view

**API Endpoints Delivered:**
```
# Assignments
GET    /api/v1/assignments
POST   /api/v1/assignments
GET    /api/v1/assignments/:id
POST   /api/v1/assignments/:id/submit
GET    /api/v1/assignments/student/:studentId

# Exams
GET    /api/v1/exams
POST   /api/v1/exams
GET    /api/v1/exams/:id
POST   /api/v1/exams/:id/results
GET    /api/v1/exams/:id/results
GET    /api/v1/results/student/:studentId
GET    /api/v1/results/export
```

---

### Phase 6: Financial Management (Week 9-10)

**Goal:** Implement fee management, expenses, and salary processing

#### Week 9 Tasks

- [ ] **Database Migrations**
  - `fee_structures`, `fee_payments`
  - `expenses`, `salaries`
  - `scholarships`, `discounts`

- [ ] **Fee Structure Management**
  - Fee structure CRUD
  - Class-based fee assignment

- [ ] **Fee Collection**
  - Payment processing
  - Receipt generation
  - Payment history

#### Week 10 Tasks

- [ ] **Expense Management**
  - Expense CRUD
  - Category-based tracking

- [ ] **Salary Management**
  - Salary processing
  - Payroll generation

- [ ] **Financial Reports**
  - Fee collection reports
  - Expense reports
  - Financial summary

**API Endpoints Delivered:**
```
# Fee Structures
GET    /api/v1/fee/structures
POST   /api/v1/fee/structures
GET    /api/v1/fee/students/:studentId
POST   /api/v1/fee/payments
GET    /api/v1/fee/payments/:id
GET    /api/v1/fee/reports

# Expenses
GET    /api/v1/expenses
POST   /api/v1/expenses
PUT    /api/v1/expenses/:id
GET    /api/v1/expenses/reports

# Salaries
GET    /api/v1/salaries
POST   /api/v1/salaries
GET    /api/v1/salaries/employee/:employeeId
```

---

### Phase 7: Communication & WebSocket (Week 11)

**Goal:** Implement messaging, notices, and real-time communication

- [ ] **Database Migrations**
  - `notices`, `messages`, `announcements`
  - `notifications`

- [ ] **WebSocket Hub Setup**
  - Hub implementation
  - Client management
  - Channel subscription

- [ ] **Notice Management**
  - Notice CRUD
  - Target audience filtering

- [ ] **Messaging System**
  - Direct messaging
  - Broadcast messages
  - Message read receipts

- [ ] **Announcement System**
  - Announcement CRUD
  - Real-time broadcast

- [ ] **Notification System**
  - Notification storage
  - Real-time delivery
  - Mark as read

**API Endpoints Delivered:**
```
# Notices
GET    /api/v1/notices
POST   /api/v1/notices
GET    /api/v1/notices/:id
DELETE /api/v1/notices/:id

# Messages
GET    /api/v1/messages
POST   /api/v1/messages
GET    /api/v1/messages/:id
POST   /api/v1/messages/broadcast

# Announcements
GET    /api/v1/announcements
POST   /api/v1/announcements

# WebSocket
WS     /ws?token=<jwt_token>
```

---

### Phase 8: Leave & Library Management (Week 12)

**Goal:** Implement leave application/approval and library system

- [ ] **Database Migrations**
  - `leave_types`, `leaves`, `leave_balances`
  - `book_categories`, `books`
  - `book_borrowings`, `library_fines`

- [ ] **Leave Management**
  - Leave type CRUD
  - Leave application
  - Leave approval/rejection
  - Leave balance tracking

- [ ] **Library Management**
  - Book category CRUD
  - Book CRUD
  - Book borrowing/returning
  - Fine management

**API Endpoints Delivered:**
```
# Leave
GET    /api/v1/leaves
POST   /api/v1/leaves
GET    /api/v1/leaves/:id
PUT    /api/v1/leaves/:id
DELETE /api/v1/leaves/:id
PATCH  /api/v1/leaves/:id/approve
PATCH  /api/v1/leaves/:id/reject
GET    /api/v1/leaves/pending
GET    /api/v1/leave-types
POST   /api/v1/leave-types
PUT    /api/v1/leave-types/:id
DELETE /api/v1/leave-types/:id
GET    /api/v1/leaves/balance/:userId

# Library
GET    /api/v1/library/books
POST   /api/v1/library/books
GET    /api/v1/library/books/:id
PUT    /api/v1/library/books/:id
DELETE /api/v1/library/books/:id
GET    /api/v1/library/books/available
GET    /api/v1/library/categories
POST   /api/v1/library/categories
PUT    /api/v1/library/categories/:id
DELETE /api/v1/library/categories/:id
POST   /api/v1/library/borrow
POST   /api/v1/library/return/:borrowId
GET    /api/v1/library/borrowed
GET    /api/v1/library/borrowed/:userId
GET    /api/v1/library/overdue
GET    /api/v1/library/fines
GET    /api/v1/library/fines/:userId
POST   /api/v1/library/fines/:id/pay
```

---

### Phase 9: Events & Calendar (Week 13)

**Goal:** Implement event management and calendar functionality

- [ ] **Database Migrations**
  - `events`, `holidays`

- [ ] **Event Management**
  - Event CRUD
  - Upcoming events

- [ ] **Holiday Management**
  - Holiday CRUD

- [ ] **Calendar Integration**
  - Combined calendar view
  - Month view API

**API Endpoints Delivered:**
```
# Events
GET    /api/v1/events
POST   /api/v1/events
GET    /api/v1/events/:id
PUT    /api/v1/events/:id
DELETE /api/v1/events/:id
GET    /api/v1/events/upcoming

# Calendar
GET    /api/v1/calendar
GET    /api/v1/calendar/month/:year/:month

# Holidays
GET    /api/v1/holidays
POST   /api/v1/holidays
PUT    /api/v1/holidays/:id
DELETE /api/v1/holidays/:id
```

---

### Phase 10: Reports & Bulk Operations (Week 14)

**Goal:** Implement comprehensive reporting and bulk operations

- [ ] **Report Generation**
  - Academic performance reports
  - Attendance summary reports
  - Fee collection reports
  - Teacher performance reports
  - Student progress reports
  - Financial summary reports

- [ ] **Bulk Operations**
  - Bulk user import (CSV/Excel)
  - Bulk student enrollment
  - Bulk attendance marking
  - Bulk export functionality
  - Bulk student promotion

**API Endpoints Delivered:**
```
# Reports
GET    /api/v1/reports/academic-performance
GET    /api/v1/reports/attendance-summary
GET    /api/v1/reports/fee-collection
GET    /api/v1/reports/teacher-performance
GET    /api/v1/reports/student-progress
GET    /api/v1/reports/financial-summary

# Bulk Operations
POST   /api/v1/users/bulk-import
POST   /api/v1/students/bulk-enroll
POST   /api/v1/attendance/bulk-mark
GET    /api/v1/users/export
GET    /api/v1/students/export
GET    /api/v1/attendance/export
GET    /api/v1/results/export
PATCH  /api/v1/users/bulk-status
PATCH  /api/v1/students/bulk-promote
```

---

### Phase 11: Polish & Production Readiness (Week 15-16)

**Goal:** Production hardening, testing, and documentation

#### Week 15 Tasks

- [ ] **API Documentation**
  - Complete Swagger/OpenAPI documentation
  - API versioning strategy
  - Postman collection

- [ ] **Testing**
  - Unit tests for services
  - Integration tests for handlers
  - Repository tests with test database
  - WebSocket tests

- [ ] **Security Hardening**
  - SQL injection prevention (parameterized queries)
  - XSS prevention
  - CSRF protection (if needed)
  - Input sanitization
  - Rate limiting fine-tuning
  - Security headers

#### Week 16 Tasks

- [ ] **Performance Optimization**
  - Database query optimization
  - Index analysis
  - Connection pooling
  - Redis caching implementation
  - Response compression (gzip)

- [ ] **DevOps**
  - Dockerfile optimization
  - Docker Compose for local development
  - CI/CD pipeline setup
  - Health check endpoints
  - Graceful shutdown implementation

- [ ] **Monitoring & Logging**
  - Structured logging
  - Request ID tracking
  - Error tracking setup
  - Metrics collection (Prometheus-ready)

---

## 8. Required Go Dependencies

```go
// go.mod
module campus-core

go 1.24

require (
    // Web Framework
    github.com/gin-gonic/gin v1.10.0
    
    // Database
    gorm.io/gorm v1.25.0
    gorm.io/driver/postgres v1.5.0
    github.com/go-redis/redis/v9 v9.0.0
    
    // Authentication
    github.com/golang-jwt/jwt/v5 v5.2.0
    golang.org/x/crypto v0.21.0  // bcrypt
    
    // WebSocket
    github.com/gorilla/websocket v1.5.0
    
    // Configuration
    github.com/spf13/viper v1.18.0
    
    // Validation
    github.com/go-playground/validator/v10 v10.18.0
    
    // UUID
    github.com/google/uuid v1.6.0
    
    // Logging
    go.uber.org/zap v1.27.0
    
    // API Documentation
    github.com/swaggo/swag v1.16.0
    github.com/swaggo/gin-swagger v1.6.0
    
    // Testing
    github.com/stretchr/testify v1.9.0
    github.com/vektra/mockery/v2 v2.42.0
    
    // Migrations
    github.com/golang-migrate/migrate/v4 v4.17.0
    
    // File handling (Excel/CSV)
    github.com/xuri/excelize/v2 v2.8.0
    
    // Rate limiting
    golang.org/x/time v0.5.0
)
```

---

## 9. Environment Variables

```env
# .env.example

# Server
SERVER_PORT=8080
SERVER_ENV=development  # development, staging, production
API_VERSION=v1

# Database
DB_HOST=localhost
DB_PORT=5432
DB_NAME=campus_core
DB_USER=postgres
DB_PASSWORD=password
DB_SSL_MODE=disable
DB_MAX_CONNECTIONS=100
DB_MAX_IDLE_CONNECTIONS=10

# Redis
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=
REDIS_DB=0

# JWT
JWT_SECRET_KEY=your-super-secret-key
JWT_ACCESS_TOKEN_EXPIRY=15m
JWT_REFRESH_TOKEN_EXPIRY=168h  # 7 days

# CORS
CORS_ALLOWED_ORIGINS=http://localhost:4200,http://localhost:3000
CORS_ALLOWED_METHODS=GET,POST,PUT,PATCH,DELETE,OPTIONS
CORS_ALLOWED_HEADERS=Content-Type,Authorization,X-Institution-ID

# Rate Limiting
RATE_LIMIT_REQUESTS=100
RATE_LIMIT_WINDOW=1m

# File Upload
MAX_UPLOAD_SIZE=10485760  # 10MB
UPLOAD_PATH=./uploads

# Logging
LOG_LEVEL=debug  # debug, info, warn, error
LOG_FORMAT=json  # json, console
```

---

## 10. Database Migration Strategy

### Migration Naming Convention

```
{version}_{description}.{up|down}.sql
000001_create_users_table.up.sql
000001_create_users_table.down.sql
```

### Migration Order

1. `000001_create_institutions_table`
2. `000002_create_users_table`
3. `000003_create_user_profiles_table`
4. `000004_create_teachers_table`
5. `000005_create_students_table`
6. `000006_create_parents_table`
7. `000007_create_parent_student_relations_table`
8. `000008_create_academic_years_table`
9. `000009_create_departments_table`
10. `000010_create_classes_table`
11. `000011_create_sections_table`
12. `000012_create_subjects_table`
13. `000013_create_timetable_table`
14. `000014_create_attendance_table`
15. `000015_create_assignments_table`
16. `000016_create_student_assignments_table`
17. `000017_create_exams_table`
18. `000018_create_exam_results_table`
19. `000019_create_fee_structures_table`
20. `000020_create_fee_payments_table`
21. `000021_create_expenses_table`
22. `000022_create_salaries_table`
23. `000023_create_notices_table`
24. `000024_create_messages_table`
25. `000025_create_announcements_table`
26. `000026_create_notifications_table`
27. `000027_create_leave_types_table`
28. `000028_create_leaves_table`
29. `000029_create_book_categories_table`
30. `000030_create_books_table`
31. `000031_create_book_borrowings_table`
32. `000032_create_library_fines_table`
33. `000033_create_events_table`
34. `000034_create_holidays_table`

---

## 11. Makefile Commands

```makefile
# Makefile

.PHONY: run build test migrate-up migrate-down swagger docker-up docker-down

# Run development server
run:
	go run cmd/server/main.go

# Build binary
build:
	go build -o bin/server cmd/server/main.go

# Run tests
test:
	go test -v ./...

# Run tests with coverage
test-coverage:
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out

# Run migrations up
migrate-up:
	migrate -path internal/database/migrations -database "postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable" up

# Run migrations down
migrate-down:
	migrate -path internal/database/migrations -database "postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable" down

# Create new migration
migrate-create:
	migrate create -ext sql -dir internal/database/migrations -seq $(name)

# Generate Swagger docs
swagger:
	swag init -g cmd/server/main.go -o docs

# Docker compose up
docker-up:
	docker-compose up -d

# Docker compose down
docker-down:
	docker-compose down

# Install dependencies
deps:
	go mod tidy
	go mod download

# Lint
lint:
	golangci-lint run

# Seed database
seed:
	go run scripts/seed.go
```

---

## 12. Summary Timeline

| Phase | Description | Duration | Weeks |
|-------|-------------|----------|-------|
| 1 | Foundation Setup | 2 weeks | 1-2 |
| 2 | Multi-Tenancy & User Management | 2 weeks | 3-4 |
| 3 | Academic Management | 2 weeks | 5-6 |
| 4 | Attendance Management | 1 week | 7 |
| 5 | Assessment Management | 1 week | 8 |
| 6 | Financial Management | 2 weeks | 9-10 |
| 7 | Communication & WebSocket | 1 week | 11 |
| 8 | Leave & Library Management | 1 week | 12 |
| 9 | Events & Calendar | 1 week | 13 |
| 10 | Reports & Bulk Operations | 1 week | 14 |
| 11 | Polish & Production Readiness | 2 weeks | 15-16 |

**Total Estimated Duration: 16 weeks**

---

## 13. Frontend-Backend Sync Points

To ensure smooth development, coordinate the following:

| Frontend Phase | Backend Phase | Sync Point |
|----------------|---------------|------------|
| Phase 1 (Foundation) | Phase 1 (Foundation) | API base URL, error formats |
| Phase 3 (Auth Module) | Phase 1-2 (Auth) | Auth endpoints ready |
| Phase 4 (Dashboard) | Phase 2 (User Mgmt) | Dashboard data APIs |
| Phase 5 (User Mgmt) | Phase 2 (User Mgmt) | User CRUD APIs |
| Phase 6 (Academic) | Phase 3 (Academic) | Academic APIs |
| Phase 7 (Attendance) | Phase 4 (Attendance) | Attendance APIs |
| Phase 8 (Assessment) | Phase 5 (Assessment) | Assessment APIs |
| Phase 9 (Finance) | Phase 6 (Finance) | Finance APIs |
| Phase 10 (Communication) | Phase 7 (WebSocket) | WebSocket + Messaging APIs |
| Phase 11 (Leave/Library) | Phase 8 (Leave/Library) | Leave & Library APIs |
| Phase 12 (Institution) | Phase 2 (Institution) | Institution APIs (early) |
| Phase 13 (Reports) | Phase 10 (Reports) | Report APIs |

---

## 14. Recommendations for Modifications

### Suggested File Modifications

1. **Update `permission-matrix.txt`**: Consider adding more granular permissions for specific features like:
   - `STUDENT_VIEW_ALL` vs `STUDENT_VIEW_OWN_CLASS` for teachers
   - `FEE_REFUND` permission for accountants
   - `AUDIT_LOG_VIEW` for admins

2. **Add missing database schema files**:
   - Create `backend/database-schema/leave-management.txt`
   - Create `backend/database-schema/library-management.txt`
   - Create `backend/database-schema/events-management.txt`

3. **Update API documentation**:
   - Add pagination parameters to all list endpoints
   - Add filter/search query parameters
   - Document error response codes for each endpoint

4. **Add new documentation files**:
   - `backend/api-error-codes.txt` - Standardized error codes
   - `backend/websocket-events.txt` - Complete WebSocket event documentation

---

*Last Updated: January 2026*
*Version: 1.0*