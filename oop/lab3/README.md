# Lab 3 — Object Serialization (variant 3: Text)

Group sequence number: **12**.
Variant: `(12 mod 5) + 1 = 3` → plain text format.

## Subject area

The hierarchy models real-world vehicles. Nine concrete classes inherit
from three intermediate group classes that share a common identification
prefix:

```
Vehicle (interface)
├── Land
│   ├── Car
│   ├── Motorcycle
│   └── Truck
├── Water
│   ├── Boat
│   └── Ship
└── Air
    ├── Airplane
    └── Helicopter
```

That is **6 + 3 = 9 classes** (excluding the two abstract ones the task
already counts).

## Open / closed for extension

Adding a new concrete class is a three-step recipe and never modifies
existing code:

1. Define a new struct embedding one of the group bases (or `commonBase`
   directly).
2. Implement the four `Vehicle` interface methods. `Fields()` reuses the
   base's helper and appends the class-specific descriptors.
3. Register a factory in the package's `init()`.

There is **no `if` / `switch` / reflection** that branches on the
concrete class anywhere in the project — neither in the serializer nor
in the editor dialog. Both walk through the `[]FieldDescriptor` returned
by the object itself.

## Build & run

```
cd lab3
go mod tidy
go run ./cmd/app
```

## File format (text, variant 3)

```
[BEGIN Car]
ID=v-001
Manufacturer=Toyota
Model=Corolla
Year=2020
MaxSpeedKmh=180
WheelCount=4
NumDoors=4
HasAirCon=true
BodyStyle=Sedan
[END]

[BEGIN Helicopter]
...
[END]
```

Newlines, carriage returns and backslashes inside string fields are
escaped with a `\n` / `\r` / `\\` sequence so any value can be
round-tripped.

## UI features

* List of all loaded vehicles
* `Add` — pick a registered type, fill in fields
* `Edit` — modify selected vehicle
* `Delete` — remove selected vehicle
* `Save…` / `Load…` — text file via standard file dialog
* `Load default` — convenience button that reads `vehicles.txt` from the
  working directory.
