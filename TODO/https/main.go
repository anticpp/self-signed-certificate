package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"log"
	"net/http"

	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

func main() {
	addr := flag.String("addr", ":4001", "Server address")
	insecure := flag.Bool("insecure", false, "Insecure connection")
	certFile := flag.String("certfile", "./cert.pem", "Certificate PEM file")
	keyFile := flag.String("keyfile", "./key.pem", "Key PEM file")
	flag.Parse()

	var mux = http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		fmt.Printf("Handle http request, URI: %v\n", req.URL.RequestURI())
		fmt.Fprintf(w, "Prouldly served with Go and HTTPS!")
	})

	srv := &http.Server{
		Addr:    *addr,
		Handler: h2c.NewHandler(mux, &http2.Server{}),
	}

	if *insecure == false {
		srv.TLSConfig = &tls.Config{
			MinVersion:               tls.VersionTLS12,
			PreferServerCipherSuites: true,
		}
	}

	var err error
	if *insecure == false {
		log.Printf("Starting server on https://localhost%v\n", *addr)
		err = srv.ListenAndServeTLS(*certFile, *keyFile)
	} else {
		log.Printf("Starting server on http://localhost%v\n", *addr)
		err = srv.ListenAndServe()
	}
	if err != nil {
		log.Fatalf("Fail to start https server: %v", err)
	}
}
