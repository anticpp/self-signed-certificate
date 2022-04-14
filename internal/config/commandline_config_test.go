package config

import (
	"reflect"
	"testing"
)

func TestCommandlineConfigParseNext(t *testing.T) {
	for _, tc := range []struct {
		args     []string
		prefix   string
		expectKV *kv
	}{
		// No prefix
		{[]string{"-cn=test-cn"}, "", &kv{"cn", "test-cn"}},
		{[]string{"--cn=test-cn"}, "", &kv{"cn", "test-cn"}},
		{[]string{"-cn", "test-cn"}, "", &kv{"cn", "test-cn"}},
		{[]string{"--cn", "test-cn"}, "", &kv{"cn", "test-cn"}},

		// With prefix
		{[]string{"-someprefix.cn=test-cn"}, "someprefix", &kv{"cn", "test-cn"}},
		{[]string{"--someprefix.cn=test-cn"}, "someprefix", &kv{"cn", "test-cn"}},
		{[]string{"-someprefix.cn", "test-cn"}, "someprefix", &kv{"cn", "test-cn"}},
		{[]string{"--someprefix.cn", "test-cn"}, "someprefix", &kv{"cn", "test-cn"}},

		// Prefix filter
		{[]string{"-cn=test-cn"}, "someprefix", nil},
		{[]string{"-cn", "test-cn"}, "someprefix", nil},
		{[]string{"-otherprefix.cn=test-cn"}, "someprefix", nil},
		{[]string{"-otherprefix.cn", "test-cn"}, "someprefix", nil},

		// Mutiple level key
		{[]string{"-key.alg=rsa"}, "", &kv{"key.alg", "rsa"}},
		{[]string{"-key.size=2048"}, "", &kv{"key.size", "2048"}},
	} {
		c := NewCommandlineConfig(tc.args, tc.prefix)
		kv, err := c.parseNext()
		if err != nil {
			t.Errorf("args [%v] fail, Parse config error: %v", tc.args, err)
			continue
		}

		if tc.expectKV == nil && kv != nil {
			t.Errorf("args [%v] fail, expectKV==nil but kv=nil\n", tc.args)
			continue
		}
		if tc.expectKV != nil && kv == nil {
			t.Errorf("args [%v] fail, expectKV!=nil but kv==nil\n", tc.args)
			continue
		}

		if tc.expectKV == nil {
			continue
		}

		if kv.key != tc.expectKV.key || kv.val != tc.expectKV.val {
			t.Errorf("args [%v] fail, kv(%v)!=expectKV(%v)\n", tc.args, kv, tc.expectKV)
			continue
		}
	}
}

func TestCommandlineConfigParse(t *testing.T) {
	for _, tc := range []struct {
		args      []string
		expectKVs []*kv
	}{
		{
			[]string{"-cn=test-cn", "-key.alg=rsa", "-key.size=2048", "-serial.attr.name=serial1", "-serial.big=1024.123", "-serial.attr.name=serial2"},
			[]*kv{
				{"cn", "test-cn"},
				{"key.alg", "rsa"},
				{"key.size", "2048"},
				{"key.size", int64(2048)}, // string can be interpreted as int64
				{"serial.big", "1024.123"},
				{"serial.big", float64(1024.123)}, // string can interpreted as float64
				{"key.alg_not_exists", nil},
				{"serial.attr.name", "serial2"}, // The last argument will be used
			},
		},
	} {
		c := NewCommandlineConfig(tc.args, "")
		err := c.Parse()
		if err != nil {
			t.Errorf("Test args [%v] fail, Parse config error: %v", tc.args, err)
			continue
		}

		for _, expectKV := range tc.expectKVs {
			v := c.Get(expectKV.key)

			var vv any
			switch expectKV.val.(type) {
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
				t.Logf("Warning: Test args %v, key \"%v\", unsupported value type: %v, you can use type constraint on expectValue\n",
					tc.args, expectKV.key, reflect.TypeOf(expectKV.val))
				continue
			} // endof switch

			if vv != expectKV.val {
				t.Errorf("Test args %v, key \"%v\", value(\"%v\")!=expectValue(\"%v\")",
					tc.args, expectKV.key, vv, expectKV.val)
				continue
			}
		} // endof for _, expectKV {}
	} // endof for _, tc {}
}
