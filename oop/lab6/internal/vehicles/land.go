package vehicles

import "oop/lab6/internal/vehicle"

// landBase carries fields shared by every road-going vehicle.
type landBase struct {
	commonBase
	MaxSpeedKmh float64
	WheelCount  int
}

// landFields concatenates common identification fields with land-only
// extras so every concrete land vehicle keeps the same prefix order.
func (l *landBase) landFields() []vehicle.FieldDescriptor {
	out := l.commonFields()
	out = append(out,
		vehicle.FloatField("MaxSpeedKmh", "Max speed, km/h", &l.MaxSpeedKmh),
		vehicle.IntField("WheelCount", "Wheel count", &l.WheelCount),
	)
	return out
}

// Car represents a passenger automobile.
type Car struct {
	landBase
	NumDoors      int
	HasAirCon     bool
	BodyStyle     string
}

func (c *Car) TypeName() string { return "Car" }
func (c *Car) Category() string { return "Land" }
func (c *Car) Summary() string  { return c.summaryPrefix("Car") }
func (c *Car) Fields() []vehicle.FieldDescriptor {
	out := c.landFields()
	out = append(out,
		vehicle.IntField("NumDoors", "Number of doors", &c.NumDoors),
		vehicle.BoolField("HasAirCon", "Air conditioner", &c.HasAirCon),
		vehicle.StringField("BodyStyle", "Body style", &c.BodyStyle),
	)
	return out
}

// Motorcycle represents a two-wheel motor vehicle.
type Motorcycle struct {
	landBase
	HasSidecar bool
	EngineCC   int
}

func (m *Motorcycle) TypeName() string { return "Motorcycle" }
func (m *Motorcycle) Category() string { return "Land" }
func (m *Motorcycle) Summary() string  { return m.summaryPrefix("Motorcycle") }
func (m *Motorcycle) Fields() []vehicle.FieldDescriptor {
	out := m.landFields()
	out = append(out,
		vehicle.BoolField("HasSidecar", "Has sidecar", &m.HasSidecar),
		vehicle.IntField("EngineCC", "Engine displacement, cc", &m.EngineCC),
	)
	return out
}

// Truck represents a heavy cargo road vehicle.
type Truck struct {
	landBase
	PayloadKg    int
	IsArticulated bool
}

func (t *Truck) TypeName() string { return "Truck" }
func (t *Truck) Category() string { return "Land" }
func (t *Truck) Summary() string  { return t.summaryPrefix("Truck") }
func (t *Truck) Fields() []vehicle.FieldDescriptor {
	out := t.landFields()
	out = append(out,
		vehicle.IntField("PayloadKg", "Payload, kg", &t.PayloadKg),
		vehicle.BoolField("IsArticulated", "Articulated", &t.IsArticulated),
	)
	return out
}

// Self-registration: each concrete class registers a factory closure
// that returns a fresh, zero-initialised pointer of itself. This is the
// only "wiring" needed to plug a class into Marshal/Unmarshal/UI.
func init() {
	vehicle.Register("Car", func() vehicle.Vehicle { return &Car{landBase: landBase{WheelCount: 4}, NumDoors: 4} })
	vehicle.Register("Motorcycle", func() vehicle.Vehicle { return &Motorcycle{landBase: landBase{WheelCount: 2}} })
	vehicle.Register("Truck", func() vehicle.Vehicle { return &Truck{landBase: landBase{WheelCount: 6}} })
}
