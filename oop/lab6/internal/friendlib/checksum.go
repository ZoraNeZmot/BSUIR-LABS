// Package friendlib emulates a plugin received from a fellow student.
// It is intentionally written with a non-conforming API so that lab 6
// has a real reason to apply the Adapter pattern: the host's pipeline
// expects (Encode, Decode) over a single []byte, while this library
// returns a sidecar metadata string in addition to the payload.
//
// The variant the friend implemented is #5 -- "Saving the checksum".
// Their interface is:
//
//	type ChecksumKeeper struct{ Algorithm string }
//	func (k *ChecksumKeeper) Wrap(payload []byte) (out []byte, sidecar string, err error)
//	func (k *ChecksumKeeper) Verify(payload []byte, sidecar string) error
//
// Wrap returns the payload unchanged together with a checksum string so
// the caller is supposed to store the sidecar separately. Verify
// recomputes the digest and reports a mismatch.
//
// Note: The package is *self-contained*. It does not import anything
// from the host; it is the adapter (in funcplugins/algos/) that bridges
// the gap.
package friendlib

import (
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"hash"
	"strings"
)

// ChecksumKeeper is the type the friend exported. It selects the digest
// algorithm via a string field, again unlike anything in the host.
type ChecksumKeeper struct {
	Algorithm string // "sha256" (default) or "md5"
}

// hasherFor maps the algorithm name to a fresh hash.Hash. New
// algorithms only need an extra entry, no consumer-side switch.
var hasherFor = map[string]func() hash.Hash{
	"":       sha256.New,
	"sha256": sha256.New,
	"md5":    md5.New,
}

// Wrap computes a hex digest of payload and returns both unchanged
// payload and the sidecar string that must travel with it.
func (k *ChecksumKeeper) Wrap(payload []byte) (out []byte, sidecar string, err error) {
	factory, ok := hasherFor[strings.ToLower(k.Algorithm)]
	if !ok {
		return nil, "", fmt.Errorf("friendlib: unsupported algorithm %q", k.Algorithm)
	}
	h := factory()
	h.Write(payload)
	return payload, hex.EncodeToString(h.Sum(nil)), nil
}

// Verify recomputes the digest of payload and compares it against the
// saved sidecar.
func (k *ChecksumKeeper) Verify(payload []byte, sidecar string) error {
	factory, ok := hasherFor[strings.ToLower(k.Algorithm)]
	if !ok {
		return fmt.Errorf("friendlib: unsupported algorithm %q", k.Algorithm)
	}
	h := factory()
	h.Write(payload)
	got := hex.EncodeToString(h.Sum(nil))
	if got != sidecar {
		return errors.New("friendlib: checksum mismatch")
	}
	return nil
}
