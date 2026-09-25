package heat

type Stack struct {
	ID          string `json:"id"`
	StackName   string `json:"stack_name"`
	StackStatus string `json:"stack_status"`
}

type CreateStackRequest struct {
	Name       string
	Template   string            // HOT/YAML as string
	Parameters map[string]string // simple key/value params
}
