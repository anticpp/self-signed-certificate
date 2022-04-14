package config

import (
	"strconv"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

// Value is an abstract object simple wrap of `any` type.
type Value struct {
	any
}

// Unmarshal Value to out structure.
// Implemented by reusing `yaml.Marshal` and `yaml.Unmarshal`.
func (v *Value) Unmarshal(out any) error {
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
func (v *Value) TryString() (string, bool) {
	vv, ok := v.any.(string)
	return vv, ok
}

// Try to convert value to int.
// Alway use int64 here, thus int8/int16/int/int32 will be assigned to int64.
// If value type is string, try to parse to int64.
// The `bool` return true if success, or false if fails.
func (v *Value) TryInt() (int64, bool) {
	var vv int64
	var ok bool

	// TODO: Destroy the boring `case`.
	switch vvv := v.any.(type) {
	case int8:
		vv = int64(vvv)
		ok = true
	case int16:
		vv = int64(vvv)
		ok = true
	case int:
		vv = int64(vvv)
		ok = true
	case int32:
		vv = int64(vvv)
		ok = true
	case int64:
		vv = int64(vvv)
		ok = true
	case string:
		var err error
		vv, err = strconv.ParseInt(vvv, 10, 64)
		if err != nil {
			ok = false
			break
		}
		ok = true
	default:
		ok = false
	}
	return vv, ok
}

// Try to convert value to float46.
// float32 will be assigned to float64.
// The `bool` return true if success, or false if fails.
func (v *Value) TryFloat() (float64, bool) {
	var vv float64
	var ok bool
	switch vvv := v.any.(type) {
	case float32:
		vv = float64(vvv)
		ok = true
	case float64:
		vv = vvv
		ok = true
	case string:
		var err error
		vv, err = strconv.ParseFloat(vvv, 64)
		if err != nil {
			ok = false
			break
		}
		ok = true
	default:
		ok = false
	}
	return vv, ok

}

// Try to convert value to boolean.
// The `bool` return true if success, or false if fails.
func (v *Value) TryBool() (bool, bool) {
	vv, ok := v.any.(bool)
	return vv, ok
}

// Try to convert value to string.
// Return `defaultValue` if fails.
func (v *Value) ToString(defaultValue string) string {
	vv, ok := v.TryString()
	if !ok {
		return defaultValue
	}
	return vv
}

// Try to convert value to int.
// Always use int64 here.
// Return `defaultValue` if fails.
func (v *Value) ToInt(defaultValue int64) int64 {
	vv, ok := v.TryInt()
	if !ok {
		return defaultValue
	}
	return vv
}

// Try to convert value to float64.
// Return `defaultValue` if fails.
func (v *Value) ToFloat(defaultValue float64) float64 {
	vv, ok := v.TryFloat()
	if !ok {
		return defaultValue
	}
	return vv
}

// Try to convert value to boolean.
// Return `defaultValue` if fails.
func (v *Value) ToBool(defaultValue bool) bool {
	vv, ok := v.TryBool()
	if !ok {
		return defaultValue
	}
	return vv
}

// Try to convert value to string.
// Panic if fails.
func (v *Value) MustString() string {
	vv, ok := v.TryString()
	if !ok {
		panic("TryString fail")
	}
	return vv
}

// Try to convert value to int.
// Always use int64.
// Panic if fails.
func (v *Value) MustInt() int64 {
	vv, ok := v.TryInt()
	if !ok {
		panic("TryInt fail")
	}
	return vv
}

// Try to convert value to float64.
// Panic if fails.
func (v *Value) MustFloat() float64 {
	vv, ok := v.TryFloat()
	if !ok {
		panic("TryFloat fail")
	}
	return vv
}

// Try to convert value to boolean.
// Panic if fails.
func (v *Value) MustBool() bool {
	vv, ok := v.TryBool()
	if !ok {
		panic("TryBool fail")
	}
	return vv
}
