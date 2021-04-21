package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

func main() {
	addr := flag.String("addr", ":4001", "HTTPS address")
	//certFile := flag.String("certfile", "./cert.pem", "Certificate PEM file")
	//keyFile := flag.String("keyfile", "./key.pem", "Key PEM file")
	flag.Parse()

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		fmt.Printf("Handle http request, URI: %v\n", req.URL.RequestURI())
		fmt.Fprintf(w, "Prouldly served with Go and HTTPS!")
	})

	h2s := &http2.Server{}
	srv := &http.Server{
		Addr: *addr,
		//Handler: mux,
		Handler: h2c.NewHandler(mux, h2s),
		/*TLSConfig: &tls.Config{
			MinVersion:               tls.VersionTLS12,
			PreferServerCipherSuites: true,
		},*/
	}

	log.Printf("Starting server on %v\n", *addr)
	//err := srv.ListenAndServeTLS(*certFile, *keyFile)
	err := srv.ListenAndServe()
	if err != nil {
		log.Fatalf("Fail to start https server: %v", err)
	}
}
