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
		key string

		expectExist bool
		expectValue any
	}{
		{key: "cn", expectValue: "test-ca"},                 // string value
		{key: "key.alg", expectValue: "rsa"},                // string value, multilple level key
		{key: "key.size", expectValue: int64(2048)},         // int value
		{key: "serial.big", expectValue: float64(1024.123)}, // float64 value
		{key: "serial.attr.name", expectValue: "seq"},       // string value, multiple level key
		{key: "isCA", expectValue: true},                    // boolean value
		{key: "key.alg_not_exists", expectValue: nil},       // Not exists key
	} {
		v = c.Get(tc.key)

		var vv any
		switch tc.expectValue.(type) {
		case string:
			vv = v.ToString("")
		case int64:
			vv = v.ToInt(0)
		case bool:
			vv = v.ToBool(false)
		case float64:
			vv = v.ToFloat(0.0)
		case nil:
			// Expect not exists
			// Set vv=nil
			vv = nil
		default:
			t.Logf("Warning: Test on key \"%v\", unsupported value type: %v, you can use type constraint on expectValue\n",
				tc.key, reflect.TypeOf(tc.expectValue))
			continue
		}
		if vv != tc.expectValue {
			t.Errorf("Fail on key \"%v\", value(\"%v\")!=expectValue(\"%v\")", tc.key, vv, tc.expectValue)
			continue
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
