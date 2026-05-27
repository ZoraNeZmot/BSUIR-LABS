// Package vehicles contains the concrete classes of the Vehicle
// hierarchy used by lab 3. Each class registers itself with the central
// registry from its own init() function so that Marshaling, Unmarshaling
// and the GUI can stay completely generic.
package vehicles

import (
	"fmt"

	"oop/lab3/internal/vehicle"
)

// commonBase carries identification fields that every vehicle owns:
// a unique identifier, the manufacturer and model strings and the
// production year. All concrete classes embed it.
type commonBase struct {
	ID           string
	Manufacturer string
	Model        string
	Year         int
}

// commonFields produces the list of FieldDescriptors for the embedded
// commonBase. Subclasses simply append their own descriptors.
func (c *commonBase) commonFields() []vehicle.FieldDescriptor {
	return []vehicle.FieldDescriptor{
		vehicle.StringField("ID", "Identifier", &c.ID),
		vehicle.StringField("Manufacturer", "Manufacturer", &c.Manufacturer),
		vehicle.StringField("Model", "Model", &c.Model),
		vehicle.IntField("Year", "Year", &c.Year),
	}
}

// summaryPrefix is a shared helper that builds the leading part of the
// short list label so subclasses can append their own suffix.
func (c *commonBase) summaryPrefix(typeName string) string {
	return fmt.Sprintf("[%s] %s %s (%d)", typeName, c.Manufacturer, c.Model, c.Year)
}
