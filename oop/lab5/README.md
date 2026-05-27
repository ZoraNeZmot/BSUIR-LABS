# Lab 5 — Plugins (functionality), variant 3 (Encryption / Decryption)

Builds on top of lab 4. The host program now exposes a generic
processing pipeline that plugins can hook into both before saving and
after loading. The variant assigned to sequence number 12 is
`(12 mod 5) + 1 = 3` → **Encryption / Decryption**.

## Two plugin systems side-by-side

| System          | Adds                | Lives in          | Format     |
|-----------------|---------------------|-------------------|------------|
| Hierarchy       | new vehicle classes | `plugins/*.json`  | JSON       |
| Functional      | encryption steps    | `funcplugins/*.json` | JSON     |

Both are loaded automatically at startup and can be reloaded or
inspected from the GUI. A new plugin file is enough — no recompilation,
no source changes.

## Built-in algorithms

Three reversible algorithms ship with the host. Each one is a tiny
strategy implementation that registers itself through `init()`:

| ID        | Display name        | Parameters     |
|-----------|---------------------|----------------|
| `xor`     | XOR Cipher          | `key`          |
| `caesar`  | Caesar (byte) Cipher| `shift`        |
| `aes-cfb` | AES-256-CFB         | `passphrase`   |

Three demo plugin descriptors instantiate these algorithms with concrete
parameter values:

* `funcplugins/xor-default.json`
* `funcplugins/caesar-7.json`
* `funcplugins/aes-strong.json`

You can drop additional descriptors that mix-and-match: e.g. a plugin
that does Caesar shift +13 and another that does AES with a different
passphrase. The "Plugin settings…" menu lets you enable/disable each
plugin and edit its parameters at runtime, satisfying the 10-point
bonus part of the task.

## Pipeline

```
Save flow:
  vehicles --Marshal--> text bytes
                         |
                         v
                    plugin #1.Encode
                         |
                         v
                    plugin #2.Encode  (in registration order)
                         |
                         v
                    raw bytes --> file

Load flow:
  file --> raw bytes
                         |
                         v
                    plugin #2.Decode  (reverse order)
                         |
                         v
                    plugin #1.Decode
                         |
                         v
                    text bytes --Unmarshal--> vehicles
```

## Build & run

```
cd lab5
go mod tidy
go run ./cmd/app
go run ./cmd/app -funcplugins=other/dir
```
