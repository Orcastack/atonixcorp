package auth

// Role constants used across the platform.
const (
	RoleAdmin     = "admin"
	RoleUser      = "user"
	RoleService   = "service"
	RoleAnalytics = "analytics"
	RoleBot       = "bot"
)

// PermissionMap maps roles to allowed actions.
var PermissionMap = map[string][]string{
	RoleAdmin: {
		"read:any",
		"write:any",
		"manage:security",
		"manage:projects",
	},
	RoleUser: {
		"read:self",
		"write:self",
	},
	RoleService: {
		"read:any",
		"write:any",
	},
	RoleAnalytics: {
		"read:analytics",
		"write:analytics",
	},
	RoleBot: {
		"read:bot",
		"write:bot",
	},
}

// HasPermission checks if a role allows a specific action.
func HasPermission(role, action string) bool {
	allowed, ok := PermissionMap[role]
	if !ok {
		return false
	}
	for _, a := range allowed {
		if a == action {
			return true
		}
	}
	return false
}
