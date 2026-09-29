package orchestration

type Orchestrator interface {
	RunWorkflow(name string, payload map[string]interface{}) error
	Reconcile() error
	EmitEvent(event Event)
}
