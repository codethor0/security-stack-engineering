package l2_identity

import (
	"sync"
	"time"
)

// IdentityGraph is the in-memory model of identities and their relationships.
type IdentityGraph struct {
	mu          sync.RWMutex
	Identities  map[string]Identity
	Relations   map[string][]string
	LastUpdate  int64
}

// Identity represents a single identity in the graph.
type Identity struct {
	ID         string            `json:"id"`
	Type       string            `json:"type"`
	Owner      string            `json:"owner"`
	Privileged bool              `json:"privileged"`
	Active     bool              `json:"active"`
	LastAuth   int64             `json:"last_auth"`
	Attributes map[string]string `json:"attributes"`
	RiskScore  float64           `json:"risk_score"`
}

// NewIdentityGraph creates an empty graph.
func NewIdentityGraph() *IdentityGraph {
	return &IdentityGraph{
		Identities: make(map[string]Identity),
		Relations:  make(map[string][]string),
		LastUpdate: time.Now().Unix(),
	}
}
