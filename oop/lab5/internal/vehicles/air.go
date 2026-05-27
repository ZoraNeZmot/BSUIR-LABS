package vehicles

import "oop/lab5/internal/vehicle"

// airBase holds fields shared by every flying vehicle.
type airBase struct {
	commonBase
	MaxAltitudeM  int
	CruiseSpeedKmh float64
}

func (a *airBase) airFields() []vehicle.FieldDescriptor {
	out := a.commonFields()
	out = append(out,
		vehicle.IntField("MaxAltitudeM", "Max altitude, m", &a.MaxAltitudeM),
		vehicle.FloatField("CruiseSpeedKmh", "Cruise speed, km/h", &a.CruiseSpeedKmh),
	)
	return out
}

// Airplane represents a fixed-wing aircraft.
type Airplane struct {
	airBase
	WingspanMeters float64
	EngineCount    int
}

func (a *Airplane) TypeName() string { return "Airplane" }
func (a *Airplane) Category() string { return "Air" }
func (a *Airplane) Summary() string  { return a.summaryPrefix("Airplane") }
func (a *Airplane) Fields() []vehicle.FieldDescriptor {
	out := a.airFields()
	out = append(out,
		vehicle.FloatField("WingspanMeters", "Wingspan, m", &a.WingspanMeters),
		vehicle.IntField("EngineCount", "Engine count", &a.EngineCount),
	)
	return out
}

// Helicopter represents a rotary-wing aircraft.
type Helicopter struct {
	airBase
	RotorDiameterM float64
	IsMilitary     bool
}

func (h *Helicopter) TypeName() string { return "Helicopter" }
func (h *Helicopter) Category() string { return "Air" }
func (h *Helicopter) Summary() string  { return h.summaryPrefix("Helicopter") }
func (h *Helicopter) Fields() []vehicle.FieldDescriptor {
	out := h.airFields()
	out = append(out,
		vehicle.FloatField("RotorDiameterM", "Rotor diameter, m", &h.RotorDiameterM),
		vehicle.BoolField("IsMilitary", "Military", &h.IsMilitary),
	)
	return out
}

func init() {
	vehicle.Register("Airplane", func() vehicle.Vehicle { return &Airplane{airBase: airBase{MaxAltitudeM: 12000}, EngineCount: 2} })
	vehicle.Register("Helicopter", func() vehicle.Vehicle { return &Helicopter{airBase: airBase{MaxAltitudeM: 5000}} })
}
