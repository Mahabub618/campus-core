# Smart School Management System (SSMS) - Implementation Plan

## Project Overview

A comprehensive multi-tenant SaaS school management system using Angular 21 with PrimeNG, featuring 6 role-based dashboards, real-time WebSocket communication, i18n support, and lazy-loaded feature modules with Angular Signals for state management.

### Tech Stack
- **Frontend**: Angular 21
- **UI Framework**: PrimeNG
- **Backend**: Golang (separate project)
- **Real-time**: WebSocket
- **State Management**: Angular Signals
- **i18n**: ngx-translate

### Key Features
- ✅ Multi-Institution Support (SaaS model)
- ✅ Role-based Access Control (6 roles)
- ✅ Real-time Notifications via WebSocket
- ✅ Mobile-responsive UI
- ✅ Multi-language Support
- ❌ Offline Support (not required)
- ❌ Print Requirements (not required)
- ❌ External Integrations (not required)

---

## 1. User Roles & Permissions

### Roles
1. **Super Admin** - Full system access, multi-institution management
2. **Admin** - School-level management
3. **Teacher** - Class, attendance, assessment management
4. **Student** - View own data, submit assignments
5. **Parent** - View child data, pay fees, communicate
6. **Accountant** - Financial operations

### Enhanced Permission Matrix

```javascript
const PERMISSIONS = {
  SUPER_ADMIN: ['*'],
  ADMIN: [
    'USER_CREATE', 'USER_UPDATE', 'USER_DELETE',
    'STUDENT_MANAGE', 'TEACHER_MANAGE', 'CLASS_MANAGE',
    'SECTION_MANAGE', 'SUBJECT_MANAGE', 'TIMETABLE_MANAGE',
    'FEE_STRUCTURE_MANAGE', 'NOTICE_PUBLISH', 'REPORT_GENERATE',
    'LEAVE_APPROVE', 'ACADEMIC_YEAR_MANAGE', 'LIBRARY_MANAGE',
    'EVENT_MANAGE'
  ],
  TEACHER: [
    'ATTENDANCE_MARK', 'ASSIGNMENT_CREATE', 'RESULT_ENTER',
    'STUDENT_PROGRESS_VIEW', 'PARENT_COMMUNICATE',
    'LEAVE_APPLY', 'EXAM_CREATE', 'RESOURCE_UPLOAD',
    'ONLINE_CLASS_CREATE', 'TIMETABLE_VIEW'
  ],
  STUDENT: [
    'PROFILE_VIEW_OWN', 'ASSIGNMENT_SUBMIT', 'RESULT_VIEW_OWN',
    'ATTENDANCE_VIEW_OWN', 'FEE_VIEW_OWN',
    'LEAVE_APPLY', 'LIBRARY_BORROW', 'EVENT_VIEW',
    'MATERIAL_DOWNLOAD', 'MESSAGE_SEND'
  ],
  PARENT: [
    'STUDENT_PROGRESS_VIEW', 'FEE_PAY', 'TEACHER_COMMUNICATE',
    'ATTENDANCE_VIEW_CHILD', 'LEAVE_APPLY_CHILD',
    'MEETING_SCHEDULE', 'EVENT_VIEW', 'NOTICE_VIEW'
  ],
  ACCOUNTANT: [
    'FEE_COLLECT', 'EXPENSE_MANAGE', 'SALARY_PROCESS',
    'FINANCIAL_REPORT_GENERATE', 'SCHOLARSHIP_MANAGE',
    'DISCOUNT_APPLY', 'INVOICE_GENERATE'
  ]
};
```

---

## 2. Complete Folder Structure

