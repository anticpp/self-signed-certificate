package main

import (
	"fmt"
	"log"

	"github.com/anticpp/self-signed-certificate/internal/config"
)

// xtls genca -config ca.yaml
// xtls genca -cn -expiry -DNSs -IPs -key.alg -key.size
// xtls gencert -config cert.yaml
// xtls gencert -cacert -cakey

type keyCfg struct {
	Alg  string `yaml:"alg,omitempty"`
	Size int    `yaml:"size,omitempty"`
}

func main() {
	c0 := config.NewYamlConfig()
	c0.UnmarshalFromFile("./conf/ca.yaml")
	fmt.Println(c0)

	var v config.Value
	var ok bool
	var err error

	v, ok = c0.Get("cn")
	if ok {
		fmt.Println("cn:", v.ToString("default-cn"))
	} else {
		fmt.Println("\"cn\" not found")
	}

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
	/*
		// Create private key
		privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
		if err != nil {
			log.Fatal(err)
		}

		// Create certificate template
		serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
		serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
		if err != nil {
			log.Fatal(err)
		}

		template := x509.Certificate{
			SerialNumber: serialNumber,
			Subject: pkix.Name{
				Organization: []string{"No Corp"},
				CommonName:   "test-ca",
			},
			DNSNames:  []string{"test-ca"},
			NotBefore: time.Now(),
			NotAfter:  time.Now().Add(3 * time.Hour),

			//KeyUsage:              x509.KeyUsageDigitalSignature,
			//ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
			BasicConstraintsValid: true,
		}

		// Create certificate data
		certBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, &privateKey.PublicKey, privateKey)
		if err != nil {
			log.Fatal(err)
		}

		certPEMBlock := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: certBytes})
		if certPEMBlock == nil {
			log.Fatal(err)
		}
		if err := ioutil.WriteFile("cacert.pem", certPEMBlock, 0644); err != nil {
			log.Fatal(err)
		}

		// Create private key data
		keyBytes, err := x509.MarshalPKCS8PrivateKey(privateKey)
		if err != nil {
			log.Fatal(err)
		}
		keyPEMBlock := pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: keyBytes})
		if keyPEMBlock == nil {
			log.Fatal(err)
		}
		if err := ioutil.WriteFile("key.pem", keyPEMBlock, 0644); err != nil {
			log.Fatal(err)
		}*/
}
