package vehicles

import "oop/lab3/internal/vehicle"

// waterBase holds fields shared by every floating vehicle.
type waterBase struct {
	commonBase
	DisplacementTons float64
	HullMaterial     string
}

// waterFields keeps common identification fields first.
func (w *waterBase) waterFields() []vehicle.FieldDescriptor {
	out := w.commonFields()
	out = append(out,
		vehicle.FloatField("DisplacementTons", "Displacement, tons", &w.DisplacementTons),
		vehicle.StringField("HullMaterial", "Hull material", &w.HullMaterial),
	)
	return out
}

// Boat represents a small recreational watercraft.
type Boat struct {
	waterBase
	HasOutboardMotor bool
	LengthMeters     float64
}

func (b *Boat) TypeName() string { return "Boat" }
func (b *Boat) Category() string { return "Water" }
func (b *Boat) Summary() string  { return b.summaryPrefix("Boat") }
func (b *Boat) Fields() []vehicle.FieldDescriptor {
	out := b.waterFields()
	out = append(out,
		vehicle.BoolField("HasOutboardMotor", "Outboard motor", &b.HasOutboardMotor),
		vehicle.FloatField("LengthMeters", "Length, m", &b.LengthMeters),
	)
	return out
}

// Ship represents a large sea-going vessel.
type Ship struct {
	waterBase
	CrewSize    int
	IsContainer bool
}

func (s *Ship) TypeName() string { return "Ship" }
func (s *Ship) Category() string { return "Water" }
func (s *Ship) Summary() string  { return s.summaryPrefix("Ship") }
func (s *Ship) Fields() []vehicle.FieldDescriptor {
	out := s.waterFields()
	out = append(out,
		vehicle.IntField("CrewSize", "Crew size", &s.CrewSize),
		vehicle.BoolField("IsContainer", "Container ship", &s.IsContainer),
	)
	return out
}

func init() {
	vehicle.Register("Boat", func() vehicle.Vehicle { return &Boat{waterBase: waterBase{HullMaterial: "Fiberglass"}} })
	vehicle.Register("Ship", func() vehicle.Vehicle { return &Ship{waterBase: waterBase{HullMaterial: "Steel"}} })
}
