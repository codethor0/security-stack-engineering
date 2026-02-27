package l0_govern

import (
	"crypto/ed25519"
	"crypto/rand"
)

// Ed25519Signer implements SigningKey with Ed25519.
type Ed25519Signer struct {
	priv ed25519.PrivateKey
}

// NewEd25519Signer creates a signer with a new keypair.
func NewEd25519Signer() (*Ed25519Signer, error) {
	_, priv, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return nil, err
	}
	return &Ed25519Signer{priv: priv}, nil
}

// Sign signs data and returns the signature.
func (s *Ed25519Signer) Sign(data []byte) ([]byte, error) {
	return ed25519.Sign(s.priv, data), nil
}
