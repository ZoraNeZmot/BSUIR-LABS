# Lab 6 — Patterns

Builds on top of lab 5. The hierarchy of vehicle classes, the JSON
descriptor plugins and the encryption pipeline are all carried over
unchanged. Lab 6 adds three explicit design-pattern demonstrations.

## 1. Adapter — `internal/friendlib` + `internal/funcplugins/algos/checksum_adapter.go`

The "friend's plugin" is an independent package (`internal/friendlib`)
that implements variant 5 of lab 5 — *Saving the checksum*. Its API is
**deliberately incompatible** with what the host expects:

```go
// What the host wants (single []byte in, single []byte out).
type Algorithm interface {
    Encode(data []byte, params map[string]string) ([]byte, error)
    Decode(data []byte, params map[string]string) ([]byte, error)
}

// What the friend offers (sidecar metadata, no params map).
type ChecksumKeeper struct{ Algorithm string }
func (k *ChecksumKeeper) Wrap(payload []byte) (out []byte, sidecar string, err error)
func (k *ChecksumKeeper) Verify(payload []byte, sidecar string) error
```

`checksumAdapter` wraps a `*friendlib.ChecksumKeeper` and rewires every
detail:

* it inlines the sidecar string into the payload as a `FRIENDCHK:<hex>:`
  header so a single byte slice carries both halves;
* it pulls the parameter out of the host's `map[string]string` and
  feeds it into the friend's struct field;
* it implements the four `Algorithm` methods and registers itself with
  the same global registry the built-in algorithms use, so it lights up
  in the GUI exactly like any other plugin.

Why Adapter is appropriate here: the existing host code cannot change
(the pipeline is fixed) and the friend's library cannot change either
(it is "third-party"). An adapter is the canonical way to reconcile two
fixed contracts.

## 2. Singleton — `internal/appsettings`

`appsettings.Instance()` returns one process-wide `*Settings` value
holding the user name (shown in the title), the last opened file (used
to seed the *Load default* dialog) and the two plugin directories.

Why Singleton is appropriate here: those settings have no natural
owner, every layer reads them, and threading the same struct through
every constructor would be needlessly verbose. The `sync.Once`
construction guarantees a single, lazily-initialised instance and the
internal `sync.RWMutex` keeps concurrent reads safe.

## 3. Observer — `internal/events`

`events.Bus` is a tiny synchronous publish/subscribe broker. The UI
hooks **three independent observers** on it:

1. the master list and statistics labels redraw,
2. a "dirty" flag is flipped on/off and reflected in the window title,
3. an append-only audit log shown at the bottom of the window records
   every mutation.

Why Observer is appropriate here: each subscriber has a different
concern and a different lifetime. With direct method calls, the
mutation sites would have to know about every consumer; with a bus,
new subscribers can attach and detach without touching the publishers.

## Running

```
cd lab6
go mod tidy
go run ./cmd/app
```

Open *Plugin settings…* and enable for example
`Friend's checksum (adapted)` together with `AES-256 (strong)` to see
the full pipeline (encrypt → checksum) at work; reload roundtrips the
file by running the inverse pipeline (checksum verify → decrypt).
