package config

import (
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

// Value is abstract return value of config.Get()
type Value struct {
	any
}

// Unmarshal Valut to out structure.
// Implemented by reusing yaml.Marshal and yaml.Unmarshal.
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

// Try to convert value to string.
// Return `defaultValue` if fails.
func (v Value) ToString(defaultValue string) string {
	vv, ok := v.TryString()
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
