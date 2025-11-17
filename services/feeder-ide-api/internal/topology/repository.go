package topology

import (
	"sync"

	"github.com/google/uuid"
)

// Repository 定義拓樸儲存介面
type Repository interface {
	Create(topology *Topology) error
	GetByID(id string) (*Topology, error)
	Update(id string, topology *Topology) error
	Delete(id string) error
	List() ([]*Topology, error)
}

// InMemoryRepository 記憶體實作（開發用）
type InMemoryRepository struct {
	mu         sync.RWMutex
	topologies map[string]*Topology
}

// NewInMemoryRepository 建立新的記憶體 repository
func NewInMemoryRepository() *InMemoryRepository {
	return &InMemoryRepository{
		topologies: make(map[string]*Topology),
	}
}

func (r *InMemoryRepository) Create(topology *Topology) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if topology.ID == "" {
		topology.ID = uuid.New().String()
	}
	r.topologies[topology.ID] = topology
	return nil
}

func (r *InMemoryRepository) GetByID(id string) (*Topology, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	topology, exists := r.topologies[id]
	if !exists {
		return nil, ErrTopologyNotFound
	}
	return topology, nil
}

func (r *InMemoryRepository) Update(id string, topology *Topology) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.topologies[id]; !exists {
		return ErrTopologyNotFound
	}
	topology.ID = id
	r.topologies[id] = topology
	return nil
}

func (r *InMemoryRepository) Delete(id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.topologies[id]; !exists {
		return ErrTopologyNotFound
	}
	delete(r.topologies, id)
	return nil
}

func (r *InMemoryRepository) List() ([]*Topology, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	topologies := make([]*Topology, 0, len(r.topologies))
	for _, topology := range r.topologies {
		topologies = append(topologies, topology)
	}
	return topologies, nil
}