```
src/
├── app/
│   ├── core/                          # Singleton services, guards, interceptors
│   │   ├── auth/
│   │   │   ├── guards/
│   │   │   │   ├── auth.guard.ts
│   │   │   │   ├── role.guard.ts
│   │   │   │   └── permission.guard.ts
│   │   │   ├── interceptors/
│   │   │   │   ├── auth.interceptor.ts
│   │   │   │   ├── error.interceptor.ts
│   │   │   │   └── loading.interceptor.ts
│   │   │   ├── services/
│   │   │   │   └── auth.service.ts
│   │   │   └── models/
│   │   │       ├── user.model.ts
│   │   │       ├── auth.model.ts
│   │   │       └── permission.model.ts
│   │   ├── services/
│   │   │   ├── api.service.ts           # Base HTTP service
│   │   │   ├── websocket.service.ts     # WebSocket connection manager
│   │   │   ├── notification.service.ts  # Real-time notifications
│   │   │   ├── storage.service.ts       # localStorage/sessionStorage abstraction
│   │   │   └── theme.service.ts         # Theme management
│   │   ├── state/
│   │   │   ├── auth.state.ts            # Authentication signals
│   │   │   ├── app.state.ts             # Global app state signals
│   │   │   └── notification.state.ts    # Notification signals
│   │   └── constants/
│   │       ├── api.constants.ts
│   │       ├── permission.constants.ts
│   │       └── route.constants.ts
│   │
│   ├── shared/                         # Reusable components, directives, pipes
│   │   ├── components/
│   │   │   ├── layouts/
│   │   │   │   ├── main-layout/        # Sidebar + header + content
│   │   │   │   ├── auth-layout/        # Login/register layout
│   │   │   │   ├── sidebar/
│   │   │   │   ├── header/
│   │   │   │   └── footer/
│   │   │   ├── ui/
│   │   │   │   ├── data-table/         # Generic table with pagination
│   │   │   │   ├── confirm-dialog/
│   │   │   │   ├── loading-spinner/
│   │   │   │   ├── empty-state/
│   │   │   │   ├── stat-card/          # Dashboard stat widgets
│   │   │   │   ├── calendar-widget/
│   │   │   │   └── notification-bell/
│   │   │   └── forms/
│   │   │       ├── form-field/
│   │   │       ├── search-input/
│   │   │       └── date-range-picker/
│   │   ├── directives/
│   │   │   ├── has-role.directive.ts
│   │   │   ├── has-permission.directive.ts
│   │   │   └── click-outside.directive.ts
│   │   ├── pipes/
│   │   │   ├── date-format.pipe.ts
│   │   │   ├── currency-format.pipe.ts
│   │   │   ├── truncate.pipe.ts
│   │   │   └── safe-html.pipe.ts
│   │   └── models/
│   │       ├── api-response.model.ts
│   │       ├── pagination.model.ts
│   │       └── common.model.ts
│   │
│   ├── features/                       # Lazy-loaded feature modules
│   │   ├── auth/
│   │   │   ├── pages/
│   │   │   │   ├── login/
│   │   │   │   ├── forgot-password/
│   │   │   │   └── reset-password/
│   │   │   └── auth.routes.ts
│   │   │
│   │   ├── dashboard/                  # Role-based dashboards
│   │   │   ├── super-admin-dashboard/
│   │   │   ├── admin-dashboard/
│   │   │   ├── teacher-dashboard/
│   │   │   ├── student-dashboard/
│   │   │   ├── parent-dashboard/
│   │   │   ├── accountant-dashboard/
│   │   │   └── dashboard.routes.ts
│   │   │
│   │   ├── institution/                # Super Admin - Multi-tenant
│   │   │   ├── pages/
│   │   │   │   ├── institution-list/
│   │   │   │   ├── institution-form/
│   │   │   │   └── institution-details/
│   │   │   ├── services/
│   │   │   │   └── institution.service.ts
│   │   │   └── institution.routes.ts
│   │   │
│   │   ├── users/                      # User Management
│   │   │   ├── pages/
│   │   │   │   ├── user-list/
│   │   │   │   ├── user-form/
│   │   │   │   ├── teacher-list/
│   │   │   │   ├── student-list/
│   │   │   │   ├── parent-list/
│   │   │   │   └── profile/
│   │   │   ├── services/
│   │   │   │   └── user.service.ts
│   │   │   └── users.routes.ts
│   │   │
│   │   ├── academic/                   # Academic Management
│   │   │   ├── pages/
│   │   │   │   ├── classes/
│   │   │   │   ├── sections/
│   │   │   │   ├── subjects/
│   │   │   │   ├── timetable/
│   │   │   │   └── academic-year/
│   │   │   ├── services/
│   │   │   │   └── academic.service.ts
│   │   │   └── academic.routes.ts
│   │   │
│   │   ├── attendance/
│   │   │   ├── pages/
│   │   │   │   ├── mark-attendance/
│   │   │   │   ├── view-attendance/
│   │   │   │   └── attendance-report/
│   │   │   ├── services/
│   │   │   │   └── attendance.service.ts
│   │   │   └── attendance.routes.ts
│   │   │
│   │   ├── assessment/
│   │   │   ├── pages/
│   │   │   │   ├── assignments/
│   │   │   │   │   ├── assignment-list/
│   │   │   │   │   ├── assignment-form/
│   │   │   │   │   └── assignment-submit/
│   │   │   │   ├── exams/
│   │   │   │   │   ├── exam-list/
│   │   │   │   │   ├── exam-form/
│   │   │   │   │   └── result-entry/
│   │   │   │   └── results/
│   │   │   │       └── result-view/
│   │   │   ├── services/
│   │   │   │   └── assessment.service.ts
│   │   │   └── assessment.routes.ts
│   │   │
│   │   ├── finance/
│   │   │   ├── pages/
│   │   │   │   ├── fee-structure/
│   │   │   │   ├── fee-collection/
│   │   │   │   ├── student-fees/
│   │   │   │   ├── expenses/
│   │   │   │   ├── salary/
│   │   │   │   └── financial-reports/
│   │   │   ├── services/
│   │   │   │   └── finance.service.ts
│   │   │   └── finance.routes.ts
│   │   │
│   │   ├── communication/
│   │   │   ├── pages/
│   │   │   │   ├── notices/
│   │   │   │   ├── messages/
│   │   │   │   └── announcements/
│   │   │   ├── services/
│   │   │   │   └── communication.service.ts
│   │   │   └── communication.routes.ts
│   │   │
│   │   ├── reports/
│   │   │   ├── pages/
│   │   │   │   ├── academic-report/
│   │   │   │   ├── attendance-report/
│   │   │   │   └── financial-report/
│   │   │   ├── services/
│   │   │   │   └── report.service.ts
│   │   │   └── reports.routes.ts
│   │   │
│   │   ├── leave/
│   │   │   ├── pages/
│   │   │   │   ├── apply-leave/
│   │   │   │   └── approve-leave/
│   │   │   ├── services/
│   │   │   │   └── leave.service.ts
│   │   │   └── leave.routes.ts
│   │   │
│   │   ├── library/
│   │   │   ├── pages/
│   │   │   │   ├── book-list/
│   │   │   │   ├── book-form/
│   │   │   │   └── borrowing/
│   │   │   ├── services/
│   │   │   │   └── library.service.ts
│   │   │   └── library.routes.ts
│   │   │
│   │   └── events/
│   │       ├── pages/
│   │       │   ├── event-list/
│   │       │   ├── event-form/
│   │       │   └── calendar/
│   │       ├── services/
│   │       │   └── events.service.ts
│   │       └── events.routes.ts
│   │
│   ├── app.ts
│   ├── app.html
│   ├── app.scss
│   ├── app.config.ts
│   └── app.routes.ts
│
├── assets/
│   ├── i18n/
│   │   ├── en.json
│   │   ├── hi.json                     # Hindi
│   │   └── es.json                     # Spanish (example)
│   ├── images/
│   └── icons/
│
├── environments/
│   ├── environment.ts
│   └── environment.prod.ts
│
├── styles/
│   ├── _variables.scss
│   ├── _mixins.scss
│   ├── _primeng-overrides.scss
│   └── _responsive.scss
│
├── index.html
├── main.ts
└── styles.scss
```

