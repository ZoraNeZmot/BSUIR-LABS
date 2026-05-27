# OOTPiSP — laboratory works (group sequence number 12)

Each lab is a self-contained Go module with its own `go.mod` and Fyne
GUI. The four projects build on top of each other but the source tree
is duplicated so that each lab can be opened, audited and graded in
isolation.

| Folder | Topic                                | Variant |
|--------|--------------------------------------|---------|
| `lab3` | Object serialization                 | 3 — Text |
| `lab4` | Plugins (hierarchy)                  | —       |
| `lab5` | Plugins (functionality)              | 3 — Encryption |
| `lab6` | Patterns                             | —       |

Variant assignment uses the formula
`(порядковый_номер mod количество_вариантов) + 1`. With sequence number
**12** and 5 variants, the result is **3** for both labs that ask for
a variant.

## Build & run

Each lab follows the same pattern:

```
cd labN
go mod tidy
go run ./cmd/app
```

The first `go mod tidy` downloads Fyne (~150 MB once); subsequent labs
reuse the Go module cache.

## Highlights per lab

* **lab3** — A 9-class Vehicle hierarchy (Car / Motorcycle / Truck /
  Boat / Ship / Airplane / Helicopter behind common Land/Water/Air
  bases). Plain-text serializer; the editor dialog and the file format
  are built from `[]FieldDescriptor` so adding a new class never
  requires `if` / `switch` / reflection on its concrete type.
* **lab4** — JSON-descriptor plugin loader for new vehicle classes.
  Drop a `*.json` file in `plugins/`, press *Reload plugins*, and the
  GUI immediately knows how to add, edit, save and load instances of
  the new type. Two demo plugins (`bicycle.json`, `submarine.json`)
  ship in the repository.
* **lab5** — Variant 3 (encryption). Three reversible algorithms
  (`xor`, `caesar`, `aes-cfb`) registered via `init()`; instantiated
  through JSON descriptors that live in `funcplugins/`. *Plugin
  settings…* exposes a generic editor for every algorithm parameter
  (10-point bonus). Save and Load run the chain forward and backward
  through `funcplugins.Encode` / `funcplugins.Decode`.
* **lab6** — Three explicit pattern demonstrations:
  * **Adapter** — `internal/funcplugins/algos/checksum_adapter.go`
    bridges the friend's library `internal/friendlib` (variant 5,
    *Save the checksum*) which has a sidecar-string API into the
    host's single-`[]byte` `Algorithm` interface.
  * **Singleton** — `internal/appsettings.Instance()` returns the
    process-wide settings struct used across UI / loader / save flow.
  * **Observer** — `internal/events.Bus` lets the master list, the
    dirty flag in the title bar and the audit log subscribe to mutation
    events independently.

## Compliance

* Comments are written in English (per project rules).
* Each lab has its own folder.
* No `if` / `switch` / reflection branches on concrete vehicle classes
  anywhere — both the serializer and the editor are completely
  generic.
* Tests:
  * `lab3`: `go test ./internal/storage/...` round-trips every built-in
    class through the text format.
  * `lab6`: `go test ./internal/funcplugins/...` round-trips the full
    `caesar -> xor -> aes -> friend-checksum` pipeline and asserts that
    the checksum adapter rejects tampered payloads.
