package config

import (
	"reflect"
	"testing"
)

func TestScalar(t *testing.T) {
	for i, tc := range []struct {
		s      string
		expect any
	}{
		{"", ""},
		{"dog", "dog"},
		{"2048", int64(2048)},
		{"-2048", int64(-2048)},
		{"1024.12", float64(1024.12)},
		{"0.12", float64(0.12)},
		{"-1024.12", float64(-1024.12)},
		{"-0.12", float64(-0.12)},
		// Quoted string
		{"\"dog\"", "dog"},
		{"'dog'", "dog"},
		{"\"2048\"", "2048"},
		{"'2048'", "2048"},
		{"\"dog'", "dog"},
		{"'dog\"", "dog"},
	} {
		ss := parseScalar(tc.s)
		ssType := reflect.TypeOf(ss)
		expectType := reflect.TypeOf(tc.expect)
		if ssType != expectType {
			t.Errorf("Test[%v] \"%v\" fail, got scalar type(\"%v\")!=expectType(\"%v\")\n", i, tc.s, ssType, expectType)
			continue
		}

		if ss != tc.expect {
			t.Errorf("Test[%v] \"%v\" fail, got scalar(\"%v\")!=expectScalarSS(\"%v\")\n", i, tc.s, ss, tc.expect)
			continue
		}
	}
}
