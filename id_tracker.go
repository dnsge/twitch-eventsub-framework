package eventsub

import (
	"context"
	"sync"
)

// IDTracker keeps track of EventSub Message IDs that have already been processed.
type IDTracker interface {
	// AddAndCheckIfDuplicate returns if the ID is a duplicate and an error.
	AddAndCheckIfDuplicate(ctx context.Context, id string) (bool, error)
}

// MapTracker uses an in-memory map to check if a notification ID is
// a duplicate.
type MapTracker struct {
	mu   sync.Mutex
	seen map[string]struct{}
}

func (m *MapTracker) AddAndCheckIfDuplicate(_ context.Context, id string) (bool, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	_, ok := m.seen[id]
	if ok {
		return true, nil
	}

	m.seen[id] = struct{}{}
	return false, nil
}

// NewMapTracker creates a new MapTracker instance which uses an in-memory map
// to check if a notification ID is a duplicate.
func NewMapTracker() *MapTracker {
	return &MapTracker{
		seen: make(map[string]struct{}),
	}
}

// TrackerFunc is a functional adapter for the IDTracker interface.
type TrackerFunc func(context.Context, string) (bool, error)

func (f TrackerFunc) AddAndCheckIfDuplicate(ctx context.Context, id string) (bool, error) {
	return f(ctx, id)
}
