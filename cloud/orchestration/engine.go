package orchestration

import (
	"orchestration/events"
	"orchestration/reconciler"
	"orchestration/workflows/compute"
)

type Engine struct {
	events     *events.Router
	reconciler *reconciler.Reconciler
}

func NewEngine() *Engine {
	return &Engine{
		events:     events.NewRouter(),
		reconciler: reconciler.NewReconciler(),
	}
}

func (e *Engine) RunWorkflow(name string, payload map[string]interface{}) error {
	switch name {
	case "compute.start":
		return compute.Start(payload["id"].(string))
	case "compute.stop":
		// ...
	}
	return nil
}

func (e *Engine) Reconcile() error {
	return e.reconciler.ReconcileInstance("instance-1")
}
