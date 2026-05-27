package algos

import (
	"bytes"
	"errors"
	"fmt"

	"oop/lab6/internal/friendlib"
	"oop/lab6/internal/funcplugins"
)

// checksumAdapter is the textbook example of the Adapter pattern.
//
// The host's plugin pipeline speaks a single, narrow interface:
//
//	Encode(data []byte, params) ([]byte, error)
//	Decode(data []byte, params) ([]byte, error)
//
// The friend's library exposes a wider, sidecar-based contract:
//
//	Wrap(payload) (out, sidecar, err)
//	Verify(payload, sidecar) error
//
// checksumAdapter owns a *friendlib.ChecksumKeeper and re-shapes its
// methods into the Algorithm contract. The encoded blob carries the
// digest as a fixed-length text prefix so a single []byte is enough to
// round-trip both halves of the friend's two-part output.
type checksumAdapter struct {
	keeper *friendlib.ChecksumKeeper
}

// header tags written before the actual payload. Using a fixed tag
// keeps the format self-describing and makes Decode robust against
// stray bytes from earlier pipeline steps.
const (
	checksumHeader = "FRIENDCHK:"
	headerSep      = ":"
)

func (a *checksumAdapter) ID() string          { return "friend-checksum" }
func (a *checksumAdapter) DisplayName() string { return "Friend's checksum (adapted)" }
func (a *checksumAdapter) Description() string {
	return "Adapter around a peer's plugin (variant 5 -- save checksum). Verifies file integrity on load."
}

func (a *checksumAdapter) Parameters() []funcplugins.ParamSpec {
	return []funcplugins.ParamSpec{
		{Name: "algorithm", Label: "Digest algorithm (sha256/md5)", Kind: funcplugins.ParamString, Default: "sha256"},
	}
}

// Encode bridges Wrap by inlining the sidecar string into the payload.
// Format:  FRIENDCHK:<digest>:<original-bytes>
func (a *checksumAdapter) Encode(data []byte, p map[string]string) ([]byte, error) {
	a.keeper.Algorithm = p["algorithm"]
	payload, sidecar, err := a.keeper.Wrap(data)
	if err != nil {
		return nil, err
	}
	out := bytes.NewBuffer(make([]byte, 0, len(payload)+len(sidecar)+len(checksumHeader)+1))
	out.WriteString(checksumHeader)
	out.WriteString(sidecar)
	out.WriteString(headerSep)
	out.Write(payload)
	return out.Bytes(), nil
}

// Decode is the inverse: pull the sidecar back out, verify, return the
// trailing payload.
func (a *checksumAdapter) Decode(data []byte, p map[string]string) ([]byte, error) {
	if !bytes.HasPrefix(data, []byte(checksumHeader)) {
		return nil, errors.New("friend-checksum: header missing")
	}
	tail := data[len(checksumHeader):]
	idx := bytes.IndexByte(tail, headerSep[0])
	if idx < 0 {
		return nil, errors.New("friend-checksum: separator missing")
	}
	sidecar := string(tail[:idx])
	payload := tail[idx+1:]
	a.keeper.Algorithm = p["algorithm"]
	if err := a.keeper.Verify(payload, sidecar); err != nil {
		return nil, fmt.Errorf("friend-checksum: %w", err)
	}
	return payload, nil
}

func init() {
	funcplugins.RegisterAlgorithm(&checksumAdapter{
		keeper: &friendlib.ChecksumKeeper{Algorithm: "sha256"},
	})
}
