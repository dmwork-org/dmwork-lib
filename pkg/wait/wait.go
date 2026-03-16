package wait

import (
	"errors"
	"fmt"
	"sync"
)

const (
	// To avoid lock contention we use an array of list struct (rw mutex & map)
	// for the id argument, we apply mod operation and uses its remainder to
	// index into the array and find the corresponding element.
	defaultListElementLength = 64
)

// ErrDuplicateID is returned when attempting to register an ID that is already registered.
var ErrDuplicateID = errors.New("duplicate id")

// Wait Wait
type Wait interface {
	// Register waits returns a chan that waits on the given ID.
	// The chan will be triggered when Trigger is called with
	// the same ID. Returns ErrDuplicateID if the ID is already registered.
	Register(id uint64) (<-chan interface{}, error)
	// Trigger triggers the waiting chans with the given ID.
	Trigger(id uint64, x interface{})
	IsRegistered(id uint64) bool
}

type list struct {
	e []listElement
}

type listElement struct {
	l sync.RWMutex
	m map[uint64]chan interface{}
}

// New creates a Wait.
func New() Wait {
	res := list{
		e: make([]listElement, defaultListElementLength),
	}
	for i := 0; i < len(res.e); i++ {
		res.e[i].m = make(map[uint64]chan interface{})
	}
	return &res
}

func (w *list) Register(id uint64) (<-chan interface{}, error) {
	idx := id % defaultListElementLength
	newCh := make(chan interface{}, 1)
	w.e[idx].l.Lock()
	defer w.e[idx].l.Unlock()
	if _, ok := w.e[idx].m[id]; ok {
		return nil, fmt.Errorf("%w: %x", ErrDuplicateID, id)
	}
	w.e[idx].m[id] = newCh
	return newCh, nil
}

func (w *list) Trigger(id uint64, x interface{}) {
	idx := id % defaultListElementLength
	w.e[idx].l.Lock()
	ch := w.e[idx].m[id]
	delete(w.e[idx].m, id)
	w.e[idx].l.Unlock()
	if ch != nil {
		ch <- x
		close(ch)
	}
}

func (w *list) IsRegistered(id uint64) bool {
	idx := id % defaultListElementLength
	w.e[idx].l.RLock()
	defer w.e[idx].l.RUnlock()
	_, ok := w.e[idx].m[id]
	return ok
}
