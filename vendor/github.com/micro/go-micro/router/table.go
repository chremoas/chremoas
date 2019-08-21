package router

import (
	"errors"
	"sync"
	"time"

	"github.com/google/uuid"
)

// table is an in memory routing table
type table struct {
	sync.RWMutex
	// routes stores service routes
	routes map[string]map[uint64]Route
	// watchers stores table watchers
	watchers map[string]*tableWatcher
}

// newtable creates a new routing table and returns it
func newTable(opts ...Option) *table {
	return &table{
		routes:   make(map[string]map[uint64]Route),
		watchers: make(map[string]*tableWatcher),
	}
}

// Create creates new route in the routing table
func (t *table) Create(r Route) error {
	service := r.Service
	sum := r.Hash()

	t.Lock()
	defer t.Unlock()

	// check if there are any routes in the table for the route destination
	if _, ok := t.routes[service]; !ok {
		t.routes[service] = make(map[uint64]Route)
		t.routes[service][sum] = r
		go t.sendEvent(&Event{Type: Create, Timestamp: time.Now(), Route: r})
		return nil
	}

	// add new route to the table for the route destination
	if _, ok := t.routes[service][sum]; !ok {
		t.routes[service][sum] = r
		go t.sendEvent(&Event{Type: Create, Timestamp: time.Now(), Route: r})
		return nil
	}

	return ErrDuplicateRoute
}

// Delete deletes the route from the routing table
func (t *table) Delete(r Route) error {
	service := r.Service
	sum := r.Hash()

	t.Lock()
	defer t.Unlock()

	if _, ok := t.routes[service]; !ok {
		return ErrRouteNotFound
	}

	delete(t.routes[service], sum)
	go t.sendEvent(&Event{Type: Delete, Timestamp: time.Now(), Route: r})

	return nil
}

// Update updates routing table with the new route
func (t *table) Update(r Route) error {
	service := r.Service
	sum := r.Hash()

	t.Lock()
	defer t.Unlock()

	// check if the route destination has any routes in the table
	if _, ok := t.routes[service]; !ok {
		t.routes[service] = make(map[uint64]Route)
		t.routes[service][sum] = r
		go t.sendEvent(&Event{Type: Create, Timestamp: time.Now(), Route: r})
		return nil
	}

	t.routes[service][sum] = r
	go t.sendEvent(&Event{Type: Update, Timestamp: time.Now(), Route: r})

	return nil
}

// List returns a list of all routes in the table
func (t *table) List() ([]Route, error) {
	t.RLock()
	defer t.RUnlock()

	var routes []Route
	for _, rmap := range t.routes {
		for _, route := range rmap {
			routes = append(routes, route)
		}
	}

	return routes, nil
}

// isMatch checks if the route matches given network and router
func isMatch(route Route, network, router string) bool {
	if network == "*" || network == route.Network {
		if router == "*" || router == route.Gateway {
			return true
		}
	}
	return false
}

// findRoutes finds all the routes for given network and router and returns them
func findRoutes(routes map[uint64]Route, network, router string) []Route {
	var results []Route
	for _, route := range routes {
		if isMatch(route, network, router) {
			results = append(results, route)
		}
	}
	return results
}

// Lookup queries routing table and returns all routes that match the lookup query
func (t *table) Query(q Query) ([]Route, error) {
	t.RLock()
	defer t.RUnlock()

	if q.Options().Service != "*" {
		if _, ok := t.routes[q.Options().Service]; !ok {
			return nil, ErrRouteNotFound
		}
		return findRoutes(t.routes[q.Options().Service], q.Options().Network, q.Options().Gateway), nil
	}

	var results []Route
	// search through all destinations
	for _, routes := range t.routes {
		results = append(results, findRoutes(routes, q.Options().Network, q.Options().Gateway)...)
	}

	return results, nil
}

// Watch returns routing table entry watcher
func (t *table) Watch(opts ...WatchOption) (Watcher, error) {
	// by default watch everything
	wopts := WatchOptions{
		Service: "*",
	}

	for _, o := range opts {
		o(&wopts)
	}

	w := &tableWatcher{
		id:      uuid.New().String(),
		opts:    wopts,
		resChan: make(chan *Event, 10),
		done:    make(chan struct{}),
	}

	// when the watcher is stopped delete it
	go func() {
		<-w.done
		t.Lock()
		delete(t.watchers, w.id)
		t.Unlock()
	}()

	// save the watcher
	t.Lock()
	t.watchers[w.id] = w
	t.Unlock()

	return w, nil
}

// sendEvent sends events to all subscribed watchers
func (t *table) sendEvent(e *Event) {
	t.RLock()
	defer t.RUnlock()

	for _, w := range t.watchers {
		select {
		case w.resChan <- e:
		case <-w.done:
		}
	}
}

var (
	// ErrRouteNotFound is returned when no route was found in the routing table
	ErrRouteNotFound = errors.New("route not found")
	// ErrDuplicateRoute is returned when the route already exists
	ErrDuplicateRoute = errors.New("duplicate route")
)