---

## 3. Core Services & Responsibilities

| Service | File | Responsibility |
|---------|------|----------------|
| **AuthService** | `core/auth/services/auth.service.ts` | Login, logout, refresh token, change/reset password, store/retrieve tokens |
| **ApiService** | `core/services/api.service.ts` | Base HTTP wrapper with `get<T>()`, `post<T>()`, `put<T>()`, `delete<T>()`, error handling |
| **WebSocketService** | `core/services/websocket.service.ts` | Connect/disconnect WebSocket, subscribe to channels, handle reconnection |
| **NotificationService** | `core/services/notification.service.ts` | Show toast (PrimeNG Toast), real-time notification handling |
| **StorageService** | `core/services/storage.service.ts` | Abstract `localStorage`/`sessionStorage`, handle token encryption |
| **ThemeService** | `core/services/theme.service.ts` | Switch PrimeNG themes, persist user preference |

---

## 4. State Management with Angular Signals

### Auth State
```typescript
// core/state/auth.state.ts
export const authState = {
  currentUser: signal<User | null>(null),
  isAuthenticated: computed(() => authState.currentUser() !== null),
  userRole: computed(() => authState.currentUser()?.role ?? null),
  permissions: computed(() => authState.currentUser()?.permissions ?? []),
  token: signal<string | null>(null),
};
```

