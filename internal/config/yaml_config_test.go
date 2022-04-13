package config

import (
	"reflect"
	"strings"
	"testing"
)

const testData = `
cn: test-ca
key:
  alg: rsa
  size: 2048
DNSs: 
  - localhost
  - test-ca
IPs:
  - 127.0.0.1
  - 192.168.1.1
expiry: 72h
isCA: true
serial:
  big: 1024.123
  attr:
    name: seq
`

func TestYamlConfigGet(t *testing.T) {
	var v *Value
	var err error

	reader := strings.NewReader(testData)
	c := NewYamlConfig()
	err = c.UnmarshalFromReader(reader)
	if err != nil {
		t.Fatalf("Parse yaml config fail, err: %v", err)
	}

	for _, tc := range []struct {
		key         string
		valueKind   reflect.Kind
		expectExist bool
		expectValue any
	}{
		{key: "cn", valueKind: reflect.String, expectExist: true, expectValue: "test-ca"},           // string value
		{key: "key.alg", valueKind: reflect.String, expectExist: true, expectValue: "rsa"},          // string value, multilple level key
		{key: "key.size", valueKind: reflect.Int, expectExist: true, expectValue: 2048},             // int value
		{key: "serial.big", valueKind: reflect.Float64, expectExist: true, expectValue: 1024.123},   // float64 value
		{key: "serial.attr.name", valueKind: reflect.String, expectExist: true, expectValue: "seq"}, // string value, multiple level key
		{key: "isCA", valueKind: reflect.Bool, expectExist: true, expectValue: true},                // boolean value
		{key: "key.alg_not_exists", valueKind: reflect.String, expectExist: false, expectValue: ""}, // Not exists key
	} {
		v = c.Get(tc.key)
		if tc.expectExist == false {
			if v != nil {
				t.Errorf("Fail on key \"%v\", expectExists false, but get return OK", tc.key)
			}
			continue
		}

		var rv any
		switch tc.valueKind {
		case reflect.String:
			rv = v.ToString("")
		case reflect.Int:
			rv = v.ToInt(0)
		case reflect.Bool:
			rv = v.ToBool(false)
		case reflect.Float64:
			rv = v.ToFloat64(0.0)
		default:
			t.Logf("Warning: Test on key \"%v\", but unknown valueKind: %v\n", tc.key, tc.valueKind)
			continue
		}
		if rv != tc.expectValue {
			t.Errorf("Fail on key \"%v\", value(\"%v\")!=expectValue(\"%v\")", tc.key, rv, tc.expectValue)
		}
	}
}

func TestYamlConfigUnmarshal(t *testing.T) {
	reader := strings.NewReader(testData)
	c := NewYamlConfig()
	err := c.UnmarshalFromReader(reader)
	if err != nil {
		t.Fatalf("Parse yaml config fail, err: %v", err)
	}
	v := c.Get("key")
	if v == nil {
		t.Fatalf("\"%v\" not found", "key")
	}

	var kc struct {
		Alg  string `yaml:"alg,omitempty"`
		Size int    `yaml:"size,omitempty"`
	}
	err = v.Unmarshal(&kc)
	if err != nil {
		t.Fatalf("Unmarshal error: %v", err)
	}

	if kc.Alg != "rsa" {
		t.Errorf("kc.Alg(\"%v\")!=expect(\"%v\")", kc.Alg, "rsa")
	}
	if kc.Size != 2048 {
		t.Errorf("kc.Size(\"%v\")!=expect(\"%v\")", kc.Alg, 2048)
	}
}

func TestYamlConfigUnmarshalArray(t *testing.T) {
	reader := strings.NewReader(testData)
	c := NewYamlConfig()
	err := c.UnmarshalFromReader(reader)
	if err != nil {
		t.Fatalf("Parse yaml config fail, err: %v", err)
	}
	v := c.Get("IPs")
	if v == nil {
		t.Fatalf("\"%v\" not found", "IPs")
	}

	var ips []string
	err = v.Unmarshal(&ips)
	if err != nil {
		t.Fatalf("Unmarshal error: %v", err)
	}

	var expectIPs = []string{
		"127.0.0.1",
		"192.168.1.1",
	}
	if len(ips) != len(expectIPs) {
		t.Errorf("len-ips(\"%v\")!=len-expect-ips(\"%v\")", len(ips), len(expectIPs))
	}
	for i := 0; i < len(ips); i++ {
		if ips[i] != expectIPs[i] {
			t.Errorf("ips[%v](%v)!=expectIPs[%v](%v)", i, ips[i], i, expectIPs[i])
			break
		}
	}
}
