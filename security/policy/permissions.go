package policy

// Action represents a secured operation in the platform.
type Action string

// Common platform actions
const (
	ActionReadAny        Action = "read:any"
	ActionWriteAny       Action = "write:any"
	ActionReadSelf       Action = "read:self"
	ActionWriteSelf      Action = "write:self"
	ActionManageSecurity Action = "manage:security"
	ActionManageProjects Action = "manage:projects"
	ActionReadAnalytics  Action = "read:analytics"
	ActionWriteAnalytics Action = "write:analytics"
	ActionReadBot        Action = "read:bot"
	ActionWriteBot       Action = "write:bot"
)

// RolePermissions maps roles to allowed actions.
var RolePermissions = map[string][]Action{
	"admin": {
		ActionReadAny,
		ActionWriteAny,
		ActionManageSecurity,
		ActionManageProjects,
	},
	"user": {
		ActionReadSelf,
		ActionWriteSelf,
	},
	"service": {
		ActionReadAny,
		ActionWriteAny,
	},
	"analytics": {
		ActionReadAnalytics,
		ActionWriteAnalytics,
	},
	"bot": {
		ActionReadBot,
		ActionWriteBot,
	},
}

// HasPermission checks if a role allows a specific action.
func HasPermission(role string, action Action) bool {
	allowed, ok := RolePermissions[role]
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