### App State
```typescript
// core/state/app.state.ts
export const appState = {
  isLoading: signal(false),
  currentInstitution: signal<Institution | null>(null),
  sidebarCollapsed: signal(false),
  currentLanguage: signal('en'),
  theme: signal('lara-light-indigo'),
};
```

### Notification State
```typescript
// core/state/notification.state.ts
export const notificationState = {
  notifications: signal<Notification[]>([]),
  unreadCount: computed(() => 
    notificationState.notifications().filter(n => !n.read).length
  ),
};
```

---

## 5. Routing Structure with Lazy Loading

```typescript
// app.routes.ts
export const routes: Routes = [
  { path: '', redirectTo: 'dashboard', pathMatch: 'full' },
  
  // Auth routes (public)
  {
    path: 'auth',
    loadChildren: () => import('./features/auth/auth.routes'),
    component: AuthLayoutComponent
  },
  
  // Protected routes
  {
    path: '',
    component: MainLayoutComponent,
    canActivate: [authGuard],
    children: [
      {
        path: 'dashboard',
        loadChildren: () => import('./features/dashboard/dashboard.routes'),
      },
      {
        path: 'institutions',
        loadChildren: () => import('./features/institution/institution.routes'),
        canActivate: [roleGuard],
        data: { roles: ['SUPER_ADMIN'] }
      },
      {
        path: 'users',
        loadChildren: () => import('./features/users/users.routes'),
        canActivate: [roleGuard],
        data: { roles: ['SUPER_ADMIN', 'ADMIN'] }
      },
      {
        path: 'academic',
        loadChildren: () => import('./features/academic/academic.routes'),
        canActivate: [roleGuard],
        data: { roles: ['ADMIN', 'TEACHER'] }
      },
      {
        path: 'attendance',
        loadChildren: () => import('./features/attendance/attendance.routes'),
      },
      {
        path: 'assessment',
        loadChildren: () => import('./features/assessment/assessment.routes'),
      },
      {
        path: 'finance',
        loadChildren: () => import('./features/finance/finance.routes'),
        canActivate: [roleGuard],
        data: { roles: ['ADMIN', 'ACCOUNTANT', 'PARENT', 'STUDENT'] }
      },
      {
        path: 'communication',
        loadChildren: () => import('./features/communication/communication.routes'),
      },
      {
        path: 'reports',
        loadChildren: () => import('./features/reports/reports.routes'),
        canActivate: [roleGuard],
        data: { roles: ['SUPER_ADMIN', 'ADMIN', 'ACCOUNTANT'] }
      },
      {
        path: 'leave',
        loadChildren: () => import('./features/leave/leave.routes'),
      },
      {
        path: 'library',
        loadChildren: () => import('./features/library/library.routes'),
      },
      {
        path: 'events',
        loadChildren: () => import('./features/events/events.routes'),
      },
      {
        path: 'settings',
        loadChildren: () => import('./features/settings/settings.routes'),
      },
      {
        path: 'profile',
        loadComponent: () => import('./features/users/pages/profile/profile.component'),
      }
    ]
  },
  
  { path: '**', redirectTo: 'dashboard' }
];
```

---

## 6. Shared Components

