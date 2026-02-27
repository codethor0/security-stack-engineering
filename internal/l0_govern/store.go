package l0_govern

import "sync"

// GovTokenStore stores governance tokens for audit/replay.
type GovTokenStore struct {
	mu     sync.RWMutex
	tokens []*GovernanceToken
}

// NewGovTokenStore creates a token store.
func NewGovTokenStore() *GovTokenStore {
	return &GovTokenStore{tokens: make([]*GovernanceToken, 0)}
}

// Store appends a token.
func (s *GovTokenStore) Store(token *GovernanceToken) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.tokens = append(s.tokens, token)
	return nil
}

// Latest returns the most recent token.
func (s *GovTokenStore) Latest() *GovernanceToken {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if len(s.tokens) == 0 {
		return nil
	}
	return s.tokens[len(s.tokens)-1]
}
