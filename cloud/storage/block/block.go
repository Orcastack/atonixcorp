package block

import "fmt"

type LocalDriver struct{}

func NewLocalDriver() *LocalDriver {
	return &LocalDriver{}
}

func (d *LocalDriver) Create(sizeGB int) (Volume, error) {
	v := Volume{
		ID:      fmt.Sprintf("vol-%d", sizeGB),
		SizeGB:  sizeGB,
		Backend: "local",
	}
	fmt.Println("Block: create volume", v.ID)
	return v, nil
}

func (d *LocalDriver) Delete(id string) error {
	fmt.Println("Block: delete volume", id)
	return nil
}

func (d *LocalDriver) Attach(volumeID, vmID string) error {
	fmt.Println("Block: attach", volumeID, "to", vmID)
	return nil
}

func (d *LocalDriver) Detach(volumeID string) error {
	fmt.Println("Block: detach", volumeID)
	return nil
}
