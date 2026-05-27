# Go Torrent Tracker (Coursework)

Minimal BitTorrent **HTTP tracker** in Go (standard library only), plus a small **seed** helper to register with the tracker and serve pieces over the peer wire protocol for demos.

## Features

- `/announce`, `/scrape`, `/healthz`
- In-memory swarms, stale-peer cleanup
- Structured logging ([logrus](https://github.com/sirupsen/logrus)): text or JSON, configurable level

## Project layout

| Path | Role |
|------|------|
| `cmd/tracker/main.go` | Tracker HTTP server |
| `cmd/seed/main.go` | Optional seeder (announce + TCP peer wire) |
| `internal/tracker/` | Swarm state, announce logic, cleanup |
| `internal/http/` | HTTP handlers |
| `internal/bencode/` | Encoder (tracker responses) + decoder (used by seed’s tracker client) |
| `internal/btwire/`, `internal/metainfo/`, `internal/trackerclient/` | Used by **seed** only |

## Run tracker

```bash
go run ./cmd/tracker
```

### Logging

| Env | Default | Meaning |
|-----|---------|--------|
| `TRACKER_LOG_LEVEL` | `info` | `debug`, `info`, `warn`, `error` |
| `TRACKER_LOG_JSON` | off | Set to `1` or `true` for JSON lines |

Flags: `-log-level`, `-log-json`

Other env/flags: `TRACKER_HOST`, `TRACKER_PORT`, `ANNOUNCE_INTERVAL_SEC`, `PEER_TIMEOUT_SEC`, `DEFAULT_NUMWANT`, `MAX_NUMWANT` (see `internal/config/config.go`).

## API examples

20-byte `info_hash` and `peer_id` in real use; placeholders below are 20 characters each.

```bash
curl "http://localhost:8080/announce?info_hash=12345678901234567890&peer_id=ABCDEFGHIJKLMNOPQRST&port=6881&uploaded=0&downloaded=0&left=100"
curl "http://localhost:8080/scrape?info_hash=12345678901234567890"
```

## Seed (optional)

Register with the tracker and serve the file named in the torrent:

```bash
go run ./cmd/seed -torrent path/to.torrent -content path/to/exact-size-file.bin -port 6882
```

Use `-tracker http://HOST:PORT/announce` if the torrent’s announce URL differs.

## Build binaries

```bash
go build -o tracker.exe ./cmd/tracker
go build -o seed.exe ./cmd/seed
```

## Tests

```bash
go test ./...
```

## Limitations

- In-memory tracker state (lost on restart)
- HTTP tracker only; no DHT / magnet / UDP tracker in this repo
