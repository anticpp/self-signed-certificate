package config

import (
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

// Value is an abstract object simple wrap of `any` type.
type Value struct {
	any
}

// Unmarshal Value to out structure.
// Implemented by reusing `yaml.Marshal` and `yaml.Unmarshal`.
func (v Value) Unmarshal(out any) error {
	var data []byte
	var err error

	data, err = yaml.Marshal(v.any)
	if err != nil {
		return errors.Wrap(err, "Marshal fail")
	}
	err = yaml.Unmarshal(data, out)
	if err != nil {
		return errors.Wrap(err, "Unmarshal fail")
	}
	return nil
}

// Try to convert value to string.
// The `bool` return true if success, or false if fails.
func (v Value) TryString() (string, bool) {
	vv, ok := v.any.(string)
	return vv, ok
}

// Try to convert value to int.
// The `bool` return true if success, or false if fails.
func (v Value) TryInt() (int, bool) {
	vv, ok := v.any.(int)
	return vv, ok
}

// Try to convert value to float46.
// The `bool` return true if success, or false if fails.
func (v Value) TryFloat64() (float64, bool) {
	vv, ok := v.any.(float64)
	return vv, ok
}

// Try to convert value to boolean.
// The `bool` return true if success, or false if fails.
func (v Value) TryBool() (bool, bool) {
	vv, ok := v.any.(bool)
	return vv, ok
}

// Try to convert value to string.
// Return `defaultValue` if fails.
func (v Value) ToString(defaultValue string) string {
	vv, ok := v.TryString()
	if !ok {
		return defaultValue
	}
	return vv
}

// Try to convert value to int.
// Return `defaultValue` if fails.
func (v Value) ToInt(defaultValue int) int {
	vv, ok := v.TryInt()
	if !ok {
		return defaultValue
	}
	return vv
}

// Try to convert value to float64.
// Return `defaultValue` if fails.
func (v Value) ToFloat64(defaultValue float64) float64 {
	vv, ok := v.TryFloat64()
	if !ok {
		return defaultValue
	}
	return vv
}

// Try to convert value to boolean.
// Return `defaultValue` if fails.
func (v Value) ToBool(defaultValue bool) bool {
	vv, ok := v.TryBool()
	if !ok {
		return defaultValue
	}
	return vv
}

// Try to convert value to string.
// Panic if fails.
func (v Value) MustString() string {
	vv, ok := v.TryString()
	if !ok {
		panic("TryString fail")
	}
	return vv
}

// Try to convert value to int.
// Panic if fails.
func (v Value) MustInt() int {
	vv, ok := v.TryInt()
	if !ok {
		panic("TryInt fail")
	}
	return vv
}

// Try to convert value to float64.
// Panic if fails.
func (v Value) MustFloat64() float64 {
	vv, ok := v.TryFloat64()
	if !ok {
		panic("TryFloat64 fail")
	}
	return vv
}

// Try to convert value to boolean.
// Panic if fails.
func (v Value) MustBool() bool {
	vv, ok := v.TryBool()
	if !ok {
		panic("TryBool fail")
	}
	return vv
}
