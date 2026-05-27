# Lab 4 — Plugins (hierarchy)

Builds on top of lab 3. The host program ships seven built-in vehicle
classes and discovers additional ones at runtime by scanning the
`plugins/` directory.

## Why JSON descriptors instead of native Go plugins

Go's standard `plugin` package only works on Linux/macOS. Because the
target platform of this lab is Windows, runtime native plugins are not
viable. The JSON descriptor approach satisfies every requirement of the
task while staying fully cross-platform:

* a new module truly extends the hierarchy (it adds a class with its
  own field set),
* all generic functions (Marshal/Unmarshal/edit dialog/master list)
  immediately recognise it,
* loading is **dynamic**: drop the file in `plugins/` and either start
  the program or hit *Reload plugins* — no rebuild, no main-module
  edit.

## Plugin format

```json
{
  "typeName": "Bicycle",
  "category": "Land",
  "summary": "[%s] %s %s (%s)",
  "summaryFields": ["TypeName", "Manufacturer", "Model", "Year"],
  "fields": [
    {"name": "ID",            "label": "Identifier",      "kind": "string"},
    {"name": "Manufacturer",  "label": "Manufacturer",    "kind": "string"},
    {"name": "Model",         "label": "Model",           "kind": "string"},
    {"name": "Year",          "label": "Year",            "kind": "int",   "default": "2024"},
    {"name": "GearCount",     "label": "Number of gears", "kind": "int",   "default": "21"},
    {"name": "FrameMaterial", "label": "Frame material",  "kind": "string"},
    {"name": "IsElectric",    "label": "Electric",        "kind": "bool"}
  ]
}
```

Supported `kind` values: `string`, `int`, `float`, `bool`.

Two example descriptors are bundled:

* `plugins/bicycle.json` — adds a Bicycle class to the Land branch.
* `plugins/submarine.json` — adds a Submarine class to the Water branch.

## Build & run

```
cd lab4
go mod tidy
go run ./cmd/app
go run ./cmd/app -plugins=other/dir   # alternative directory
```

Use the *Add plugin file…* button to load a descriptor that lives
outside the configured directory.
