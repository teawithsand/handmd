package nonce

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"io"
)

type Generator interface {
	GenerateNonce(ctx context.Context) (nce string, err error)
}

type HexGenerator struct {
	BytesLength int
	RNG         io.Reader
}

func (hng *HexGenerator) GenerateNonce(ctx context.Context) (nce string, err error) {
	rng := hng.RNG
	if rng == nil {
		rng = rand.Reader
	}
	buf := make([]byte, hng.BytesLength)
	_, err = io.ReadFull(rng, buf)
	if err != nil {
		return
	}

	nce = hex.EncodeToString(buf)
	return
}
