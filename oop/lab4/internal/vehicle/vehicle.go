package vehicle

// Vehicle is the abstract contract for every node in the hierarchy.
// Adding a new concrete class only requires implementing this interface
// and registering a factory in the registry: no existing code must
// change, no if/else, switch or reflection is involved.
type Vehicle interface {
	// TypeName returns the unique tag that identifies the concrete
	// class both in the registry and in the serialized text file.
	TypeName() string

	// Category returns a human-readable group ("Land", "Water", "Air")
	// used by the UI to organise the available types.
	Category() string

	// Summary returns a short label rendered in the master list.
	Summary() string

	// Fields returns the list of editable fields backing the object.
	// Field order is preserved between Marshal/Unmarshal.
	Fields() []FieldDescriptor
}
