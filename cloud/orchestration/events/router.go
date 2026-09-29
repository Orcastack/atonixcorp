package events

type Event struct {
	Type      string
	Message   string
	Timestamp int64
}

type Router struct{}

func NewRouter() *Router {
	return &Router{}
}

func (r *Router) Emit(e Event) {
	// send to DB, logs, portal, metrics
}
