package main

import (
	"crypto/tls"
	"crypto/x509"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func main() {
	addr := flag.String("addr", "localhost:4001", "HTTPS Server address")
	caCertFile := flag.String("ca-certfile", "cert.pem", "trusted CA certificate")
	flag.Parse()

	pemCert, err := ioutil.ReadFile(*caCertFile)
	if err != nil {
		log.Fatal(err)
	}
	certPool := x509.NewCertPool()
	if ok := certPool.AppendCertsFromPEM(pemCert); !ok {
		log.Fatalf("unable to parse cert from %v", *caCertFile)
	}

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				RootCAs: certPool,
			},
		},
	}

	r, err := client.Get("https://" + *addr)
	if err != nil {
		log.Fatal(err)
	}
	defer r.Body.Close()

	html, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%v\n", r.Status)
	fmt.Println(string(html))
}
