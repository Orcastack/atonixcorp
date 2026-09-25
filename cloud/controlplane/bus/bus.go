package bus

type EventType string

type Event struct {
	Type EventType
	Data any
}

type Handler func(Event)

type Bus struct {
	handlers map[EventType][]Handler
}

func NewBus() *Bus {
	return &Bus{handlers: make(map[EventType][]Handler)}
}

func (b *Bus) Subscribe(t EventType, h Handler) {
	b.handlers[t] = append(b.handlers[t], h)
}

func (b *Bus) Publish(e Event) {
	if hs, ok := b.handlers[e.Type]; ok {
		for _, h := range hs {
			h(e)
		}
	}
}
