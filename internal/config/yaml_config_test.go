package config

import (
	"fmt"
	"strings"
	"testing"
)

const testData = `
cn: test-ca
names:
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
`

func TestYamlConfigReadString(t *testing.T) {
	var v Value
	var ok bool
	var err error

	reader := strings.NewReader(testData)
	c := NewYamlConfig()
	err = c.UnmarshalFromReader(reader)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(c)

	v, ok = c.Get("cn")
	if ok {
		fmt.Println("cn:", v.ToString("default-cn"))
	} else {
		fmt.Println("\"cn\" not found")
	}
	/*
		v, ok = c0.Get("key.alg")
		if ok {
			fmt.Println("key.alg:", v.ToString("default-alg"))
		} else {
			fmt.Println("\"key.alg\" not found")
		}

		v, ok = c0.Get("key.alg_not_exists")
		if ok {
			fmt.Println("key.alg_not_exists:", v.ToString("default-alg"))
		} else {
			fmt.Println("\"key.alg_not_exists\" not found")
		}

		v, ok = c0.Get("key")
		if ok {
			fmt.Println("key:", v)
		} else {
			fmt.Println("\"key\" not found")
		}

		var kc keyCfg
		err = v.Unmarshal(&kc)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(kc)
	*/
}
