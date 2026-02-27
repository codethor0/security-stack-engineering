package orchestrator

import "sync"

// InMemoryPersistence stores evidence in memory for PoC.
type InMemoryPersistence struct {
	mu      sync.RWMutex
	records []EvidenceRecord
}

// NewInMemoryPersistence creates an in-memory evidence store.
func NewInMemoryPersistence() *InMemoryPersistence {
	return &InMemoryPersistence{
		records: make([]EvidenceRecord, 0),
	}
}

// StoreEvidence appends a record.
func (p *InMemoryPersistence) StoreEvidence(record EvidenceRecord) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.records = append(p.records, record)
	return nil
}

// QueryEvidence returns records for a layer since the given timestamp.
func (p *InMemoryPersistence) QueryEvidence(layer string, since int64) ([]EvidenceRecord, error) {
	p.mu.RLock()
	defer p.mu.RUnlock()
	var out []EvidenceRecord
	for _, r := range p.records {
		if r.Layer == layer && r.Timestamp >= since {
			out = append(out, r)
		}
	}
	return out, nil
}

// AllEvidence returns all stored records (for tests and L8 reporting).
func (p *InMemoryPersistence) AllEvidence() []EvidenceRecord {
	p.mu.RLock()
	defer p.mu.RUnlock()
	out := make([]EvidenceRecord, len(p.records))
	copy(out, p.records)
	return out
}
