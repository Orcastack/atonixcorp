package block

type Volume struct {
	ID         string
	SizeGB     int
	Backend    string
	AttachedTo string
}
