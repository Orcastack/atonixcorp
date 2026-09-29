package registry

// AccessControl defines permissions for datasets.
type AccessControl struct {
	datasetPermissions map[string][]string // datasetID -> list of users
}

func NewAccessControl() *AccessControl {
	return &AccessControl{
		datasetPermissions: make(map[string][]string),
	}
}

func (ac *AccessControl) Grant(datasetID, user string) {
	ac.datasetPermissions[datasetID] = append(ac.datasetPermissions[datasetID], user)
}

func (ac *AccessControl) Revoke(datasetID, user string) {
	users := ac.datasetPermissions[datasetID]
	filtered := []string{}
	for _, u := range users {
		if u != user {
			filtered = append(filtered, u)
		}
	}
	ac.datasetPermissions[datasetID] = filtered
}

func (ac *AccessControl) Allowed(datasetID, user string) bool {
	for _, u := range ac.datasetPermissions[datasetID] {
		if u == user {
			return true
		}
	}
	return false
}
