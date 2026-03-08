package driver

import "sync"

// Repository defines the interface for driver storage operations.
type Repository interface {
	Save(d *Driver)
	FindAll() []*Driver
	FindByID(id string) (*Driver, bool)
}

// InMemoryRepository is a thread-safe in-memory implementation of Repository.
type InMemoryRepository struct {
	mu      sync.RWMutex
	drivers map[string]*Driver
}

// NewInMemoryRepository creates a new InMemoryRepository.
func NewInMemoryRepository() *InMemoryRepository {
	return &InMemoryRepository{
		drivers: make(map[string]*Driver),
	}
}

// Save stores or updates a driver's record.
func (r *InMemoryRepository) Save(d *Driver) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.drivers[d.ID] = d
}

// FindAll returns all stored drivers.
func (r *InMemoryRepository) FindAll() []*Driver {
	r.mu.RLock()
	defer r.mu.RUnlock()
	result := make([]*Driver, 0, len(r.drivers))
	for _, d := range r.drivers {
		result = append(result, d)
	}
	return result
}

// FindByID returns a driver by their ID.
func (r *InMemoryRepository) FindByID(id string) (*Driver, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	d, ok := r.drivers[id]
	return d, ok
}
