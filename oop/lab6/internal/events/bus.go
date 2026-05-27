// Package events is the lab-6 demonstration of the OBSERVER design
// pattern.
//
// Why an observer? Several independent sub-systems care about changes
// to the in-memory vehicle list:
//
//   - the master list widget needs to redraw,
//   - the audit log records every mutation,
//   - the title bar shows a small "*" when there are unsaved changes.
//
// Tying them together with direct method calls would couple the UI to
// every observer; an event bus lets each side stay independent and
// turn-on or off subscriptions at will.
package events

import "sync"

// Kind enumerates the events the host emits. New events are added by
// extending this constant set; observers branch on the small closed
// enum, never on object types, so the no-if-on-class rule from earlier
// labs still holds.
type Kind int

const (
	VehicleAdded Kind = iota
	VehicleEdited
	VehicleRemoved
	ListReplaced
	PluginsChanged
	SaveCompleted
	LoadCompleted
)

// Event is what observers receive. Payload carries an optional context
// value (typically the affected index or a message string).
type Event struct {
	Kind    Kind
	Payload any
}

// Listener is the observer signature.
type Listener func(Event)

// Bus is a tiny synchronous publish/subscribe broker.
type Bus struct {
	mu        sync.RWMutex
	listeners []Listener
}

// NewBus returns a fresh bus. Most apps use a single instance.
func NewBus() *Bus { return &Bus{} }

// Subscribe registers a listener. The returned function detaches it.
func (b *Bus) Subscribe(l Listener) (unsubscribe func()) {
	b.mu.Lock()
	defer b.mu.Unlock()
	id := len(b.listeners)
	b.listeners = append(b.listeners, l)
	return func() {
		b.mu.Lock()
		defer b.mu.Unlock()
		if id < len(b.listeners) {
			b.listeners[id] = nil
		}
	}
}

// Publish fans the event out to every active listener. Listeners run
// in the caller's goroutine, which is fine for a Fyne UI app.
func (b *Bus) Publish(e Event) {
	b.mu.RLock()
	snapshot := make([]Listener, len(b.listeners))
	copy(snapshot, b.listeners)
	b.mu.RUnlock()
	for _, l := range snapshot {
		if l != nil {
			l(e)
		}
	}
}
