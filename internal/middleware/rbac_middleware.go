package middleware

import (
	"campus-core/internal/models"
	"campus-core/internal/utils"

	"github.com/gin-gonic/gin"
)

// RequireRole returns a middleware that checks if the user has one of the required roles
func RequireRole(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole := GetUserRole(c)
		if userRole == "" {
			utils.Error(c, 401, utils.ErrTokenMissing)
			c.Abort()
			return
		}

		// Super Admin has access to everything
		if userRole == models.RoleSuperAdmin {
			c.Next()
			return
		}

		// Check if user has one of the required roles
		for _, role := range roles {
			if userRole == role {
				c.Next()
				return
			}
		}

		utils.Error(c, 403, utils.ErrRoleNotAllowed)
		c.Abort()
	}
}

// RequirePermission returns a middleware that checks if the user has all required permissions
func RequirePermission(permissions ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userPerms := GetUserPermissions(c)

		// Super Admin has all permissions
		if contains(userPerms, "*") {
			c.Next()
			return
		}

		// Check all required permissions
		for _, required := range permissions {
			if !contains(userPerms, required) {
				utils.Error(c, 403, utils.ErrInsufficientPermissions)
				c.Abort()
				return
			}
		}

		c.Next()
	}
}

// RequireAnyPermission returns a middleware that checks if the user has at least one of the permissions
func RequireAnyPermission(permissions ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userPerms := GetUserPermissions(c)

		// Super Admin has all permissions
		if contains(userPerms, "*") {
			c.Next()
			return
		}

		// Check if user has at least one permission
		for _, required := range permissions {
			if contains(userPerms, required) {
				c.Next()
				return
			}
		}

		utils.Error(c, 403, utils.ErrInsufficientPermissions)
		c.Abort()
	}
}

// RequireSuperAdmin returns a middleware that only allows super admins
func RequireSuperAdmin() gin.HandlerFunc {
	return RequireRole(models.RoleSuperAdmin)
}

// RequireAdmin returns a middleware that allows admins and super admins
func RequireAdmin() gin.HandlerFunc {
	return RequireRole(models.RoleSuperAdmin, models.RoleAdmin)
}

// RequireTeacher returns a middleware that allows teachers, admins, and super admins
func RequireTeacher() gin.HandlerFunc {
	return RequireRole(models.RoleSuperAdmin, models.RoleAdmin, models.RoleTeacher)
}

// RequireStaff returns a middleware that allows all staff (not students/parents)
func RequireStaff() gin.HandlerFunc {
	return RequireRole(models.RoleSuperAdmin, models.RoleAdmin, models.RoleTeacher, models.RoleAccountant)
}

// contains checks if a slice contains a string
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

// RolePermissions maps roles to their permissions
var RolePermissions = map[string][]string{
	models.RoleSuperAdmin: {"*"},
	models.RoleAdmin: {
		"USER_CREATE", "USER_UPDATE", "USER_DELETE", "USER_VIEW",
		"STUDENT_MANAGE", "TEACHER_MANAGE", "CLASS_MANAGE",
		"SECTION_MANAGE", "SUBJECT_MANAGE", "DEPARTMENT_MANAGE",
		"ACADEMIC_YEAR_MANAGE", "TIMETABLE_MANAGE",
		"FEE_STRUCTURE_MANAGE",
		"NOTICE_PUBLISH", "ANNOUNCEMENT_CREATE",
		"REPORT_GENERATE",
		"LEAVE_APPROVE",
		"LIBRARY_MANAGE",
		"EVENT_MANAGE",
	},
	models.RoleTeacher: {
		"ATTENDANCE_MARK", "ATTENDANCE_VIEW",
		"ASSIGNMENT_CREATE", "ASSIGNMENT_GRADE",
		"EXAM_CREATE", "RESULT_ENTER",
		"STUDENT_PROGRESS_VIEW",
		"PARENT_COMMUNICATE", "MESSAGE_SEND",
		"LEAVE_APPLY",
		"RESOURCE_UPLOAD", "MATERIAL_UPLOAD",
		"TIMETABLE_VIEW",
		"ONLINE_CLASS_CREATE",
	},
	models.RoleStudent: {
		"PROFILE_VIEW_OWN", "PROFILE_UPDATE_OWN",
		"ASSIGNMENT_VIEW", "ASSIGNMENT_SUBMIT",
		"RESULT_VIEW_OWN",
		"ATTENDANCE_VIEW_OWN",
		"FEE_VIEW_OWN",
		"LEAVE_APPLY",
		"LIBRARY_BORROW", "LIBRARY_VIEW",
		"EVENT_VIEW",
		"MATERIAL_DOWNLOAD",
		"MESSAGE_SEND", "NOTICE_VIEW",
	},
	models.RoleParent: {
		"STUDENT_PROGRESS_VIEW",
		"FEE_PAY", "FEE_VIEW_CHILD",
		"TEACHER_COMMUNICATE", "MESSAGE_SEND",
		"ATTENDANCE_VIEW_CHILD",
		"LEAVE_APPLY_CHILD",
		"MEETING_SCHEDULE",
		"EVENT_VIEW", "NOTICE_VIEW",
	},
	models.RoleAccountant: {
		"FEE_COLLECT", "FEE_VIEW_ALL", "FEE_STRUCTURE_VIEW",
		"EXPENSE_MANAGE", "EXPENSE_CREATE", "EXPENSE_VIEW",
		"SALARY_PROCESS", "SALARY_VIEW",
		"FINANCIAL_REPORT_GENERATE",
		"INVOICE_GENERATE",
		"SCHOLARSHIP_MANAGE", "DISCOUNT_APPLY",
	},
}

// GetPermissionsForRole returns the permissions for a given role
func GetPermissionsForRole(role string) []string {
	if perms, ok := RolePermissions[role]; ok {
		return perms
	}
	return []string{}
}