| Component | Purpose | PrimeNG Components Used |
|-----------|---------|------------------------|
| `MainLayoutComponent` | Page wrapper with sidebar/header | `p-sidebar`, `p-menubar` |
| `AuthLayoutComponent` | Login/register page layout | Basic layout |
| `SidebarComponent` | Dynamic navigation menu by role | `p-panelMenu`, `p-menu` |
| `HeaderComponent` | Top bar with user, notifications, lang | `p-toolbar`, `p-avatar`, `p-dropdown` |
| `DataTableComponent` | Reusable table with sort/filter/page | `p-table`, `p-paginator` |
| `StatCardComponent` | Dashboard stat widget | `p-card` |
| `ConfirmDialogComponent` | Confirm actions | `p-confirmDialog` |
| `LoadingSpinnerComponent` | Global loading overlay | `p-progressSpinner` |
| `NotificationBellComponent` | Header notification dropdown | `p-badge`, `p-overlayPanel` |
| `FormFieldComponent` | Form input wrapper with errors | `p-inputText`, `p-message` |
| `EmptyStateComponent` | No data placeholder | Custom |
| `CalendarWidgetComponent` | Dashboard calendar | `p-calendar` |
| `SearchInputComponent` | Search with debounce | `p-inputText`, `p-iconField` |
| `DateRangePickerComponent` | Date range selection | `p-calendar` |

---

## 7. Feature Modules Breakdown

| Module | Pages | Role Access | Key API Endpoints |
|--------|-------|-------------|-------------------|
| **Auth** | Login, Forgot/Reset Password | Public | `/auth/*` |
| **Dashboard** | 6 role-specific dashboards | All roles | Various aggregated data |
| **Institution** | List, Create, Edit, Details | Super Admin | `/institutions/*` |
| **Users** | User List, Forms, Profile | Admin, Super Admin | `/users/*`, `/teachers/*`, `/students/*` |
| **Academic** | Classes, Sections, Subjects, Timetable | Admin, Teacher | `/classes/*`, `/subjects/*`, `/timetable/*` |
| **Attendance** | Mark, View, Report | Teacher, Student, Parent | `/attendance/*` |
| **Assessment** | Assignments, Exams, Results | Teacher, Student | `/assignments/*`, `/exams/*`, `/results/*` |
| **Finance** | Fees, Expenses, Salary, Reports | Accountant, Admin, Parent | `/fee/*`, `/expenses/*`, `/salaries/*` |
| **Communication** | Notices, Messages, Announcements | All roles | `/notices/*`, `/messages/*` |
| **Reports** | Academic, Attendance, Financial | Admin, Accountant | `/reports/*` |
| **Leave** | Apply, Approve | Teacher, Student, Admin | `/leaves/*` |
| **Library** | Books, Borrowing | Admin, Student, Teacher | `/library/*` |
| **Events** | Events, Calendar | All roles | `/events/*` |

---

## 8. WebSocket Implementation

### Service Structure
```typescript
@Injectable({ providedIn: 'root' })
export class WebSocketService {
  private socket: WebSocket | null = null;
  private reconnectAttempts = 0;
  private maxReconnectAttempts = 5;
  
  readonly connectionStatus = signal<'connected' | 'disconnected' | 'connecting'>('disconnected');
  readonly messages$ = new Subject<WebSocketMessage>();

  connect(token: string): void { /* ... */ }
  disconnect(): void { /* ... */ }
  subscribe(channel: string): void { /* ... */ }
  send(message: WebSocketMessage): void { /* ... */ }
  private handleReconnect(): void { /* ... */ }
}
```

### WebSocket Channels
- `notifications` - Real-time notifications
- `messages:{userId}` - Direct messages
- `announcements:{institutionId}` - Institution-wide announcements
- `attendance:{classId}` - Live attendance updates

### Integration Points
- Connect after successful login
- Disconnect on logout
- Auto-reconnect with exponential backoff
- Update `notificationState` on incoming messages

---

## 9. i18n Setup

### Installation
```bash
npm install @ngx-translate/core @ngx-translate/http-loader
```

### Configuration
```typescript
// app.config.ts
export function HttpLoaderFactory(http: HttpClient) {
  return new TranslateHttpLoader(http, './assets/i18n/', '.json');
}

provideTranslateModule({
  loader: { provide: TranslateLoader, useFactory: HttpLoaderFactory, deps: [HttpClient] }
})
```

### Translation File Structure
```json
// assets/i18n/en.json
{
  "common": { "save": "Save", "cancel": "Cancel", "delete": "Delete" },
  "auth": { "login": "Login", "logout": "Logout" },
  "dashboard": { "welcome": "Welcome, {{name}}" },
  "attendance": { "mark": "Mark Attendance", "present": "Present" }
}
```

### Usage
- Template: `{{ 'common.save' | translate }}`
- Binding: `[label]="'common.save' | translate"`
- Language switcher in HeaderComponent with persisted preference

