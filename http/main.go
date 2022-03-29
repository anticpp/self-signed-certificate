package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
)

func runServer(prog string, args []string) error {
	var addr string
	var insecure bool
	var certFile string
	var keyFile string

	fg := flag.NewFlagSet(prog, flag.ExitOnError)
	fg.StringVar(&addr, "addr", ":4433", "Server address")
	fg.BoolVar(&insecure, "insecure", false, "Insecure connection")
	fg.StringVar(&certFile, "certfile", "./server.crt", "Certificate file")
	fg.StringVar(&keyFile, "keyfile", "./server.key", "Key file")
	fg.Parse(args)

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		log.Printf("%v %v\n", req.Method, req.URL.Path)
		fmt.Fprintf(w, "Prouldly served with Go and HTTPS!\n")
	})

	srv := &http.Server{
		Addr:    addr,
		Handler: mux,
	}

	var err error
	if insecure == false {
		log.Printf("Listening on https://%v\n", addr)
		srv.TLSConfig = &tls.Config{
			MinVersion: tls.VersionTLS12,
		}
		err = srv.ListenAndServeTLS(certFile, keyFile)
	} else {
		log.Printf("Listening on http://%v\n", addr)
		err = srv.ListenAndServe()
	}

	if err != nil {
		log.Fatalf("Fail to start https server: %v", err)
	}

	return nil
}

func runClient(prog string, args []string) error {

	return nil
}

func main() {
	prog := os.Args[0]
	if len(os.Args) < 2 {
		fmt.Printf("Usage: %v [client|server]\n", prog)
		return
	}
	mode := os.Args[1]
	args := os.Args[2:]
	//fmt.Printf("mode: %v, args: %v\n", mode, args)

	subProg := fmt.Sprintf("%v %v", prog, mode)
	if mode == "server" {
		runServer(subProg, args)
	} else {
		runClient(subProg, args)
	}

	// Server
	/*
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
	*/

	// Client
	/*	addr := flag.String("addr", "localhost:4001", "HTTPS Server address")
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

	*/

}
