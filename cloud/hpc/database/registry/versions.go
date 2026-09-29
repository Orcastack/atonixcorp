package registry

import "fmt"

// Version represents a single dataset version.
type Version struct {
	Number    int
	Timestamp string
	Notes     string
}

type VersionManager struct {
	versions map[string][]Version
}

func NewVersionManager() *VersionManager {
	return &VersionManager{
		versions: make(map[string][]Version),
	}
}

func (vm *VersionManager) AddVersion(datasetID string, v Version) {
	vm.versions[datasetID] = append(vm.versions[datasetID], v)
	fmt.Printf("Version added: %s v%d\n", datasetID, v.Number)
}

func (vm *VersionManager) GetVersions(datasetID string) []Version {
	return vm.versions[datasetID]
}