---

## 10. Required Dependencies

```json
{
  "dependencies": {
    "primeng": "^18.0.0",
    "primeicons": "^7.0.0",
    "@primeng/themes": "^18.0.0",
    "@ngx-translate/core": "^16.0.0",
    "@ngx-translate/http-loader": "^9.0.0"
  }
}
```

---

## 11. Development Phases

### Phase 1: Foundation (Week 1-2)
- [ ] Install PrimeNG and configure themes
- [ ] Set up environment files with API base URL
- [ ] Create Core module structure and constants
- [ ] Implement `AuthService`, `ApiService`, `StorageService`
- [ ] Create auth interceptors (JWT, Error, Loading)
- [ ] Implement auth guards (auth, role, permission)
- [ ] Set up Angular Signals state (auth, app, notification)
- [ ] Configure i18n with ngx-translate
- [ ] Create base translation files (en.json)

### Phase 2: Layouts & Shared Components (Week 2-3)
- [ ] Build `AuthLayoutComponent`
- [ ] Build `MainLayoutComponent` with sidebar/header
- [ ] Create `SidebarComponent` with role-based menu
- [ ] Create `HeaderComponent` with user menu, language switch
- [ ] Build shared UI components
- [ ] Implement `hasRole` and `hasPermission` directives
- [ ] Create shared pipes
- [ ] Set up responsive SCSS mixins

### Phase 3: Auth Module (Week 3)
- [ ] Login page with form validation
- [ ] Forgot password page
- [ ] Reset password page
- [ ] Token refresh logic
- [ ] Redirect to role-based dashboard after login

### Phase 4: Dashboard Module (Week 4)
- [ ] Super Admin Dashboard
- [ ] Admin Dashboard
- [ ] Teacher Dashboard
- [ ] Student Dashboard
- [ ] Parent Dashboard
- [ ] Accountant Dashboard

### Phase 5: User Management (Week 5)
- [ ] User list with filters and pagination
- [ ] User create/edit form
- [ ] Teacher, Student, Parent lists
- [ ] Profile view/edit

### Phase 6: Academic Management (Week 5-6)
- [ ] Class management (CRUD)
- [ ] Section management
- [ ] Subject management
- [ ] Academic year setup
- [ ] Timetable builder UI

### Phase 7: Attendance Module (Week 6)
- [ ] Mark attendance (batch selection)
- [ ] View attendance
- [ ] Attendance reports

### Phase 8: Assessment Module (Week 7)
- [ ] Assignment management
- [ ] Assignment submission (Student)
- [ ] Exam management
- [ ] Result entry and viewing

### Phase 9: Finance Module (Week 8)
- [ ] Fee structure management
- [ ] Fee collection UI
- [ ] Student fee view/payment
- [ ] Expense management
- [ ] Salary processing
- [ ] Financial reports

### Phase 10: Communication Module (Week 9)
- [ ] Notice board
- [ ] Messaging system
- [ ] Announcements
- [ ] WebSocket integration

### Phase 11: Additional Modules (Week 10)
- [ ] Leave management
- [ ] Library management
- [ ] Events and calendar

### Phase 12: Institution Management (Week 10-11)
- [ ] Institution list (Super Admin)
- [ ] Create/edit institution
- [ ] Institution details
- [ ] Admin assignment

### Phase 13: Reports Module (Week 11)
- [ ] Academic performance reports
- [ ] Attendance summary reports
- [ ] Financial reports

### Phase 14: Polish & Testing (Week 12)
- [ ] Responsive design testing
- [ ] Error handling improvements
- [ ] Loading states and skeleton screens
- [ ] Accessibility improvements
- [ ] Performance optimization
- [ ] Cross-browser testing

---

## 12. TypeScript Path Aliases

```json
// tsconfig.json
{
  "compilerOptions": {
    "paths": {
      "@core/*": ["src/app/core/*"],
      "@shared/*": ["src/app/shared/*"],
      "@features/*": ["src/app/features/*"],
      "@env/*": ["src/environments/*"]
    }
  }
}
```

---

## 13. Recommended PrimeNG Theme

**Lara Light Indigo** (default) with option to switch to **Lara Dark Indigo** for dark mode toggle.

---

*Last Updated: January 2026*

