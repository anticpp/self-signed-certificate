package config

import (
	"fmt"
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

		if tc.expectKV != nil && (kv.key != tc.expectKV.key || kv.val != tc.expectKV.val) {
			t.Errorf("args [%v] fail, kv(%v)!=expectKV(%v)\n", tc.args, kv, tc.expectKV)
			continue
		}
	}
}

func TestCommandlineConfigParse(t *testing.T) {
	args := []string{
		"-cn=test-cn",
		"-key.alg=rsa",
		"-key.size=2048",
	}
	c := NewCommandlineConfig(args, "")
	err := c.Parse()
	if err != nil {
		t.Fatalf("args [%v] fail, Parse config error: %v", args, err)
	}
	fmt.Println(c)
}
