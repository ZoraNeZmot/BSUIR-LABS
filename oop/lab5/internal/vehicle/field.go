// Package vehicle defines the abstract Vehicle contract used by every
// concrete class. The hierarchy never relies on if/else, switch or
// reflection: each class advertises its own set of editable fields
// through a slice of FieldDescriptor closures.
package vehicle

import (
	"fmt"
	"strconv"
)

// FieldKind enumerates supported primitive field types. UI and serializer
// branch on it generically; concrete classes never have to know about it.
type FieldKind int

const (
	KindString FieldKind = iota
	KindInt
	KindFloat
	KindBool
)

// FieldDescriptor is a generic getter/setter pair around a single struct
// field. The owning class produces these in Fields(); the UI and the
// serializer iterate them without ever knowing the concrete type.
type FieldDescriptor struct {
	Name  string
	Label string
	Kind  FieldKind
	Get   func() string
	Set   func(string) error
}

// StringField builds a FieldDescriptor bound to a *string.
func StringField(name, label string, ptr *string) FieldDescriptor {
	return FieldDescriptor{
		Name:  name,
		Label: label,
		Kind:  KindString,
		Get:   func() string { return *ptr },
		Set:   func(s string) error { *ptr = s; return nil },
	}
}

// IntField builds a FieldDescriptor bound to a *int.
func IntField(name, label string, ptr *int) FieldDescriptor {
	return FieldDescriptor{
		Name:  name,
		Label: label,
		Kind:  KindInt,
		Get:   func() string { return strconv.Itoa(*ptr) },
		Set: func(s string) error {
			v, err := strconv.Atoi(s)
			if err != nil {
				return fmt.Errorf("field %q: %w", name, err)
			}
			*ptr = v
			return nil
		},
	}
}

// FloatField builds a FieldDescriptor bound to a *float64.
func FloatField(name, label string, ptr *float64) FieldDescriptor {
	return FieldDescriptor{
		Name:  name,
		Label: label,
		Kind:  KindFloat,
		Get:   func() string { return strconv.FormatFloat(*ptr, 'f', -1, 64) },
		Set: func(s string) error {
			v, err := strconv.ParseFloat(s, 64)
			if err != nil {
				return fmt.Errorf("field %q: %w", name, err)
			}
			*ptr = v
			return nil
		},
	}
}

// BoolField builds a FieldDescriptor bound to a *bool.
func BoolField(name, label string, ptr *bool) FieldDescriptor {
	return FieldDescriptor{
		Name:  name,
		Label: label,
		Kind:  KindBool,
		Get:   func() string { return strconv.FormatBool(*ptr) },
		Set: func(s string) error {
			v, err := strconv.ParseBool(s)
			if err != nil {
				return fmt.Errorf("field %q: %w", name, err)
			}
			*ptr = v
			return nil
		},
	}
}
