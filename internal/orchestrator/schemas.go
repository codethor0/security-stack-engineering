package orchestrator

import (
	"crypto/rand"
	"encoding/hex"
)

// GenerateID returns a random hex string for message/entity IDs.
func GenerateID() string {
	b := make([]byte, 8)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}
